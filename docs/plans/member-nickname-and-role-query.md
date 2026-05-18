# 成员表增加 user_nickname 冗余字段 & 按角色 + 昵称模糊查询

## 1. 目标

前端 member list 需要支持三种查询场景：

| 场景 | 条件                        |
|------|-----------------------------|
| 1    | 仅按昵称模糊搜索              |
| 2    | 仅按**单一**职位筛选          |
| 3    | 按**单一**职位 + 昵称模糊搜索  |

目前 `t_member` 的角色使用 8 列可空时间戳存储，无 `user_nickname` 列，不能直接在 member 表上完成模糊昵称查询。

---

## 2. Schema 变更（migrations/）

### 2.1 删除 user 表无用 trgm 索引

全代码库无任何对 `t_user.nickname` 的 `ILIKE` / `similarity` 模糊查询。`nickname` 的 `UNIQUE` 约束已隐式创建唯一 B-tree 索引，足以满足去重检查。该 GIN trgm 索引从未被使用。

```sql
-- down
CREATE INDEX IF NOT EXISTS "trgm_idx_user_nickname"
    ON "t_user" USING gin ("nickname" gin_trgm_ops);

-- up
DROP INDEX IF EXISTS "trgm_idx_user_nickname";
```

### 2.2 给 member 表增加 user_nickname 列

```sql
-- up
ALTER TABLE "t_member"
    ADD COLUMN "user_nickname" TEXT;

-- 回填已有数据（以当前 user 表 nickname 为准）
UPDATE "t_member" AS m
    SET "user_nickname" = u."nickname"
FROM "t_user" AS u
WHERE m."user_id" = u."id"
  AND m."user_nickname" IS NULL;

-- 后续 migration 将 user_nickname 设为 NOT NULL（等数据安全后再做）
ALTER TABLE "t_member"
    ALTER COLUMN "user_nickname" SET NOT NULL;

-- down
ALTER TABLE "t_member"
    DROP COLUMN IF EXISTS "user_nickname";
```

### 2.3 精简角色部分索引（8 个）

现有索引形如：

```sql
CREATE INDEX "idx_member_team_translator"
    ON "t_member" ("team_id", "assigned_translator_at")
    WHERE "assigned_translator_at" IS NOT NULL;
```

其中第二列 `assigned_xxx_at` 是冗余的：所有查询都会带 `team_id = ?`，而 `WHERE assigned_xxx_at IS NOT NULL` 已在索引过滤条件中。该列从不用于排序，无实际加速价值。

改为仅保留 `team_id` 一列：

```sql
-- 删除旧索引
DROP INDEX IF EXISTS "idx_member_team_raw_provider";
DROP INDEX IF EXISTS "idx_member_team_translator";
DROP INDEX IF EXISTS "idx_member_team_proofreader";
DROP INDEX IF EXISTS "idx_member_team_typesetter";
DROP INDEX IF EXISTS "idx_member_team_redrawer";
DROP INDEX IF EXISTS "idx_member_team_reviewer";
DROP INDEX IF EXISTS "idx_member_team_publisher";
DROP INDEX IF EXISTS "idx_member_team_admin";

-- 重建为单列 + WHERE 部分索引
CREATE INDEX IF NOT EXISTS "idx_member_team_raw_provider"
    ON "t_member" ("team_id")
    WHERE "assigned_raw_provider_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_translator"
    ON "t_member" ("team_id")
    WHERE "assigned_translator_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_proofreader"
    ON "t_member" ("team_id")
    WHERE "assigned_proofreader_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_typesetter"
    ON "t_member" ("team_id")
    WHERE "assigned_typesetter_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_redrawer"
    ON "t_member" ("team_id")
    WHERE "assigned_redrawer_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_reviewer"
    ON "t_member" ("team_id")
    WHERE "assigned_reviewer_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_publisher"
    ON "t_member" ("team_id")
    WHERE "assigned_publisher_at" IS NOT NULL;
CREATE INDEX IF NOT EXISTS "idx_member_team_admin"
    ON "t_member" ("team_id")
    WHERE "assigned_admin_at" IS NOT NULL;
```

### 2.4 新建 user_nickname 的 GIN trgm 索引

仅此一个 GIN 索引，配合上述 8 个 B-tree 部分索引，PostgreSQL 通过 BitmapAnd 自动合并处理三种场景：

```sql
-- up
CREATE INDEX IF NOT EXISTS "trgm_idx_member_user_nickname"
    ON "t_member" USING gin ("user_nickname" gin_trgm_ops);

-- down
DROP INDEX IF EXISTS "trgm_idx_member_user_nickname";
```

查询计划示意：

```
场景 1: user_nickname ILIKE '%keyword%'
  → 仅 trgm_idx_member_user_nickname

场景 2: team_id = ? AND assigned_translator_at IS NOT NULL
  → 仅 idx_member_team_translator

场景 3: team_id = ? AND assigned_translator_at IS NOT NULL
         AND user_nickname ILIKE '%keyword%'
  → BitmapAnd(idx_member_team_translator, trgm_idx_member_user_nickname)
```

### 2.5 同步修改原 migration 文件中的默认 member 插入

`20260425163027_create-member-table.up.sql` 中的默认 super-admin member 插入需增加 `user_nickname` 值，或在回填 migration 中覆盖。

