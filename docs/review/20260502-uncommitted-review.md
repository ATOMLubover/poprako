# 2026-05-02 未提交改动 review

`go test ./...`（在 `refactor/` 下）通过，但以下问题仍成立。

## High

### 1. `t_chapter` 失去“每个 comic 只能有一个 pinned chapter”的数据库约束

- 锚点：`refactor/migrations/20260427000003_create-chapter-table.up.sql`、`refactor/internal/infra/repo/chapter.go`
- 现状：migration 只保留了普通索引 `("comic_id", "pinned")`；`Create` / `Update` 仍靠“先清旧 pinned，再写新 pinned”维持唯一性。
- 问题：并发创建/置顶时可以出现同一 `comic_id` 下多条 `pinned = TRUE`。`FindPinnedByComicId` 与 comic workflow filter 都默认 pinned 唯一，这个前提已经被破坏。

### 2. `TeamRepo` 未过滤 `deleted_at`，软删除 team 仍可分配新 workset index

- 锚点：`refactor/migrations/20260425162956_create-team-table.up.sql`、`refactor/internal/infra/repo/team.go`、`refactor/internal/app/impl/workset.go`
- 现状：`t_team` 仍然保留 `deleted_at`，但 `GetById` / `IncrementWorksetNextIndex` 都只按 `id` 查询。
- 问题：软删除 team 仍可被 `Create` 路径继续使用；这既违反 soft-delete 查询约定，也会让已删除 team 继续产生新的 `workset`。

### 3. `Comic.Remove` / `Workset.Remove` 直接删 assignment/chapter，但没有补发移除事件

- 锚点：`refactor/internal/app/impl/comic.go`、`refactor/internal/app/impl/workset.go`、`refactor/main.go`、`refactor/internal/infra/event/assignment_chapter.go`
- 现状：删除流程直接调用 `assignmentRepo.DeleteByChapterId` 与 `chapterRepo.Remove`。
- 问题：`AssignmentRemovedEv` / `ChapterRemovedEv` 没有发出；但 `main.go` 仍注册了对应的 user-stats handler。删 comic/workset 时，用户统计会静默漂移。

## Medium

### 4. 级联删除遍历没有分页循环，还带固定上限

- 锚点：`refactor/internal/app/impl/comic.go`、`refactor/internal/app/impl/chapter.go`
- 现状：`comic.Remove` 对 chapter/page 直接写死 `Limit: 10000`；`chapter.Remove` 对 assignments 写死 `Limit: 500`。
- 问题：超过上限的记录仍会被 FK 级联删掉，但对应的 page image OSS 删除消息、`assignedUserIds` 收集不会完整，副作用不完整。

### 5. `chapter port import` 对“未分配用户”返回 500，而不是 403

- 锚点：`refactor/internal/app/impl/chapter_port.go`、`refactor/internal/infra/repo/assignment.go`、`refactor/internal/api/http/chapter_port.go`
- 现状：`Import` 先调 `assignmentRepo.GetByChapterUserId`；not found 被直接当成 server error。
- 问题：未被分配的调用方本应命中“仅翻译/校对可导入”的业务拒绝，现在会错误返回 500，和接口文档也不一致。

### 6. 语义约定漂移：现在是“硬删除”，但接口仍叫 `Remove`

- 锚点：`refactor/internal/app/{chapter,comic,workset}.go`、`refactor/internal/domain/repo/{chapter,comic,workset}.go`、`refactor/internal/infra/repo/{chapter,comic,workset}.go`
- 现状：当前实现已经是物理 `DELETE`，但 app/repo 契约仍叫 `Remove`。
- 问题：这和你在 `AGENTS.md` 里固定的语义冲突：`Delete` = hard delete，`Remove` = soft delete。名称已经开始误导契约阅读者。

## Style / Contract Drift

### 7. 新增 app/http 代码有明显风格与契约漂移

- 锚点：
  - `refactor/internal/app/comic.go`
  - `refactor/internal/app/impl/comic.go`
  - `refactor/internal/app/impl/workset.go`
  - `refactor/internal/app/impl/chapter.go`
  - `refactor/internal/api/http/comic.go`
  - `refactor/internal/api/http/chapter_port.go`
- 现状：
  - `ComicApp` 新增公开方法未按 `app-layer-style` 使用多行签名；
  - 多个新改的大型事务函数缺少 step comments / 语句留白；
  - Swagger body 仍声明 `comic_id` / `chapter_id`，但 handler 实际从 path 回填并覆盖 body 字段。
- 问题：这批改动虽然能编译，但已经偏离你通过 skill 明确给出的 app-layer / HTTP 契约风格，继续迭代会越修越散。
