FROM alpine:latest

ENV API_PORT 8000
ENV DB_URI mongodb://localhost:27017

RUN apk --no-cache add ca-certificates

RUN adduser -D cedrusservice
USER cedrusservice

WORKDIR /home/cedrusservice

COPY ./bin/cedrusservice .
COPY ./assets/static ./static

EXPOSE $API_PORT

CMD ["./cedrusservice"]
