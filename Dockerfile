# syntax=docker/dockerfile:1
FROM golang:1.19 AS builder
WORKDIR /go/src/github.com/mig-elgt/chatws/
ADD . .
RUN go test --cover -v ./...
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/chatws/main.go

FROM alpine:latest  
RUN apk --no-cache --update add ca-certificates
COPY --from=builder /bin/app /usr/local/bin/app
RUN chmod +x /usr/local/bin/app
CMD ["app"]
