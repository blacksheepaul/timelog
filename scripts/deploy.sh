#!/bin/bash
set -e

# 检查必需的环境变量
if [ -z "$ARTIFACT_PATH" ]; then
    echo "Error: ARTIFACT_PATH is not set"
    exit 1
fi

# 配置变量（可通过环境变量覆盖）
DEPLOY_DIR="${DEPLOY_DIR:-/var/www/timelog}"
TEMP_DIR=$(mktemp -d)

echo "Deploying timelog..."

# 1. 解压到临时目录
echo "Extracting artifact..."
tar -xzf "$ARTIFACT_PATH" -C "$TEMP_DIR"

# 2. 停止服务
echo "Stopping service..."
sudo systemctl stop timelog.service

# 3. 替换文件
echo "Replacing files..."
rm -rf "$DEPLOY_DIR/main"
mv "$TEMP_DIR/main.linux" "$DEPLOY_DIR/main"
chmod +x "$DEPLOY_DIR/main"

# 4. 重启服务
echo "Starting service..."
sudo systemctl start timelog.service

# 5. 清理临时文件
echo "Cleaning up..."
rm -rf "$TEMP_DIR"
rm -f "$ARTIFACT_PATH"

# 验证服务状态
sleep 2
if sudo systemctl is-active --quiet timelog.service; then
    echo "Deployment completed successfully!"
else
    echo "Error: Service failed to start"
    sudo systemctl status timelog.service
    exit 1
fi
