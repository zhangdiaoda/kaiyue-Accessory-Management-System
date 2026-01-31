# 构建阶段
FROM node:20-alpine AS build
WORKDIR /app

# 复制package文件并安装依赖
COPY package*.json ./
RUN npm ci --registry=https://registry.npmmirror.com

# 复制源代码并构建
COPY . .
RUN npm run build

# 运行阶段
FROM nginx:1.25-alpine

# 复制自定义 nginx 配置以实现 API 代理
COPY nginx.conf /etc/nginx/nginx.conf

# 复制构建产物到 nginx 目录
COPY --from=build /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
