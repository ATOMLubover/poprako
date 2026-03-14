# Info 层重构进度

## 背景

参照已完成的 workset 链路，对 user、team、member 域的 info 层进行统一重构。  
核心目标：将原本硬编码的 JOIN 查询改造为由 `includes[]` 参数驱动的可选嵌套查询机制，  
并统一各域的信息模型（value、model、entity），彻底消除散乱分裂的返回类型。

---

## Workset 链路重构参照模式

| 层次             | 文件                 | 核心变更                                                                             |
| ---------------- | -------------------- | ------------------------------------------------------------------------------------ |
| domain/service   | `workset_include.go` | 纯函数 `ResolveWorksetListIncludeSpec`，输出 `WorksetListIncludeSpec{NeedTeam bool}` |
| value            | `workset.go`         | `WorksetInfo.Team *TeamInfo`（可选指针）                                             |
| model            | `model/workset.go`   | `WorksetInfo.Team *model.TeamInfo`（可选指针）                                       |
| query_option     | `workset.go`         | `WorksetQuery().IncludeTeamInfo()`                                                   |
| entity           | `entity/workset.go`  | `WorksetInfoRow` 含 team 别名列，`ToWorksetInfo` 按零值判断是否填充 Team             |
| infra/repository | `workset.go`         | 统一 `List(executor, options...)`                                                    |
| app              | `workset.go`         | 先解析 spec，再按需追加 query option，最后注入 presigned URL                         |

---

## 从 App 层角度：需修改的 Domain 汇总

### 1. Member Domain（主要工作）

**问题：**

- `ListMembers` 目前硬编码 LEFT JOIN user（无论前端是否需要用户信息）
- `ListMyMembers` 目前硬编码 LEFT JOIN team（无论前端是否需要组信息）
- 存在两个功能重复的 repo 方法：`ListProfiles` 与 `ListProfilesWithUserInfo`
- `MemberProfile`（内嵌 UserInfo）和 `MemberWithTeamInfo` 是两个分裂的返回类型，无法统一表达
- 字段命名存在 typo：`AssignedPublishererAt`（多一个 r）

**App 层接口变更：**

```go
// 变更前
ListMyMembers(scope, currentUserID string) ([]value.MemberWithTeamInfo, error)
ListMembers(scope, currentUserID string, args value.ListTeamMemberArgs) ([]value.MemberProfile, error)

// 变更后
ListMyMembers(scope, currentUserID string, args value.ListMyMemberArgs) ([]value.MemberInfo, error)
ListMembers(scope, currentUserID string, args value.ListTeamMemberArgs) ([]value.MemberInfo, error)
// ListTeamMemberArgs 新增 Includes []string
```

**新增 includes 支持：**

- `includes[]=user` → 在成员列表中携带用户信息（头像、昵称等）
- `includes[]=team` → 在"我的成员"列表中携带汉化组信息

---

### 2. Team Domain（已接入 include）

**问题：**

- `ListAllTeams` 和 `ListMyTeams` 没有任何 args，无法支持 includes 扩展
- 缺少 include spec 基础设施

**App 层接口变更：**

```go
// 变更前
ListAllTeams(scope, currentUserID string) ([]value.TeamInfo, error)
ListMyTeams(scope, currentUserID string) ([]value.TeamInfo, error)

// 变更后
ListAllTeams(scope, currentUserID string, args value.ListTeamArgs) ([]value.TeamInfo, error)
ListMyTeams(scope, currentUserID string, args value.ListMyTeamArgs) ([]value.TeamInfo, error)
```

**新增 includes 支持：**

- `includes[]=member` → 在 team 列表中携带成员信息
- `includes[]=member.user` → 在成员信息中继续携带用户信息

---

### 3. User Domain（已接入 include）

`GetUser` / `GetMyUser` 增加 query args 并支持 include 规格解析。

**新增 includes 支持：**

- `includes[]=member` → 在 user 详情中携带成员信息
- `includes[]=member.team` → 在成员信息中继续携带所属 team 信息

---

## 各域变更文件列表

### Member Domain（全量重构）

