#!/bin/bash

echo "=== Docker Compose 服务更新脚本 ==="
echo "此脚本将更新所有使用 latest 标签的镜像"
echo ""

# 检查 docker-compose.yml 是否存在
if [ ! -f "docker-compose.yml" ]; then
    echo "错误: 在当前目录中未找到 docker-compose.yml 文件"
    exit 1
fi

echo "步骤 1: 停止当前运行的服务..."
docker-compose down

echo ""
echo "步骤 2: 拉取最新的镜像..."
docker-compose pull

echo ""
echo "步骤 3: 重新创建并启动服务..."
docker-compose up -d

echo ""
echo "步骤 4: 验证服务状态..."
docker-compose ps

echo ""
echo "步骤 5: 显示服务日志 (最近20行)..."
docker-compose logs --tail=20

echo ""
echo "=== 更新完成 ==="
echo "您可以通过访问 http://localhost:3000 来验证服务是否正常运行"