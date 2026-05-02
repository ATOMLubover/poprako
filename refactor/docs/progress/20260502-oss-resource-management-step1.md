# OSS 资源治理修复进度（Step 1）

## 执行日期

- 2026-05-02

## 本步开始前验证

- 已先执行基线验证：`go test ./...`（在 `refactor/` 下）通过。
- 已核对 `refactor/docs/plan/oss-resource-management.md` 以及现有迁移脚本与实现状态。
- 已确认 `refactor/docs/progress/` 原先不存在，本步已创建并开始记录。

## 本步完成内容

### 1) P0：补齐 OSS Worker 后台协程

- 新增：`refactor/internal/infra/worker/oss_worker.go`
- 在 `main.go` 中接入启动：
  - 使用 `signal.NotifyContext` 获取运行上下文；
  - `go worker_infra.NewOssWorker(ossMsgRepo, ossClient).Run(runCx)` 启动后台循环。
- 实现 3 个循环：
  - `consumeLoop`（5min）
  - `reclaimLoop`（5min）
  - `purgeLoop`（24h）
- `create` 消息未过期时，按计划使用：
  - `MarkPending(id, "not yet expired", expireAt)`

### 2) P1：接入 `OssMsgSvc.SavePendingDel`

- 新增：`internal/domain/svc/oss_msg.go` 中 `SavePendingDel`。
- 替换 app 层直接 `ossMsgRepo.SavePendingDel` 的调用：
  - `chapter.Remove`
  - `page.RemoveByChapterId`
  - 新增级联删除流程中的 `comic.Remove` / `workset.Remove`

### 3) P1：实现 Comic/Workset 级联删除（物理删除 + OSS 入队）

- `comic.Remove`：
  - 遍历 chapter/page，页面图片入队 OSS 删除；
  - 删除 pages / assignments / chapters；
  - 漫画封面入队 OSS 删除；
  - 物理删除 comic；
  - 更新 workset comic_count。
- `workset.Remove`：
  - 遍历 workset 下 comics，逐层执行章节/页面 OSS 入队与物理删除；
  - 物理删除 workset。

### 4) P1：移除 t_workset / t_comic / t_chapter 软删除并接入 next_index

- 按要求仅修改现有 migration schema 定义，未新增任何 alter 脚本。
- 修改 migration：
  - `20260427000001_create-workset-table.up.sql`
  - `20260427000002_create-comic-table.up.sql`
  - `20260427000003_create-chapter-table.up.sql`
- 变更点：
  - 删除 `deleted_at` 列；
  - partial index 调整为普通唯一/普通索引（或仅保留 `is_completed` 过滤）；
  - 保留 `next_index` 列作为分配来源。
- repo/app 接入 next_index 分配：
  - `TeamRepo.IncrementWorksetNextIndex`
  - `WorksetRepo.IncrementComicNextIndex`
  - `ComicRepo.IncrementChapterNextIndex`
  - `Create` 流程从 `Count(*)` 改为上述原子分配。
- repo 查询与删除语义同步改为硬删除：
  - 移除 `deleted_at IS NULL` 条件；
  - `Remove` 改为物理 `DELETE`。

## 本步结束验证

- 已执行：`go test ./...`（`refactor/`）通过。
- 已执行检索确认：Go 实现中已无 `deleted_at` 引用。

## 下步计划

- 继续按计划推进：
  - P2：用户头像更新时旧 key 清理防呆逻辑；
  - P3：补齐 Team/User 软删除策略与 Comic 封面上传完整流程（按计划优先级推进）。
