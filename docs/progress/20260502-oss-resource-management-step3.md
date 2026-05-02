# OSS 资源治理修复进度（Step 3）

## 执行日期

- 2026-05-02

## 本步开始前状态

- 上一步后代码处于不可编译状态，`go test ./...` 失败。
- 主要错误点集中在 comic 封面流程：
  - `app_util.SafeGetUrl` 不存在；
  - `ComicApp` 新增方法后，`comic_log` 装饰器未对齐；
  - `main.go` 中 `NewComicApp` 参数未补齐 `ossSigner`。

## 本步完成内容

### 1) 修复 `comic` 封面流程编译与依赖问题

- 文件：`refactor/internal/app/impl/comic.go`
- 调整：
  - `GetById` 的封面 URL 生成改为 `a.ossSigner.GenGetUrl(...)`；
  - `ResvCover` 的上传 URL 生成改为 `a.ossSigner.GenPutUrl(...)`；
  - 删除错误依赖 `app_util.SafeGetUrl`。

### 2) 补齐 `ComicApp` 装饰器方法集

- 文件：`refactor/internal/app/impl/comic_log.go`
- 新增转发方法：
  - `ResvCover`
  - `MarkCoverUploaded`
- 与 `internal/app/comic.go` 接口保持一致，恢复接口完整实现。

### 3) 补齐 `main` 构造参数

- 文件：`refactor/main.go`
- `NewComicApp(...)` 调用增加 `ossClient`（`oss_iface.Signer`）参数，修复构造签名不一致。

### 4) 接入 comic 封面 HTTP API

- 文件：
  - `refactor/internal/api/http/comic.go`
  - `refactor/internal/api/http/http.go`
- 新增接口：
  - `POST /api/v1/comic/{comic_id}/cover`
  - `POST /api/v1/comic/{comic_id}/cover/confirm`
- 同步增加 swagger 注释。

## 格式化与校验

- 已执行 `gofmt`：
  - `internal/app/impl/comic.go`
  - `internal/app/impl/comic_log.go`
  - `internal/api/http/comic.go`
  - `internal/api/http/http.go`
  - `main.go`
- 已执行 `go test ./...`（`refactor/`）并通过。

## 当前剩余项（待确认）

- 计划文档中的 P3 仍有一项“User/Team 软删除独立处理”。
- 当前 `refactor` 的 user/team app-repo contract 尚未暴露明确删除用例，本项需要先明确目标行为与接口边界后再实施。
