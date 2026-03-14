# Copilot 遵守宪章

能不使用指针，则绝对不允许使用指针，否则会引入过多的空值检查。

必须严格按照用户指定的风格进行代码编写，不能随意更改用户的代码风格，否则会严重破坏代码规范和优雅设计。

所有 value 对象、model 对象都必须使用 New 构造函数进行构造，禁止直接使用 struct literal 进行构造，否则会导致代码风格不统一，增加代码维护难度。只有 Args 类对象不需要 New 构造函数，因为这类都是直接由 ReadJSON 反序列化构造的，不需要 New 构造函数，否则会导致代码冗余和不必要的复杂性。

value 对象是 API 层和 APP 层的数据传输对象，model 对象是领域模型对象（用于 repository 和 APP 层通信）。

权限校验使用 Perm 开头而不是 CheckPermission 的函数，因为使用字符串而不是类型的权限管理方式是老版本即将被重构的内容。

每次修改后，必须使用 get_errors MCP 检查代码是否符合规范、有静态错误。

无须更新 swagger 文档，这是自动生成的。但是 handler 上的 godoc 需要更新，以保持文档的准确性和完整性。

除了 (wa worksetApplication) 这种关联对象，其他的参数一律不得使用简写，必须使用全称，否则会导致代码可读性差，增加理解难度。

所有的数据的 single source of truth 全部在于 migrations/ 下的数据库 SQL 脚本。任何 Go 代码中的字段都可能因为 SQL 脚本改变而失效或缺失。

大部分情况下，你需要的格式化代码都在 justfile 中。只有测试 build 不应该使用 just build，它是最终打包为 Linux 可执行文件时才使用。
