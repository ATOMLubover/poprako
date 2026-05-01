---
description: "Use when writing repo mocks or svc/app unit tests under refactor/internal, including mock structs, accessor helpers, test fixtures, helper naming, and test structure conventions."
name: "Refactor Mock & Test Style"
applyTo:
  - "refactor/internal/domain/repo/mock/**/*.go"
  - "refactor/internal/domain/svc/test/**/*.go"
  - "refactor/internal/app/**/test/**/*.go"
---

# Refactor Mock & Test Style

这份 instruction 约束 `refactor/internal` 下的 mock 构建与 svc/app 单元测试。

它是在共享 Go 宪法之上的测试专项补充。

## 核心定位

- mock 是领域 repo 接口的纯内存替身，用于隔离 svc/app 逻辑与 infra 细节
- 单元测试不启动数据库、不启动 HTTP server、不 import infra 包
- 测试验证的是领域逻辑（校验、排序、并发语义、副作用）而非 mock 自身
- mock 不应包含任何测试断言逻辑；断言一律留在 test 文件中

## 目录职责

- `repo/mock` 放 repo 接口的内存实现，包名 `repo_mock`
- `svc/test` 放 domain service 的单元测试，包名 `test`
- `app/**/test` 放 app 层编排逻辑的单元测试，包名 `test`

## Mock 构建规则

### 结构体与构造函数

- mock 实现 struct 必须私有（小写开头），例如 `type mockUnitRepo struct {}`
- 构造函数命名 `NewMock<InterfaceName>`，返回对应的 repo 接口类型而非私有 struct
- 构造函数接受初始数据切片或 nil，内部按主键构建 `map[string]*aggr.Xxx` 索引
- mock 不包含 `sync.Mutex`、channel 或其他并发原语，除非被测逻辑明确依赖并发语义

### 公开访问器

- 为测试代码提供公开的访问器辅助函数，命名 `ListMock<Xxx>By<Key>` / `GetMock<Xxx>ById`
- 访问器接受 repo 接口类型（非私有 struct），内部做 type assertion 到私有 mock struct
- type assertion 失败时返回 nil / 空切片，不 panic
- 访问器必须返回 **clone**，绝不暴露 mock 内部 map 中的真实指针，防止测试代码意外修改 mock 内部状态
- clone 方式：`clone := *original; return &clone`
- 切片返回前需要按业务语义排序（例如按 `Index` 排序），避免依赖 map 的迭代序

### 接口方法实现

- 所有 repo 接口方法在私有 struct 上实现
- `Create` 方法逐字段显式拷贝 `*Cre` → `*Aggr`，不使用反射或自动化映射
- `Save` 语义为 **全量覆盖 PUT**：
  - 若记录已存在，逐字段覆盖所有可写字段
  - 若记录不存在，回退到 `Create` 逻辑（调用自身的 `Create` 方法）
  - Save 不保留旧值、不做 partial merge，因为领域语义就是 PUT
- `Delete` 物理删除（delete from map），非软删除
- `Reindex` 等批量操作逐一处理，顺序执行
- 所有方法返回的 `RepoErr` 默认返回 `nil`；需要模拟错误时通过扩展 mock struct 注入

### 字段覆盖清单

- 实现 mock 方法时必须覆盖 aggr/Cre/Save 结构体中的所有字段，逐一赋值
- 新增领域字段时，mock 方法也必须同步更新
- 不写 `// 省略其他字段` 这类简写注释

## 测试构建规则

### 包与依赖

- 测试包名使用扁平的 `test`（非 `svc_test` 或 `svc`）
- 不引入 testify、goconvey、gomock 等第三方断言/ mock 框架
- 断言全部使用 `t.Helper()` + `t.Fatalf`
- import 分组：标准库 → 空行 → `poprako-s/...` 项目包

### Helper 命名与职责

