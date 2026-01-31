# Docker Hub 推送脚本

# 1. 设置变量
$username = "zhangyu991014"
$repoName = "kaiyuewms"

# 2. 登录 Docker Hub (如果尚未登录)
# docker login 
# (注意：为了安全，建议手动在终端运行 docker login，此处注释掉以免卡住)

# 3. 标记镜像
$backendImage = "$username/$repoName-backend:latest"
$frontendImage = "$username/$repoName-frontend:latest"

Write-Host "正在标记镜像..." -ForegroundColor Cyan
docker tag warehouse-backend:latest $backendImage
docker tag warehouse-frontend:latest $frontendImage

# 4. 推送镜像
Write-Host "正在推送到 Docker Hub (这可能需要几分钟)..." -ForegroundColor Cyan

Write-Host "推送到 $backendImage ..."
docker push $backendImage

Write-Host "推送到 $frontendImage ..."
docker push $frontendImage

Write-Host "✅ 推送完成！" -ForegroundColor Green
Write-Host "您现在可以在任何服务器上使用以下命令拉取并运行："
Write-Host "docker run -d -p 8090:80 $frontendImage"
