#!/bin/bash

APP_NAME=gin-fataMorgana-dev
LOG_FILE=logs/dev_app.log
START_CMD="go run main.go -c config/config.yaml > $LOG_FILE 2>&1 &"
CONFIG_FILE=config/config.yaml

get_db_info() {
  DB_HOST=$(grep '^  host:' $CONFIG_FILE | head -1 | sed 's/.*host:[ ]*"\{0,1\}\([^"]*\)"\{0,1\}.*/\1/')
  DB_USER=$(grep '^  username:' $CONFIG_FILE | head -1 | sed 's/.*username:[ ]*"\{0,1\}\([^"]*\)"\{0,1\}.*/\1/')
  DB_PASS=$(grep '^  password:' $CONFIG_FILE | head -1 | sed 's/.*password:[ ]*"\{0,1\}\([^"]*\)"\{0,1\}.*/\1/')
  DB_NAME=$(grep '^  dbname:' $CONFIG_FILE | head -1 | sed 's/.*dbname:[ ]*"\{0,1\}\([^"]*\)"\{0,1\}.*/\1/')
  echo "[dev] 当前环境: 本地开发"
  echo "[dev] 数据库: $DB_HOST ($DB_NAME) 用户: $DB_USER 密码: $DB_PASS"
}

case "$1" in
  start)
    get_db_info
    echo "[dev] 启动本地开发服务..."
    eval $START_CMD
    sleep 1
    pgrep -fl "go run main.go"
    ;;
  stop)
    echo "[dev] 停止本地开发服务..."
    pkill -f "go run main.go" && echo "已停止开发服务" || echo "未找到开发服务进程"
    ;;
  restart)
    $0 stop
    sleep 1
    $0 start
    ;;
  log)
    echo "[dev] 查看本地开发日志 (Ctrl+C 退出)"
    tail -f $LOG_FILE
    ;;
  clean)
    echo "[dev] 清理本地日志..."
    rm -f $LOG_FILE
    echo "已清理 $LOG_FILE"
    ;;
  *)
    echo "用法: $0 {start|stop|restart|log|clean}"
    ;;
esac 