---

## 3. Go 代码变更

### 3.1 三层数据变更概览

```
domain/model/aggr/member.go   → Member 增加 UserNickname；MemberCre 增加 UserNickname
domain/model/query/member.go  → ListMemberOpt 增加 UserNicknameKeyword *string
domain/repo/member.go         → 新增 UpdateUserNickname(userId, name) RepoErr

infra/repo/entity/member.go   → MemberRow / MemberCreRow 增加 UserNickname 列映射
infra/repo/member.go          → List 增加 user_nickname ILIKE 条件
                              → 实现 UpdateUserNickname
                              → Create 填充 user_nickname

domain/svc/member.go          → NewMemberCre 增加 userNickname 参数
domain/svc/user.go            → NewUserReg 显式返回或携带 nickname（供后续 member 创建使用）

app/val/member.go             → MemberVal 增加 UserNickname
                              → ListMemberByTeamArgs 增加 UserNicknameKeyword
                              → CreateMemberArgs 不需要改（user_nickname 从 user 表获取）

app/impl/member_util.go       → asmMemberVal 映射新字段
app/impl/member.go            → ListByTeam 等将 args.UserNicknameKeyword 传入 ListMemberOpt
app/impl/user.go              → Update 包裹事务，同步更新 member 表 user_nickname
app/impl/user.go              → Register 中创建 member 时传入 user nickName

api/http/member.go            → parseMemberQuery 增加 user_nickname_keyword 参数解析
```

### 3.2 关键一致性问题：事务覆盖

用户修改昵称时（`userAppImpl.Update`），**当前无事务包裹**，需改为：

```go
func (a *userAppImpl) Update(cx context.Context, args *val.UserUpdArgs) app_res.AppRes[app_res.None] {
    // 验证参数...

    re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
        userRepo := prov.UserRepo()
        memberRepo := prov.MemberRepo()

        if err := userRepo.Update(&aggr.UserUpd{Id: args.Id, Qid: args.Qid, Name: args.Name}); err != nil {
            if repo_infra.IsNotFound(err) {
                return app_res.Reject[app_res.None](app_res.BadRequest, "用户不存在"), app_res.DefErr()
            }
            if repo_infra.IsDupKey(err) {
                return app_res.Reject[app_res.None](app_res.Conflict, "qid 或昵称已被使用"), app_res.DefErr()
            }
            return app_res.Reject[app_res.None](app_res.ServerError, "更新用户信息失败"), err
        }

        if err := memberRepo.UpdateUserNickname(args.Id, args.Name); err != nil {
            return app_res.Reject[app_res.None](app_res.ServerError, "更新用户信息失败"), err
        }

        return app_res.Accept(&app_res.None{}), nil
    })
    // ...
}
```

用户注册时（`userAppImpl.Register`）已在事务内创建 member，直接传入 `userReg.Nickname` 即可。

### 3.3 Repo 层 ILIKE 查询实现

```go
// member.go — List 方法中
if opt.UserNicknameKeyword != nil && *opt.UserNicknameKeyword != "" {
    keyword := strings.TrimSpace(*opt.UserNicknameKeyword)
    query = query.Where("user_nickname ILIKE ?", "%"+keyword+"%")
}
```

### 3.4 Role 列名映射

entity 层已有 `column` tag 常量，role → 列名映射可放在 `entity/member.go`：

```go
var MemberRoleToColumn = map[enum.Role]string{
    enum.RoleRawProvider: "assigned_raw_provider_at",
    enum.RoleTranslator:  "assigned_translator_at",
    // ...
}
```

当 `ListMemberOpt` 增加 `Role *enum.Role` 字段时，repo 层据此拼接 `WHERE <col> IS NOT NULL`。

---

## 4. 实施顺序

| 步骤 | 内容                                      | 影响范围        |
|------|-------------------------------------------|-----------------|
| 1    | 新建 migration：drop user trgm index       | 无业务影响       |
| 2    | 新建 migration：add user_nickname + 回填    | 向后兼容（NULL） |
| 3    | 新建 migration：重建 8 个角色部分索引         | 删除+重建       |
| 4    | 新建 migration：创建 trgm_idx_member_user_nickname | 新增索引  |
| 5    | domain/aggr + domain/query + domain/repo 接口变更 | 编译断点        |
| 6    | entity + infra repo 实现                    | 编译断点        |
| 7    | domain/svc 适配                             | 编译断点        |
| 8    | app/val + app/impl 适配                     | 编译断点        |
| 9    | api/http 参数解析                           | 编译断点        |
| 10   | 新建 migration：user_nickname SET NOT NULL  | 确认无 NULL 后执行 |
| 11   | 修改原 migration 默认 member 插入语句        | 仅影响新部署环境  |

---

## 5. 不做的

- **不修改** `t_member` 现有 8 个时间戳列（数据模型不变）
- **不修改** `RoleMask` / `TimedRoles` 领域设计
- **不动** `t_user.nickname` 的 UNIQUE 约束或 B-tree 唯一索引
- 本次不实现 `ListMemberOpt.Role` 字段（即按角色筛选的 query option），仅在 DB 层完成索引准备。如需可后续再加，改动量很小
