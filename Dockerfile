FROM golang:1.16.4-alpine3.13 AS build
WORKDIR /src/
COPY . go.* /src/
RUN ls -a
RUN CGO_ENABLED=0 go build -o /target


FROM alpine:3 as certs
RUN apk --no-cache add ca-certificates


FROM scratch
COPY --from=build /target /target
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/target"]