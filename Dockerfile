# 多阶段构建：编译与运行分离，CGO_ENABLED=0 便于交叉编译（SQLite 用纯 Go 驱动）
FROM golang:1.25-alpine AS build
WORKDIR /src
# 先拉取依赖，利用层缓存
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server .

# 运行镜像：alpine 自带 busybox wget，可供 healthcheck 使用
FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/server /app/server
COPY manifest/config /app/manifest/config
RUN mkdir -p /app/data
ENV DB_DSN=/app/data/studyplanet.db
ENV SERVER_PORT=8080
EXPOSE 8080
# 优雅退出
STOPSIGNAL SIGTERM
CMD ["/app/server"]
