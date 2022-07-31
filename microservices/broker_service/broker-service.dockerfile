# get the base Go image
#FROM golang:1.18-alpine as builder
#
## make a directory for our app
#RUN mkdir /app
#
## copy everything from the current directory to /app in Docker
#COPY . /app
#
## set the working directory
#WORKDIR /app
#
## build our Go code -- name the binary brokerApp
#RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
#
## make our binary executable, just in case
#RUN chmod +x /app/brokerApp

# build a tiny Docker image that will execute our program
FROM alpine:latest

RUN mkdir /app

#COPY --from=builder /app/brokerApp /app
COPY brokerApp /app

CMD [ "/app/brokerApp" ]