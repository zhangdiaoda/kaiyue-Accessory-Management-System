# 仓储管理系统 Android APK项目生成脚本

Write-Host "1. 安装 Capacitor 依赖..." -ForegroundColor Cyan
cd frontend
npm install @capacitor/core @capacitor/cli @capacitor/android
if ($LASTEXITCODE -ne 0) { Write-Error "依赖安装失败"; exit 1 }

Write-Host "2. 初始化 Capacitor 项目..." -ForegroundColor Cyan
# 如果尚未初始化
if (-not (Test-Path "capacitor.config.json") -and -not (Test-Path "capacitor.config.ts")) {
    npx cap init warehouse com.example.warehouse --web-dir dist
}

Write-Host "3. 添加 Android 平台支持..." -ForegroundColor Cyan
if (-not (Test-Path "android")) {
    npx cap add android
}

Write-Host "4. 构建前端资源..." -ForegroundColor Cyan
npm run build
if ($LASTEXITCODE -ne 0) { Write-Error "前端构建失败"; exit 1 }

Write-Host "5. 同步资源到 Android 项目..." -ForegroundColor Cyan
npx cap sync

Write-Host "✅ Android 项目生成完毕！" -ForegroundColor Green
Write-Host "请执行以下命令打开 Android Studio 进行打包："
Write-Host "cd frontend"
Write-Host "npx cap open android"
