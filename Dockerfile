FROM golang:1.17 AS build

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o s3-proxy


# real image
FROM scratch
WORKDIR /app
ENV PORT 8080
EXPOSE $PORT
CMD ["/app/s3-proxy"]
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert.pem
COPY --from=build /build/s3-proxy /app/s3-proxy
