# 远程部署快速指南

本文档总结了将本项目部署到远程服务器的最小可行流程（手动或 CI 可复用）。

## 概览

部署方案采用：在构建端打包镜像并传到目标机，目标机通过 `docker load` 加载镜像，先停止旧应用、启动数据库容器、运行一次性迁移容器，最后重建应用容器。此流程已在仓库中实现为 GitHub Actions 工作流，并且现在提供了可直接复用的本地/服务端脚本。

主要参考文件：

- docker-compose 生产配置： docker/compose.prod.yml
- 增量迁移脚本： docker/prod-database-migrate.sh
- GitHub Actions 发布流程： .github/workflows/deploy-prod.yml
- 本地生产镜像 Dockerfile： docker/poprako-s-main/Dockerfile

## 前置条件

- 目标服务器已安装 Docker 与 docker compose（Docker >= 20.x，Compose v2）。
- 你可以通过 SSH 访问目标服务器并有权限运行 docker。
- 远程机器至少开放应用端口（默认 8080）或通过反向代理（NGINX/Caddy）做前端代理与 HTTPS。

## 快速手动部署步骤（在本地执行）

先说明一个常见错误：`docker build` 必须带构建上下文目录，通常是最后的 `.`。你遇到的 `docker buildx build requires 1 argument`，就是因为命令被截断或漏掉了上下文。

推荐方式是直接使用脚本或 `just` 命令。

### 本地一键打包并上传

先设置环境变量：

```sh
export SERVER_USER=youruser
export SERVER_HOST=your.host.example
export DEPLOY_ROOT=/opt/poprako-s
export IMAGE_TAG=$(git rev-parse --short=12 HEAD)
export TARGET_PLATFORM=linux/amd64
```

确保仓库根目录 `.env` 已填好生产运行所需变量。`package-release` 现在会读取本地 `.env`，规范化后自动上传到服务器 `${DEPLOY_ROOT}/shared/.env`。

生产 env 模板可参考：`docker/.env.prod.sample`

然后执行任一方式：

```sh
sh scripts/package-release.sh
```

```sh
just package-release
```

这个脚本会完成以下事情：

- 构建 `poprako-s-main:${IMAGE_TAG}`
- 构建 `poprako-s-database:${IMAGE_TAG}`
- 默认按 `TARGET_PLATFORM=linux/amd64` 构建，适配常见 Ubuntu x86_64 服务器
- 导出两个镜像到 `dist/`
- 打包迁移文件与生产 compose 文件
- 上传发布包到 `${DEPLOY_ROOT}/releases/${IMAGE_TAG}`
- 上传并规范化服务器运行时 env 到 `${DEPLOY_ROOT}/shared/.env`
- 上传服务端切换脚本到 `${DEPLOY_ROOT}/shared/bin/remote-switch-release.sh`

### 服务端切换到新版本

完成打包上传后，直接执行：

```sh
IMAGE_TAG=${IMAGE_TAG} DEPLOY_ROOT=${DEPLOY_ROOT} sh /opt/poprako-s/shared/bin/remote-switch-release.sh
```

如果服务器本身也保留了这份仓库工作树，也可以在仓库根目录执行：

```sh
IMAGE_TAG=${IMAGE_TAG} DEPLOY_ROOT=${DEPLOY_ROOT} just switch-release
```

服务端脚本会完成以下事情：

- 加载上传好的主服务和数据库镜像
- 解压迁移包到 `${DEPLOY_ROOT}/shared`
- 更新 `${DEPLOY_ROOT}/shared/.env` 里的 `IMAGE_TAG`
- 停止旧的 `prod-main-server`
- 强制重建 `prod-postgres`，避免旧架构容器残留
- 执行 `prod-db-migrate`
- 启动新的 `prod-main-server`

### 手动拆分步骤

1. 在本地生成 IMAGE_TAG（示例使用 git sha）并构建镜像：

```sh
IMAGE_TAG=$(git rev-parse --short=12 HEAD)
mkdir -p dist
docker build --platform linux/amd64 -f docker/poprako-s-main/Dockerfile -t poprako-s-main:${IMAGE_TAG} .
docker build --platform linux/amd64 -f docker/poprako-s-database/Dockerfile -t poprako-s-database:${IMAGE_TAG} .
docker save poprako-s-main:${IMAGE_TAG} | gzip > dist/poprako-s-main-${IMAGE_TAG}.tar.gz
docker save poprako-s-database:${IMAGE_TAG} | gzip > dist/poprako-s-database-${IMAGE_TAG}.tar.gz
tar -czf dist/poprako-s-migrations-${IMAGE_TAG}.tar.gz migrations docker/prod-database-migrate.sh docker/compose.prod.yml
```

如果你只执行了其中的 `docker save ... > dist/...` 这一行，而没有先创建 `dist/` 目录，就会得到 `zsh: no such file or directory`。要么先执行 `mkdir -p dist`，要么整段一起复制执行。

