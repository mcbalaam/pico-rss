FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app/pico-rss .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates && adduser -D -H appuser
COPY --from=build /app/pico-rss /app/pico-rss
WORKDIR /app
ENV RSS_TARGET_DIR=/data
USER appuser
ENTRYPOINT ["/app/pico-rss"]
