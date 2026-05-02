# OSS 资源管理与软删除去除方案

## 一、问题诊断

### 1.1 OSS Worker 后台协程缺失

对比旧项目 `internal/app/worker/oss_worker.go`，重构项目缺少后台协程来消费 `t_oss_message` 中的 pending 消息。

`OssMsgRepo` 的 `ClaimPending`、`MarkCompleted`、`MarkPending`、`ResetStuck`、`CleanCompleted` 以及 `oss_iface.Cleaner.DelBatch` 均已实现，但没有消费循环把它们串起来。

**后果：** `SavePendingDel` 入队的 OSS 删除消息永远不会被执行，远程对象泄漏；`SavePendingCre` 的过期消息不会被清理；completed 消息无限增长。

### 1.2 删除操作级联缺失

**Comic.Remove：** 只软删除 comic + 更新计数，不处理其下 chapter/page/assignment 及关联 OSS 资源。

**Workset.Remove：** 同理，不处理其下 comic/chapter/page 及 OSS 资源。

### 1.3 资源层级与 OSS 关联

```
Team (OSS: avatar_key)
  └── Workset
        └── Comic (OSS: cover_key)
              └── Chapter
                    └── Page (OSS: image_key)
                          └── Unit

User (OSS: avatar_key)
```

`t_page` 和 `t_unit` 已使用硬删除 + `ON DELETE CASCADE`，无需改动。

## 二、软删除去除：`next_index` 替代方案

### 2.1 当前使用软删除的原因

`t_workset`、`t_comic`、`t_chapter` 使用软删除的核心原因：子资源创建时需要通过 `COUNT(*) WHERE deleted_at IS NULL` 获取下一个 index。如果硬删除，COUNT 会减少，新资源可能复用已删除资源的 index，导致"不重不漏"被破坏。

### 2.2 替代方案：`UPDATE parent SET next_index = next_index + 1 RETURNING next_index`

Schema 中已预留了 `next_index` 列：

| 父表        | 列名                 | 用途                  |
| ----------- | -------------------- | --------------------- |
| `t_team`    | `workset_next_index` | 分配 workset 的 index |
| `t_workset` | `comic_next_index`   | 分配 comic 的 index   |
| `t_comic`   | `chapter_next_index` | 分配 chapter 的 index |

**创建流程（以 comic 为例）：**

```
事务内:
1. SELECT ... FOR UPDATE 锁 workset 行
2. index := UPDATE t_workset SET comic_next_index = comic_next_index + 1
   WHERE id = ? RETURNING comic_next_index
3. INSERT INTO t_comic (..., index) VALUES (..., index)
提交事务
```

若事务失败回滚，UPDATE 自动回滚，next_index 不变。

**删除：** 直接物理 DELETE。空洞不处理——index 仅用于排序展示，不需要连续。

**为什么比 COUNT 好：**

- UPDATE ... RETURNING 是单行操作，COUNT 需要扫索引
- UPDATE 自带行锁，COUNT + INSERT 之间有竞态窗口需要额外加锁
- next_index 单调递增，绝不会与已有 index 冲突

### 2.3 逐表改造

#### t_workset

- 去除 `deleted_at` 列
- `uidx_workset_team_id_index` 从 partial unique → 普通 `UNIQUE (team_id, index)`
- `Create` 中 `COUNT(*)` → `teamRepo.IncrementWorksetNextIndex(teamId)`
- `Remove` → 物理 DELETE（含级联，见第三章）
- 所有查询去除 `deleted_at IS NULL`

#### t_comic

- 去除 `deleted_at` 列
- `uidx_comic_workset_id_index` 和其余 partial index → 普通 index
- `Create` 中 `comicRepo.Count(...)` → `worksetRepo.IncrementComicNextIndex(worksetId)`
- `Remove` → 物理 DELETE（含级联，见第三章）
- 所有查询去除 `deleted_at IS NULL`

#### t_chapter

- 去除 `deleted_at` 列
- 三个 partial index → 普通 index
- `Create` 中 `chapterRepo.Count(...)` → `comicRepo.IncrementChapterNextIndex(comicId)`
- `Remove` → 物理 DELETE（含级联）
- 所有查询去除 `deleted_at IS NULL`

#### t_user / t_team

`deleted_at` 与 index 分配无关，保护的是 `qid`/`nickname`/`name` 唯一性。去除软删除后可改为删除时追加 `_deleted_{timestamp}` 后缀释放原名，或保留软删除。此项独立于本期核心目标。

#### t_page / t_unit

已使用硬删除 + `ON DELETE CASCADE`，无需改动。

### 2.4 Index 空洞说明

删除资源后其 index 永久废弃，`next_index` 不复用。这导致列表中存在空洞（如 workset 的 comics 为 index 0, 1, 3，2 已被删）。

这是**可接受的**：index 仅用于排序和展示位置，不是主键，不需要连续。用户看到 0, 1, 3 不影响功能。

## 三、级联删除方案

### 3.1 核心原则

