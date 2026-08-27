PORT ?= 18080

## 拉取/整理依赖
tidy:
	go mod tidy

## 本地编译（CGO 关闭，纯 Go SQLite）
build:
	CGO_ENABLED=0 go build -o bin/server .

## 本地运行（数据写入 ./data/grade5.db）
run: build
	DB_DSN=data/grade5.db ./bin/server

## 静态检查
vet:
	go vet ./...

## 容器构建
docker-build:
	docker compose build

## 部署并启动（构建 + 后台运行）
deploy:
	docker compose up -d --build

## 查看日志
logs:
	docker compose logs -f

## 停止并移除容器（保留数据卷）
down:
	docker compose down

.PHONY: tidy build run vet docker-build deploy logs down
