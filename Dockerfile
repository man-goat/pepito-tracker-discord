from golang:alpine as builder
workdir /blarg
COPY *.go go.mod go.sum .
RUN CGO_ENABLED=0 go build -o /ppbot

from scratch as final
ARG CONFIG_FILE=./config.json
WORKDIR /app
COPY --from=builder /ppbot ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ${CONFIG_FILE} /app/config.json
CMD [ "/app/ppbot", "/app/config.json"]
