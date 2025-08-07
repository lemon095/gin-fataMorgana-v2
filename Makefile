# Makefile for Gin-FataMorgana
.PHONY: help build run test clean deploy dev prod stop logs status

# 变量定义
PROJECT_NAME := gin-fataMorgana
DOCKER_COMPOSE := docker-compose
GO := go

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help: ## 显示帮助信息
	@echo "Gin-FataMorgana 项目管理工具"
	@echo ""
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 开发相关
build: ## 构建项目
	@echo "构建项目..."
	$(GO) build -o $(PROJECT_NAME) main.go

run: ## 运行项目（开发模式）
	@echo "运行项目..."
	$(GO) run main.go

test: ## 运行测试
	@echo "运行测试..."
	$(GO) test ./...

clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -f $(PROJECT_NAME)
	rm -rf logs/*

# Docker相关
docker-build: ## 构建Docker镜像
	@echo "构建Docker镜像..."
	$(DOCKER_COMPOSE) build --no-cache

docker-up: ## 启动Docker服务
	@echo "启动Docker服务..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## 停止Docker服务
	@echo "停止Docker服务..."
	$(DOCKER_COMPOSE) down

docker-restart: ## 重启Docker服务
	@echo "重启Docker服务..."
	$(DOCKER_COMPOSE) restart

docker-logs: ## 查看Docker日志
	@echo "查看Docker日志..."
	$(DOCKER_COMPOSE) logs -f

docker-status: ## 查看Docker服务状态
	@echo "查看Docker服务状态..."
	$(DOCKER_COMPOSE) ps

# 部署相关
deploy: ## 一键部署（开发环境）
	@echo "开始部署..."
	@chmod +x deploy.sh
	./deploy.sh dev

deploy-prod: ## 一键部署（生产环境）
	@echo "开始生产环境部署..."
	@chmod +x deploy.sh
	./deploy.sh prod

# 数据库相关
db-migrate: ## 数据库迁移
	@echo "执行数据库迁移..."
	$(GO) run cmd/migrate/main.go

db-seed: ## 数据库种子数据
	@echo "插入种子数据..."
	@chmod +x init_admin.sh
	./init_admin.sh

db-check-index: ## 检测并创建缺失的索引
	@echo "检测并创建缺失的索引..."
	$(GO) run cmd/migrate/main.go -check-index

db-show-index: ## 显示当前数据库的所有索引
	@echo "显示当前数据库索引..."
	$(GO) run cmd/migrate/main.go -show-index

# 健康检查
health: ## 健康检查
	@echo "执行健康检查..."
	curl -f http://localhost:9002/health || echo "健康检查失败"

# 开发工具
dev-setup: ## 开发环境设置
	@echo "设置开发环境..."
	go mod tidy
	go mod download
	@chmod +x *.sh

# 生产工具
prod-setup: ## 生产环境设置
	@echo "设置生产环境..."
	mkdir -p logs
	mkdir -p docker/nginx/ssl
	mkdir -p docker/mysql
	cp config/config.example.yaml config/config.yaml

# 监控相关
monitor: ## 监控服务状态
	@echo "监控服务状态..."
	@echo "应用状态:"
	curl -s http://localhost:9002/health | jq '.' 2>/dev/null || echo "应用未响应"
	@echo ""
	@echo "Docker服务状态:"
	$(DOCKER_COMPOSE) ps

# 备份相关
backup: ## 备份数据
	@echo "备份数据..."
	@mkdir -p backups
	$(DOCKER_COMPOSE) exec mysql mysqldump -u root -proot123456 future_v2 > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

# 恢复相关
restore: ## 恢复数据（需要指定备份文件）
	@echo "恢复数据..."
	@if [ -z "$(file)" ]; then echo "请指定备份文件: make restore file=backups/backup_20240101_120000.sql"; exit 1; fi
	$(DOCKER_COMPOSE) exec -T mysql mysql -u root -proot123456 future_v2 < $(file)

# 安全相关
security-check: ## 安全检查
	@echo "执行安全检查..."
	@echo "检查配置文件权限..."
	ls -la config/
	@echo "检查日志文件权限..."
	ls -la logs/
	@echo "检查Docker容器安全..."
	$(DOCKER_COMPOSE) ps

# 性能相关
performance-test: ## 性能测试
	@echo "执行性能测试..."
	@echo "性能测试功能已移除，请使用其他工具进行性能测试"

# 清理所有
clean-all: ## 清理所有（包括Docker）
	@echo "清理所有..."
	$(DOCKER_COMPOSE) down -v
	rm -f $(PROJECT_NAME)
	rm -rf logs/*
	rm -rf backups/*
	docker system prune -f

# 快速启动开发环境
dev: ## 快速启动开发环境
	@echo "快速启动开发环境..."
	$(MAKE) dev-setup
	$(MAKE) docker-build
	$(MAKE) docker-up
	@echo "等待服务启动..."
	@sleep 10
	$(MAKE) health

# 快速启动生产环境
prod: ## 快速启动生产环境
	@echo "快速启动生产环境..."
	$(MAKE) prod-setup
	$(MAKE) deploy-prod

# 停止所有服务
stop: ## 停止所有服务
	@echo "停止所有服务..."
	$(DOCKER_COMPOSE) down

# 查看日志
logs: ## 查看日志
	@echo "查看日志..."
	$(DOCKER_COMPOSE) logs -f

# 查看状态
status: ## 查看状态
	@echo "查看状态..."
	$(MAKE) docker-status
	@echo ""
	$(MAKE) health 