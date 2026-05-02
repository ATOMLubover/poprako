# 权限清单 (Refactor)

> 仅统计 `refactor/` 下实现。括号内为文件:行号。统一最后提交: 2026-04-29。

## 角色

| 角色 | 位  | 常量                | 域         | 备注                                           |
| ---- | --- | ------------------- | ---------- | ---------------------------------------------- |
| 图源 | 0   | `RoleRawProvider`   | Team+Chap  | 可推 `upload_complete`                         |
| 翻译 | 1   | `RoleTranslator`    | Team+Chap  | 可推 `translate_start\|complete`               |
| 校对 | 2   | `RoleProofreader`   | Team+Chap  | 可推 `proofread_start\|complete`               |
| 嵌字 | 3   | `RoleTypesetter`    | Team+Chap  | 可推 `typeset_start\|complete`                 |
| 美工 | 4   | `RoleRedrawer`      | Team+Chap  |                                                |
| 监修 | 5   | `RoleReviewer`      | Team+Chap  | 章级:创建/更新/删除分配, 可推**任意** workflow |
| 发布 | 6   | `RolePublisher`     | Team+Chap  | 可推 `publish_complete`                        |
| 管理 | 7   | `RoleAdmin`         | **仅Team** | 章分配不支持; 管作品集/漫画/章CRUD             |
| 超管 | —   | `User.IsSuperAdmin` | 全局       | bool字段, 创建/删除汉化组                      |

## 资源×操作

| 资源       | 操作              | 鉴权                               | Svc方法                                          |
| ---------- | ----------------- | ---------------------------------- | ------------------------------------------------ |
| Workset    | List              | 团队成员                           | `CanListWorkset` (`svc/workset.go:53`)           |
| Workset    | Create/Upd/Remove | Team Admin                         | `CanAdminWorkset` (`svc/workset.go:34`)          |
| Comic      | List              | 团队成员                           | `CanListComic` (`svc/comic.go:23`)               |
| Comic      | Create/Upd/Remove | Team Admin                         | `CanAdminComic` (`svc/comic.go:42`)              |
| Chapter    | List/GetPinned    | 团队成员                           | `CanListChapter` (`svc/chapter.go:23`)           |
| Chapter    | Create/Remove     | Team Admin                         | `CanAdminChapter` (`svc/chapter.go:42`)          |
| Chapter    | Upd(元数据)       | Team Admin                         | `CanAdminChapter` (`svc/chapter.go:42`)          |
| Chapter    | Workflow 推进     | **章级角色**                       | `CanTransiteWorkflow` (`svc/chapter.go:xx`)      |
| Assignment | ListByChapter     | **团队成员**                       | 链查 `chapter→comic→workset→team→member`         |
| Assignment | ListByUser        | 已认证                             | 无(自己的分配)                                   |
| Assignment | Upsert            | Chap Reviewer + 目标用户有Team角色 | `CanReviewAssignment` + `CanTakeAssignmentRoles` |
| Assignment | Delete            | Chap Reviewer                      | `CanReviewAssignment`                            |

## Workflow → 章角色映射

> Reviewer(监修)可推动所有阶段。

| Transition                               | 需要角色    |
| ---------------------------------------- | ----------- |
| `upload_complete`                        | RawProvider |
| `translate_start` / `translate_complete` | Translator  |
| `proofread_start` / `proofread_complete` | Proofreader |
| `typeset_start` / `typeset_complete`     | Typesetter  |
| `review_complete`                        | Reviewer    |
| `publish_complete`                       | Publisher   |

## 关键文件

| 文件                                  | 最后提交   |
| ------------------------------------- | ---------- |
| `refactor/.../svc/workset.go`         | 2026-04-29 |
| `refactor/.../svc/comic.go`           | 2026-04-29 |
| `refactor/.../svc/chapter.go`         | 2026-04-29 |
| `refactor/.../svc/assignment.go`      | 2026-04-29 |
| `refactor/.../app/impl/workset.go`    | 2026-04-29 |
| `refactor/.../app/impl/comic.go`      | 2026-04-29 |
| `refactor/.../app/impl/chapter.go`    | 2026-04-29 |
| `refactor/.../app/impl/assignment.go` | 2026-04-29 |
| `refactor/.../model/enum/role.go`     | 2026-04-27 |
| `refactor/.../model/aggr/role.go`     | 2026-04-27 |
