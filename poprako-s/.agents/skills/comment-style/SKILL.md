---
name: comment-style
description: the appropriate comment style for all Go code in workspace
---

# comment-style

在编写任何 struct、interface 和 function 的注释时，请使用以下格式：

```go
// StructName 代表 ...
type StructName struct {
    // FieldName 表示 ... 的 ...
    FieldName FieldType
}

// InterfaceName 指定了 ... 的行为
type InterfaceName interface {
    // MethodName 查找/加密/...
    MethodName() ReturnType
}

// 函数同上
// 注意，函数的注释不需要解释 **内部** 的具体做法
func FunctionName() ReturnType {
    // 首先，根据 ... 获取 ...

    // 然后，处理 ...

    // 最后，返回 ...
}
```

## When to use

任何 Go 文件都必须遵循这套注释风格，以确保代码的可读性和一致性。

## Instructions

None.
