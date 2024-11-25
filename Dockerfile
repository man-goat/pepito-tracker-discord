from golang:alpine as builder
workdir /blarg
COPY * .
RUN CGO_ENABLED=0 go build -o /ppbot

from scratch as final
WORKDIR /app
COPY --from=builder /ppbot ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY config.json ./
CMD [ "/app/ppbot" ]
