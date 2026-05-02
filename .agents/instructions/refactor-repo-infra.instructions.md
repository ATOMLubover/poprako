---
description: "Use when editing refactor repo infra files under refactor/internal/infra/repo, including GORM entity structs, preload incl patterns, transaction helpers, query optimization, 仓储实现, and repo infra style."
name: "Refactor Repo Infra Style"
applyTo:
  - "refactor/internal/infra/repo/*.go"
  - "refactor/internal/infra/repo/**/*.go"
---
# Refactor Repo Infra Style

这份 instruction 只约束 `refactor/internal/infra/repo` 及其 `entity` 子目录

它是在共享 Go 宪法之上的 repo infra 专项补充

## 核心原则

- `repo_infra` 只负责持久化实现与查询装配 不承载 app 层编排 也不承载 domain 业务规则
- 一切 schema 认知以 `refactor/migrations` 为准 老项目只能提供行为参考 不能覆盖表结构事实
- GORM 交互必须建立在强类型 `entity` struct 之上 不允许使用 `map[string]any` `[]map[string]any` 或其他弱类型载体绕过建模
- 优先为单个操作定义专用 row struct 而不是复用一个臃肿的全字段 struct
- 每个查询与写入都要尽量只触达最小必要字段数量 通过专用 row struct 与显式列选择控制带宽和扫描范围

## `entity` 建模规则

- `entity` 是 GORM 的唯一直接交互层 repo 方法读取和写入的数据形状都应先落到 `entity` struct
- 每个持久化概念维持独立文件 不把多个无关聚合揉进同一文件
- 同一文件内按 struct 分区组织 先定义 struct 再紧跟它的 `TableName` 转换函数 构造函数 与相关 helper
- 禁止把新的相关函数一股脑追加到文件最后 必须跟在所属 struct 或所属小主题之后
- 读取模型 写入模型 更新模型 统计模型 预加载模型可以拆成不同 struct 例如 `UserRow` `UserCredsRow` `UserRegRow`
- `ToXxxAggr` 和 `NewXxxRowFromAggr` 这类转换逻辑必须留在 `entity` 包内 不外泄到 repo 文件
- `TableName` 与表常量要显式声明并与 row struct 共置 让表归属一眼可见
- 需要表示未 preload 的关联时 关联字段优先使用指针类型 让 未加载 与 已加载但为空 的语义保持可区分
- `ToXxxAggr` 需要对可选关联保持 nil safe 未 preload 不应被当作错误

## 查询与写入规则

- 每个 repo 方法都从稳定的 base query 开始 通常是 `r.gdb.Table(entity.XxxTable)` 再按固定顺序叠加条件 include 与执行动作
- 默认显式指定表来源 不依赖 GORM 从 model 名推断表名
- 查询返回什么 就定义什么 row struct 不为了省事扫描多余字段
- `Exist` `Count` 之类方法只查询标量或最小结果 不回表装配完整聚合
- 新建 更新 删除 都要使用强类型 row struct 或强类型 helper 不使用 `Save` 这类写入范围不透明的 API
- 如果更新需要保留零值语义 仍然禁止回退到 `map[string]any` 应通过专用 update row struct 搭配显式 `Select` 或精确 `UpdateColumn` `UpdateColumns` 完成
- 一个操作需要统一时间戳时 在构造 row 的地方一次性生成 `now` 再分发给相关字段 不在多个位置各自取时间
- repo 方法对外返回 domain aggregate 标量结果或 `repo_iface` 约定的错误类型 不向包外暴露 GORM row 或 `gorm.DB`

## `incl` 与 preload 规则

- belongs to 或其他常用 include 场景必须通过强类型 `Incl` 枚举表达 例如 `enum.MemberIncl`
- repo 方法接收 `inc ...enum.XxxIncl` 这类参数时 只能在 repo 内部把它映射为明确的 preload 行为 不能把裸字符串 preload 名泄漏到调用方
- preload 开关的语义要稳定 可预测 同一个 `Incl` 取值始终映射同一组关联
- 如果多个 repo 共享同一类 preload 或查询片段 优先抽成 repo 内可复用的强类型 helper 而不是复制粘贴 switch
- preload 需要自定义列选择 排序 或子条件时 也应收敛为 typed helper 不在业务方法中铺开重复细节

## 强制个人约定

- 若聚合有 belongs to 关联字段（例如 `Workset.Team`）repo 必须实现 typed include 到 preload 的映射 不能声明字段但不支持 include
- `PUT` 语义更新时 可空字段即使为 `nil` 也要被显式 `Select` 并写入 `NULL` 禁止把 `nil` 当作“跳过列更新”
- `List` 查询的排序必须严格匹配业务默认排序 并在 `refactor/migrations` 提供与过滤列 + 排序列一致的复合索引
- 删除命名固定：`Delete` 只能实现硬删除 `Remove` 只能实现软删除
- 已有常见缩写必须统一使用 例如 `Desc`/`desc` 禁止回退到 `Description`/`description`（SQL 字符串除外）

## util 与复用规则

- repo infra 已有的 util 与 helper 应优先复用 例如事务提取 重复错误判断 公共查询片段 与 include 装配
- 只有在现有 helper 无法表达意图时 才新增 helper 新 helper 必须职责单一 命名具体 输入输出强类型
- 一段查询拼装逻辑在两个及以上 repo 中重复出现时 应考虑上提为 repo infra helper
- helper 负责消除样板 不能把查询真实语义隐藏成难以追踪的黑盒

## 文件内排布规则

- repo 文件按 repo struct 分区组织 一般顺序是 repo struct `TxnXxxRepo` `NewXxxRepo` 查询方法 写入方法 本 repo 专属 helper
- 如果某个 helper 只服务于一个 struct 或一个方法族 要尽量贴近该区域放置 不要沉到底部堆积
- `entity` 文件也遵循同样原则 以 struct 为中心组织 而不是以后追加顺序组织
- 新增代码时 优先维护已有分区的可读性 宁可微调局部顺序 也不要破坏文件分层

## 命名与注释规则

- 沿用项目里的个人缩写习惯 例如 `gdb` `cx` `txn` 不要为了迎合 Go 常见写法改成 `db` `ctx`
- 每个 struct 每个字段 每个函数 每个方法 每个常量都要求完整注释
- 注释里的标识符必须使用反引号包裹 例如 `UserRow` `GetById` `gdb`
- 这层的注释要求是完整覆盖 不能因为函数很短或字段含义看起来明显就省略

## 建议补充模式

- 当一个 repo 方法同时承担过滤 include 分页 排序时 优先先抽出稳定的小型 query builder helper 再组装主流程
- 对于 create update 场景 优先使用操作专属 row struct 来锁定允许写入的列集合 防止未来字段扩张后误写
- 对于 preload 关联 如果只需要极少列 应为被 preload 的关联也准备更轻量的 row 形状 或在 preload 回调里显式收窄列集合
- repo infra 中的转换函数只做数据映射 不在 `ToXxxAggr` 内做额外查询 远程调用 或业务决策
- 对重复出现的错误语义 提供类似 `IsDupKey` 的 helper 让上层判断建立在明确语义上 而不是散落的底层错误比较