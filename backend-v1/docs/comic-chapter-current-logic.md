# Comic & Chapter 当前实现逻辑完整描述

> 本文档对 Comic 和 Chapter 两个聚合根在各层的实现进行完整、精确的描述，作为状态机重构的前置参考。

---

## 1. 数据模型（Single Source of Truth: migrations/）

### 1.1 comic_table

| 列名                                     | 类型                    | 说明                                                     |
| ---------------------------------------- | ----------------------- | -------------------------------------------------------- |
| id                                       | TEXT PK                 | 漫画 ID                                                  |
| workset_id                               | TEXT FK → workset_table | 所属工作集                                               |
| index                                    | INTEGER                 | 在工作集中的序号（0-based），UNIQUE(workset_id, index)   |
| title                                    | TEXT                    | 漫画标题                                                 |
| author                                   | TEXT                    | 作者名                                                   |
| composed_title                           | TEXT                    | 组合标题格式：`【index】[author] title`，用于模糊搜索    |
| description                              | TEXT                    | 漫画描述                                                 |
| chapter_count                            | INTEGER                 | 章节计数（冗余副本，由 `SyncLatestChapterReplica` 同步） |
| latest_uploaded_at ~ latest_published_at | TIMESTAMPTZ (9列)       | **最新章节的工作流状态副本**，从最新章节复制而来         |
| creator_id                               | TEXT FK → user_table    | 创建者                                                   |
| last_active_at                           | TIMESTAMPTZ             | 最后活跃时间                                             |
| created_at / updated_at / deleted_at     | TIMESTAMPTZ             | 软删除                                                   |

**关键索引**：composed*title 上建有 `gin_trgm_ops` 索引用于模糊搜索；每种 latest*\*\_at 状态列均有独立索引。

### 1.2 chapter_table

| 列名                                 | 类型                  | 说明                                                    |
| ------------------------------------ | --------------------- | ------------------------------------------------------- |
| id                                   | TEXT PK               | 章节 ID                                                 |
| comic_id                             | TEXT FK → comic_table | 所属漫画                                                |
| index                                | INTEGER               | 在漫画中的序号（0-based），UNIQUE(comic_id, index DESC) |
| subtitle                             | TEXT                  | 章节副标题                                              |
| page_count                           | INTEGER               | 页面计数（冗余）                                        |
| total_unit_count                     | INTEGER               | 翻译单元总数（冗余）                                    |
| translated_unit_count                | INTEGER               | 已翻译单元数（冗余）                                    |
| proofread_unit_count                 | INTEGER               | 已校对单元数（冗余）                                    |
| uploaded_at                          | TIMESTAMPTZ           | 上传完成时间                                            |
| transalating_at                      | TIMESTAMPTZ           | 翻译开始时间                                            |
| translated_at                        | TIMESTAMPTZ           | 翻译完成时间                                            |
| proofreading_at                      | TIMESTAMPTZ           | 校对开始时间                                            |
| proofread_at                         | TIMESTAMPTZ           | 校对完成时间                                            |
| typesetting_at                       | TIMESTAMPTZ           | 嵌字开始时间                                            |
| typeset_at                           | TIMESTAMPTZ           | 嵌字完成时间                                            |
| reviewed_at                          | TIMESTAMPTZ           | 审核完成时间                                            |
| published_at                         | TIMESTAMPTZ           | 发布完成时间                                            |
| creator_id                           | TEXT FK → user_table  | 创建者                                                  |
| created_at / updated_at / deleted_at | TIMESTAMPTZ           | 软删除                                                  |

---

## 2. 领域层（domain/）

### 2.1 Workflow & WorkflowStatus

定义于 `model/workflow.go`：

- **Workflow**（工作流分类）共 6 种 + 1 特殊值：
  - `WorkflowUploading`、`WorkflowTranslating`、`WorkflowProofreading`、`WorkflowTypesetting`、`WorkflowReviewing`、`WorkflowPublishing`
  - `WorkflowNone`（标记不检查）

- **WorkflowStatus**（工作流状态）共 3 种 + 1 特殊值：
  - `WorkflowPending`（待处理）
  - `WorkflowInProgress`（进行中）
  - `WorkflowCompleted`（已完成）
  - `WorkflowUnset`（不更新/不筛选）

- **合法组合规则** (`IsValidWorkflowCombination`)：
  - Upload / Review / Publish：仅 `pending | completed`（二态）
  - Translate / Proofread / Typeset：`pending | in_progress | completed`（三态）

