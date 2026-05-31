# 权限清单 (Legacy 对齐)

> 统计当前 `internal/` 实现。权限逻辑统一内化在 `internal/domain/svc/`，`app` 层仅做参数校验和编排。

## 角色

| 角色 | 位  | 常量                | 域         | 备注                                           |
| ---- | --- | ------------------- | ---------- | ---------------------------------------------- |
| 图源 | 0   | `RoleRawProvider`   | Team+Chap  | 可推 `upload_complete`；可回退 `upload_revert`                                                                       |
| 翻译 | 1   | `RoleTranslator`    | Team+Chap  | 可推 `translate_start\|complete`；可回退 `translate_start_revert\|translate_revert`                                 |
| 校对 | 2   | `RoleProofreader`   | Team+Chap  | 可推 `proofread_start\|complete`；可回退 `proofread_start_revert\|proofread_revert` + translate 系列               |
| 嵌字 | 3   | `RoleTypesetter`    | Team+Chap  | 可推 `typeset_start\|complete`；可回退 `typeset_start_revert\|typeset_revert`                                       |
| 美工 | 4   | `RoleRedrawer`      | Team+Chap  |                                                                                                                      |
| 监修 | 5   | `RoleReviewer`      | Team+Chap  | 章级:创建/更新/删除分配, 可推**任意** workflow；可回退**任意** workflow（`publish_revert` 不存在，发布不可回退）    |
| 发布 | 6   | `RolePublisher`     | Team+Chap  | 可推 `publish_complete`；发布不可回退                                                                               |
| 管理 | 7   | `RoleAdmin`         | **仅Team** | 章分配不支持; 管作品集/漫画/章CRUD             |
| 超管 | —   | `User.IsSuperAdmin` | 全局       | bool字段, 创建/删除汉化组                      |

## 资源×操作

| 资源 | 操作 | 鉴权 | Svc 方法 |
| --- | --- | --- | --- |
| Workset | List | 团队成员 | `WorksetSvc.CanListWorkset` |
| Workset | Create/Upd/Delete | Team Admin | `WorksetSvc.CanAdminWorkset` |
| Comic | List | 团队成员 | `ComicSvc.CanListComic` |
| Comic | Create/Upd/Delete | Team Admin | `ComicSvc.CanAdminComic` |
| Chapter | List/GetPinned/GetById | 团队成员 | `ChapterSvc.CanListChapter` |
| Chapter | Create/Upd/Delete | Team Admin | `ChapterSvc.CanAdminChapter` |
| Chapter | Workflow 推进 | 章级角色（Reviewer 全通） | `ChapterSvc.CanTransiteWorkflow` |
| Chapter | Workflow 回退 | 章级角色（见下表）；Publish 不可回退 | `ChapterSvc.CanRevertWorkflow` |
| Page | ListByChapter | 团队成员；若团队成员校验失败，允许章节 assignment 回退访问（legacy 兼容） | `PageSvc.CanListByChapter` |
| Page | ResvChapterPages | 仅图源或监修 | `PageSvc.CanResvPages` |
| Page | MarkImageUploaded | 仅图源 | `PageSvc.CanMarkImageUploaded` |
| Page | DeleteByChapterId | Team Admin | `ChapterSvc.CanAdminChapter` |
| Unit | ListByPage | 仅当前章节参与者 | `UnitSvc.CanListPageUnits` |
| Unit | SaveByPage | 仅当前章节翻译或校对 | `UnitSvc.CanEditPageUnits` |
| Assignment | ListByChapter | 团队成员；若团队成员校验失败，允许章节 assignment 回退访问（legacy 兼容） | `AssignmentSvc.CanListByChapter` |
| Assignment | ListByUser | 已认证（只看自己） | 无 |
| Assignment | Upsert/Delete | 仅章节监修；目标用户需具备对应 team 角色 | `AssignmentSvc.CanReviewAssignment` + `AssignmentSvc.CanTakeAssignmentRoles` |
| Assignment Invitation | Create/Delete/Accept | 仅章节监修；角色合法性校验不允许 `RoleAdmin` | `AssignmentSvc.CanReviewAssignment` + `AssignmentInvSvc.NewAssignmentInvCre` |
| Member | ListByTeam | 团队成员 | `MemberSvc.CanListMember` |
| Member | UpdateRole/Delete | Team Admin | `MemberSvc.CanAdminMember` |
| Member Invitation | List | 团队成员 | `MemberInvSvc.CanListMemberInv` |
| Member Invitation | Create/Delete/UpdateRole | Team Admin | `MemberInvSvc.CanAdminMemberInv` |
| Team | List/Create（全团队维度） | Super Admin | `TeamSvc.CanListTeam` |
| Team | Update/Delete/Avatar 操作 | Team Admin | `TeamSvc.CanAdminTeam` |

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

| 文件 |
| --- |
| `internal/domain/svc/workset.go` |
| `internal/domain/svc/comic.go` |
| `internal/domain/svc/chapter.go` |
| `internal/domain/svc/page.go` |
| `internal/domain/svc/unit.go` |
| `internal/domain/svc/assignment.go` |
| `internal/domain/svc/member.go` |
| `internal/domain/svc/member_inv.go` |
| `internal/domain/svc/team.go` |
| `internal/app/impl/workset.go` |
| `internal/app/impl/comic.go` |
| `internal/app/impl/chapter.go` |
| `internal/app/impl/page.go` |
| `internal/app/impl/unit.go` |
| `internal/app/impl/assignment.go` |
