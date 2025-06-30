#!/bin/bash

APP_NAME=gin-fataMorgana
LOG_FILE=logs/app.log
START_CMD="./gin-fataMorgana -c config/config.prod.yaml > $LOG_FILE 2>&1 &"
CONFIG_FILE=config/config.prod.yaml

get_db_info() {
  DB_HOST=$(grep 'host:' $CONFIG_FILE | grep -v '#' | head -1 | awk '{print $2}' | tr -d '"')
  DB_USER=$(grep 'username:' $CONFIG_FILE | grep -v '#' | head -1 | awk '{print $2}' | tr -d '"')
  DB_PASS=$(grep 'password:' $CONFIG_FILE | grep -v '#' | head -1 | awk '{print $2}' | tr -d '"')
  DB_NAME=$(grep 'dbname:' $CONFIG_FILE | grep -v '#' | head -1 | awk '{print $2}' | tr -d '"')
  echo "[prod] 当前环境: 生产环境"
  echo "[prod] 数据库: $DB_HOST ($DB_NAME) 用户: $DB_USER 密码: $DB_PASS"
}

case "$1" in
  start)
    get_db_info
    echo "[prod] 启动生产服务..."
    eval $START_CMD
    sleep 1
    pgrep -fl $APP_NAME
    ;;
  stop)
    echo "[prod] 停止生产服务..."
    pkill -f $APP_NAME && echo "已停止生产服务" || echo "未找到生产服务进程"
    ;;
  restart)
    $0 stop
    sleep 1
    $0 start
    ;;
  log)
    echo "[prod] 查看生产日志 (Ctrl+C 退出)"
    tail -f $LOG_FILE
    ;;
  update)
    echo "[prod] 拉取最新代码并重启..."
    git pull
    go build -o $APP_NAME main.go
    $0 restart
    ;;
  clean)
    echo "[prod] 清理生产日志..."
    rm -f $LOG_FILE
    echo "已清理 $LOG_FILE"
    ;;
  *)
    echo "用法: $0 {start|stop|restart|log|update|clean}"
    ;;
esac 