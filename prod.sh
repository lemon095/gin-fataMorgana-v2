#!/bin/bash

# 生产环境管理脚本
# 包含部署、更新、监控、日志管理等功能

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 检查是否为root用户
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "此脚本需要root权限运行"
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    echo "生产环境管理脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  deploy      - 部署应用"
    echo "  update      - 更新应用"
    echo "  restart     - 重启应用"
    echo "  status      - 查看应用状态"
    echo "  logs        - 查看应用日志"
    echo "  clean-logs  - 清理日志文件"
    echo "  clean-docker- 清理Docker空间"
    echo "  disk-space  - 检查磁盘空间"
    echo "  health      - 健康检查"
    echo "  backup      - 备份数据"
    echo "  help        - 显示此帮助信息"
    echo ""
}

# 检查磁盘空间
check_disk_space() {
    log_info "检查磁盘空间使用情况..."
    
    echo "📊 磁盘空间使用情况："
    df -h | grep -E "(Filesystem|/dev/)"
    
    echo ""
    echo "📁 大文件/目录检查："
    du -sh /* 2>/dev/null | sort -hr | head -10
    
    echo ""
    echo "🐳 Docker空间使用情况："
    docker system df
    
    # 检查是否空间不足
    USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
    if [ "$USAGE" -gt 90 ]; then
        log_warn "磁盘空间使用率超过90%，建议清理空间"
        return 1
    else
        log_info "磁盘空间使用正常"
        return 0
    fi
}

# 清理日志文件
clean_logs() {
    log_info "开始清理日志文件..."
    
    # 清理系统日志（保留7天）
    log_info "清理系统日志..."
    journalctl --vacuum-time=7d
    
    # 清理Docker日志
    log_info "清理Docker日志..."
    find /var/lib/docker/containers -name "*.log" -exec truncate -s 0 {} \; 2>/dev/null || true
    
    # 清理应用日志
    log_info "清理应用日志..."
    if [ -d "./logs" ]; then
        find ./logs -name "*.log" -mtime +7 -delete 2>/dev/null || true
        log_info "已清理7天前的应用日志"
    fi
    
    # 清理临时文件
    log_info "清理临时文件..."
    find /tmp -type f -mtime +3 -delete 2>/dev/null || true
    
    log_info "日志清理完成"
}

# 清理Docker空间
clean_docker() {
    log_info "开始清理Docker空间..."
    
    # 显示清理前的状态
    log_info "清理前的Docker使用情况："
    docker system df
    
    # 清理未使用的资源
    log_info "清理未使用的镜像..."
    docker image prune -f
    
    log_info "清理已停止的容器..."
    docker container prune -f
    
    log_info "清理未使用的网络..."
    docker network prune -f
    
    log_info "清理构建缓存..."
    docker builder prune -f
    
    log_info "清理未使用的卷..."
    docker volume prune -f
    
    # 显示清理后的状态
    log_info "清理后的Docker使用情况："
    docker system df
    
    log_info "Docker空间清理完成"
}

# 部署应用
deploy() {
    log_info "开始部署应用..."
    
    # 检查磁盘空间
    if ! check_disk_space; then
        log_warn "磁盘空间不足，建议先清理空间"
        read -p "是否继续部署？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "部署已取消"
            return 1
        fi
    fi
    
    # 停止现有容器
    log_info "停止现有容器..."
    docker compose down || true
    
    # 构建新镜像
    log_info "构建新镜像..."
    docker compose build --no-cache
    
    # 启动服务
    log_info "启动服务..."
    docker compose up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 10
    
    # 健康检查
    health_check
    
    log_info "部署完成"
}

# 更新应用
update() {
    log_info "开始更新应用..."
    
    # 备份当前状态
    log_info "备份当前状态..."
    docker compose ps > backup_status.txt 2>/dev/null || true
    
    # 拉取最新代码
    log_info "拉取最新代码..."
    git fetch origin
    git stash push -m "Auto stash before update $(date)" || true
    git pull origin main || git pull origin master || true
    
    # 如果有本地更改，尝试恢复
    if git stash list | grep -q "Auto stash before update"; then
        log_info "恢复本地更改..."
        git stash pop || true
    fi
    
    # 重新构建和部署
    log_info "重新构建应用..."
    docker compose build --no-cache
    
    log_info "重启服务..."
    docker compose down
    docker compose up -d
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 10
    
    # 健康检查
    health_check
    
    log_info "更新完成"
}

# 重启应用
restart() {
    log_info "重启应用..."
    docker compose restart
    log_info "应用重启完成"
}

# 查看状态
status() {
    log_info "应用状态："
    docker compose ps
    
    echo ""
    log_info "服务健康状态："
    health_check
}

# 查看日志
logs() {
    log_info "显示应用日志..."
    docker compose logs -f --tail=100
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    # 检查容器状态
    if docker-compose ps | grep -q "Up"; then
        log_info "✅ 容器运行正常"
    else
        log_error "❌ 容器运行异常"
        return 1
    fi
    
    # 检查应用健康端点
    if curl -f -s http://localhost:9001/health > /dev/null; then
        log_info "✅ 应用健康检查通过"
    else
        log_warn "⚠️  应用健康检查失败"
        return 1
    fi
    
    # 检查数据库连接
    if docker exec mysql57 mysqladmin ping -h localhost -u root -p'shgytywe!#%65926328' > /dev/null 2>&1; then
        log_info "✅ 数据库连接正常"
    else
        log_warn "⚠️  数据库连接异常"
    fi
    
    # 检查Redis连接
    if docker exec redis redis-cli ping > /dev/null 2>&1; then
        log_info "✅ Redis连接正常"
    else
        log_warn "⚠️  Redis连接异常"
    fi
}

# 备份数据
backup() {
    log_info "开始备份数据..."
    
    BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    # 备份数据库
    log_info "备份数据库..."
    docker exec mysql57 mysqldump -u root -p'shgytywe!#%65926328' future > "$BACKUP_DIR/database.sql" 2>/dev/null || log_warn "数据库备份失败"
    
    # 备份配置文件
    log_info "备份配置文件..."
    cp config.yaml "$BACKUP_DIR/" 2>/dev/null || true
    cp docker-compose.yml "$BACKUP_DIR/" 2>/dev/null || true
    
    # 备份日志
    log_info "备份日志文件..."
    tar -czf "$BACKUP_DIR/logs.tar.gz" logs/ 2>/dev/null || log_warn "日志备份失败"
    
    log_info "备份完成，位置: $BACKUP_DIR"
}

# 设置定时任务
setup_cron() {
    log_info "设置定时任务..."
    
    # 创建定时任务脚本
    cat > /etc/cron.d/gin-fataMorgana << EOF
# 每天凌晨2点清理日志
0 2 * * * root cd /gin-fataMorgana && ./prod.sh clean-logs

# 每天凌晨3点检查磁盘空间
0 3 * * * root cd /gin-fataMorgana && ./prod.sh disk-space

# 每周日凌晨1点清理Docker空间
0 1 * * 0 root cd /gin-fataMorgana && ./prod.sh clean-docker

# 每天凌晨4点备份数据
0 4 * * * root cd /gin-fataMorgana && ./prod.sh backup
EOF
    
    # 重启cron服务
    systemctl restart crond
    
    log_info "定时任务设置完成"
}

# 主函数
main() {
    check_root
    
    case "${1:-help}" in
        deploy)
            deploy
            ;;
        update)
            update
            ;;
        restart)
            restart
            ;;
        status)
            status
            ;;
        logs)
            logs
            ;;
        clean-logs)
            clean_logs
            ;;
        clean-docker)
            clean_docker
            ;;
        disk-space)
            check_disk_space
            ;;
        health)
            health_check
            ;;
        backup)
            backup
            ;;
        setup-cron)
            setup_cron
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 