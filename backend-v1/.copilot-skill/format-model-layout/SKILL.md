# Skill: format-model-layout

目的
- 指导 AI 和开发者按项目约定格式化 `value` 层与 `model` 层的 Go 源码文件。
- 强制顺序：`info`、`list`、`create`、`update`（严格按此顺序排列）。
- 强制构造函数约定：`value` 层不需要 `New*` 构造函数；`model` 层必须提供 `New*` 构造函数。

范围
- Workspace 作用域（供本仓库内所有模块使用）。
- 主要面向领域对象（例如 `Comic`、`Chapter`、`Page` 等）。

何时使用
- 新建或重构 `value` / `model` 相关的 Go 文件时。
- 自动化格式化代码或由 AI 生成 model/value 文件时校验约定。

步骤（按顺序）
1. 确定文件目标：是 `value`（DTO）还是 `model`（领域模型）。
2. 对于 `value` 文件，按严格顺序声明类型：
   - `Info`：单个对象详细视图（通常用于获取详情）。
   - `List`：用于列表返回的简化结构或分页项。
   - `Create`：用于创建请求的入参结构（仅包含必要字段）。
   - `Update`：用于更新请求的入参结构（允许零值/可选）。
3. 对于 `model` 文件：先声明领域类型/字段（按自然语义分组），并必须实现 `New*` 构造函数（例如 `NewComic(...)`），构造函数负责初始化不为零值的字段与内嵌聚合。
4. 构造函数签名与行为规范：
   - 名称：`New<TypeName>`，返回 `*TypeName` 或 `TypeName`（本仓库偏好，按现有代码风格，一般返回指针）。
   - 参数：必须为构造该领域对象所需的核心字段（不要直接传入 repository/DB 的类型）。
   - 验证：在构造函数中尽可能进行参数校验或明确文档说明由调用方保证。
5. 指针使用规则：能不使用指针则不使用指针（尤其在 `value` 层避免指针）。
6. 代码注释与 godoc：handler/godoc 在变更时同步更新（仓库约定）。

决策点与分支逻辑
- 可选字段如何表达？
  - `Create`：只包含必需字段。
  - `Update`：可采用指针或 `option` 类型表达“是否修改”，但优先使用项目现有 `option.go` 约定（参见 `util/option.go`）。
- 字段默认值与零值：构造函数（`model` 层）负责设置默认值；`value` 层保持轻量，不封装业务默认。

质量标准（完成检查）
- 文件中类型顺序满足 `info` → `list` → `create` → `update`（对 `value` 文件强制）。
- `value` 文件中不存在 `New*` 构造函数。
- `model` 文件包含 `New*` 构造函数并返回合适类型。
- 遵循仓库的指针使用约定（尽量避免指针，特殊场景除外）。
- 变更通过 `go vet`/`go build`（本地或 CI）无误，且用 `get_errors` 检查静态错误。

示例（Comic）

value/comic.go（示例布局，伪代码）

- `ComicInfo`（Info）: 详细字段（id、title、author、summary、publish_date 等）
- `ComicListItem`（List）: 列表项字段（id、title、cover、status）
- `CreateComic`（Create）: 创建入参（title、author、cover 等必填）
- `UpdateComic`（Update）: 更新入参（可选字段，按项目约定使用 `option` 或指针）

model/comic.go（示例）

- `type Comic struct { ... }`（领域模型字段，可能包含聚合与行为）
- `func NewComic(id ID, title string, author string, ...) *Comic { ... }`（构造函数，负责必要初始化）

AI 使用范式（Prompt 模板）
- "请根据本仓库的 `format-model-layout` 规则，将下面的字段列表格式化为 `value` 层的 Go 文件，严格按 `info`, `list`, `create`, `update` 顺序：<字段列表>。不要添加任何 `New*` 构造函数。"
- "请为领域对象 `<Name>` 生成 `model` 层的 Go 文件，包含类型定义与 `New<Name>` 构造函数。构造函数应初始化必要默认值并返回 `*<Name>`。"

迭代与澄清点
1. 保存草稿并在仓库中运行 `get_errors` 进行静态检查。  
2. 若构造函数返回值（指针或值）在仓库风格中不统一，请指出当前文件或模块的偏好。  
3. 如有特定字段（例如时间戳、UUID、状态字典）需遵循迁移脚本或全局约定，请提供示例或引用对应 SQL/迁移文件路径。

示例快速验证清单
- [ ] `value` 文件类型顺序正确
- [ ] `value` 文件没有 `New*` 构造函数
- [ ] `model` 文件包含 `New*` 构造函数
- [ ] 无不必要的指针使用
- [ ] 通过 `go build` 与 `get_errors` 检查

示例 prompts to try
- "Format `Comic` value types per `format-model-layout` and include `ComicInfo`, `ComicListItem`, `CreateComic`, `UpdateComic`."
- "Generate `model` for `Comic` with `NewComic(title, author, ...)` constructor and default timestamps."

---

备注
- 该 skill 遵循仓库现有约定（例如不滥用指针与 `New` 构造函数的差异）。
- 如需将该规则作为自动化检查（pre-commit / linter），可进一步扩展为脚本或 GitHub Action。