| 文件                                                        | 变更说明                                                                                                                            |
| ----------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `internal/domain/model/user.go`                             | 新增 `NewUserInfo` 构造函数（补齐缺失）                                                                                             |
| `internal/domain/model/member.go`                           | 新增统一模型 `MemberWithInfo`（替代 `MemberWithUserInfo` + `MemberWithTeamInfo`），修正 typo，更新 `NewMemberUpdate` 签名           |
| `internal/domain/repository/member.go`                      | 接口用 `List` 替换 `ListProfiles` + `ListWithTeamInfo`；`GetProfile` 返回类型改为 `model.MemberWithInfo`                            |
| `internal/domain/service/member_include.go`                 | **新建**：`ResolveMemberListIncludeSpec` / `ResolveMyMemberListIncludeSpec`                                                         |
| `internal/infrastructure/repository/entity/member.go`       | 新增 `MemberWithInfoRow`（合并 user + team 别名列）；新增 `ToMemberWithInfo` 转换函数                                               |
| `internal/infrastructure/repository/query_option/member.go` | 新增 `MemberQuery().IncludeUserInfo()` 和 `MemberQuery().IncludeTeamInfo()`                                                         |
| `internal/infrastructure/repository/member.go`              | 实现统一 `List`（无硬编码 JOIN），更新 `GetProfile` 实现，移除重复方法                                                              |
| `internal/value/member.go`                                  | 以统一 `MemberInfo` 替换 `MemberProfile` + `MemberWithTeamInfo`；新增 `ListMyMemberArgs`；`ListTeamMemberArgs` 增加 `Includes` 字段 |
| `internal/application/member.go`                            | 更新 `ListMembers`、`ListMyMembers`（含 include spec 驱动）、`UpdateMemberRole`、`RemoveMember`                                     |
| `internal/api/http/member.go`                               | 更新 `ListMyMembers` handler（读取 query args），更新 godoc                                                                         |

### Team Domain（已接入 include）

| 文件                                      | 变更说明                                                                          |
| ----------------------------------------- | --------------------------------------------------------------------------------- |
| `internal/value/team.go`                  | 新增 `ListTeamArgs`、`ListMyTeamArgs`（含 `Includes []string`）                   |
| `internal/domain/service/team_include.go` | **新建**：`ResolveTeamListIncludeSpec`                                            |
| `internal/application/team.go`            | 接口和实现更新 `ListAllTeams` / `ListMyTeams` 签名，按 include 加载成员与成员用户 |
| `internal/api/http/team.go`               | 更新 handlers，读取 query args                                                    |

### User Domain（已接入 include）

| 文件                                      | 变更说明                                                             |
| ----------------------------------------- | -------------------------------------------------------------------- |
| `internal/domain/service/user_include.go` | **新建**：`ResolveUserDetailIncludeSpec`                             |
| `internal/value/user.go`                  | 新增 `GetUserArgs` / `GetMyUserArgs`，`UserInfo` 增加 `Members` 字段 |
| `internal/application/user.go`            | `GetUser` / `GetMyUser` 支持 include 编排，按需加载成员及成员团队    |
| `internal/api/http/user.go`               | `GetUserByID` / `GetMyUser` 读取 includes[] query                    |

---

## 重构进度

| 任务                            | 状态                               |
| ------------------------------- | ---------------------------------- |
| 创建本文档                      | ✅ 完成                            |
| model/user.go: NewUserInfo      | ✅ 完成                            |
| model/member.go: MemberWithInfo | ✅ 完成                            |
| service/member_include.go       | ✅ 完成                            |
| domain/repository/member.go     | ✅ 完成                            |
| entity/member.go                | ✅ 完成                            |
| query_option/member.go          | ✅ 完成                            |
| infra/repository/member.go      | ✅ 完成                            |
| value/member.go                 | ✅ 完成                            |
| application/member.go           | ✅ 完成                            |
| api/http/member.go              | ✅ 完成                            |
| value/team.go                   | ✅ 完成                            |
| application/team.go             | ✅ 完成                            |
| api/http/team.go                | ✅ 完成                            |
| 编译验证（Member/Team/User）    | ✅ 完成（`go build ./...` 无错误） |

---

## Comic Domain 重构（第二阶段）

**背景：**

- `comic_table` SQL 中有 `workset_id` 列（FK），但原 Go 实体用的是 `team_id`（错误）
- 需修正字段映射，并支持通过 `includes[]=workset` 按需嵌套工作集信息

