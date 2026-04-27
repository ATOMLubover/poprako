---
description: "Use when editing refactor app-layer files under refactor/internal/app, including app iface definitions, app impl files, log decorators, app result handling, logger context propagation, 参数清洗, and app-layer style."
name: "Refactor App Layer Style"
applyTo:
  - "refactor/internal/app/*.go"
  - "refactor/internal/app/impl/*.go"
  - "refactor/internal/app/val/*.go"
  - "refactor/internal/app/res/*.go"
---
# Refactor App Layer Style

这份 instruction 只约束 `refactor/internal/app` 及其 `impl` `val` `res` 子目录

它是在共享 Go 宪法之上的 app 层专项补充

## 核心职责

- `app_iface` 负责声明面向上层的用例契约 不承载实现细节
- `app_impl` 负责组织业务流程 编排 domain service repo event ext 等依赖
- `log app impl` 是 `inner app` 的装饰器 负责追加日志上下文 记录入参 并在转发前做最基本的数据清洗 以降低 `inner app` 因空上下文或脏输入直接崩溃的风险
- `app` 层是流程编排层 可以跨多个 repo 与 service 协作 但不应下沉持久化细节到自身

## 包与目录规则

- `refactor/internal/app/*.go` 使用 `package app_iface`
- `refactor/internal/app/impl/*.go` 使用 `package app_impl`
- app interface 与 app implementation 必须分包 不能把接口和实现混在同一个 package
- 一个 app interface 一个文件 例如 `user.go` `user_stats.go`
- 一个 impl struct 一个文件 例如 `userAppImpl` 放在 `impl/user.go` `userLogAppImpl` 放在 `impl/user_log.go`
- 纯 helper 文件允许独立存在 但必须按作用域命名 例如模块专属 helper 放在 `user_util.go` 包级共享 helper 放在 `util.go`
- 不允许在一个 impl 文件中同时堆放多个无关 impl struct

## 接口与返回值规则

- app interface 是上层唯一依赖的 app 入口 构造函数必须返回 `app_iface.XxxApp`
- app 方法第一个参数永远是 `cx context.Context`
- refactor app 层统一使用 `res.AppRes[T]` 表达结果 不走 `(*Res, error)` 双返回值风格
- 成功返回统一使用 `res.Accept`
- 拒绝返回统一使用 `res.Reject`
- 面向业务的数据载体优先落在 `val` 包中 app 方法不要直接向上暴露 domain aggregate
- `Accept` 中的数据通常传入 `*val.Xxx` 或 `*res.None` 保持返回形状统一

## `inner app impl` 规则

- `inner app impl` 只处理业务流程与依赖协作 不负责给 logger 增加请求参数等调用上下文
- `inner app impl` 获取 logger 的唯一入口是 `takeLgr(cx)`
- 除 `err` 或极少数必要诊断字段外 不要在 `inner app impl` 内继续向 logger 注入大段上下文 调用级上下文由 `log app impl` 提前写入
- 需要事务时 在 app 层开启事务 并在事务闭包内部通过 `repo_infra.TxnXxxRepo(cx)` 取到事务态 repo
- app 层负责决定 BadRequest 与 ServerError 这类业务响应语义 但不要把底层错误原文直接暴露给最终消息
- app 层允许组合多个 repo service ext event 完成一个用例 但不要在这里写 domain 规则本身

## 强制个人约定

- 当前用户 id 参数和变量统一使用 `currUid` 命名 禁止使用 `currUserId`
- 事务编排风格统一对齐 `userAppImpl` 已有写法 使用本地 `errCode` + 直接返回真实错误的流程 禁止构造哨兵业务错误（例如 `errNotAdmin`）
- 在事务闭包里拿到事务态 repo 后 局部变量命名使用业务名本身 例如 `memberRepo` `worksetRepo` 禁止添加 `txn` 前缀
- 当 use-case 是 `PUT` 语义时 app 层必须构造全量更新输入 不能把 `nil` 可空字段解释为“跳过更新”
- 删除语义命名固定：`Delete` 表示硬删除 `Remove` 表示软删除
- 已有常见缩写必须统一使用 例如 `Desc`/`desc` 禁止回退到 `Description`/`description`（SQL 字符串除外）
- 当 app 方法除 `cx` 和 `currUid` 外存在超过一个业务参数时 必须封装为 `val` args 结构体
- 所有 `List` 类 app 接口必须带分页参数 并通过 `val.ListXxxArgs` 显式传入
- `inner app impl` 的参数校验与参数归一化必须收敛到同模块 `*_util.go` helper 主流程方法只保留编排与业务步骤

