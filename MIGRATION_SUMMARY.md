# Gin-FataMorgana v1 到 v2 迁移总结

## 概述

本项目已从 Gin-FataMorgana v1 成功迁移到 v2 版本，主要修改包括端口号、数据库名称、容器名称、服务名称和 API 版本。

## 主要修改内容

### 1. 端口号修改

- **原端口**: 9001
- **新端口**: 9002
- **修改文件**:
  - `config.yaml` - 服务器配置
  - `docker-compose.yml` - 容器端口映射
  - `Dockerfile` - 暴露端口和健康检查
  - `main.go` - 默认端口配置
  - `config/config.go` - 默认端口配置
  - `Makefile` - 健康检查 URL
  - 所有文档文件中的端口引用

### 2. 数据库名称修改

- **原数据库**: `future`
- **新数据库**: `future_v2`
- **修改文件**:
  - `config.yaml` - 数据库配置
  - `docker-compose.yml` - 环境变量
  - `database/migrations/01_create_all_tables.sql` - 数据库使用语句
  - `docker/mysql/init.sql` - 数据库创建脚本
  - `Makefile` - 备份和恢复命令
  - `database/migrations/README.md` - 迁移说明

### 3. 容器名称修改

- **原容器名**: `gin-fataMorgana-app`
- **新容器名**: `gin-fataMorgana-app-v2`
- **修改文件**:
  - `docker-compose.yml` - 容器名称配置

### 4. 服务名称修改

- **原服务名**: `gin-fataMorgana`
- **新服务名**: `gin-fataMorgana-v2`
- **修改文件**:
  - `docker-compose.yml` - 服务名称
  - `docker/nginx/nginx.conf` - 上游服务器配置

### 5. 镜像名称修改

- **原镜像**: `gin-fatamorgana:latest`
- **新镜像**: `gin-fatamorgana-v2:latest`
- **修改文件**:
  - `docker-compose.yml` - 镜像名称

### 6. API 版本修改

- **原版本**: `/api/v1`
- **新版本**: `/api/v2`
- **修改文件**:
  - `main.go` - API 路由组和 Swagger 配置
  - 所有控制器文件中的 Swagger 注释
  - 所有文档文件中的 API 路径
  - `middleware/rate_limit.go` - 限流路径

### 7. 二进制文件名修改

- **原文件名**: `gin-fataMorgana`
- **新文件名**: `gin-fataMorgana-v2`
- **修改文件**:
  - `Dockerfile` - 构建和启动命令

## 修改的文件列表

### 配置文件

- `config.yaml`
- `docker-compose.yml`
- `Dockerfile`
- `go.mod` (无需修改)

### 源代码文件

- `main.go`
- `config/config.go`
- `middleware/rate_limit.go`
- 所有控制器文件中的 Swagger 注释

### 数据库文件

- `database/migrations/01_create_all_tables.sql`
- `docker/mysql/init.sql`
- `database/migrations/README.md`

### 文档文件

- `README.md`
- `README_DATABASE.md`
- `docs/API_DOCUMENTATION.md`
- `docs/PROJECT_OVERVIEW.md`
- `docs/swagger.json`
- `docs/swagger.yaml`
- `CURRENCY_API_README.md`

### 构建和部署文件

- `Makefile`

## 部署说明

### 1. 构建新镜像

```bash
make docker-build
```

### 2. 启动新服务

```bash
make docker-up
```

### 3. 检查服务状态

```bash
make status
```

### 4. 健康检查

```bash
curl http://localhost:9002/health
```

## 注意事项

1. **数据库迁移**: 新版本使用 `future_v2` 数据库，需要确保该数据库存在
2. **端口冲突**: 确保 9002 端口未被其他服务占用
3. **API 兼容性**: 所有 API 路径已从 `/api/v1` 更改为 `/api/v2`
4. **文档更新**: 所有 API 文档已更新为 v2 版本
5. **容器隔离**: 新版本使用独立的容器名称，可以与 v1 版本并行运行

## 回滚方案

如果需要回滚到 v1 版本：

1. 停止 v2 服务: `docker-compose down`
2. 恢复配置文件到 v1 版本
3. 重新构建和启动 v1 服务

## 验证清单

- [ ] 服务在 9002 端口正常启动
- [ ] 数据库连接正常
- [ ] API 接口响应正常
- [ ] Swagger 文档可访问
- [ ] 健康检查通过
- [ ] 所有功能测试通过