### 2.2 model/comic.go

- **ComicInfo**：读模型，包含所有列字段。`Workset`、`Creator` 为可选关联填充字段。
- **ComicCreation**：写模型（创建），包含 `WorksetID, Index, Title, Author, Description, CreatorID`。通过 `NewComicCreation` 构造。
- **ComicUpdate**：写模型（更新），包含 `ID, Title, Author, Description, ComposedTitle`。通过 `NewComicUpdate` 构造。
- **ComposeComicTitle(index, author, title)**：合成格式化标题。

> Comic 自身无工作流状态字段，其 `latest_*_at` 全部由 `SyncLatestChapterReplica` 从最新章节同步而来。

### 2.3 model/chapter.go

- **ChapterInfo**：读模型。包含所有列字段（含 9 个状态时间戳指针）以及可选填充的 `Comic`、`Creator`。
- **ChapterStats**：统计子模型，仅含 `TotalUnitCount / TranslatedUnitCount / ProofreadUnitCount`。
- **ChapterCreation**：写模型（创建），包含 `ComicID, Index, Subtitle, CreatorID`。
- **ChapterUpdate**：写模型（更新），**核心状态转换逻辑在 `NewChapterUpdate` 中实现**。

#### `NewChapterUpdate` 的状态转换规则

该函数接收当前 ChapterInfo 和 6 个 `*WorkflowStatus` 参数，根据规则计算每个时间戳的新值：

| 工作流    | pending                                   | in_progress                               | completed                                | unset（nil） |
| --------- | ----------------------------------------- | ----------------------------------------- | ---------------------------------------- | ------------ |
| Upload    | uploaded_at = nil                         | —                                         | uploaded_at = now                        | 保留当前值   |
| Translate | translating_at = nil, translated_at = nil | translating_at = now, translated_at = nil | 保留 translating_at, translated_at = now | 保留当前值   |
| Proofread | proofreading_at = nil, proofread_at = nil | proofreading_at = now, proofread_at = nil | 保留 proofreading_at, proofread_at = now | 保留当前值   |
| Typeset   | typesetting_at = nil, typeset_at = nil    | typesetting_at = now, typeset_at = nil    | 保留 typesetting_at, typeset_at = now    | 保留当前值   |
| Review    | reviewed_at = nil                         | —                                         | reviewed_at = now                        | 保留当前值   |
| Publish   | published_at = nil                        | —                                         | published_at = now                       | 保留当前值   |

**注意**：所有时间戳均使用 `*time.Time`（nil 表示未到达该状态），状态通过"时间戳是否为nil"来推理。

### 2.4 model/permission.go（Comic/Chapter 相关权限）

#### Comic 权限

| 权限            | 签名                                                                      | 规则                                         |
| --------------- | ------------------------------------------------------------------------- | -------------------------------------------- |
| PermComicList   | `(userID, teamID, OnLoadMemberInfo)`                                      | 只要是汉化组成员                             |
| PermComicCreate | `(userID, teamID, OnLoadMemberInfo)`                                      | 汉化组管理员                                 |
| PermComicUpdate | `(userID, comicID, OnLoadComicInfo, OnLoadWorksetInfo, OnLoadMemberInfo)` | 通过 comic→workset→team 链路确认汉化组管理员 |
| PermComicDelete | 同 Update                                                                 | 同 Update                                    |

#### Chapter 权限

| 权限              | 签名                                                                      | 规则                                   |
| ----------------- | ------------------------------------------------------------------------- | -------------------------------------- |
| PermChapterList   | `(userID, comicID, OnLoadMemberInfo, OnLoadComicInfo, OnLoadWorksetInfo)` | 只要是汉化组成员                       |
| PermChapterCreate | 同 List 签名                                                              | 汉化组管理员                           |
| PermChapterDelete | 同 List 签名                                                              | 汉化组管理员                           |
| PermChapterUpdate | `(userID, chapterID, []Workflow, OnLoadAssignmentInfo)`                   | **基于 Assignment 的细粒度工作流权限** |

`PermChapterUpdate` 的核心逻辑：

1. 通过 `OnLoadAssignmentInfo(chapterID, userID)` 加载当前用户在该章节的分工信息
2. 若用户是 **Reviewer**（监修），无条件允许更新所有工作流
3. 否则逐个检查要更新的工作流，用户必须持有对应角色：
   - Upload → RoleRawProvider
   - Translate → RoleTranslator
   - Proofread → RoleProofreader
   - Typeset → RoleTypesetter
   - Publish → RolePublisher

