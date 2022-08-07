# 使用 alpine 构建 Linux 环境镜像
FROM alpine:latest

# 创建工作目录
RUN mkdir /app

# 将本地构建的可执行二进制文件 拷贝到 /app
COPY authApp /app

# 设置容器启动默认执行命令
CMD ["/app/authApp"]
