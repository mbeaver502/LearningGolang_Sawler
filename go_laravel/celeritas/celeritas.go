package celeritas

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

// Celeritas is the main type.
type Celeritas struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
}

type config struct {
	port     string
	renderer string
}

// New sets up a new Celeritas value.
func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers",
			"migrations",
			"views",
			"data",
			"public",
			"tmp",
			"logs",
			"middleware",
		},
	}

	err := c.Init(pathConfig)
	if err != nil {
		return err
	}

	err = c.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	err = godotenv.Load(fmt.Sprintf("%s/.env", rootPath))
	if err != nil {
		return err
	}

	infoLog, errorLog := c.startLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errorLog

	c.Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		c.Debug = true
	}

	c.Version = version
	c.RootPath = rootPath
	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	return nil
}

// Init initializes expected folder structure.
func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		// create folder if it does not exist
		err := c.CreateDirIfNotExists(fmt.Sprintf("%s/%s", root, path))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Celeritas) checkDotEnv(path string) error {
	err := c.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
