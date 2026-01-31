#!/bin/bash
# 仓储管理系统 - 一键部署脚本

set -e

echo "=================================="
echo "仓储管理系统 - 自动部署脚本"
echo "=================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查Docker和Docker Compose
echo -e "${YELLOW}[1/6]${NC} 检查依赖..."
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker未安装，请先安装Docker${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: Docker Compose未安装${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker和Docker Compose已安装${NC}"

# 停止现有容器
echo -e "${YELLOW}[2/6]${NC} 停止现有容器..."
docker-compose down 2>/dev/null || true
echo -e "${GREEN}✓ 现有容器已停止${NC}"

# 拉取最新镜像
echo -e "${YELLOW}[3/6]${NC} 拉取最新Docker镜像..."
docker-compose pull
echo -e "${GREEN}✓ 镜像拉取完成${NC}"

# 询问是否重置数据库
echo ""
read -p "是否重置数据库（删除所有数据）？[y/N]: " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}[4/6]${NC} 删除数据卷..."
    docker-compose down -v
    echo -e "${GREEN}✓ 数据卷已删除，数据库将重新初始化${NC}"
else
    echo -e "${GREEN}✓ 保留现有数据${NC}"
fi

# 启动服务
echo -e "${YELLOW}[5/6]${NC} 启动服务..."
docker-compose up -d
echo -e "${GREEN}✓ 服务启动中...${NC}"

# 等待服务就绪
echo -e "${YELLOW}[6/6]${NC} 等待服务就绪..."
echo "等待MySQL初始化（约30秒）..."
sleep 30

# 检查服务状态
echo ""
echo "=================================="
echo "服务状态："
echo "=================================="
docker-compose ps

# 显示访问信息
echo ""
echo "=================================="
echo -e "${GREEN}部署完成！${NC}"
echo "=================================="
echo ""
echo "访问信息："
echo "  前端页面: http://$(hostname -I | awk '{print $1}'):8090"
echo "  后端API:  http://$(hostname -I | awk '{print $1}'):8091"
echo ""
echo "默认账户："
echo "  超级管理员: admin / admin123"
echo "  仓库管理员: warehouse / admin123"
echo ""
echo "查看日志："
echo "  docker-compose logs -f"
echo ""
echo "停止服务："
echo "  docker-compose down"
echo ""
