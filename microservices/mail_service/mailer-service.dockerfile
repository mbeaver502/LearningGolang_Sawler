# build a tiny Docker image that will execute our program
FROM alpine:latest

RUN mkdir /app

COPY mailerApp /app
COPY templates /templates

CMD [ "/app/mailerApp" ]