顶层资源删除时，在同一事务内逐层入队所有 OSS 删除消息，然后逐层物理删除 DB 记录。事务提交后，OSS Worker 异步执行远程清理。Page/Unit 的级联删除依赖 `ON DELETE CASCADE`，不在应用层显式处理。

### 3.2 Comic.Remove

```
事务内:
1. 查询 comic，锁 workset 行
2. 权限校验
3. 遍历 comic 下所有 chapters:
   a. 遍历 chapter 下所有 pages:
      - 若 image_key 非空 → ossMsgRepo.SavePendingDel(OssResPageImage, pageId, [imageKey])
   b. DELETE pages WHERE chapter_id = ? — unit 由 FK CASCADE 自动删除
   c. DELETE assignments WHERE chapter_id = ?
   d. DELETE chapter WHERE id = ?
4. 若 comic.cover_key 非空 → ossMsgRepo.SavePendingDel(OssResComicCover, comicId, [coverKey])
5. DELETE comic
6. workset.comic_count -= 1
提交事务
```

### 3.3 Chapter.Remove

```
事务内:
1. 查询 chapter，锁 comic 行
2. 权限校验
3. 遍历 chapter 下所有 pages:
   - 若 image_key 非空 → ossMsgRepo.SavePendingDel
4. DELETE pages WHERE chapter_id = ? — unit 由 FK CASCADE 自动删除
5. DELETE assignments WHERE chapter_id = ?
6. DELETE chapter
7. comic.chapter_count -= 1
8. 发布 ChapterRemovedEv
提交事务
```

### 3.4 Workset.Remove

```
事务内:
1. 查询 workset，锁 team 行
2. 权限校验
3. 遍历 workset 下所有 comics:
   a. 对每个 comic 执行 Comic.Remove 级联逻辑（章节/页面/OSS入队）
   b. DELETE comic
4. DELETE workset
提交事务
```

### 3.5 Page.RemoveByChapterId

现有实现正确处理了 OSS 入队 + page 硬删除 + unit CASCADE。无需改动。

## 四、OSS Worker 实现

### 4.1 架构

```
refactor/internal/infra/worker/oss_worker.go
```

依赖 `repo_iface.OssMsgRepo` + `oss_iface.Cleaner`。

### 4.2 三个循环

| 循环          | 间隔 | 职责                                                                                 |
| ------------- | ---- | ------------------------------------------------------------------------------------ |
| `consumeLoop` | 5min | Claim pending del → `DelBatch(objKeys)` → MarkCompleted / MarkPending                |
|               |      | Claim pending cre → 仅处理已过期 → `DelBatch(objKeys)` → MarkCompleted / MarkPending |
| `reclaimLoop` | 5min | `ResetStuck`：processing 超 15min → 重置为 pending                                   |
| `purgeLoop`   | 24h  | `CleanCompleted`：删除超过 24h 的 completed 消息                                     |

### 4.3 Create 消息处理修正

旧 worker 对未过期 create 消息返回 nil（成功）导致被误标 completed。

**修正：** 未过期时 `MarkPending(id, "not yet expired", expireAt)`，visible_at 设为 expire_at，过期前不再被重复 claim。

### 4.4 main.go 启动方式

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()
go worker.NewOssWorker(ossMsgRepo, ossClient).Run(ctx)
```

## 五、其他补充

### 5.1 OssMsgSvc 补充 SavePendingDel

当前调用方直接调 `ossMsgRepo.SavePendingDel`，绕过 domain service。补充 `OssMsgSvc.SavePendingDel` 统一入口。

### 5.2 头像/封面更新时的旧文件清理

Avatar 及 Cover 类 OSS 资源的 key 是固定的（如 `user_avatar/{userId}.{ext}`），`ResvAvatar` 生成的新 key 与 DB 中旧 key 通常相同。若不加判断直接入队 `SavePendingDel`，会删除刚刚上传的文件。

**防呆逻辑：** 入队删除前比对旧 key 与新 key，仅当两者不同时才入队：

```go
// ResvAvatar 事务内:
oldKey := user.AvatarKey
newKey := fmt.Sprintf("user_avatar/%s.%s", args.UserId, args.FileExt)
if oldKey != "" && oldKey != newKey {
    ossMsgRepo.SavePendingDel(OssResUserAvatar, args.UserId, []string{oldKey})
}
```

当前 key 生成规则固定，旧 key 与新 key 恒相等，此逻辑不会触发。保留作为未来 key 规则变更时的安全防护。

### 5.3 Comic 封面上传流程

Migration 已定义 `cover_key` + `cover_uploaded`，app 层待实现。

## 六、实施优先级

| 优先级 | 任务                                                             |
| ------ | ---------------------------------------------------------------- |
| **P0** | 实现 OSS Worker 后台协程                                         |
| **P1** | 实现 Comic.Remove / Workset.Remove 级联删除（需要处理 OSS 资源） |
| **P1** | 去除 t_workset/t_comic/t_chapter 软删除 + 接入 next_index        |
| **P1** | OssMsgSvc 补充 SavePendingDel                                    |
| **P2** | User 换头像清理旧头像                                            |
| **P3** | User/Team 软删除独立处理                                         |
| **P3** | Comic 封面上传流程                                               |