### 2.5 Repository 接口

#### ComicRepository

```
List / Get / Count / Create / Update / Delete（常规 CRUD）
LockByWorksetID                   — 悲观锁，创建时防并发
GetLatestChapterFirstPageOSSKey   — 获取封面（最新章节第一页 OSS Key）
SyncLatestChapterReplica          — 将最新章节状态同步到 comic 表
```

#### ChapterRepository

```
List / Get / Count / Create / Update / Delete（常规 CRUD）
LockByID / LockByComicID          — 悲观锁
GetStatsByID                       — 获取章节统计
UpdateStats                        — 仅更新三个统计字段（供 unit save 事务使用）
```

### 2.6 Domain Service

- `ResolveComicListIncludeSpec`：解析 includes[] 参数，决定是否需要 JOIN Workset / Creator
- `ResolveChapterListIncludeSpec`：解析 includes[] 参数，决定是否需要 JOIN Creator

---

## 3. 应用层（application/）

### 3.1 ComicApplication

依赖：`OSSClient, UserRepository, MemberRepository, WorksetRepository, ComicRepository`

#### ListComics

1. 参数验证
2. 通过 `WorksetRepository.Get` 获取 teamID
3. **鉴权**：`PermComicList(userID, teamID, ...)`
4. 解析 includes 规格
5. 构建查询选项（workset 筛选、排序 last_active_at DESC、分页、模糊搜索、最新章节状态筛选）
6. 执行查询、组装返回

#### CreateComic

1. 参数验证
2. 获取 teamID → **鉴权** PermComicCreate
3. 开事务 → `LockByWorksetID` → `Count` → 计算新 index → `Create`
4. 提交事务

#### UpdateComic

1. 参数验证
2. 获取现有漫画信息 → **鉴权** PermComicUpdate
3. 开事务 → `Update`（更新 title, author, description, composed_title）
4. 提交事务

#### DeleteComic

1. 验证存在 → **鉴权** PermComicDelete
2. 直接 `Delete`（无事务）

#### GetComicCover

1. 获取漫画 → 获取工作集 → **鉴权** PermComicList
2. `GetLatestChapterFirstPageOSSKey` → OSS 预签名 URL

### 3.2 ChapterApplication

依赖：`OSSClient, UserRepository, MemberRepository, WorksetRepository, ComicRepository, ChapterRepository, AssignmentRepository, PageRepository`

#### CreateComicChapter

1. 参数验证 → **鉴权** PermChapterCreate
2. 开事务 → `LockByComicID` → `Count` → 计算新 index → `Create`
3. **`SyncLatestChapterReplica`**（将新章节状态同步到漫画表）
4. 提交事务

#### ListComicChapters

1. 参数验证 → **鉴权** PermChapterList
2. 解析 includes → 构建查询（按 index DESC 排序、分页）→ 组装返回

#### UpdateChapter（**最复杂的操作**）

1. 参数验证
2. 开事务 → `LockByID` → `Get`（获取当前状态）
3. 收集要更新的工作流列表 → **鉴权** `PermChapterUpdate(userID, chapterID, workflows, ...)`
4. 调用 `NewChapterUpdate` 计算状态转换结果
5. **检测"上传完成"事件**：`targetChapter.UploadedAt == nil && chapterUpdate.UploadedAt != nil`
6. `ChapterRepository.Update` → `SyncLatestChapterReplica`
7. 若触发"上传完成"事件 → `handleChapterUploaded`：
   - 查询该章节所有 Assignment
   - 对每个 Assignment 的用户：确保 user_stats 存在 → `IncrementStats(pending: -1, completed: +1)`
   - 设置 afterCommitTask = 清理页面 OSS 资源
8. 提交事务
9. 执行 afterCommitTask（异步 goroutine）

#### handleChapterUploaded（上传完成后的副作用）

- 遍历章节所有 Assignment，更新每个用户的统计（pending_count -= 1, completed_count += 1）
- 设置后提交任务：清理章节所有页面的 OSS 资源

#### cleanupChapterPagesAfterUploaded

- 异步执行：查询该章节所有页面 → 逐个删除 OSS 对象 → 逐个删除页面记录
- 此操作在事务外执行，失败仅记录日志

#### DeleteComicChapter

