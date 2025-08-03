FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

#RUN go test -v ./internal/services/tests/...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o runner cmd/sso/main.go

FROM scratch
COPY --from=builder /app/ .

EXPOSE 4053

CMD ["/runner"]