FROM alpine:latest

ENV DB_URI mongodb://localhost:27017

RUN adduser -D emailservice && \
    apk --no-cache add ca-certificates

USER emailservice

WORKDIR /home/emailservice

COPY ./bin/emailservice .
COPY ./assets/email ./email

CMD ["./emailservice"]
