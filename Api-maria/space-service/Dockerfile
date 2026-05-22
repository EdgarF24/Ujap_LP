# ---- Builder Stage ----
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod only — go.sum will be generated during build
COPY go.mod ./

# Download all dependencies listed in go.mod (creates go.sum)
RUN go mod download

# Copy all source code
COPY . .

# Install swag and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o ./docs

# Build with -mod=mod so Go can update go.sum if needed
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -a -installsuffix cgo -o space-service ./cmd/main.go

# ---- Final Stage ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /root/

COPY --from=builder /app/space-service .

# Install wait utility for DB readiness
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

EXPOSE 8002

CMD /wait && ./space-service