| 文件                                                       | 变更说明                                                                                                                                           |
| ---------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/domain/model/comic.go`                           | `ComicInfo.TeamID` → `WorksetID`；新增 `Workset *WorksetInfo` 可选字段；`ComicCreation.TeamID` → `WorksetID`                                       |
| `internal/infrastructure/repository/entity/comic.go`       | 实体列 `workset_id`；加入 workset 别名列；`ToComicInfo` 改为 `ToComicWithInfo`                                                                     |
| `internal/infrastructure/repository/query_option/comic.go` | `FilterByTeamID` → `FilterByWorksetID`；新增 `IncludeWorksetInfo()`                                                                                |
| `internal/domain/repository/comic.go`                      | `LockByTeamID` → `LockByWorksetID`                                                                                                                 |
| `internal/domain/service/comic_include.go`                 | **新建**：`ResolveComicListIncludeSpec`                                                                                                            |
| `internal/infrastructure/repository/comic.go`              | 始终 JOIN workset 获取 team_id；`LockByWorksetID`；`Create` 用 `WorksetID`                                                                         |
| `internal/value/comic.go`                                  | `ListTeamComicArgs` → `ListComicArgs`（`WorksetID`，`Includes[]string`）；`CreateComicArgs.TeamID` → `WorksetID`；`NewComicInfoFromModel` 补全字段 |
| `internal/application/comic.go`                            | `ListTeamComics` → `ListComics`；注入 `worksetRepository` 依赖；uses include spec                                                                  |
| `internal/api/http/comic.go`                               | `ListTeamComics` → `ListComics`；更新 godoc                                                                                                        |
| `main.go`                                                  | `NewComicApplication` 增加 `worksetRepository` 参数                                                                                                |

---

## Chapter Domain 重构（第二阶段）

**背景：**

- `ChapterDetail` 未携带创建者信息，前端无法显示
- 支持通过 `includes[]=creator` 按需嵌套创建者用户信息

| 文件                                                         | 变更说明                                                                                 |
| ------------------------------------------------------------ | ---------------------------------------------------------------------------------------- |
| `internal/domain/model/chapter.go`                           | `ChapterDetail` 增加 `Creator *UserInfo`（可选）                                         |
| `internal/infrastructure/repository/entity/page.go`          | 新增 `ChapterWithInfoRow`（含 creator 别名列）；新增 `ToChapterWithInfo`                 |
| `internal/infrastructure/repository/query_option/chapter.go` | 新增 `ChapterQuery().IncludeCreatorInfo()`                                               |
| `internal/domain/service/chapter_include.go`                 | **新建**：`ResolveChapterListIncludeSpec`                                                |
| `internal/infrastructure/repository/chapter.go`              | `List`/`Get` 改用 `ChapterWithInfoRow`                                                   |
| `internal/value/chapter.go`                                  | `ListChapterArgs` 增加 `Includes []string`；`NewChapterInfoFromModel` 补全 `CreatorInfo` |
| `internal/application/chapter.go`                            | `NewChapterApplication` 增加 `ossClient` 依赖；`ListComicChapters` 支持 include spec     |
| `internal/api/http/chapter.go`                               | 更新 godoc，增加 `includes[]` 参数说明                                                   |
| `main.go`                                                    | `NewChapterApplication` 增加 `ossClient` 参数                                            |

---

## Assignment Domain 重构（第二阶段）

**背景：**

- DB 中 `assignment_table` 和 `member_table` 均有 `assigned_redrawer_at` 列，Go 代码中缺失
- `ListWithChapterInfo` JOIN 中引用了已不存在的 `comic_table.team_id`（应通过 workset 获取）
- `ListChapterAssignmentArgs` 需要支持 `includes[]` 参数

| 文件                                                            | 变更说明                                                                                                                                                                                                                                 |
| --------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/domain/model/assignment.go`                           | `AssignmentInfo`/`AssignmentCreation`/`AssignmentUpdate`/`AssignmentWithUserInfo`/`AssignmentWithChapterInfo` 均补全 `AssignedRedrawerAt *time.Time`                                                                                     |
| `internal/infrastructure/repository/entity/assignment.go`       | `AssignmentInfoRow`/`AssignmentInsertRow` 补全 `assigned_redrawer_at`；`ToAssignmentWithUserInfo`/`ToAssignmentWithChapterInfo` 更新；`AssignmentWithChapterAndComicRow` 新增 `comic_workset_id`，通过 JOIN workset_table 获取 `team_id` |
| `internal/infrastructure/repository/query_option/assignment.go` | 无变更（预留扩展）                                                                                                                                                                                                                       |
| `internal/domain/service/assignment_include.go`                 | **新建**：`ResolveAssignmentListIncludeSpec`                                                                                                                                                                                             |
| `internal/infrastructure/repository/assignment.go`              | `ListWithChapterInfo` SELECT 增加 `comic_workset_id`，JOIN `workset_table`；`Create`/`Update` 补全 `assigned_redrawer_at`                                                                                                                |
| `internal/value/assignment.go`                                  | `AssignmentWithUserInfo`/`AssignmentWithChapterInfo` 补全 `AssignedRedrawerAt *int64`；`ListChapterAssignmentArgs` 增加 `Includes []string`                                                                                              |
| `internal/api/http/assignment.go`                               | 更新 godoc，增加 `includes[]`/`offset`/`limit` 参数说明                                                                                                                                                                                  |

---

## 全量编译验证

| 任务                   | 状态                               |
| ---------------------- | ---------------------------------- |
| Comic domain 重构      | ✅ 完成                            |
| Chapter domain 重构    | ✅ 完成                            |
| Assignment domain 重构 | ✅ 完成                            |
| 最终编译验证           | ✅ 完成（`go build ./...` 无错误） |
