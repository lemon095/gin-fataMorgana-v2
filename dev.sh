#!/bin/bash

# 本地开发环境管理脚本
# 使用方法: ./dev.sh [start|stop|restart|log|status|clean]

APP_NAME="gin-fataMorgana-dev"
LOG_FILE="logs/dev_app.log"
PID_FILE="logs/dev.pid"

# 检查配置文件是否存在
check_config() {
    if [ ! -f "config.yaml" ]; then
        echo "❌ 配置文件 config.yaml 不存在"
        exit 1
    fi
}

# 检查数据库连接
check_database() {
    echo "🔍 检查MySQL连接..."
    if ! mysql -h 127.0.0.1 -P 3306 -u root -proot -e "SELECT 1;" > /dev/null 2>&1; then
        echo "❌ MySQL连接失败，请确保MySQL服务正在运行"
        echo "💡 启动MySQL命令: brew services start mysql (macOS) 或 sudo systemctl start mysql (Linux)"
        exit 1
    fi
    echo "✅ MySQL连接成功"

    echo "🔍 检查Redis连接..."
    if ! redis-cli -h 127.0.0.1 -p 6379 ping > /dev/null 2>&1; then
        echo "❌ Redis连接失败，请确保Redis服务正在运行"
        echo "💡 启动Redis命令: brew services start redis (macOS) 或 sudo systemctl start redis (Linux)"
        exit 1
    fi
    echo "✅ Redis连接成功"

    echo "🔍 检查数据库是否存在..."
    if ! mysql -h 127.0.0.1 -P 3306 -u root -proot -e "USE future;" > /dev/null 2>&1; then
        echo "⚠️  数据库 'future' 不存在，正在创建..."
        mysql -h 127.0.0.1 -P 3306 -u root -proot -e "CREATE DATABASE IF NOT EXISTS future CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        echo "✅ 数据库创建成功"
    fi
}

# 创建必要目录
create_dirs() {
    echo "📁 创建日志目录..."
    mkdir -p logs
}

# 显示配置信息
show_config() {
    echo "📋 当前配置信息:"
    echo "  - 服务端口: 9001"
    echo "  - 数据库: 127.0.0.1:3306 (future)"
    echo "  - Redis: 127.0.0.1:6379"
    echo "  - 模式: debug"
    echo "  - 日志文件: $LOG_FILE"
}

# 启动服务
start_service() {
    echo "🚀 启动本地开发服务..."
    
    # 检查是否已经运行
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p $PID > /dev/null 2>&1; then
            echo "⚠️  服务已经在运行 (PID: $PID)"
            return
        else
            echo "🧹 清理无效的PID文件"
            rm -f "$PID_FILE"
        fi
    fi

    # 启动服务
    nohup go run main.go > "$LOG_FILE" 2>&1 &
    PID=$!
    echo $PID > "$PID_FILE"
    
    sleep 2
    
    # 检查是否启动成功
    if ps -p $PID > /dev/null 2>&1; then
        echo "✅ 服务启动成功 (PID: $PID)"
        echo "📍 服务地址: http://localhost:9001"
        echo "📝 日志文件: $LOG_FILE"
    else
        echo "❌ 服务启动失败"
        rm -f "$PID_FILE"
        exit 1
    fi
}

# 停止服务
stop_service() {
    echo "🛑 停止本地开发服务..."
    
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p $PID > /dev/null 2>&1; then
            kill $PID
            echo "✅ 服务已停止 (PID: $PID)"
        else
            echo "⚠️  服务未运行"
        fi
        rm -f "$PID_FILE"
    else
        # 如果没有PID文件，尝试通过进程名停止
        PIDS=$(pgrep -f "go run main.go")
        if [ -n "$PIDS" ]; then
            echo "🔍 发现运行中的进程: $PIDS"
            kill $PIDS
            echo "✅ 服务已停止"
        else
            echo "⚠️  未找到运行中的服务"
        fi
    fi
}

# 重启服务
restart_service() {
    echo "🔄 重启本地开发服务..."
    stop_service
    sleep 2
    start_service
}

# 查看日志
show_logs() {
    echo "📊 查看服务日志 (Ctrl+C 退出)"
    echo "----------------------------------------"
    if [ -f "$LOG_FILE" ]; then
        tail -f "$LOG_FILE"
    else
        echo "❌ 日志文件不存在: $LOG_FILE"
    fi
}

# 查看状态
show_status() {
    echo "🔍 服务状态检查..."
    
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p $PID > /dev/null 2>&1; then
            echo "✅ 服务正在运行 (PID: $PID)"
            echo "📍 服务地址: http://localhost:9001"
            
            # 检查端口
            if lsof -i :9001 > /dev/null 2>&1; then
                echo "✅ 端口 9001 正在监听"
            else
                echo "❌ 端口 9001 未监听"
            fi
            
            # 检查健康状态
            if curl -s http://localhost:9001/health > /dev/null 2>&1; then
                echo "✅ 健康检查通过"
            else
                echo "❌ 健康检查失败"
            fi
        else
            echo "❌ 服务未运行 (PID文件存在但进程不存在)"
            rm -f "$PID_FILE"
        fi
    else
        echo "❌ 服务未运行 (无PID文件)"
    fi
}

# 清理
clean_service() {
    echo "🧹 清理本地开发环境..."
    
    # 停止服务
    stop_service
    
    # 清理日志
    if [ -f "$LOG_FILE" ]; then
        rm -f "$LOG_FILE"
        echo "✅ 已清理日志文件"
    fi
    
    # 清理PID文件
    if [ -f "$PID_FILE" ]; then
        rm -f "$PID_FILE"
        echo "✅ 已清理PID文件"
    fi
    
    echo "✅ 清理完成"
}

# 主逻辑
case "$1" in
    start)
        check_config
        check_database
        create_dirs
        show_config
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        check_config
        check_database
        create_dirs
        restart_service
        ;;
    log)
        show_logs
        ;;
    status)
        show_status
        ;;
    clean)
        clean_service
        ;;
    *)
        echo "❓ 使用方法: $0 {start|stop|restart|log|status|clean}"
        echo ""
        echo "命令说明:"
        echo "  start   - 启动本地开发服务"
        echo "  stop    - 停止本地开发服务"
        echo "  restart - 重启本地开发服务"
        echo "  log     - 查看服务日志"
        echo "  status  - 查看服务状态"
        echo "  clean   - 清理服务文件"
        echo ""
        echo "📝 配置说明:"
        echo "  - 服务端口: 9001"
        echo "  - 数据库: 127.0.0.1:3306 (future)"
        echo "  - Redis: 127.0.0.1:6379"
        echo "  - 模式: debug"
        exit 1
        ;;
esac 