---
name: domain-service-style
description: |
  记录 domain service 代码规范
  适用于 internal/domain/service 下全部实现与重构任务
---

# domain service style

## 目标

- 所有 domain service 保持无内禀状态
- 所有 repo db 依赖由方法参数显式传入
- 所有 interface 与实现函数声明使用分行风格
- 注释不出现任何句号

## 强制规则

### 1 无内禀状态

- `xxxServiceImpl` 必须是空结构体 `struct{}`
- 不允许在 service struct 中保存 repo db cfg client 等字段
- `NewXxxService` 不接收依赖参数

### 2 显式传入依赖

- 业务方法如果需要 repo 必须通过参数传入
- 示例 `NewCreation(mr repo.MemberRepo ...)`
- 不允许 `s.memberRepo` 这类字段访问

### 3 声明分行

- interface 内每个方法参数必须逐行
- 实现方法参数必须逐行
- 返回值写在独立行结尾

### 4 注释风格

- struct interface function 注释必须有
- 方法体关键步骤前必须有注释
- 注释中禁止 `。` 与 `.`

### 5 空行风格

- 任意两个可执行语句之间保留空行
- if 之前与之后保持一行空白

## 校验清单

- 是否仍有 `type xxxServiceImpl struct {` 且存在字段
- 是否仍有构造函数接收 repo 参数
- 是否仍有单行接口方法声明
- 是否仍有单行实现方法声明
- 注释中是否出现 `。` 或 `.`

## 推荐命令

- `rg -n "type\\s+[a-zA-Z0-9_]+serviceImpl\\s+struct\\s*\\{" internal/domain/service`
- `rg -n "^\\s*//.*[\\.|。]" internal/domain/service`
- `rg -n "^func\\s+\\([^)]*\\)\\s+[A-Za-z0-9_]+\\([^)]*\\)\\s*\\{" internal/domain/service`

## 示例请求

- "按 domain-service-style 重构 internal/domain/service 全部文件"
- "检查所有 service 是否无内禀状态 并给出违规点"
- "把所有 service 的接口与实现函数声明改成分行"
