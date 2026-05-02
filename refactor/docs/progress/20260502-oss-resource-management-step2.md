# OSS 资源治理修复进度（Step 2）

## 执行日期

- 2026-05-02

## 本步开始前验证

- 已先执行 `go test ./...`（`refactor/`）作为修改前基线检查，结果通过。

## 本步目标

- 对齐 `refactor/docs/plan/oss-resource-management.md` 的 P2：
  - 用户头像更新时旧文件清理防呆逻辑。

## 本步完成内容

### 1) User 头像预留流程增加旧 key 比对防呆

- 文件：`refactor/internal/app/impl/user.go`
- 变更点（`ResvAvatar` 事务内）：
  1. 在 `PrefillAvatarKey` 之前读取当前用户，拿到 `oldKey := user.AvatarKey`。
  2. 计算新 key：`user_avatar/{userId}.{ext}`（保持现有规则）。
  3. 仅当 `oldKey != "" && oldKey != newKey` 时，调用：
     - `ossMsgSvc.SavePendingDel(ossMsgRepo, OssResUserAvatar, userId, []string{oldKey})`
  4. 继续原有 `SavePendingCre` 逻辑。

### 2) 行为说明

- 在当前固定 key 规则下，旧 key 与新 key 通常相同，删除入队通常不会触发。
- 该逻辑作为未来 key 规则变化时的安全防护，避免误删刚上传的新对象。

## 本步结束验证

- 已执行 `gofmt`：
  - `refactor/internal/app/impl/user.go`
- 已执行 `go test ./...`（`refactor/`），结果通过。

## 下一步建议

- 按计划继续 P3：
  - User/Team 软删除独立策略；
  - Comic 封面上传流程补齐（与已完成的 OSS worker 与删除策略联动）。