1. 获取章节 → **鉴权** PermChapterDelete
2. 开事务 → `LockByID` → 重新 `Get` → `Delete` → `SyncLatestChapterReplica`
3. 提交事务

---

## 4. API 层（api/http/）

### Comic 路由

| 方法   | 路径                     | Handler       | 说明                           |
| ------ | ------------------------ | ------------- | ------------------------------ |
| GET    | /comics                  | ListComics    | 查询参数绑定 ListComicArgs     |
| POST   | /comics                  | CreateComic   | JSON Body 绑定 CreateComicArgs |
| PUT    | /comics/{comic_id}       | PatchComic    | 校验路径参数与 Body ID 一致    |
| DELETE | /comics/{comic_id}       | DeleteComic   |                                |
| GET    | /comics/{comic_id}/cover | GetComicCover | 返回 ComicCoverResult          |

### Chapter 路由

| 方法   | 路径                   | Handler            | 说明                             |
| ------ | ---------------------- | ------------------ | -------------------------------- |
| GET    | /chapters              | ListComicChapters  | 查询参数绑定 ListChapterArgs     |
| POST   | /chapters              | CreateComicChapter | JSON Body 绑定 CreateChapterArgs |
| PATCH  | /chapters/{chapter_id} | UpdateChapter      | chapter_id 从路径注入            |
| DELETE | /chapters/{chapter_id} | DeleteComicChapter |                                  |

---

## 5. 值对象层（value/）

- `ComicInfo`：API 返回用 DTO，时间戳为 `int64`（UnixMilli）
- `ChapterInfo`：API 返回用 DTO，状态时间戳为 `*int64`（nil 表示未到达）
- `ListComicArgs`：支持 6 种状态筛选、模糊搜索、分页、includes
- `ListChapterArgs`：支持分页、includes
- `CreateComicArgs` / `UpdateComicArgs`：含字段长度验证
- `CreateChapterArgs` / `UpdateChapterArgs`：UpdateChapterArgs 含 6 个可选状态字段，每个都验证合法组合

---

## 6. 基础设施层（infrastructure/repository/）

### SyncLatestChapterReplica 实现

核心 SQL 逻辑：

```sql
UPDATE comic_table SET
  chapter_count = (SELECT COUNT(*) FROM chapter_table WHERE ...),
  latest_uploaded_at = latest_chapter.uploaded_at,
  -- ... 其余 8 个状态字段
FROM (子查询: chapter_count) AS chapter_stats
LEFT JOIN LATERAL (子查询: 最新章节的状态时间戳, ORDER BY index DESC LIMIT 1) AS latest_chapter ON TRUE
WHERE comic.id = ? AND comic.deleted_at IS NULL
```

**含义**：每次章节状态变更、创建或删除后，将 comic 表中的 `latest_*_at` 字段更新为 index 最大的未删除章节的对应状态。

---

## 7. 当前架构的问题总结

### 7.1 状态管理散落在多处

- 状态转换规则硬编码在 `NewChapterUpdate` 函数中
- 状态的合法性验证分散在 `IsValidWorkflowCombination`（领域层）和 `UpdateChapterArgs.Validate`（值对象层）
- 无法表达"状态 A 只能从状态 B 转换"的前置条件约束

### 7.2 副作用（领域事件）与业务逻辑耦合

- `handleChapterUploaded` 逻辑直接嵌入 UpdateChapter 方法内部
- "上传完成"事件的检测是通过比较前后 UploadedAt 指针实现的 ad-hoc 检查
- 未来若需新增事件（如"翻译完成通知"、"发布完成触发推送"），必须不断在 UpdateChapter 中堆积 if 分支

### 7.3 Comic 的"最新章节副本"同步是命令式的

- 每个修改章节状态的操作都必须手动调用 `SyncLatestChapterReplica`
- 遗漏调用将导致 comic 表状态不一致
- 这本质上是一个"投影"，应由事件驱动自动维护

### 7.4 状态无前置条件校验

- 当前实现允许跳过状态（例如直接从 pending 跳到 completed 而无需经过 in_progress）
- 无法强制工作流顺序（例如必须先 upload 完成才能开始 translate）

### 7.5 事务边界与副作用混合

- `afterCommitTask` 模式虽然解决了"提交后执行"的问题，但属于 ad-hoc 方案
- 多个事件叠加时（如一次更新同时触发多个工作流完成），管理将变得困难
