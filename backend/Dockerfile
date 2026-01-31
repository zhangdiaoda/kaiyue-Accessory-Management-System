# 直接运行阶段 (二进制已在宿主机编译完成)
FROM alpine:latest

WORKDIR /app

# 安装基础依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 复制宿主机编译好的二进制文件
# 注意：需要在宿主机先执行 set GOOS=linux && go build -o warehouse-server-linux ./cmd/server
COPY warehouse-server-linux ./warehouse-server
COPY config/ ./config/

RUN chmod +x ./warehouse-server

EXPOSE 8080

CMD ["./warehouse-server"]