## `log app impl` 规则

- 每个需要日志包装的 app 都应有自己的 `log app impl`
- `log app impl` 是装饰器 只做防御性清洗 日志增强 与调用转发 不承载业务判断
- 进入方法后先兜底 `cx` 如果调用方传入 nil 则转成 `context.Background()`
- 然后从 `cx` 提取 `lgr` 若不存在则回退到 `zap.L()`
- 然后把当前方法的关键入参写入 logger 优先记录稳定 可读 安全的字段
- 参数是结构体指针时 可以记录 `zap.Any("args", args)` 但前提是该对象不会泄露敏感信息 也不会因脏值导致下游崩溃
- 如果原始参数可能包含敏感信息 过深嵌套 或需要归一化后才能安全转发 应在 `log app impl` 先做清洗或裁剪 再注入日志并转发给 `inner app`
- 完成 logger 增强后 必须通过 `saveLgr(cx, lgr)` 写回上下文 再调用 `inner`
- `log app impl` 不要跳过 `inner` 直接复制业务逻辑

## `cx` 与 `lgr` 规则

- `cx` 在 app 层主要承担传递带键值对的 `lgr` 这一职责
- `takeLgr` 与 `saveLgr` 这类上下文 helper 应集中放在 `impl/util.go` 这类共享 helper 文件中
- app 方法之间透传 `cx` 时 不要随意替换整个上下文 只在必要时基于原 `cx` 增强 logger
- 若 `cx` 为 nil `takeLgr` 应提供稳定回退 避免装饰器缺失时出现空指针路径

## 文件内排布规则

- interface 文件按 interface 为中心组织 先 interface 再其紧邻注释与少量同主题声明
- impl 文件按 struct 分区组织 一般顺序是 struct 构造函数 方法
- `log app impl` 文件也遵循同样顺序 先 decorator struct 再构造函数 再各方法
- helper 文件中的函数应按服务对象分组 不要把新 helper 无脑追加到文件底部
- 任何新代码都要优先维护模块局部可读性 而不是追求追加最省事

## 依赖与构造规则

- `NewXxxApp` 负责创建真实 app 实现
- `NewXxxLogApp` 负责创建日志装饰器
- 构造函数必须显式校验依赖 缺失依赖时应立刻失败 通常使用 `zap.L().Panic`
- 真实 app struct 只持有完成该用例编排所需的最小依赖集合 不要因为未来可能会用就提前注入
- 如果一个 helper 只被某个模块使用 应优先放到该模块的 `*_util.go` 中 而不是扩张包级共享 util

## 建议补充模式

- 对外暴露的方法按业务能力切片 不要让一个 app interface 演化成无边界的大杂烩
- 如果一个模块同时存在真实 impl 与 log decorator 其方法集必须严格对齐 避免装饰器漏转发某些方法
- 涉及用户输入的日志记录 默认先考虑脱敏 例如密码 令牌 签名串 上传凭据等不应直接进入日志
- 参数清洗应优先做成轻量 helper 例如规范 nil slice nil map 默认值字段 或裁剪过大的日志载荷
- 从 domain aggregate 组装 `val` 的转换逻辑 应优先收敛到 `*_util.go` 中 例如 `asmUserVal`
- `res.Reject` 的 `msg` 应保持稳定 友好 可直接面向用户 不把内部 repo service ext 错误细节透出
- 若一个流程包含明显的多个业务步骤 `inner app impl` 方法内部应按步骤写注释 并保持步骤之间的空行
