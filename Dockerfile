FROM golang:1.12-alpine3.9 AS build
RUN apk add upx git 

COPY go.mod go.sum main.go /src/
WORKDIR /src

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/tibber-exporter 
RUN upx /bin/tibber-exporter -9

FROM scratch
LABEL maintainer="Gunnar Inge G. Sortland <gunnaringe@gmail.com>"
EXPOSE 9501
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/tibber-exporter /tibber-exporter
ENTRYPOINT ["/tibber-exporter"]
CMD ["--help"]
