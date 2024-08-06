#!/bin/bash

# chmod +x ./auto.sh

# 先决条件说明（可选）
# 0. 需要在本机执行以下语句，将密钥同步到云端服务器，这样可以无需密码验证(这部分不执行也行，每次需要输入服务器密码)，直接进行同步
# 1. 本机生成密码（公钥）
#    ssh-keygen -t rsa
#    ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
# 2. 将密钥同步至云端服务器
#    ssh-copy-id root@122.9.41.45
# 3. 此时会需要输入云端服务器密码
# 4. 在云端服务器上及本机上安装rsync
#    在 Ubuntu 或 Debian 系统上安装 rsync
#    在 Ubuntu 或 Debian 系统上，可以使用 apt-get 包管理器安装 rsync。运行以下命令：
#    sudo apt-get update
#    sudo apt-get install rsync
#    在 CentOS 或 Red Hat 系统上安装 rsync
#    在 CentOS 或 Red Hat 系统上，可以使用 yum 包管理器安装 rsync。运行以下命令：
#    sudo yum install rsync
#    在 macOS 上安装 rsync
#    在 macOS 上，可以使用 Homebrew 或 MacPorts 安装 rsync。如果你使用 Homebrew，可以运行以下命令：
#    brew install rsync

# 设置变量
PROJECT_NAME="renhao_go2"
SERVER_ADDR="101.89.144.155"
SOURCE_DIR="/Users/kealuya/mywork/my_git/info_project/back/$PROJECT_NAME/"  
DEST_DIR="root@$SERVER_ADDR:/root/$PROJECT_NAME"
PORT=9009

# 获取当前时间
start_time=$(date +%s)

# 开始同步代码到云服务器
echo "开始同步代码到云服务器..."
rsync -avz $SOURCE_DIR $DEST_DIR
echo "代码同步完成。"

# 开始在云服务器上构建和运行 Docker 镜像
echo "开始在云服务器上构建和运行 Docker 镜像..."
ssh root@$SERVER_ADDR << EOF
echo "进入目录 /root/${PROJECT_NAME}"
cd /root/$PROJECT_NAME

echo "停止容器 ${PROJECT_NAME}_service（如果存在）..."
docker stop $PROJECT_NAME_service || true

echo "删除容器 ${PROJECT_NAME}_service（如果存在）..."
docker rm $PROJECT_NAME_service || true

echo "删除镜像 ${PROJECT_NAME}_image（如果存在）..."
docker image rm $PROJECT_NAME_image || true

echo "构建 Docker 镜像 ${PROJECT_NAME}_image..."
docker build -t $PROJECT_NAME_image -f ./Dockerfile . || { echo "Docker 镜像构建失败"; exit 1; }

echo "运行 Docker 容器 ${PROJECT_NAME}_service..."
docker run -itd -p $PORT:$PORT --name $PROJECT_NAME_service $PROJECT_NAME_image || { echo "Docker 容器运行失败"; exit 1; }

echo "Docker 容器 ${PROJECT_NAME}_service 运行成功。"
EOF

# 获取当前时间
end_time=$(date +%s)
duration=$(( end_time - start_time ))

echo "脚本执行完成。总耗时: ${duration}秒"