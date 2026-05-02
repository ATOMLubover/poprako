---
name: refactor-user-style
description: >
  Use whenever editing Go code under `refactor/`, especially when naming local
  variables, shaping transaction flows, or aligning app/http/repo code to the
  repo owner's personal style. Load this before refactoring a slice, fixing
  review feedback, or whenever a draft introduces placeholder names like `res`,
  `txErr`, `cm`, `ws`, `ch`, `v`, or bare `vals`.
---

# Refactor 用户个人风格对齐

## 作用范围

- 本 skill 只补充 `refactor/` 下的个人命名与排版偏好。
- 它不覆盖仓库的硬规则：先服从 `common-code-style-constitution`，再服从
  `AGENTS.md` 里的 `currUid`、`Delete` / `Remove` 等固定约定。

## 金标准文件

在命名拿不准时，先对照下面这些文件，不要自己发明风格：

- `refactor/internal/app/impl/user.go`
- `refactor/internal/app/impl/user_log.go`
- `refactor/internal/api/http/user.go`
- `refactor/internal/infra/repo/user.go`

## 核心原则

选择“最短但仍然有语义”的名字。

用户能接受仓库里已经稳定存在的短名，例如：

- `re`：仅用于 `app_res.AppRes[...]` 这类 app 层临时返回值
- `err`
- `lgr`
- `ev`
- `currUid`

用户不能接受把语义抹掉的占位名，即使它们很常见。

## 禁忌命名

以下命名默认禁止，除非代码本身就是第三方约定或生成代码：

- `res`
- `resp`
- `ret`
- `tmp`
- `txErr`
- `txnErr`
- `cm`
- `ws`
- `ch`
- 单独的 `v`
- 不带领域前缀的 `vals`

这些名字的问题不是“太短”，而是“看不出它到底装的是什么”。

## 推荐命名方式

### 单个领域对象

- `comic`
- `chapter`
- `workset`
- `assignment`
- `invitation`
- `page`
- `user`

不要写成 `cm` / `ch` / `ws`。

### 单个 app 输出值

- `comicVal`
- `chapterVal`
- `worksetVal`
- `exportVal`
- `importRes`

不要写成 `v`。

### 集合值

- `comicVals`
- `chapterVals`
- `worksetVals`
- `assignments`
- `pages`
- `existingUnits`

不要写成裸 `vals`，除非上下文里真的没有更具体的领域名可用。

### 事务结果

`RunWithTxn` 的既定写法就是：

```go
re, err := repo_iface.RunWithTxn(...)
if err != nil {
    ...
    return re
}

return re
```

不要发明：

- `res, txErr := ...`
- `result, txnErr := ...`

### GORM 链式返回值

如果必须接 `*gorm.DB` 返回句柄，按操作命名：

- `queryRe`
- `updRe`
- `claimRe`

不要统一叫 `res`。

## 结构偏好

- 一个局部变量如果承载明确领域语义，就直接用领域名。
- 只有 `app_res` 这类仓库里已经稳定存在的抽象结果，才使用 `re`。
- 组装值对象时，变量名要体现最终返回的东西，而不是写成 `v`。
- 同一函数里不要同时出现多套命名体系，例如一边是 `comic`，另一边又是 `cm`。

## 例子

### Bad

```go
res, txErr := repo_iface.RunWithTxn(...)

cm, err := comicRepo.GetById(comicId)

v := asmComicVal(cm)

vals := make([]val.ComicVal, len(comics))
```

### Good

```go
re, err := repo_iface.RunWithTxn(...)

comic, err := comicRepo.GetById(comicId)

comicVal := asmComicVal(comic)

comicVals := make([]val.ComicVal, len(comics))
```

## 自检清单

改完后至少人工复查这些点：

1. 有没有出现 `txErr` / `txnErr` / `res` 这种占位名。
2. 有没有把 `comic` / `chapter` / `workset` 缩成 `cm` / `ch` / `ws`。
3. 有没有单独的 `v` 或无语义的 `vals`。
4. `RunWithTxn` 是否回到了 `re, err := ...` 的既有写法。
5. 新代码是否看起来像是从 `user` slice 自然长出来的，而不是另一套作者风格。
