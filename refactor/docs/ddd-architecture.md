# DDD 架构文件组织

```txt
➜  refactor git:(refactor/simplify) ✗  tree .
.
├── docs
│   ├── ddd-architecture.md
│   ├── perm-list.md
│   └── todos.md
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── http
│   │   │   ├── assignment.go
│   │   │   ├── assignment_invitation.go
│   │   │   ├── auth.go
│   │   │   ├── chapter.go
│   │   │   ├── comic.go
│   │   │   ├── http.go
│   │   │   ├── middleware
│   │   │   │   ├── auth.go
│   │   │   │   └── log.go
│   │   │   ├── res
│   │   │   │   └── res.go
│   │   │   ├── sys_mail.go
│   │   │   ├── team.go
│   │   │   ├── user.go
│   │   │   ├── util.go
│   │   │   └── workset.go
│   │   └── state
│   │       └── state.go
│   ├── app
│   │   ├── assignment.go
│   │   ├── assignment_invitation.go
│   │   ├── chapter.go
│   │   ├── comic.go
│   │   ├── impl
│   │   │   ├── assignment.go
│   │   │   ├── assignment_invitation.go
│   │   │   ├── assignment_invitation_log.go
│   │   │   ├── assignment_log.go
│   │   │   ├── assignment_util.go
│   │   │   ├── chapter.go
│   │   │   ├── chapter_log.go
│   │   │   ├── chapter_util.go
│   │   │   ├── comic.go
│   │   │   ├── comic_log.go
│   │   │   ├── comic_util.go
│   │   │   ├── sys_mail.go
│   │   │   ├── sys_mail_log.go
│   │   │   ├── sys_mail_util.go
│   │   │   ├── team.go
│   │   │   ├── team_log.go
│   │   │   ├── team_util.go
│   │   │   ├── user.go
│   │   │   ├── user_log.go
│   │   │   ├── user_stats.go
│   │   │   ├── user_stats_log.go
│   │   │   ├── user_util.go
│   │   │   ├── workset.go
│   │   │   ├── workset_log.go
│   │   │   └── workset_util.go
│   │   ├── res
│   │   │   ├── enum.go
│   │   │   └── res.go
│   │   ├── sys_mail.go
│   │   ├── team.go
│   │   ├── user.go
│   │   ├── user_stats.go
│   │   ├── util
│   │   │   └── util.go
│   │   ├── val
│   │   │   ├── assignment.go
│   │   │   ├── assignment_invitation.go
│   │   │   ├── chapter.go
│   │   │   ├── comic.go
│   │   │   ├── sys_mail.go
│   │   │   ├── team.go
│   │   │   ├── user.go
│   │   │   ├── user_stats.go
│   │   │   └── workset.go
│   │   └── workset.go
│   ├── cfg
│   │   ├── cfg.go
│   │   ├── db.go
│   │   ├── env.go
│   │   └── http.go
│   ├── domain
│   │   ├── ext
│   │   │   ├── oss
│   │   │   │   ├── cleaner.go
│   │   │   │   ├── client.go
│   │   │   │   └── signer.go
│   │   │   └── token
│   │   │       └── parser.go
│   │   ├── model
│   │   │   ├── aggr
│   │   │   │   ├── assignment.go
│   │   │   │   ├── assignment_inv.go
│   │   │   │   ├── chapter.go
│   │   │   │   ├── comic.go
│   │   │   │   ├── member.go
│   │   │   │   ├── member_inv.go
│   │   │   │   ├── oss_msg.go
│   │   │   │   ├── role.go
│   │   │   │   ├── role_iface.go
│   │   │   │   ├── sys_mail.go
│   │   │   │   ├── team.go
│   │   │   │   ├── user.go
│   │   │   │   ├── user_stats.go
│   │   │   │   └── workset.go
│   │   │   ├── enum
│   │   │   │   ├── chapter.go
│   │   │   │   ├── comic.go
│   │   │   │   ├── member.go
│   │   │   │   ├── oss_op.go
│   │   │   │   ├── role.go
│   │   │   │   ├── workflow.go
│   │   │   │   └── workset.go
│   │   │   ├── event
│   │   │   │   ├── assignment.go
│   │   │   │   ├── chapter.go
│   │   │   │   ├── typ.go
│   │   │   │   └── user.go
│   │   │   └── query
│   │   │       ├── assignment.go
│   │   │       ├── assignment_inv.go
│   │   │       ├── chapter.go
│   │   │       ├── comic.go
│   │   │       ├── member.go
│   │   │       ├── member_inv.go
│   │   │       ├── pagi.go
│   │   │       ├── sys_mail.go
│   │   │       └── workset.go
│   │   ├── repo
│   │   │   ├── assignment.go
│   │   │   ├── assignment_inv.go
│   │   │   ├── chapter.go
│   │   │   ├── comic.go
│   │   │   ├── err.go
│   │   │   ├── member.go
│   │   │   ├── member_inv.go
│   │   │   ├── oss_msg.go
│   │   │   ├── prov.go
│   │   │   ├── sys_mail.go
│   │   │   ├── team.go
│   │   │   ├── txn_ctrl.go
│   │   │   ├── user.go
│   │   │   ├── user_stats.go
│   │   │   ├── util.go
│   │   │   └── workset.go
│   │   └── svc
│   │       ├── assignment.go
│   │       ├── assignment_inv.go
│   │       ├── chapter.go
│   │       ├── comic.go
│   │       ├── member.go
│   │       ├── oss_msg.go
│   │       ├── res
│   │       │   ├── enum.go
│   │       │   └── res.go
│   │       ├── user.go
│   │       ├── util.go
│   │       └── workset.go
│   ├── event
│   │   ├── base.go
│   │   ├── bus.go
│   │   ├── event.go
│   │   ├── handler.go
│   │   └── typ.go
│   ├── infra
│   │   ├── event
│   │   │   ├── assignment_chapter.go
│   │   │   ├── bus.go
│   │   │   ├── user.go
│   │   │   └── work_unit.go
│   │   ├── ext
│   │   │   ├── oss
│   │   │   │   ├── oss.go
│   │   │   │   ├── platform.go
│   │   │   │   ├── r2_client.go
│   │   │   │   └── util.go
│   │   │   └── token
│   │   │       ├── claims.go
│   │   │       └── parser.go
│   │   └── repo
│   │       ├── assignment.go
│   │       ├── assignment_inv.go
│   │       ├── chapter.go
│   │       ├── comic.go
│   │       ├── entity
│   │       │   ├── assignment.go
│   │       │   ├── assignment_inv.go
│   │       │   ├── chapter.go
│   │       │   ├── comic.go
│   │       │   ├── member.go
│   │       │   ├── member_inv.go
│   │       │   ├── oss_msg.go
│   │       │   ├── sys_mail.go
│   │       │   ├── team.go
│   │       │   ├── user.go
│   │       │   ├── user_stats.go
│   │       │   └── workset.go
│   │       ├── err.go
│   │       ├── err_classifier.go
│   │       ├── member.go
│   │       ├── member_inv.go
│   │       ├── oss_msg.go
│   │       ├── prov.go
│   │       ├── repo.go
│   │       ├── sys_mail.go
│   │       ├── team.go
│   │       ├── txn_ctrl.go
│   │       ├── user.go
│   │       ├── user_stats.go
│   │       └── workset.go
│   └── lgr
│       └── lgr.go
├── justfile
├── main.go
├── migrations
│   ├── 20260425162940_enable-features.down.sql
│   ├── 20260425162940_enable-features.up.sql
│   ├── 20260425162943_create-user-table.down.sql
│   ├── 20260425162943_create-user-table.up.sql
│   ├── 20260425162956_create-team-table.down.sql
│   ├── 20260425162956_create-team-table.up.sql
│   ├── 20260425163018_create-member-invitation-table.down.sql
│   ├── 20260425163018_create-member-invitation-table.up.sql
│   ├── 20260425163027_create-member-table.down.sql
│   ├── 20260425163027_create-member-table.up.sql
│   ├── 20260425163102_create-oss-message-table.down.sql
│   ├── 20260425163102_create-oss-message-table.up.sql
│   ├── 20260426101603_create-system-mail-table.down.sql
│   ├── 20260426101603_create-system-mail-table.up.sql
│   ├── 20260427000001_create-workset-table.down.sql
│   ├── 20260427000001_create-workset-table.up.sql
│   ├── 20260427000002_create-comic-table.down.sql
│   ├── 20260427000002_create-comic-table.up.sql
│   ├── 20260427000003_create-chapter-table.down.sql
│   ├── 20260427000003_create-chapter-table.up.sql
│   ├── 20260428000001_create-assignment-invitation-table.down.sql
│   ├── 20260428000001_create-assignment-invitation-table.up.sql
│   ├── 20260428000002_create-assignment-table.down.sql
│   ├── 20260428000002_create-assignment-table.up.sql
│   ├── 20260428000003_create-user-stats-table.down.sql
│   ├── 20260428000003_create-user-stats-table.up.sql
│   ├── 20260428000004_create-page-table.down.sql
│   ├── 20260428000004_create-page-table.up.sql
│   ├── 20260428000005_create-unit-table.down.sql
│   └── 20260428000005_create-unit-table.up.sql
└── pkg
    └── util
        └── id.go
```