2. 上传到远程：

```sh
SERVER_USER=youruser
SERVER_HOST=your.host.example
DEPLOY_ROOT=/opt/poprako-s
ssh ${SERVER_USER}@${SERVER_HOST} "mkdir -p ${DEPLOY_ROOT}/releases/${IMAGE_TAG} ${DEPLOY_ROOT}/shared"
scp dist/poprako-s-main-${IMAGE_TAG}.tar.gz ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_ROOT}/releases/${IMAGE_TAG}/
scp dist/poprako-s-database-${IMAGE_TAG}.tar.gz ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_ROOT}/releases/${IMAGE_TAG}/
scp dist/poprako-s-migrations-${IMAGE_TAG}.tar.gz ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_ROOT}/releases/${IMAGE_TAG}/
```

3. 在远程加载镜像与准备目录（在远程 shell 执行）：

```sh
DEPLOY_ROOT=/opt/poprako-s
IMAGE_TAG=... # 与上面一致
RELEASE_DIR=${DEPLOY_ROOT}/releases/${IMAGE_TAG}
SHARED_DIR=${DEPLOY_ROOT}/shared

docker load -i ${RELEASE_DIR}/poprako-s-database-${IMAGE_TAG}.tar.gz
docker load -i ${RELEASE_DIR}/poprako-s-main-${IMAGE_TAG}.tar.gz

tar -xzf ${RELEASE_DIR}/poprako-s-migrations-${IMAGE_TAG}.tar.gz -C ${SHARED_DIR}
```

4. 在远程创建运行时 env（至少包含以下变量）：

如果你使用 `just package-release` 或 `sh scripts/package-release.sh`，这一步通常不需要手动做；脚本会基于本地 `.env` 自动上传。

```
IMAGE_TAG=${IMAGE_TAG}
APP_ENV=prod
POSTGRES_PASSWORD=your_database_password
DATABASE_USER=poprako_s
DATABASE_NAME=db_poprako_s
DATABASE_PASSWORD=your_database_password
DATABASE_HOST=prod-postgres
DATABASE_PORT=5432
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION_HOURS=336
OSS_PLATFORM=r2
R2_ACCOUNT_ID=your_r2_account_id
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_BUCKET_NAME=your_r2_bucket_name
R2_REGION=auto
R2_CUSTOM_DOMAIN=your_r2_custom_domain
```

将上面内容写入 ${SHARED_DIR}/.env 并设置仅 owner 可读（chmod 600）。

5. 使用 production compose 文件启动服务：

```sh
docker compose -f ${SHARED_DIR}/docker/compose.prod.yml --env-file ${SHARED_DIR}/.env up -d --wait prod-postgres
docker compose -f ${SHARED_DIR}/docker/compose.prod.yml --env-file ${SHARED_DIR}/.env run --rm prod-db-migrate
docker compose -f ${SHARED_DIR}/docker/compose.prod.yml --env-file ${SHARED_DIR}/.env up -d --force-recreate prod-main-server
docker compose -f ${SHARED_DIR}/docker/compose.prod.yml --env-file ${SHARED_DIR}/.env ps
```

## 验证

- 查看容器状态： docker ps 或 docker compose -f ${SHARED_DIR}/docker/compose.prod.yml --env-file ${SHARED_DIR}/.env ps
- 查看应用日志： docker logs -f prod-main-server（或用 compose logs）
- 测试 HTTP：curl http://localhost:8080/ 或通过已配置的反向代理访问外网地址

## 回滚（最小可行）

1. 在 shared/.env 把 IMAGE_TAG 改为此前稳定的 tag。
2. 执行：

```sh
docker compose -f ${SHARED_DIR}/docker/compose.prod.yml --env-file ${SHARED_DIR}/.env up -d --force-recreate prod-main-server
```

注意：回滚只替换应用镜像，不能回退已经运行的数据库迁移；如果新迁移对 schema 有不可逆改动，请在部署前确保有数据库备份策略。

## 风险与注意事项

- 迁移校验： docker/prod-database-migrate.sh 会校验已记录迁移的 checksum，已应用的 SQL 被修改会导致部署失败 — 请在变更 migration 文件后谨慎操作。
- 数据安全：务必保护 ${SHARED_DIR}/.env，并为数据库做备份（快照/pg_dump）。
- 端口与防火墙：默认暴露 8080，生产环境建议在前端放一层反向代理并启用 TLS。
- 持久化：Postgres 使用 Docker 卷 `poprako_s_prod_db_data`，除非你手动删除卷，否则数据会保留。

## 脚本入口

- 本地打包上传：`scripts/package-release.sh`
- 服务端切换版本：`scripts/remote-switch-release.sh`
- just 包装命令：`just package-release`，以及在服务器存在仓库工作树时可用的 `just switch-release`

---
