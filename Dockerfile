FROM alpine:latest
WORKDIR /app

# 复制可执行文件和配置
COPY main.arm .
COPY config.yml .
COPY config-test.yml .
COPY dev.db .
COPY entrypoint.sh .

# 赋予启动脚本执行权限
RUN chmod +x entrypoint.sh

# 设置 ENV 变量，默认 dev
ENV ENV=dev

# 启动时用 entrypoint.sh 选择配置
ENTRYPOINT ["/app/entrypoint.sh"]
