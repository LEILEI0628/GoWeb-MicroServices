# 基础镜像
FROM ubuntu:20.04

COPY goweb /app/goweb

# 设置工作目录
WORKDIR /app

ENTRYPOINT ["/app/goweb"]

# Build命令：docker build -t goweb:v0.0.1 .