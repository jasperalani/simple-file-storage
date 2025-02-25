FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .

# Install any dependencies your API needs
RUN go mod download
# Build with static linking to avoid glibc dependency issues
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/api-server /app/
RUN mkdir -p /app/uploads

EXPOSE 8080
CMD ["/app/api-server"]