# Team Announcement + Comment 全链路实现计划

## 1. 目标与已确认约束

本计划覆盖两个 team-scoped 功能：`announcement`（公告）和 `comment`（留言/寄语）。

已确认的业务约束如下：

1. `announcement`
   - 创建：仅 `team admin`
   - 列表：`team member` 可查看
   - 典型场景：进入工作区后拉取最新三条（由调用方传 `limit=3`）

2. `comment`
   - 创建：`team member` 可创建
   - 列表：`team member` 可查看
   - 展示：team 维度留言板列表

3. 通用约束
   - 两个功能都只做 `create + list`，本期不做 `update/delete`
   - `val` 层允许并默认携带 `User` 嵌入
   - 接口默认 `include user`，HTTP 层不暴露 `includes` 参数

---

## 2. 范围边界

本期会实现：

- migration 对齐（仅在必要时微调）
- domain model / query / enum / repo contract
- infra entity / repo / provider
- app iface / app impl / log wrapper / util validator
- api http handler + route 注册 + state + main DI
- swagger 文档生成

本期不实现：

- update / delete 链路
- count 接口
- event bus / async worker / 通知推送

---

## 3. 现状差距（实现前需先修）

1. `internal/app/val/announcement.go` 的 `Id` 字段 `json` tag 未闭合，当前会导致编译失败。
2. `internal/domain/model/aggr/announcement.go` 目前为空结构体，无法承载 create/list。
3. `comment` 还未落 aggregate、repo、app、http 等主干链路。

---

## 4. 数据与领域设计

### 4.1 Aggregate 设计

新增/完善以下聚合：

- `aggr.Announcement`
  - `Id`, `TeamId`, `UserId`, `User`, `Title`, `Content`, `CreatedAt`
- `aggr.AnnouncementCre`
  - `Id`, `TeamId`, `UserId`, `Title`, `Content`
- `aggr.Comment`
  - `Id`, `TeamId`, `UserId`, `User`, `Content`, `CreatedAt`
- `aggr.CommentCre`
  - `Id`, `TeamId`, `UserId`, `Content`

`User` 关系默认预加载，向上层直接提供作者信息。

### 4.2 Query 设计

新增：

- `query.ListAnnouncementOpt`
- `query.ListCommentOpt`

字段保持最小化：`TeamId` + `Pagi`。

### 4.3 Enum 设计

新增：

- `enum.AnnouncementIncl`
  - `AnnouncementInclUser = "user"`
- `enum.CommentIncl`
  - `CommentInclUser = "user"`

说明：HTTP 不透传 includes，但 repo 仍保留 typed include 扩展点；app 内部固定传 `InclUser`。

### 4.4 Domain Service 设计

新增：

- `svc.AnnouncementSvc`
  - `CanListAnnouncement`（member）
  - `CanCreateAnnouncement`（admin）
  - `NewAnnouncementCre`
- `svc.CommentSvc`
  - `CanListComment`（member）
  - `CanCreateComment`（member）
  - `NewCommentCre`

权限判定直接复用现有 `MemberRepo` + `ErrClsf` 范式。

---

## 5. Repo 与 Infra 方案

### 5.1 Domain Repo Contract

新增接口文件：

- `internal/domain/repo/announcement.go`
- `internal/domain/repo/comment.go`

方法最小集合：

- `List(...)`
- `Create(...)`

### 5.2 Provider 扩展

扩展：

- `internal/domain/repo/prov.go`
  - `AnnouncementRepo() AnnouncementRepo`
  - `CommentRepo() CommentRepo`
- `internal/infra/repo/prov.go`
  - 对应实现 `NewAnnouncementRepo` / `NewCommentRepo`

### 5.3 Entity 映射

新增：

- `internal/infra/repo/entity/announcement.go`
- `internal/infra/repo/entity/comment.go`

每个文件包含：

- table 常量
- read row（含 `User *UserRow` 关系）
- create row
- `To...Aggr()`
- `New...CreRowFromAggr()`

