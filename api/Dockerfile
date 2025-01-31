FROM golang:1.23.5-alpine AS builder
WORKDIR /app
RUN adduser -D apiuser
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/api-bin cmd/api/main.go

FROM golang:1.23.5-alpine
COPY --from=builder /app/api-bin /app/api-bin
COPY --from=builder /etc/passwd /etc/passwd
USER apiuser
EXPOSE 8080
CMD ["/app/api-bin"]
