#!/bin/bash

# 一键部署脚本
# 使用方法: ./deploy.sh [dev|prod]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 环境变量
ENV=${1:-dev}
PROJECT_NAME="gin-fataMorgana"
DOCKER_COMPOSE_FILE="docker-compose.yml"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查部署依赖..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    # 检查Docker服务状态
    if ! docker info &> /dev/null; then
        log_error "Docker服务未启动，请先启动Docker服务"
        exit 1
    fi
    
    log_success "依赖检查完成"
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    mkdir -p logs
    mkdir -p docker/nginx/ssl
    mkdir -p docker/mysql
    
    log_success "目录创建完成"
}

# 检查配置文件
check_config() {
    log_info "检查配置文件..."
    
    if [ ! -f "config/config.yaml" ]; then
        log_warning "配置文件不存在，创建默认配置..."
        cp config/config.example.yaml config/config.yaml
        log_success "默认配置文件已创建，请根据实际情况修改"
    fi
    
    log_success "配置文件检查完成"
}

# 构建镜像
build_images() {
    log_info "构建Docker镜像..."
    
    docker-compose build --no-cache
    
    log_success "镜像构建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    if [ "$ENV" = "prod" ]; then
        log_info "生产环境部署..."
        docker-compose -f $DOCKER_COMPOSE_FILE up -d
    else
        log_info "开发环境部署..."
        docker-compose -f $DOCKER_COMPOSE_FILE up -d
    fi
    
    log_success "服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."
    
    # 等待MySQL就绪
    log_info "等待MySQL就绪..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if docker-compose exec -T mysql mysqladmin ping -h localhost --silent; then
            log_success "MySQL已就绪"
            break
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "MySQL启动超时"
        exit 1
    fi
    
    # 等待Redis就绪
    log_info "等待Redis就绪..."
    timeout=30
    while [ $timeout -gt 0 ]; do
        if docker-compose exec -T redis redis-cli ping &> /dev/null; then
            log_success "Redis已就绪"
            break
        fi
        sleep 1
        timeout=$((timeout - 1))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "Redis启动超时"
        exit 1
    fi
    
    # 等待应用就绪
    log_info "等待应用就绪..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if curl -f http://localhost:8080/health &> /dev/null; then
            log_success "应用已就绪"
            break
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    
    if [ $timeout -le 0 ]; then
        log_error "应用启动超时"
        exit 1
    fi
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    # 检查应用健康状态
    if curl -f http://localhost:8080/health &> /dev/null; then
        log_success "应用健康检查通过"
    else
        log_error "应用健康检查失败"
        exit 1
    fi
    
    # 检查数据库连接
    if docker-compose exec -T mysql mysqladmin ping -h localhost --silent; then
        log_success "数据库健康检查通过"
    else
        log_error "数据库健康检查失败"
        exit 1
    fi
    
    # 检查Redis连接
    if docker-compose exec -T redis redis-cli ping &> /dev/null; then
        log_success "Redis健康检查通过"
    else
        log_error "Redis健康检查失败"
        exit 1
    fi
}

# 显示服务状态
show_status() {
    log_info "显示服务状态..."
    
    echo ""
    echo "=========================================="
    echo "          服务部署状态"
    echo "=========================================="
    echo "环境: $ENV"
    echo "项目: $PROJECT_NAME"
    echo ""
    
    docker-compose ps
    
    echo ""
    echo "=========================================="
    echo "          访问地址"
    echo "=========================================="
    echo "应用地址: http://localhost:8080"
    echo "Nginx地址: http://localhost:80"
    echo "健康检查: http://localhost:8080/health"
    echo ""
    echo "=========================================="
    echo "          默认账户"
    echo "=========================================="
    echo "管理员邮箱: admin@example.com"
    echo "管理员密码: admin123"
    echo "邀请码: ADMIN1"
    echo ""
    echo "=========================================="
    echo "          常用命令"
    echo "=========================================="
    echo "查看日志: docker-compose logs -f"
    echo "停止服务: docker-compose down"
    echo "重启服务: docker-compose restart"
    echo "更新服务: ./deploy.sh $ENV"
    echo ""
}

# 清理函数
cleanup() {
    log_info "清理临时文件..."
    # 可以在这里添加清理逻辑
}

# 主函数
main() {
    echo -e "${GREEN}==========================================${NC}"
    echo -e "${GREEN}        Gin-FataMorgana 一键部署${NC}"
    echo -e "${GREEN}==========================================${NC}"
    echo ""
    
    # 设置错误处理
    trap cleanup EXIT
    
    # 执行部署步骤
    check_dependencies
    create_directories
    check_config
    build_images
    start_services
    wait_for_services
    health_check
    show_status
    
    echo -e "${GREEN}==========================================${NC}"
    echo -e "${GREEN}           部署完成！${NC}"
    echo -e "${GREEN}==========================================${NC}"
}

# 执行主函数
main "$@" 