### 5.4 Repo 实现

新增：

- `internal/infra/repo/announcement.go`
- `internal/infra/repo/comment.go`

规则：

1. `List`
   - `WHERE team_id = ?`
   - `ORDER BY created_at DESC`
   - 支持 `offset/limit`
   - app 内默认携带 `InclUser`

2. `Create`
   - 插入 create row
   - 按 `id` 回查并返回完整 aggregate（包含 `User`）

---

## 6. App 层方案

### 6.1 App Interface

新增：

- `internal/app/announcement.go`
- `internal/app/comment.go`

每个接口仅暴露：

- `List(cx, currUid, args)`
- `Create(cx, currUid, args)`

### 6.2 Val 设计

完善/新增：

- `internal/app/val/announcement.go`
  - 修复 tag
  - `AnnouncementVal`（含 `User *UserVal`）
  - `ListAnnouncementArgs`
  - `CreateAnnouncementArgs`
  - `AnnouncementCreatedRes`
- `internal/app/val/comment.go`
  - `CommentVal` 增加 `User *UserVal`
  - `ListCommentArgs`
  - `CreateCommentArgs`
  - `CommentCreatedRes`

### 6.3 App Impl

新增：

- `internal/app/impl/announcement.go`
- `internal/app/impl/announcement_util.go`
- `internal/app/impl/announcement_log.go`
- `internal/app/impl/comment.go`
- `internal/app/impl/comment_util.go`
- `internal/app/impl/comment_log.go`

实现模式：

1. `List`
   - 参数校验（team_id + 分页）
   - 权限校验
   - repo list（固定 incl user）
   - aggregate -> val 组装

2. `Create`
   - 参数校验
   - `RunWithTxn`
   - 权限校验
   - svc 构建 cre
   - repo create
   - 返回 `id`

分页统一使用 `app_util.ClampOffsetLimit`。

---

## 7. HTTP / State / Main 接线

### 7.1 HTTP Handler

新增：

- `internal/api/http/announcement.go`
- `internal/api/http/comment.go`

路由：

- `GET /api/v1/announcements`
- `POST /api/v1/announcements`
- `GET /api/v1/comments`
- `POST /api/v1/comments`

说明：

- 不增加 `includes` query 参数
- 继续沿用当前 auth、`HttpRes` 包裹、错误映射范式

### 7.2 State 扩展

更新 `internal/api/state/state.go`：

- 增加 `AnnouncementApp`
- 增加 `CommentApp`
- `NewAppState(...)` 参数与构造同步扩展

### 7.3 Main DI

更新 `main.go`：

- 构造 `announcementRepo` / `commentRepo`
- 构造 `announcementSvc` / `commentSvc`
- 构造 `announcementApp` / `commentApp` 及其 log 包装
- 注入 `state.NewAppState(...)`

---

## 8. Migration 与索引策略

现有 migration 已满足核心需求：

- 表结构完整
- team + created_at desc 索引已具备

本期仅在发现编译或运行映射不一致时做最小修正，不追加额外业务字段。

---

## 9. 实施顺序

1. 修 `val/announcement.go` 编译问题（tag）
2. 完成 `aggr/query/enum`
3. 完成 `domain/repo` 与 `prov` 扩展
4. 完成 `infra/entity/repo/prov`
5. 完成 `domain/svc`
6. 完成 `app iface + impl + util + log`
7. 完成 `http + route + state + main`
8. 生成 swagger（`just swag`）
9. 编译与回归（`go test ./...`）

---

## 10. Done Criteria

满足以下条件即视为完成：

1. `announcement/comment` 的 create/list 全链路可编译并可运行。
2. 权限语义符合已确认规则：
   - announcement create=admin, list=member
   - comment create=member, list=member
3. 两类 list 返回中都包含 `User` 嵌入信息。
4. HTTP 不暴露 includes 参数，但服务端默认加载 `User`。
5. swagger 文档与实际路由一致并生成成功。
