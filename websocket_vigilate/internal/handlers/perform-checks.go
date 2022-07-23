package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tsawler/vigilate/internal/models"
)

const (
	HTTP    = 1
	HTTPS   = 2
	SSLCert = 3
)

type jsonResponse struct {
	Ok            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
}

// ScheduleCheck performs a scheduled check on a host service by id.
func (repo *DBRepo) ScheduledCheck(hostServiceID int) {
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	host, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	newStatus, msg := repo.testServiceForHost(host, hs)

	if newStatus != hs.Status {
		repo.updateHostServiceStatusCount(host, hs, newStatus, msg)
	}
}

func (repo *DBRepo) updateHostServiceStatusCount(h models.Host, hs models.HostService, newStatus string, msg string) {
	// update host service record in database with status and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()

	err := repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	pending, healthy, warning, problem, err := repo.DB.GetAllServiceStatusCounts()
	if err != nil {
		log.Println(err)
		return
	}

	data := make(map[string]string)
	data["healthy_count"] = strconv.Itoa(healthy)
	data["pending_count"] = strconv.Itoa(pending)
	data["problem_count"] = strconv.Itoa(problem)
	data["warning_count"] = strconv.Itoa(warning)

	repo.broadcastMessage("public-channel", "host-service-count-changed", data)

	log.Println("New status is", newStatus, "and message is", msg)
}

func (repo *DBRepo) broadcastMessage(channel, event string, data map[string]string) {
	err := app.WsClient.Trigger(channel, event, data)
	if err != nil {
		log.Println(err)
	}
}

// TestCheck checks the status of a host service.
func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParamFromCtx(r.Context(), "id"))
	oldStatus := chi.URLParam(r, "oldStatus")

	ok := true

	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		ok = false
	}

	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		ok = false
	}

	newStatus, msg := repo.testServiceForHost(h, hs)

	// update host service in database with status and last check
	hs.Status = newStatus
	hs.UpdatedAt = time.Now()
	hs.LastCheck = time.Now()

	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		ok = false
	}

	// broadcast service status changed event

	var resp jsonResponse

	if ok {
		resp = jsonResponse{
			Ok:            ok,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			LastCheck:     time.Now(),
		}
	} else {
		resp.Ok = false
		resp.Message = "Something went wrong."
	}

	out, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (string, string) {
	var msg, newStatus string

	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = testHTTPForHost(h.URL)
	}

	// if the host service status has changed, broadcast to all clients
	if newStatus != hs.Status {
		data := make(map[string]string)

		data["host_id"] = strconv.Itoa(hs.HostID)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		data["host_name"] = h.HostName
		data["service_name"] = hs.Service.ServiceName
		data["icon"] = hs.Service.Icon
		data["status"] = newStatus
		data["message"] = fmt.Sprintf("%s on %s reports %s", hs.Service.ServiceName, h.HostName, newStatus)
		data["last_check"] = time.Now().Format("2006-01-02 15:04:06")

		repo.broadcastMessage("public-channel", "host-service-status-changed", data)

	}

	// TODO - if appropriate, send email or SMS message

	return newStatus, msg
}

func testHTTPForHost(url string) (string, string) {
	url = strings.TrimSuffix(url, "/")

	url = strings.Replace(url, "https://", "http://", -1)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}