# DDD 架构文件组织

```bash
➜  poprako-s-adv git:(refactor/simplify) tree re
factor
refactor
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── http
│   │   │   ├── auth.go
│   │   │   ├── http.go
│   │   │   ├── middleware
│   │   │   │   ├── auth.go
│   │   │   │   └── log.go
│   │   │   ├── res
│   │   │   │   └── res.go
│   │   │   ├── user.go
│   │   │   └── util.go
│   │   └── state
│   │       └── state.go
│   ├── app
│   │   ├── impl
│   │   │   ├── user.go
│   │   │   ├── user_util.go
│   │   │   └── util.go
│   │   ├── res
│   │   │   ├── enum.go
│   │   │   └── res.go
│   │   ├── user.go
│   │   ├── user_stats.go
│   │   └── val
│   │       ├── user.go
│   │       └── user_stats.go
│   ├── cfg
│   │   ├── cfg.go
│   │   └── env.go
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
│   │   │   │   ├── member.go
│   │   │   │   ├── member_inv.go
│   │   │   │   ├── oss_msg.go
│   │   │   │   ├── role.go
│   │   │   │   ├── role_iface.go
│   │   │   │   ├── sys_mail.go
│   │   │   │   ├── team.go
│   │   │   │   ├── user.go
│   │   │   │   └── user_stats.go
│   │   │   ├── enum
│   │   │   │   ├── member.go
│   │   │   │   ├── oss_op.go
│   │   │   │   └── role.go
│   │   │   ├── event
│   │   │   │   ├── typ.go
│   │   │   │   └── user.go
│   │   │   └── query
│   │   │       ├── member.go
│   │   │       ├── member_inv.go
│   │   │       └── pagi.go
│   │   ├── repo
│   │   │   ├── err.go
│   │   │   ├── member.go
│   │   │   ├── member_inv.go
│   │   │   ├── oss_msg.go
│   │   │   ├── sys_mail.go
│   │   │   ├── team.go
│   │   │   ├── txn_ctrl.go
│   │   │   ├── user.go
│   │   │   └── user_stats.go
│   │   └── svc
│   │       ├── member.go
│   │       ├── oss_msg.go
│   │       └── user.go
│   ├── event
│   │   ├── base.go
│   │   ├── bus.go
│   │   ├── event.go
│   │   ├── handler.go
│   │   └── typ.go
│   └── infra
│       ├── event
│       │   └── user.go
│       ├── ext
│       │   └── token
│       │       ├── claims.go
│       │       └── parser.go
│       └── repo
│           ├── entity
│           │   ├── member.go
│           │   ├── team.go
│           │   └── user.go
│           ├── err.go
│           ├── member.go
│           ├── member_inv.go
│           ├── oss_msg.go
│           ├── txn_ctrl.go
│           └── user.go
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
│   └── 20260425163102_create-oss-message-table.up.sql
└── pkg
    └── util
        └── id.go
```