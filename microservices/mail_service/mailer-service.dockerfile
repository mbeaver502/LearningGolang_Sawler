# build a tiny Docker image that will execute our program
FROM alpine:latest

RUN mkdir /app

COPY mailerApp /app

CMD [ "/app/mailerApp" ]