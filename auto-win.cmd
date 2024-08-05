@echo off

:: 设置变量
set PROJECT_NAME=renhao_go2
set SERVER_ADDR=101.89.144.155
set SOURCE_DIR=C:\myCoding\renhao_go2\main.go
set DEST_DIR=root@%SERVER_ADDR%:/root/%PROJECT_NAME%
set PORT=9006



:: 开始同步代码到云服务器
echo 开始同步代码到云服务器...
scp -r %SOURCE_DIR% %DEST_DIR%
echo 代码同步完成。

:: 开始在云服务器上构建和运行 Docker 镜像
echo 开始在云服务器上构建和运行 Docker 镜像...
ssh root@%SERVER_ADDR% "cd /root/%PROJECT_NAME% && echo '进入目录 /root/%PROJECT_NAME%' && echo '停止容器 %PROJECT_NAME%_service（如果存在）...' && docker stop %PROJECT_NAME%_service || true && echo '删除容器 %PROJECT_NAME%_service（如果存在）...' && docker rm %PROJECT_NAME%_service || true && echo '删除镜像 %PROJECT_NAME%_image（如果存在）...' && docker image rm %PROJECT_NAME%_image || true && echo '构建 Docker 镜像 %PROJECT_NAME%_image...' && docker build -t %PROJECT_NAME%_image -f ./Dockerfile . || (echo 'Docker 镜像构建失败' && exit 1) && echo '运行 Docker 容器 %PROJECT_NAME%_service...' && docker run -itd -p %PORT%:%PORT% --name %PROJECT_NAME%_service %PROJECT_NAME%_image || (echo 'Docker 容器运行失败' && exit 1) && echo 'Docker 容器 %PROJECT_NAME%_service 运行成功。'"



echo 脚本执行完成。