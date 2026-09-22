FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app/go-rss .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates && adduser -D -H appuser
COPY --from=build /app/go-rss /app/go-rss
WORKDIR /app
USER appuser
ENTRYPOINT ["/app/go-rss"]
