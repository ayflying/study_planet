# 构建上下文为仓库根目录；Go 服务端源码位于 server/。
FROM golang:1.25-alpine AS build
WORKDIR /src

COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./

ARG APP_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X studyplanet/internal/handler.AppVersion=${APP_VERSION}" -o /out/server .

# Alpine 自带 busybox wget，可供 Compose healthcheck 使用。
FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/server /app/server
COPY server/manifest/config /app/manifest/config
COPY VERSION /app/VERSION
RUN mkdir -p /app/data
ENV DB_DSN=/app/data/studyplanet.db
ENV SERVER_PORT=8080
EXPOSE 8080
STOPSIGNAL SIGTERM
CMD ["/app/server"]
