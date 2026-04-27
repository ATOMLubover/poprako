---
description: "Use when editing refactor domain-layer files under refactor/internal/domain, including aggr models, enum and incl types, domain events, query option structs, repo interfaces, domain services, ext abstractions, and domain-layer style."
name: "Refactor Domain Layer Style"
applyTo:
  - "refactor/internal/domain/**/*.go"
---
# Refactor Domain Layer Style

这份 instruction 只约束 `refactor/internal/domain` 及其全部子目录

它是在共享 Go 宪法之上的 domain 层专项补充

## 核心定位

- domain 层负责表达业务语义 不变式 领域行为 领域事件 和抽象契约
- domain 层是 refactor 架构的语义中心 不是 DTO 仓库 也不是持久化细节仓库
- schema 事实仍以 `refactor/migrations` 为准 但 domain 不直接承载列名 SQL 或 GORM 细节
- app 层负责编排 infra 层负责实现 因此 domain 不应反向依赖 `internal/app` `internal/api` `internal/infra`

## 目录职责

- `model/aggr` 放聚合与与该聚合强相关的领域结构 允许一个文件中包含同一业务概念下的多个 struct 或 type
- `model/enum` 放稳定的领域枚举与 include 类型 不把这类闭集常量散落到 `aggr` `repo` 或 `app`
- `model/event` 放领域事件常量与事件 payload 类型 使用 `event_impl` 作为包名
- `model/query` 放 repo 查询所需的强类型过滤与分页参数 不把查询条件零散铺到 repo 方法签名里
- `repo` 放 `repo_iface` 包下的抽象契约与共享错误类型 只描述领域需要的读写语义
- `svc` 放无状态 domain service 用于生成 校验 归并 或补全领域对象
- `ext/*` 放对外部能力的抽象接口 例如 `oss_iface` `token_iface`

## `aggr` 规则

- `aggr` 是充血模型 允许并鼓励在结构体上定义领域行为 不应退化成只有字段的贫血 struct
- 同一业务概念下的主聚合 创建参数 更新参数 凭证结构 辅助类型 可以共置在同一个文件中 例如 `User` `UserCreds` `UserReg` `UserUpd`
- 共享领域行为优先抽成 `aggr` 内的辅助类型或小接口 例如 `RoleMask` `TimedRoles` `WithRoles` 而不是过早上提到 app 或 infra
- 领域对象允许嵌入 `event_iface.EvBase` 来积累待发布事件 事件的产生应贴近业务动作本身
- `aggr` 中的方法应返回领域语义结果 不应返回 app 层结果类型 也不应了解 repo 或 GORM 实现细节
- `aggr` 可以依赖标准库与必要的纯算法库 但不应依赖 app infra api zap gorm 这类外层实现细节

## `enum` `query` `event` 规则

- 枚举 常量范围 include 选项等闭集语义必须放入 `model/enum` 不使用裸字符串和裸整数在多层漂移
- repo include 这类选择器优先定义成强类型枚举 例如 `MemberIncl`
- 查询过滤 分页 批量读取参数优先定义成 `model/query` 下的 option struct 例如 `ListMemberOpt` `PagiOpt`
- 可选查询条件优先使用指针字段表达 缺失与零值必须可区分
- 领域事件常量集中放在 `model/event/typ.go`
- 具体事件 payload 与构造函数按主题拆文件 例如用户相关事件留在 `model/event/user.go`
- 事件构造函数优先返回 `event_iface.Event` 而不是把外层依赖暴露给调用方

## `repo_iface` 规则

- `repo` 包只定义抽象契约 不描述 SQL GORM preload join upsert 等 infra 细节
- repo 方法签名使用强类型参数与返回值 优先接收 `aggr` `enum` `query` 下的类型
- repo 返回的核心数据应是领域聚合或标量 不返回 entity row 或 DB 句柄
- repo 错误统一走 `RepoErr` 这一抽象边界 由 infra 负责把具体错误映射进来
- include 需求通过 `enum.XxxIncl` 表达 过滤需求通过 `query.XxxOpt` 表达 不在方法签名中堆积大量布尔值与可空标量
- `txn` 相关抽象可以存在于 domain repo 中 但只保留调用语义 例如 `TxnCtrl` `TxnFn` 不泄露底层事务实现方式
- repo 文件按概念一文件一接口组织 共享类型单独放在 `err.go` `txn_ctrl.go`

## `svc` 规则

- domain service 必须保持无内禀状态 不保存 repo db cfg client logger 等字段
- `NewXxxSvc` 不接收 repo db client 等依赖
- service 的依赖应通过方法参数显式传入 或只依赖纯函数与纯工具
- service 负责生成 补全 校验 或组合领域对象 例如生成 `UserReg` `MemberCre` 或保存待处理的 oss message 语义
- service 可以调用不引入 infra 语义的纯工具 例如 `pkg/util.GenId`
- service 不直接持久化对象 不承担 app 层事务编排与响应语义
- 若某个逻辑本质上是单个聚合的自有行为 优先留在 `aggr` 方法中 而不是放进 service

## `ext` 规则

- `ext` 目录只定义外部能力抽象 不包含任何实现逻辑
- 接口命名优先表达能力本身 例如 `Signer` `Cleaner` `Parser` `Client`
- 更大的接口可以通过组合更小的接口形成 例如 `Client` 嵌入 `Signer` 与 `Cleaner`
- `ext` 接口只描述 domain 真正需要的能力表面 不为了迁就某个具体 SDK 暴露多余方法

## 文件与组织规则

- domain 文件命名保持 `snake_case`
- `aggr` 文件按业务概念聚类 而不是按 struct 种类机械拆分
- `enum` `event` `query` 均按主题拆分 不要把无关常量和 option 堆在一个总文件中
- 当多个聚合共享同一组规则时 优先抽出概念化的辅助类型 文件仍放在最贴近语义的目录下 例如角色系统相关类型放在 `aggr/role.go` 与 `enum/role.go`
- 新增领域类型时 优先归位到正确子目录 不要图省事塞进现有无关文件

## 边界约束

- domain 不返回 `val` `res.AppRes` 或任何 app 层结果包装
- domain 不依赖 `context.Context` `*zap.Logger` 作为普遍机制 只有契约本身确实需要时 才在抽象层显式出现 例如 `TxnFn`
- domain 不出现 GORM tag SQL 语句 HTTP 语义 路由 DTO 序列化标签等外层细节
- repo 接口不泄露 infra helper 名称 service 不泄露 app 编排意图 aggr 不泄露数据库结构细节

## 建议补充模式

- 对有明确语义的标量优先使用命名类型或别名 而不是长时间停留在裸 `string` `int` `uint32` 上 例如 `RoleMask` `SignedToken`
- 当多个聚合共享同一套操作协议时 优先定义小接口表达该协议 例如 `WithRoles`
- 创建与更新这类领域意图 若只服务某个聚合 家族对象应贴近该聚合定义 而不是搬到 app 层
- 领域事件应由触发该业务动作的聚合或 service 在最靠近业务语义的位置产生 不要拖到 infra 层临时拼接
- repo 的 list 类方法默认优先接收 query object 与 typed include 而不是不断膨胀方法参数列表
- domain service 若只做对象构造和校验 应保持轻量清晰 不要演化成承接一切逻辑的杂糅层