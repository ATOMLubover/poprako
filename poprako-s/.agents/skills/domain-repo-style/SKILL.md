---
name: domain-repo-style
description: 为本仓库 DDD（domain/repo）接口总结命名、注释、缩写与审查要点
---

# domain-repo-style

## 目的

本技能（Skill）总结并规范本仓库 `internal/domain/repo` 下的接口风格，包含：命名约定、注释风格、缩写优先级、文件与包约定、接口语义建议以及 PR 审查清单 目的是统一风格、降低沟通成本，并明确哪些信息应保留于 domain 层，哪些应由 infra/实现层负责

## 适用范围

- 范围：`internal/domain/repo` 下的所有接口类型（`*Repo`, `*Mgr` 等）以及与之直接关联的 model 定义
- 场景：代码审查、接口定义、生成样板注释、自动化校验规则（linters / CI 检查提示）

## 总则（快速概览）

- 导出标识符使用 PascalCase（首字母大写）
- 接口命名采用角色后缀：仓库用 `Repo`，管理器用 `Mgr`（保持现有风格一致）
- 常见缩写使用大写惯例（如 `ID`, `QQ`, `OSS`）；项目特定缩写保留且在团队内统一（如 `Txn` 在当前代码中被接受）
- 注释使用中文，公开符号必须有文档注释，注释首词应以符号名开头（符合 Go 文档注释习惯）
- domain 层接口只表达领域意图与不变式，尽量避免包含具体的持久化实现细节（例如 SQL 语句、`UPDATE ON CONFLICT` 之类的实现提示应放到 infra/实现层或单独注释为“实现说明”）

## 详细约定

### 命名（Naming）

- 接口命名
  - 形如 `UserRepo`, `UserStatsRepo`, `TxnMgr` 保留 `Repo` 与 `Mgr` 后缀，避免同时出现 `Manager` 与 `Mgr` 两种写法
  - 接口名首字母大写导出（domain 层通常导出给 infra 实现）
- 方法命名
  - 使用动词或动词短语：`Get`, `GetByID`, `GetOrCreate`, `Create`, `Reg`（现有项目使用 `Reg` 表示注册），`Update`, `Patch`, `Delete`, `Confirm`
  - 对于标识符（如 ID）使用全大写缩写：`GetByID`（不是 `GetById`）
  - 对于第三方或产品专有缩写（如 `QQ`）保持原样：`GetCredsByQQ`
- 参数与变量
  - 参数使用驼峰小写（camelCase）：`userID string`、`avatarOSSKey string` 当上下文明确时短名 `id` 可接受，但在可能混淆时使用 `userID`/`orderID` 等更具体的名字
- 文件命名
  - 使用小写并可含下划线（项目中为 `user.go`, `user_stats.go`, `tx_mgr.go`），保持现有风格一致即可

### 缩写偏好（Acronyms & Abbreviations）

- 标准缩写（按 Go 社区惯例）：`ID`, `URL`, `HTTP`, `UUID` 等全部大写
- 项目中出现且需保留的缩写：
  - `QQ` → `QQ`
  - `OSS` → `OSS`
  - `Txn` → `Txn`（项目当前使用 `TxnMgr`，建议将 `Txn` 视为项目约定；若愿意统一为 `Tx` 可在团队内决定并替换）
  - `Reg` → `Reg`（当前用法表示 Register；可选：改为 `Register` 以提高可读性，但会引入命名变更）
- 规则：在复合标识符中保留缩写的大写：`GetByID`, `PreFillAvatarOSSKey`

### 注释风格（Docs / Comments）

- 语言：中文；对外接口需中文文档注释以便团队阅读
- 格式：导出类型与方法的注释应以符号名开头（Go doc 推荐样式），例如：

```go
// UserRepo 是用户仓库的接口
type UserRepo interface {
    // GetByID 根据用户 ID 返回用户信息；若不存在返回 (nil, nil) 或 (nil, ErrNotFound)
    GetByID(id string) (*model.UserInfo, error)
}
```

- 内容要点：
  - 描述功能语义（做什么）、重要语义（幂等性、事务边界、并发/锁的预期行为）和错误语义（什么时候会返回 error）
  - 避免在 domain 注释中包含低层持久化实现细节（如 SQL 语句、特定 DB 行为） 如果有必要，可用“实现说明（infra）：...”单独标注并放在 infra/实现处
  - 注释语句要完整、标点规范，避免半个括号、未闭合的注释

### 接口语义建议（Domain semantics）

- `GetOrCreate`：应明确返回值语义（是否保证创建后返回非 nil、并说明并发冲突的期望处理方式）
- `Patch` vs `Update`：
  - `Patch` 用于部分更新（merge/patch），`Update` 用于全量替换或具有明确的替换语义 PR 中应标注 `Patch` 的可空字段语义
- 事务：domain 层可以在注释中说明“此操作需在事务中执行”或“调用方应保证事务”，但不要写出具体事务实现方式

### 返回值约定

- 返回 `(*T, error)` 为首选，error 始终为最后一个返回值
- 若函数可能返回“未找到”，在注释中明确是否返回 `(nil, nil)` 或 `ErrNotFound`

## 示例（推荐写法）

```go
// UserStatsRepo 负责用户统计信息的读取与更新
// 语义：实现需保证对同一 userID 的并发写入具有合理的并发控制（或在 infra 层保证幂等）
type UserStatsRepo interface {
    // GetOrCreate 返回用户统计信息；若不存在则创建并返回（调用方无需关心创建细节）
    // 若创建/读取失败返回非空 error
    GetOrCreate(userID string) (*model.UserStats, error)
    // Patch 对用户统计信息执行部分更新，仅修改非零/非空字段
    Patch(stats *model.UserStats) error
}
```

## PR 审查清单（Reviewer Quick-Checklist）

- [ ] 导出类型/方法有中文注释，注释首词以符号名开头
- [ ] 注释只描述语义，不泄露实现细节（SQL/DB 特定语句移至 infra）
- [ ] 接口/方法命名遵循 `Repo`/`Mgr` 后缀约定
- [ ] 缩写使用一致（`ID`, `QQ`, `OSS`, `Txn` 等）
- [ ] 方法返回值顺序正确：数据, error

## 已确认的团队策略（根据项目维护者回复）

1. `Reg` 与 `Txn` 为项目内常见且接受的缩写：保留 `Reg` 表示注册（`Reg`）以及 `Txn` 用于事务管理（`TxnMgr`） 一般原则：若存在常见缩写（例如 `ID`, `QQ`, `OSS`, `Reg`, `Txn`）则使用缩写；否则使用全称（例如 `Assignment` 使用 `Assignment`）
2. 允许在 domain 层注释中出现实现提示（implementation hint），以便说明领域意图或对 infra 实现者提出建议，例如 `可能需要 UPDATE ON CONFLICT` 实现提示应以明确前缀标注为“实现提示（infra）”或放到 infra 目录的实现注释中，以避免与领域语义混淆
3. 强制文件命名使用下划线风格（snake_case），例如 `user_stats.go`；不允许使用驼峰风格文件名以避免导入不便

以上规则已被记录进本技能，并将作为后续自动检查与 PR 审查的参考标准

## 使用示例（可以对话中直接调用的提示）

- "请根据本技能检查 `internal/domain/repo`，列出所有不符合命名或注释规范的接口 "
- "把 `internal/domain/repo` 下的接口注释统一改为本技能推荐的格式并生成 PR 描述草稿 "

## 保存位置

技能文件保存在：`.agents/skills/domain-repo-style/SKILL.md`
