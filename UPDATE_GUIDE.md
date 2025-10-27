# Docker Compose 服务更新指南

## 概述

当使用 `docker-compose up -d` 启动的服务中，如果 `latest` 标签的官方镜像更新了，您需要按照以下步骤来更新服务。

## 快速更新方法

1. 将 `update-services.sh` 脚本上传到服务器上的 `docker-compose.yml` 文件所在目录
2. 给脚本添加执行权限：`chmod +x update-services.sh`
3. 运行脚本：`./update-services.sh`

## 手动更新步骤

如果您不想使用脚本，也可以手动执行以下命令：

### 1. 停止当前运行的服务
```bash
docker-compose down
```

### 2. 拉取最新的镜像
```bash
docker-compose pull
```

### 3. 重新创建并启动服务
```bash
docker-compose up -d
```

### 4. 验证服务状态
```bash
docker-compose ps
```

### 5. 查看服务日志（可选）
```bash
docker-compose logs --tail=20
```

## 注意事项

1. **数据持久化**：您的 `docker-compose.yml` 已经正确配置了数据卷（`mysql_data`、`redis_data` 等），所以更新容器不会丢失数据。

2. **停机时间**：更新过程中会有短暂的停机时间，通常在几分钟内完成。

3. **备份建议**：在更新前，建议备份重要数据：
   ```bash
   # 备份数据库
   docker exec mysql mysqldump -u root -p123456 new-api > backup.sql
   ```

4. **回滚方案**：如果新版本有问题，可以回滚到之前的镜像：
   ```bash
   # 查看本地镜像历史
   docker images calciumion/new-api
   
   # 使用特定标签启动（如果有的话）
   # 修改 docker-compose.yml 中的 image: calciumion/new-api:latest 为具体版本号
   docker-compose up -d
   ```

## 自动化建议

如果您希望自动化这个过程，可以考虑：

1. **设置定时任务**：使用 cron 定期检查并更新
2. **监控工具**：使用 Watchtower 等工具自动监控和更新容器
3. **CI/CD**：集成到持续部署流程中

## 使用 Watchtower 自动更新（可选）

如果您希望自动监控和更新容器，可以添加 Watchtower 服务：

```yaml
# 在 docker-compose.yml 中添加
watchtower:
  image: containrrr/watchtower
  container_name: watchtower
  restart: always
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  command: --interval 3600 --cleanup  # 每小时检查一次
```

这样 Watchtower 会每小时检查一次镜像更新，并自动更新容器。