# 仓储管理系统 Docker构建脚本

Write-Host "1. 开始构建后端镜像..." -ForegroundColor Cyan
cd backend
docker build -t warehouse-backend:latest .
if ($LASTEXITCODE -ne 0) { Write-Error "后端构建失败"; exit 1 }
cd ..

Write-Host "2. 开始构建前端镜像..." -ForegroundColor Cyan
cd frontend
docker build -t warehouse-frontend:latest .
if ($LASTEXITCODE -ne 0) { Write-Error "前端构建失败"; exit 1 }
cd ..

Write-Host "3. 启动服务..." -ForegroundColor Cyan
docker-compose up -d

Write-Host "✅ 部署完成！" -ForegroundColor Green
Write-Host "访问地址: http://localhost:8090"
