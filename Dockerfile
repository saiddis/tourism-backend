FROM golang:1.25-alpine3.22 AS builder

COPY . /app
WORKDIR /app

# Toggle CGO based on your app requirement. CGO_ENABLED=1 for enabling CGO
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/tourism-backend ./cmd/main.go
# Use below if using vendor
# RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-s -w -extldflags "-static"' -o /app/appbin *.go

FROM alpine:3.22
# LABEL MAINTAINER="Author <author@example.com>"

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /app /app

WORKDIR /app

EXPOSE 8080
EXPOSE 8443

CMD ["./tourism-backend"]
