# 构建上下文为仓库根目录；GoFrame 同时托管 Vue 构建产物与 API。
FROM node:22-alpine AS client-build
WORKDIR /client
COPY client/package.json client/package-lock.json ./
RUN npm ci
COPY client/ ./
RUN npm run build

FROM golang:1.25-alpine AS build
WORKDIR /src

COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
COPY --from=client-build /client/dist ./client-dist

ARG APP_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X studyplanet/internal/service/studyplanet.AppVersion=${APP_VERSION}" -o /out/server .

# Alpine 自带 busybox wget，可供 Compose healthcheck 使用。
FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/server /app/server
COPY --from=client-build /client/dist /app/client
COPY server/manifest/config /app/manifest/config
COPY VERSION /app/VERSION
RUN mkdir -p /app/data
ENV DB_DSN=/app/data/studyplanet.db
ENV SERVER_PORT=8080
EXPOSE 8080
STOPSIGNAL SIGTERM
CMD ["/app/server"]
