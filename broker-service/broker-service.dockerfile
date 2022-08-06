# go 基础镜像
FROM golang:1.18-alpine as builder

# 创建 app 工作目录
RUN mkdir /app

# 拷贝当前目录内容到app
COPY . /app

# 设置工作目录
WORKDIR /app

# 构建 Go 程序（CGO_ENABLED=0 禁用CGO，项目没用到，减少构建体积）
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# 添加可执行权限
RUN chmod +x /app/brokerApp

# 使用 alpine 构建 Linux 环境镜像
FROM alpine:latest

# 创建工作目录
RUN mkdir /app

# 将 builder 拷贝到 /app
COPY --from=builder /app/brokerApp /app

# 设置容器启动默认执行命令
CMD ["/app/brokerApp"]
