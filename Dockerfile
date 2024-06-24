FROM golang:1.22-alpine3.19 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM alpine:3.19
WORKDIR /

COPY --from=builder /app/config.yaml /config.yaml

EXPOSE 8281
CMD ["/app/server"]