- `mustXxx` — 调用 `t.Helper()`，失败时 `t.Fatalf`。用于必须成功的操作（如 `mustApplyOpsAccept`）
- `assertXxx` — 调用 `t.Helper()`，失败时 `t.Fatalf`。用于值比较断言（如 `assertOrderEqual`、`assertStringPtrEqual`）
- 辅助构造器如 `strPtr(v string) *string` 不标 must/assert 前缀，直接命名
- 持多个参数时，`t *testing.T` 总是第一参数
- 每个 helper 功能单一，不把多个不相关断言合并进一个 helper

### 测试函数命名

格式：`Test<StructOrInterface>_<Scenario>`

- 对 svc struct 的方法测试：`Test<Struct><Method>_<场景描述>`，例如 `TestUnitSvcApplyOps_RejectsCandOrderMissingSubmittedSaveId`
- 对独立函数测试：`Test<FuncName>_<场景描述>`
- 场景描述用 CamelCase 英文，命中具体行为而非笼统分类

### Fixtures 与内联数据

- 测试数据全部内联构造在测试函数体内，不使用 `testdata/` 外部文件
- 短页面 id 用常量 `const pageId = "p1"` 而非全局变量
- `initUnits` 用 `[]*aggr.Xxx{...}` 切片字面量内联
- diff / ops 用 struct 字面量内联，字段按结构体定义顺序排列
- 明显无关的字段不填（用零值），只在需要验证语义时显式填值

### 测试结构范式

每个测试函数按以下顺序组织：

1. 声明局部常量（pageId、id 片段等）
2. 用 `repo_mock.NewMockXxxRepo(...)` 构造 repo（带上初始数据）
3. 内联构造被测输入（diff、参数等）
4. 调用被测函数/方法
5. 通过 mock 访问器读取副作用结果
6. 用 assert helper 或直接 `t.Fatalf` 做断言

### Reject 路径验证

- 除了断言返回值是 reject 之外，还必须验证 repo 状态未被修改
- 通过 mock 访问器重新读取数据，与初始数据做逐条比对
- 用专门的 `mustRejectXxx` helper 封装上述两步

### 并发/多用户语义测试

- 多客户端编辑场景需显式构造两次 `ApplyOps`（或被测方法）
- 第一次提交模拟客户端 A，第二次提交模拟客户端 B（从同一初始快照出发）
- 验证点：
  - B 的字段值覆盖了 A（later submitter wins）
  - B 的删除生效（不受 A 的 save 影响）
  - A 新增、B 未知的 unit 保留且字段为 A 当时的值
  - 最终排序符合 B 的 CandOrder，B 未知的 unit 聚集在相邻位置

### 用 mock 访问器而非直接访问

- 测试代码不 import mock 包来做 type assertion
- 所有对 mock 内部状态的检查必须通过公开访问器函数进行
- 访问器的排序约定（如 `ListMockUnitsByPage` 按 Index 排序）是测试可依赖的稳定契约

## 错误注入（待扩展）

- 当需要模拟 repo 返回错误时，在 mock struct 中添加可选的 `err` 字段
- 构造函数不自动注入错误，错误由测试在调用前通过公开方法设置
- 每个方法的错误注入可以独立控制（例如 `SetCreateErr`、`SetListErr`）

## 强制个人约定

- mock 访问器永远返回 clone，一个测试修改返回值不应污染另一个测试的 mock 状态
- `Save` mock 行为必须与领域指令中的 PUT 全量覆盖语义一致：先尝试直接覆盖，不存在则 create
- 测试中 `t.Fatalf` 的消息格式用 `"got=%v want=%v"` 或 `"got=%q want=%q"` 的键值对风格
- 删除后验证用 `!= nil` 而非 `== nil` 放在断言中表达"预期已删除、不应存在"
- 所有字符串断言使用 `%q` 而非 `%s` 以便区分空字符串与带引号内容
- 测试 panic 视为 bug，不使用 `t.Run` 子测试——每个场景一个独立 Test 函数
- 测试函数声明后紧跟一段注释（`// ...`）说明该测试验证的具体行为
- `diff.CandOrder` 中的字符串序列断言使用 `assertOrderEqual` 专用 helper，不做逐个 `if got[i] != want[i]` 的裸写
