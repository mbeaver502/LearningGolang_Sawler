# build a tiny Docker image that will execute our program
FROM alpine:latest

RUN mkdir /app

COPY loggerApp /app

CMD [ "/app/loggerApp" ]