# 仓储管理系统 - 一键部署脚本 (Windows PowerShell)

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "仓储管理系统 - 自动部署脚本" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# 检查Docker
Write-Host "[1/6] 检查依赖..." -ForegroundColor Yellow
try {
    docker --version | Out-Null
    docker-compose --version | Out-Null
    Write-Host "✓ Docker和Docker Compose已安装" -ForegroundColor Green
}
catch {
    Write-Host "错误: Docker未安装或未启动" -ForegroundColor Red
    exit 1
}

# 停止现有容器
Write-Host "[2/6] 停止现有容器..." -ForegroundColor Yellow
docker-compose down 2>$null
Write-Host "✓ 现有容器已停止" -ForegroundColor Green

# 拉取最新镜像
Write-Host "[3/6] 拉取最新Docker镜像..." -ForegroundColor Yellow
docker-compose pull
Write-Host "✓ 镜像拉取完成" -ForegroundColor Green

# 询问是否重置数据库
Write-Host ""
$reset = Read-Host "是否重置数据库（删除所有数据）？[y/N]"
if ($reset -eq 'y' -or $reset -eq 'Y') {
    Write-Host "[4/6] 删除数据卷..." -ForegroundColor Yellow
    docker-compose down -v
    Write-Host "✓ 数据卷已删除，数据库将重新初始化" -ForegroundColor Green
}
else {
    Write-Host "✓ 保留现有数据" -ForegroundColor Green
}

# 启动服务
Write-Host "[5/6] 启动服务..." -ForegroundColor Yellow
docker-compose up -d
Write-Host "✓ 服务启动中..." -ForegroundColor Green

# 等待服务就绪
Write-Host "[6/6] 等待服务就绪..." -ForegroundColor Yellow
Write-Host "等待MySQL初始化（约30秒）..."
Start-Sleep -Seconds 30

# 检查服务状态
Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "服务状态：" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
docker-compose ps

# 显示访问信息
Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "部署完成！" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "访问信息："
Write-Host "  前端页面: http://localhost:8090"
Write-Host "  后端API:  http://localhost:8091"
Write-Host ""
Write-Host "默认账户："
Write-Host "  超级管理员: admin / admin123"
Write-Host "  仓库管理员: warehouse / admin123"
Write-Host ""
Write-Host "查看日志："
Write-Host "  docker-compose logs -f"
Write-Host ""
Write-Host "停止服务："
Write-Host "  docker-compose down"
Write-Host ""
