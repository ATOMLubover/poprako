---
name: app-layer-style
description: |
  记录 `app` 层代码规范与实现约定，适用于本仓库中所有 `app` 层实现。
  使用时机：新增或审查 `app` 层代码、编写 `app` 层实现模板、代码审查时校验一致性。
---

# App 层代码规范（仓库级 Skill）

## 目的

- 统一 `app` 层与接口层的对接契约与实现风格，减少各模块之间的差异。
- 保证 `app` 层入参与返回值均使用 `val` 包中定义、带序列化标签的结构体。
- 明确 `app` 层错误语义，使 API 层能直接向客户端返回合理的错误信息（Bad Request）。

## 何时使用

- 编写或审查 `internal/app/*` 下的实现。
- 创建新的 `app` 子模块（如 `user`, `team`, `comic` 等）。
- 在编写中间层包装（如日志包装器）时参考实现样式。

## 约定

- 所有公开的 `app` 接口方法：
  - 第一个参数必须是 `context.Context`，命名为 `cx`。
  - 所有输入/输出结构体必须来自 `poprako-s/internal/app/val` 包。
  - 只要是 `val` 包中的结构体，签名中一律使用指针类型（`*val.XxxArgs` / `*val.XxxRes` / `*val.XxxInfo`），避免不必要的值拷贝。
  - 函数签名参数采用多行格式（与现有 `Login` 风格一致），例如：

    func (a *SomeApp) DoSomething(
    cx context.Context,
    args *val.DoSomethingArgs,
    ) (\*val.DoSomethingRes, error)

- 每个 `app` 模块应提供两种实现：
  1. 真实业务实现（例如 `userAppImpl`），负责处理核心逻辑；
  2. 日志或包装实现（例如 `logUserAppImpl`），负责初始化/注入日志并将调用转发给真实实现。

- 构造函数命名：
  - `New<Name>App()` 返回 `UserApp` 风格的接口实现，或 `NewLog<Name>App()` 返回日志包装实现。

- 日志注入与上下文处理：
  - 日志记录器 `*zap.Logger` 可通过包装实现注入并使用 `injectLgr(cx, lgr)` 将带上下文的 logger 放入 `context.Context` 中以便下游使用。

- 参数验证：
  - `app` 层负责对输入参数进行边界/格式校验（例如长度、必填），并返回可直接展示给客户端的错误字符串。
  - 校验失败应返回 `error`，其 `Error()` 字符串会被 API 层映射为 Bad Request 的 `message` 字段。
  - 错误信息不得包含敏感信息（如密码、密钥、完整个人数据）。

- 注释与排版：
  - `app` 层实现中的每一个逻辑步骤都应有简短注释，说明该步骤的业务意图，而不是重复代码字面含义。
  - 任意两个可执行语句之间必须保留空行；不要把日志、赋值、返回等语句紧贴写在一起。
  - 特别注意 `zap.L().Error(...)`、变量赋值、`return` 这类连续语句之间也必须留空行。
  - 完成后应使用 `rg` 或等效搜索工具复查是否仍存在“语句黏连”的位置。

- 错误语义：
  - `app` 层的错误表示客户端错误（Bad Request），不要在错误中包含内部堆栈或敏感上下文。
  - 对于真正的服务器/系统错误（数据库连接、第三方服务不可用等），`app` 层可以返回包装后的错误供上层区分（建议使用 sentinel 或 errors.Is 进行判断），但默认交付给客户端的 message 应为不暴露内部细节的友好提示。

## 示例模板

- 接口定义（放在 `internal/app/<module>.go`）

  type SomeApp interface {
  DoSomething(
  cx context.Context,
  args *val.DoSomethingArgs,
  ) (*val.DoSomethingRes, error)
  }

- 真实实现骨架：

  type someAppImpl struct {
  // 注入仓库、外部依赖等
  }

  func NewSomeApp(dep SomeDep) SomeApp {
  // 校验依赖。
  if dep == nil {
  panic("dep is nil")
  }

  // 返回真实实现。
  return &someAppImpl{}
  }

- 日志包装实现骨架：

  type logSomeAppImpl struct {
  lgr \*zap.Logger
  app SomeApp
  }

  func NewLogSomeApp(lgr \*zap.Logger, app SomeApp) SomeApp {
  // 校验依赖。
  if lgr == nil || app == nil {
  panic("invalid dependency")
  }

  // 返回日志包装实现。
  return &logSomeAppImpl{lgr: lgr, app: app}
  }

  func (a *logSomeAppImpl) DoSomething(
  cx context.Context,
  args *val.DoSomethingArgs,
  ) (\*val.DoSomethingRes, error) {
  // 校验包装器是否可用。
  if a == nil || a.app == nil {
  return nil, errors.New("invalid app")
  }

  // 注入日志上下文。
  cx = injectLgr(cx, a.lgr.With(zap.Any("args", args)))

  // 转发调用。
  return a.app.DoSomething(cx, args)

  }

## 校验清单（PR 审查）

- 方法签名是否以 `context.Context` 开头且使用 `cx` 命名。
- 所有输入输出类型来自 `val` 包并带有序列化标签，且是否统一使用指针。
- 函数参数是否为多行格式（每个参数独立一行）。
- 是否存在对应的日志包装实现（如果需要日志注入）。
- 错误返回是否遵循“不包含敏感信息且可作为客户端 message”的规范。
- 每个逻辑步骤前是否存在说明性注释。
- 任意两个可执行语句之间是否保留空行。

## Example prompts

- "请基于这个 SKILL 校验 `internal/app/user.go` 中的方法签名是否符合约定，并指出不符合项。"
- "为新模块 `subscription` 生成 `SomeApp` 接口与 `log` 包装实现的骨架代码。"
- "请用 `rg` 检查 `internal/app/user.go` 中是否还有任意两个语句之间没有空行的位置。"
- "请按这个 SKILL 为 `team app` 生成带逐步注释和指针签名的完整实现骨架。"

## 备注

- 该 SKILL 为仓库级别约定，建议放在 `.agents/skills/app-layer-style/SKILL.md` 或 `.github/skills/app-layer-style/SKILL.md`，以便团队成员在审查或实现时统一参考。
