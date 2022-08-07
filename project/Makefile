## 定义基础变量
FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

## up: 在后台启动所有容器
up:
	@echo "Starting docker containers..."
	docker-compose up -d
	@echo "Docker containers SUCCESS"

## up_build: 停止 docker-compose, 重新构建
up_build: build_broker build_auth
	@echo "Stopping docker containers"
	docker-compose down
	@echo "Building images..."
	docker-compose up --build -d
	@echo "Docker images build SUCCESS"

## down: 停止 docker-compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done"

## build_broker: 构建 Go 可执行文件
build_broker:
	@echo "Building broker ..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "SUCCESS"

## build_auth: 构建 Go 可执行文件
build_auth:
	@echo "Building auth ..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "SUCCESS"


## build_front: 构建前端可执行文件
build_front:
	@echo "Building front ..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "SUCCESS"

## start: 启动前端服务
start: build_front
	cd ../front-end && ./${FRONT_END_BINARY}

## stop: 停止前端服务
stop:
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Done"
