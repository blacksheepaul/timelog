#!/bin/sh
# entrypoint.sh: 根据 ENV 变量选择配置文件
set -e

if [ "$ENV" = "prod" ]; then
    cp -f config.yml config.active.yml
    echo "[entrypoint] Using production config.yml"
elif [ "$ENV" = "test" ] || [ "$ENV" = "dev" ]; then
    cp -f config-test.yml config.active.yml
    echo "[entrypoint] Using config-test.yml (test/dev)"
else
    cp -f config.yml config.active.yml
    echo "[entrypoint] Using default config.yml"
fi

exec ./main.arm -config config.active.yml
