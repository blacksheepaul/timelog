#### 需要配置的 GitHub Secrets

- DEPLOY_HOST ：目标服务器地址
- DEPLOY_USER ：SSH 用户名
- DEPLOY_PORT ：SSH 端口（可选，默认 22）
- DEPLOY_SSH_KEY ：私钥内容
- DEPLOY_SCRIPT ：目标服务器上的部署脚本绝对路径（例如 /opt/timelog/deploy.sh ）
- DEPLOY_PATH ：远端保存产物的目录（例如 /opt/timelog/releases ）

#### 部署脚本可用变量

- RELEASE_TAG ：tag 名
- ARTIFACT_NAME ：tgz 文件名
- ARTIFACT_PATH ：远端完整路径