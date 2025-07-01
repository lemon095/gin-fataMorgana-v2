#!/bin/bash

# 生产环境管理脚本
# 使用方法: ./prod.sh [start|stop|restart|logs|status|update|backup|clean]

COMPOSE_FILE="docker-compose.yml"
SERVICE_NAME="gin-fataMorgana"

# 拉取最新代码
pull_latest_code() {
    echo "📥 拉取最新代码..."
    
    # 检查是否在git仓库中
    if [ ! -d ".git" ]; then
        echo "⚠️  当前目录不是git仓库，跳过代码拉取"
        return 0
    fi
    
    # 检查git命令是否可用
    if ! command -v git &> /dev/null; then
        echo "⚠️  git命令不可用，跳过代码拉取"
        return 0
    fi
    
    # 获取当前分支
    CURRENT_BRANCH=$(git branch --show-current)
    echo "📍 当前分支: $CURRENT_BRANCH"
    
    # 保存当前工作目录
    PWD_BACKUP=$(pwd)
    
    # 拉取最新代码
    if git pull origin $CURRENT_BRANCH; then
        echo "✅ 代码拉取成功"
        
        # 检查是否有新提交
        if git log --oneline -1 | grep -q "$(git rev-parse HEAD)"; then
            echo "📝 代码已是最新版本"
        else
            echo "🔄 检测到新代码，需要重新构建镜像"
            return 1
        fi
    else
        echo "❌ 代码拉取失败"
        return 1
    fi
    
    return 0
}

case "$1" in
    start)
        echo "🚀 启动生产环境服务..."
        
        # 拉取最新代码
        pull_latest_code
        NEED_REBUILD=$?
        
        # 如果有新代码，重新构建镜像
        if [ $NEED_REBUILD -eq 1 ]; then
            echo "🔨 重新构建镜像..."
            docker compose -f $COMPOSE_FILE build --no-cache
        fi
        
        docker compose -f $COMPOSE_FILE up -d
        echo "✅ 启动完成！访问地址: http://localhost:9001"
        echo "💡 请确保MySQL和Redis服务已启动并可访问"
        ;;
    stop)
        echo "🛑 停止生产环境服务..."
        docker compose -f $COMPOSE_FILE down
        echo "✅ 停止完成！"
        ;;
    restart)
        echo "🔄 重启生产环境服务..."
        
        # 拉取最新代码
        pull_latest_code
        NEED_REBUILD=$?
        
        # 停止服务
        docker compose -f $COMPOSE_FILE down
        
        # 如果有新代码，重新构建镜像
        if [ $NEED_REBUILD -eq 1 ]; then
            echo "🔨 重新构建镜像..."
            docker compose -f $COMPOSE_FILE build --no-cache
        fi
        
        # 启动服务
        docker compose -f $COMPOSE_FILE up -d
        echo "✅ 重启完成！"
        ;;
    logs)
        echo "📊 查看服务日志..."
        docker compose -f $COMPOSE_FILE logs -f
        ;;
    status)
        echo "🔍 查看服务状态..."
        docker compose -f $COMPOSE_FILE ps
        echo ""
        echo "📊 容器资源使用情况:"
        docker stats --no-stream
        ;;
    update)
        echo "🔄 零停机更新服务..."
        
        # 拉取最新代码
        pull_latest_code
        NEED_REBUILD=$?
        
        # 检查当前服务状态
        if ! docker compose -f $COMPOSE_FILE ps | grep -q "Up"; then
            echo "⚠️  服务未运行，直接启动新版本..."
            docker compose -f $COMPOSE_FILE down
            
            # 如果有新代码，重新构建镜像
            if [ $NEED_REBUILD -eq 1 ]; then
                echo "🔨 重新构建镜像..."
                docker compose -f $COMPOSE_FILE build --no-cache
            fi
            
            docker compose -f $COMPOSE_FILE up -d
            echo "✅ 更新完成！"
            exit 0
        fi
        
        echo "📋 当前服务状态:"
        docker compose -f $COMPOSE_FILE ps
        
        # 如果有新代码，重新构建镜像
        if [ $NEED_REBUILD -eq 1 ]; then
            echo "🔨 构建新版本镜像..."
            docker compose -f $COMPOSE_FILE build --no-cache
        fi
        
        # 使用Docker Compose的滚动更新功能
        echo "🚀 执行零停机更新..."
        docker compose -f $COMPOSE_FILE up -d --force-recreate
        
        # 等待服务启动
        echo "⏳ 等待服务启动..."
        sleep 15
        
        # 检查服务健康状态
        echo "🔍 检查服务健康状态..."
        if curl -s http://localhost:9001/health > /dev/null 2>&1; then
            echo "✅ 服务健康检查通过"
            echo "✅ 零停机更新完成！"
            echo "📍 服务地址: http://localhost:9001"
        else
            echo "❌ 服务健康检查失败，请检查日志"
            echo "💡 查看日志: ./prod.sh logs"
            exit 1
        fi
        ;;
    backup)
        echo "💾 备份数据库..."
        echo "⚠️  请手动备份您的MySQL数据"
        echo "💡 示例命令: docker exec your-mysql-container mysqldump -u root -proot future > backup_$(date +%Y%m%d_%H%M%S).sql"
        ;;
    clean)
        echo "🧹 清理未使用的Docker资源..."
        docker system prune -f
        docker volume prune -f
        echo "✅ 清理完成！"
        ;;
    *)
        echo "❓ 使用方法: $0 {start|stop|restart|logs|status|update|backup|clean}"
        echo ""
        echo "命令说明:"
        echo "  start   - 启动生产环境服务（自动拉取代码）"
        echo "  stop    - 停止生产环境服务"
        echo "  restart - 重启生产环境服务（自动拉取代码）"
        echo "  logs    - 查看服务日志"
        echo "  status  - 查看服务状态"
        echo "  update  - 零停机更新服务（自动拉取代码）"
        echo "  backup  - 数据库备份提示"
        echo "  clean   - 清理Docker资源"
        echo ""
        echo "📝 配置说明:"
        echo "  - 服务端口: 9001"
        echo "  - MySQL连接: 172.31.46.166:3306"
        echo "  - Redis连接: 172.31.46.166:6379"
        echo "  - 模式: release"
        echo ""
        echo "⚠️  注意: 请确保MySQL和Redis服务已启动并可访问"
        echo "💡 提示: start/restart/update 命令会自动拉取最新代码"
        exit 1
        ;;
esac 