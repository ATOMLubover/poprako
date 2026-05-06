---
description: "Use when editing refactor API HTTP files under refactor/internal/api/http, including Iris handlers, swagger godoc, auth middleware, HttpRes response wrapping, request context helpers, and API HTTP style."
name: "Refactor Api Http Style"
applyTo:
  - "internal/api/http/*.go"
  - "internal/api/http/**/*.go"
  - "refactor/internal/api/http/*.go"
  - "refactor/internal/api/http/**/*.go"
---
# Refactor Api Http Style

这份 instruction 只约束 `refactor/internal/api/http` 及其子目录

它是在共享 Go 宪法之上的 API HTTP 层专项补充

## 核心定位

- API HTTP 层是薄适配层 只负责路由绑定 鉴权接入 请求解析 调用 app 层 与 HTTP 响应翻译
- HTTP handler 不承载业务规则 不直接操作 repo domain service 或 infra entity
- 对外响应形状以 `res.HttpRes` 为准 不以 app 层的 `val` 或 `res.AppRes` 裸露给客户端

## handler 规则

- 每个 route handler function 都必须有完整 godoc 注释 作为 swagger 文档生成源
- handler 工厂签名统一为 `func Xxx(st *state.AppState) iris.Handler`
- handler 方法体只做四类事 读取 HTTP 输入 取当前用户 调用 app 翻译为 HTTP 响应
- 路径参数 查询参数 body cookie header 的读取与基础校验留在 handler 层完成
- 若某个 app 入参结构体字段已声明 `url:"..."` tag 则必须使用统一 query 绑定流程自动注入 不允许继续手写 `Params().Get` 或 `URLParam*` 再手动拼装该字段
- body 解析失败 缺少必要路径参数 缺少当前用户等 HTTP 级错误 直接在 handler 层返回 `res.Reject`
- handler 不直接写 `cx.JSON` `cx.WriteString` 这类原始输出 优先统一走 `res.Accept` 与 `res.Reject`

## 路由语义规则

- 返回资源集合时 路由必须使用资源自身集合路径 例如 `/api/v1/assignments` `/api/v1/chapters` `/api/v1/members`
- 返回单资源时 路由必须使用资源自身 id 例如 `/api/v1/assignments/{assignment_id}` `/api/v1/chapters/{chapter_id}`
- 当集合接口需要按其他资源筛选时 外部资源 id 只能放 query 参数 例如 `?chapter_id=` `?team_id=` `?workset_id=` 不允许写成 `/assignments/chapters/{chapter_id}` 这类 foreign-key path
- 路径最后一段只允许是资源本身复数名 资源本身 id 或资源本身操作名 不允许在资源路径直接插入其他资源字段
- swagger `@Param` 必须与真实输入来源一致 若输入来自 query 则必须写 `query` 不允许误写成 `path`

## swagger godoc 规则

- 每个公开 route handler 的 godoc 至少要包含 `@Summary` `@Description` `@Tags` `@Produce` `@Router`
- 需要 body 的接口补 `@Accept json` 与对应 `@Param body body ...`
- 有路径参数 查询参数或其他输入时 必须逐项写清 `@Param`
- 需要鉴权的接口必须写 `@Security ApiKeyAuth`
- 若运行时同时接受 `Authorization` header 与 `authorization` cookie 必须在 godoc 中明确写出 不能只让读者从 middleware 猜测
- 当前 swagger 全局安全定义注册的是 header 模式 因此 cookie 兜底至少要写进 `@Description` 或额外参数说明中
- godoc 中的成功与失败返回必须描述真实传输层形状 不能把裸 `val.Xxx` 误写成最终响应对象

## 响应文档规则

- 运行时只要返回 JSON 外层都应以 `res.HttpRes` 为准 因此 swagger 注释必须围绕 `res.HttpRes` 描述
- 对成功响应 优先使用能表达包裹体的 swagger 写法 例如 `res.HttpRes` 或工具链支持时的 `res.HttpRes{data=val.Xxx}` 形式
- 对失败响应 也要标注为 `res.HttpRes` 包裹的错误消息 而不是裸字符串或裸 app 错误
- 若某个成功分支实际传入 `res.Accept(cx, code, nil)` 导致当前 helper 只写状态码不写 JSON body godoc 必须如实反映这一事实 不能伪造有 `data` 包裹体
- handler 对 app 返回值做二次文案转换时 swagger 里的失败消息应描述传输结构 而不是假设透传 app 层原消息

## 鉴权与上下文规则

- 需要登录态的路由应放在启用了 `middleware.Auth()` 的路由组下
- 这类 handler 内获取当前用户时 应统一走 `takeCurrUid`
- `newReqCx` 是 app 层调用上下文的唯一标准入口 用它注入 request scoped logger 与 request id
- 不要在各 handler 内重复手写 `context.WithValue` 或重复拼接 logger
- middleware 与 util 之间共享的上下文键必须集中定义 不能在多个文件里散落 magic string
- 运行时若 cookie 优先于 header 这类行为存在 必须在文档或注释中明确 让接口使用方知道真实优先级

## 强制个人约定

- 当前用户 id 局部变量统一命名为 `currUid` 并原样传给 app 层 禁止使用 `currUserId` 等其他命名
- 对 `List` 接口必须显式接收分页输入（例如 `offset` `limit`）并传递到 app 的 `val.ListXxxArgs`
- 当需求被明确要求为“闭环”或“全链路”时 实现范围必须覆盖到 HTTP 层 至少包含 route 注册 handler godoc `AppState` 注入与 `main` 中的 app 构造接线 仅完成 domain app repo 不算完成

## 文件与组织规则

- 按资源主题拆文件 例如 `auth.go` `user.go` `team.go`
- 通用 helper 放在 `util.go`
- 通用响应结构放在 `res/`
- middleware 放在 `middleware/` 且只承载横切逻辑 例如鉴权与请求日志
- 新增 route 时 优先放进对应资源文件 不把无关 handler 堆到一个总入口文件中

## 建议补充模式

- handler 对 app 的 reject 分支应优先使用 `re.Code()` 作为 HTTP 状态来源 保持 app 与 transport 的失败码一致
- swagger `@Description` 应写清认证来源 返回包裹体与关键前置条件 不只写一行重复 `@Summary` 的空描述
- 认证信息的读取 键名存取 与当前用户提取应共用 helper 与常量 避免 `middleware` 与 `util` 各写一份字符串导致错位
- 对需要 swagger 暴露的响应 若 `res.HttpRes` 的信息不够表达真实 `data` 形状 应考虑补充文档专用 response type 或使用包裹体字段覆盖写法
- 中间件与 helper 也应遵守完整 godoc 要求 因为这层经常直接影响对外 API 行为
