# build-application-closed-loop

Purpose

- 指导如何为一个资源（以 `user` 为例）构建完整的闭环：从数据库迁移脚本、domain model、repository、application/service、value（DTO）、到 HTTP API（handler/路由）及测试。

Scope

- 面向本仓库的 Go 服务；示例以 `user` 资源为中心。适用于 `internal/value/*`、`domain/model/*`、`internal/application/*`、`infrastructure/repository/*`、`internal/api/http/*`、`migrations/*`。

High-level steps (closed-loop)

1. 数据库层（Migration）
   - 编写 SQL migration 脚本（up/down）。遵循命名规范：`YYYYMMDDHHMMSS_resource-table.up.sql`。
   - 确保 SQL 是单向幂等的，并包含必要的索引与约束。
   - 示例： [migrations/20260306101212_user-table.up.sql](migrations/20260306101212_user-table.up.sql)

2. Domain 层（Model）
   - 在 `domain/model` 中定义充血 model，封装核心不变量与领域方法。
   - 提供 `NewXXX(...)` 构造函数（若构造可能失败则返回 error）。
   - 文件示例： [domain/model/user.go](domain/model/user.go)

3. Repository 层
   - 定义 repository 接口（domain/repository）以描述需要的持久化操作。
   - 在 `infrastructure/repository` 中实现接口（与 gorm/DB 交互）。实现要支持事务上下文传递。
   - 保持接口与实现分离以便测试替换。

4. Application / Service 层
   - 在 `internal/application` 中实现面向用例的服务函数，协调多个 repository 的调用与事务（在此层统一管理事务）。
   - 服务函数使用 domain model、返回 domain errors，且保持纯粹的业务流程逻辑。
   - 文件示例： [internal/application/user.go](internal/application/user.go)

5. Value / DTO 层
   - 在 `internal/value` 中定义 API 层的 DTO（请求/响应/Args/Result）。
   - 使用 `NewXxxFromModel` 的转换函数将 `model` 转换为 `value`。遵循构造器策略：model 使用 `New` 构造，value 的转换函数负责映射。
   - 文件示例： [internal/value/workset.go](internal/value/workset.go)（参考模式），以及 [internal/value/user.go](internal/value/user.go)

6. API 层（HTTP handler）
   - 在 `internal/api/http` 中实现 handler：解析请求到 `Args`，调用 `application` 服务，处理错误并返回 `value` 结果。
   - Handler 应只负责 HTTP 细节（状态码、序列化、鉴权中间件），业务逻辑委托给 application 层。
   - 文件示例： [internal/api/http/user.go](internal/api/http/user.go)

7. 测试与验证
   - 单元测试：
     - domain model 的不变量与方法测试。
     - repository 的实现（使用 SQLite 或 mock DB）测试 CRUD。
     - application 服务的集成测试（在事务回滚的测试上下文中运行）。
     - handler 的端到端测试（使用 httptest）。
   - 数据迁移脚本应包含在 CI 的迁移验证步骤中（可对临时 DB 应用 up/down）。

Design principles and rules

- 单一构造器：模型与关键对象通过 `New...` 构造，避免外部直接 struct literal。
- Args 为例外：`Args` 类型可直接从 JSON 反序列化，但必须实现 `Validate() error`。
- PUT 更新语义：在 update Args 中，`nil` 表示置 NULL，非 nil 表示更新为该值，空字符串视作有效值。
- 最小指针策略：除非字段为可选或需要变更语义，否则采用值类型。
- 转换函数放在 `value` 包侧：如 `NewUserFromModel(m model.User) value.User`，以减少循环依赖。
- 事务在 Application 层统一管理：Application 层开始/提交/回滚事务，并把事务上下文传给 repository。

Concrete `user` example mapping (quick reference)

- Migration: [migrations/20260306101212_user-table.up.sql](migrations/20260306101212_user-table.up.sql)
- Domain model: [domain/model/user.go](domain/model/user.go)
- Application service: [internal/application/user.go](internal/application/user.go)
- Value/DTO: [internal/value/user.go](internal/value/user.go)
- HTTP handler: [internal/api/http/user.go](internal/api/http/user.go)

Checklist when implementing a new resource

- [ ] 设计并添加 migrations up/down
- [ ] 在 `domain/model` 创建 model + `New` 构造器
- [ ] 在 `domain/repository` 定义接口
- [ ] 在 `infrastructure/repository` 实现接口（支持 tx 上下文）
- [ ] 在 `internal/application` 添加服务，包含事务边界
- [ ] 在 `internal/value` 添加 DTO + `NewXxxFromModel` 转换
- [ ] 在 `internal/api/http` 添加 handler、路由与 godoc
- [ ] 添加单元与集成测试，CI 中运行 migration 验证

Common pitfalls and mitigations

- 循环依赖：通过将转换函数放在 `value` 包并传递最小必要的字段来避免。
- 事务泄漏：在 Application 层显式处理事务并确保 defer 回滚在异常时触发。
- 不一致的错误语义：在 application 层统一转换 domain error 到明确的 HTTP status

Suggested prompts to run this skill

- "Scaffold closed-loop for resource `project` using this skill."
- "Review `internal/application/user.go` and point out missing transaction or repository calls."

Outcome

- 提供一套可重复的、可测试的工作流程来实现资源的全链路闭环（SQL→model→repo→service→value→API）。

Next actions

- 我可以把仓库中以 `user` 为例的各层文件做自动扫描并列出哪些转换/构造/测试缺失。要我继续吗？
