# 部署文档

## 目录
- [环境要求](#环境要求)
- [Docker部署（推荐）](#docker部署推荐)
- [手动部署](#手动部署)
- [常见问题](#常见问题)

---

## 环境要求

### Docker部署
- Docker 20.10+
- Docker Compose 2.0+
- 2GB+ 可用内存
- 10GB+ 可用磁盘空间

### 手动部署
- Go 1.21+
- Node.js 20+
- MySQL 8.0+
- Redis 7+
- Nginx 1.25+

---

## Docker部署（推荐）

### 1. 准备部署文件

```bash
# 进入项目目录
cd c:/Users/admin/Desktop/仓储管理

# 查看目录结构
ls
# 应该包含: backend/, frontend/, docker-compose.yml, .env等
```

### 2. 配置环境变量

编辑`.env`文件，设置密码：

```env
# MySQL配置
MYSQL_ROOT_PASSWORD=your_strong_root_password
MYSQL_PASSWORD=your_warehouse_password

# Redis配置  
REDIS_PASSWORD=your_redis_password

# JWT密钥（建议64位随机字符串）
JWT_SECRET=your_random_jwt_secret_key_here
```

### 3. 启动服务

```bash
# 启动所有服务
docker-compose -p warehouse up -d

# 查看服务状态
docker-compose -p warehouse ps

# 查看日志
docker-compose -p warehouse logs -f
```

### 4. 验证部署

**检查容器状态：**
```bash
docker-compose -p warehouse ps
```

所有服务应显示"Up"状态：
- warehouse-mysql
- warehouse-redis
- warehouse-backend
- warehouse-frontend

**访问系统：**
- 前端管理后台: http://localhost
- 后端API: http://localhost:8080
- 默认账号: admin / admin123

### 5. 停止服务

```bash
# 停止所有服务
docker-compose -p warehouse stop

# 停止并删除容器
docker-compose -p warehouse down

# 停止并删除容器+数据卷（⚠️ 会删除数据库数据）
docker-compose -p warehouse down -v
```

---

## 手动部署

### 1. 部署MySQL

```bash
# 使用Docker部署MySQL
docker run -d \
  --name warehouse-mysql \
  -e MYSQL_ROOT_PASSWORD=rootpassword \
  -e MYSQL_DATABASE=warehouse \
  -e MYSQL_USER=warehouse \
  -e MYSQL_PASSWORD=Warehouse@2026 \
  -p 3306:3306 \
  mysql:8.0

# 导入初始化SQL
docker exec -i warehouse-mysql mysql -uroot -prootpassword warehouse < init-db/01-schema.sql
```

### 2. 部署Redis

```bash
docker run -d \
  --name warehouse-redis \
  -p 6379:6379 \
  redis:7-alpine redis-server --requirepass Redis@2026
```

### 3. 部署后端

```bash
cd backend

# 下载依赖
go mod tidy

# 编译
go build -o warehouse-server cmd/server/main.go

# 运行（或使用systemd服务）
./warehouse-server
```

**systemd服务配置** (`/etc/systemd/system/warehouse-backend.service`):
```ini
[Unit]
Description=Warehouse Management Backend
After=network.target mysql.service redis.service

[Service]
Type=simple
User=warehouse
WorkingDirectory=/opt/warehouse/backend
ExecStart=/opt/warehouse/backend/warehouse-server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 4. 部署前端

```bash
cd frontend

# 安装依赖
npm install

# 构建生产版本
npm run build

# 输出到 dist/ 目录
```

**Nginx配置** (`/etc/nginx/sites-available/warehouse`):
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/warehouse/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

启动Nginx:
```bash
sudo ln -s /etc/nginx/sites-available/warehouse /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 生产环境优化

### 1. 数据库优化

**MySQL配置优化** (`my.cnf`):
```ini
[mysqld]
max_connections=500
innodb_buffer_pool_size=2G
innodb_log_file_size=512M
slow_query_log=1
slow_query_log_file=/var/log/mysql/slow.log
long_query_time=2
```

### 2. 后端优化

**编译优化**:
```bash
# 编译时优化
CGO_ENABLED=0 go build -ldflags="-s -w" -o warehouse-server cmd/server/main.go

# 使用upx压缩（可选）
upx --best warehouse-server
```

### 3. 前端优化

**Vite构建优化** (`vite.config.js`):
```javascript
export default defineConfig({
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'vue-vendor': ['vue', 'vue-router', 'pinia']
        }
      }
    }
  }
})
```

---

## 备份与恢复

### 数据库备份

```bash
# 备份
docker exec warehouse-mysql mysqldump -uwarehouse -pWarehouse@2026 warehouse > backup_$(date +%Y%m%d).sql

# 恢复
docker exec -i warehouse-mysql mysql -uwarehouse -pWarehouse@2026 warehouse < backup_20260125.sql
```

### Docker卷备份

```bash
# 备份MySQL数据卷
docker run --rm -v warehouse_mysql-data:/data -v $(pwd):/backup alpine tar czf /backup/mysql-data.tar.gz /data

# 恢复
docker run --rm -v warehouse_mysql-data:/data -v $(pwd):/backup alpine tar xzf /backup/mysql-data.tar.gz -C /
```

---

## 监控与日志

### 查看日志

```bash
# 后端日志
docker-compose -p warehouse logs -f backend

# 前端日志
docker-compose -p warehouse logs -f frontend

# MySQL日志
docker-compose -p warehouse logs -f mysql
```

### 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看后端进程
docker exec warehouse-backend ps aux
```

---

## 常见问题

### Q: 容器启动失败？
**A**: 检查端口占用
```bash
netstat -ano | findstr :8080
netstat -ano | findstr :3306
netstat -ano | findstr :6379
```

### Q: 数据库连接超时？
**A**: 等待MySQL健康检查完成（约30秒），查看日志：
```bash
docker-compose -p warehouse logs mysql
```

### Q: 前端访问404？
**A**: 确认Nginx配置正确，检查前端构建产物：
```bash
ls frontend/dist
```

### Q: JWT Token过期？
**A**: 修改后端配置 `config/config.yaml`:
```yaml
jwt:
  expiration: 86400  # 24小时
```

### Q: 如何修改端口？
**A**: 编辑 `docker-compose.yml`:
```yaml
frontend:
  ports:
    - "8888:80"  # 修改为8888端口
```

---

## 安全建议

1. **修改默认密码** - 首次登录后立即修改admin密码
2. **使用强密码** - MySQL、Redis密码至少16位
3. **启用HTTPS** - 生产环境使用SSL证书
4. **防火墙配置** - 仅开放必要端口(80/443)
5. **定期备份** - 每日备份数据库
6. **更新依赖** - 定期更新Docker镜像和依赖包

---

## 技术支持

如有问题，请查看：
- [项目README](README.md)
- [GitHub Issues](issues)
- 或联系管理员
