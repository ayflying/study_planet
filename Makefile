PORT ?= 18180

## 拉取/整理服务端依赖
tidy:
	cd server && go mod tidy

## 本地编译服务端（CGO 关闭，纯 Go SQLite）
build-server:
	mkdir -p bin
	cd server && CGO_ENABLED=0 go build -o ../bin/server .

## 本地运行服务端 API
run-server: build-server
	DB_DSN=data/grade5.db ./bin/server

## 静态检查服务端
vet:
	cd server && go vet ./...

## 本地构建双容器
build-images:
	docker compose build

## 启动双容器（客户端入口 + 服务端 API）
deploy:
	docker compose up -d --build

## 查看服务日志
logs:
	docker compose logs -f

## 停止并移除容器（保留数据卷）
down:
	docker compose down

.PHONY: tidy build-server run-server vet build-images deploy logs down
