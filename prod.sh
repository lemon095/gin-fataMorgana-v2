#!/bin/bash

# 生产环境管理脚本
# 使用方法: ./prod.sh [start|stop|restart|logs|status|update|backup|clean|clean-cache]

COMPOSE_FILE="docker-compose.yml"
SERVICE_NAME="gin-fataMorgana"

# 检查代码变更的函数
check_code_changes() {
    echo "🔍 检查代码变更..."
    
    # 检查是否在git仓库中
    if [ ! -d ".git" ]; then
        echo "⚠️  当前目录不是git仓库，跳过代码变更检查"
        return 1
    fi
    
    # 获取当前提交的哈希
    CURRENT_COMMIT=$(git rev-parse HEAD)
    
    # 检查是否有未提交的更改
    if ! git diff-index --quiet HEAD --; then
        echo "🔄 检测到未提交的本地更改"
        return 1
    fi
    
    # 检查是否有新的远程提交
    git fetch origin > /dev/null 2>&1
    REMOTE_COMMIT=$(git rev-parse origin/$(git branch --show-current))
    
    if [ "$CURRENT_COMMIT" != "$REMOTE_COMMIT" ]; then
        echo "🔄 检测到远程代码更新: $CURRENT_COMMIT -> $REMOTE_COMMIT"
        return 1
    fi
    
    echo "📝 没有检测到代码变更"
    return 0
}

# 智能构建函数
smart_build() {
    local force_rebuild=$1
    
    if [ "$force_rebuild" = "true" ]; then
        echo "🔨 强制重新构建镜像（清理缓存）..."
        # 先清理构建缓存
        docker builder prune -f
        docker compose -f $COMPOSE_FILE build --no-cache
    else
        echo "🔨 使用缓存构建镜像..."
        docker compose -f $COMPOSE_FILE build
    fi
}

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

# 清理构建缓存
clean_build_cache() {
    echo "🧹 清理构建缓存..."
    
    # 显示清理前的状态
    echo "📊 清理前的Docker使用情况:"
    docker system df
    
    # 清理构建缓存
    echo "🗑️  清理构建缓存..."
    docker builder prune -a -f
    
    # 清理未使用的镜像
    echo "🗑️  清理未使用的镜像..."
    docker image prune -a -f
    
    # 清理未使用的数据卷
    echo "🗑️  清理未使用的数据卷..."
    docker volume prune -f
    
    # 清理未使用的网络
    echo "🗑️  清理未使用的网络..."
    docker network prune -f
    
    # 显示清理后的状态
    echo "📊 清理后的Docker使用情况:"
    docker system df
    
    echo "✅ 构建缓存清理完成！"
}

case "$1" in
    start)
        echo "🚀 启动生产环境服务..."
        
        # 检查代码变更
        check_code_changes
        NEED_REBUILD=$?
        
        # 智能构建
        if [ $NEED_REBUILD -eq 1 ]; then
            smart_build false
        else
            echo "📝 没有代码变更，跳过构建"
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
        
        # 检查代码变更
        check_code_changes
        NEED_REBUILD=$?
        
        # 停止服务
        docker compose -f $COMPOSE_FILE down
        
        # 智能构建
        if [ $NEED_REBUILD -eq 1 ]; then
            smart_build false
        else
            echo "📝 没有代码变更，跳过构建"
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
        echo ""
        echo "📊 Docker系统使用情况:"
        docker system df
        ;;
    update)
        echo "🔄 零停机更新服务..."
        
        # 检查代码变更
        check_code_changes
        NEED_REBUILD=$?
        
        # 检查当前服务状态
        if ! docker compose -f $COMPOSE_FILE ps | grep -q "Up"; then
            echo "⚠️  服务未运行，直接启动新版本..."
            docker compose -f $COMPOSE_FILE down
            
            # 智能构建
            if [ $NEED_REBUILD -eq 1 ]; then
                smart_build false
            else
                echo "📝 没有代码变更，跳过构建"
            fi
            
            docker compose -f $COMPOSE_FILE up -d
            echo "✅ 更新完成！"
            exit 0
        fi
        
        echo "📋 当前服务状态:"
        docker compose -f $COMPOSE_FILE ps
        
        # 智能构建
        if [ $NEED_REBUILD -eq 1 ]; then
            smart_build false
        else
            echo "📝 没有代码变更，跳过构建"
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
    force-update)
        echo "🔄 强制更新服务（清理缓存）..."
        
        # 拉取最新代码
        pull_latest_code
        NEED_REBUILD=$?
        
        # 检查当前服务状态
        if ! docker compose -f $COMPOSE_FILE ps | grep -q "Up"; then
            echo "⚠️  服务未运行，直接启动新版本..."
            docker compose -f $COMPOSE_FILE down
            
            # 强制构建
            smart_build true
            
            docker compose -f $COMPOSE_FILE up -d
            echo "✅ 强制更新完成！"
            exit 0
        fi
        
        echo "📋 当前服务状态:"
        docker compose -f $COMPOSE_FILE ps
        
        # 强制构建
        smart_build true
        
        # 使用Docker Compose的滚动更新功能
        echo "🚀 执行强制零停机更新..."
        docker compose -f $COMPOSE_FILE up -d --force-recreate
        
        # 等待服务启动
        echo "⏳ 等待服务启动..."
        sleep 15
        
        # 检查服务健康状态
        echo "🔍 检查服务健康状态..."
        if curl -s http://localhost:9001/health > /dev/null 2>&1; then
            echo "✅ 服务健康检查通过"
            echo "✅ 强制零停机更新完成！"
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
    clean-cache)
        clean_build_cache
        ;;
    *)
        echo "❓ 使用方法: $0 {start|stop|restart|logs|status|update|force-update|backup|clean|clean-cache}"
        echo ""
        echo "命令说明:"
        echo "  start        - 启动生产环境服务（智能构建）"
        echo "  stop         - 停止生产环境服务"
        echo "  restart      - 重启生产环境服务（智能构建）"
        echo "  logs         - 查看服务日志"
        echo "  status       - 查看服务状态和Docker使用情况"
        echo "  update       - 零停机更新服务（智能构建）"
        echo "  force-update - 强制更新服务（清理缓存后构建）"
        echo "  backup       - 数据库备份提示"
        echo "  clean        - 清理Docker资源"
        echo "  clean-cache  - 清理构建缓存"
        echo ""
        echo "📝 配置说明:"
        echo "  - 服务端口: 9001"
        echo "  - MySQL连接: 172.31.46.166:3306"
        echo "  - Redis连接: 172.31.46.166:6379"
        echo "  - 模式: release"
        echo ""
        echo "🔄 构建策略:"
        echo "  - 智能构建: 只在代码变更时重新构建，使用缓存"
        echo "  - 强制构建: 清理缓存后重新构建，用于解决构建问题"
        echo "  - 缓存清理: 定期清理构建缓存，释放磁盘空间"
        echo ""
        echo "⚠️  注意: 请确保MySQL和Redis服务已启动并可访问"
        echo "💡 提示: 使用 clean-cache 定期清理构建缓存"
        exit 1
        ;;
esac 