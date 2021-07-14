FROM golang:1.16.5-alpine as build

RUN apk add --no-cache ca-certificates  && \
    update-ca-certificates

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -a -ldflags="-s -w" -installsuffix cgo -o app

FROM scratch

COPY --from=build /app/app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]
