# Workset Include 底层逻辑说明

## 1. CheckResourceNeccessary 的正确语义

文件：`internal/value/include.go`

`CheckResourceNeccessary(rawOption, options...)` 的约束如下：

- 第一个参数 `rawOption` 必须是用户传入的 include 原字符串，例如：`team`、`chapter.comic`
- 第二个参数 `options...` 由应用侧定义为“允许的嵌套路径”，按层级从浅到深传入
- 匹配规则：
  - `rawOption` 会按 `.` 拆分
  - 每一段必须与 `options` 对应层级相等
  - `rawOption` 深度不能超过 `options` 深度
  - 任一段为空或不相等则返回 `false`

示例：

- `CheckResourceNeccessary("team", "team") == true`
- `CheckResourceNeccessary("chapter.comic", "chapter", "comic") == true`
- `CheckResourceNeccessary("chapter.creator_id", "chapter", "comic") == false`
- `CheckResourceNeccessary("team.extra", "team") == false`

## 2. includes 解析职责下沉到 service 纯函数

文件：`internal/domain/service/workset_include.go`

新增纯函数：

- `ResolveWorksetListIncludeSpec(includes []string) WorksetListIncludeSpec`

输出结构：

- `NeedTeam bool`

该函数只负责“解释 include 请求”，不触发 IO，不依赖 repository，实现可测试的纯逻辑层。

## 3. app 层职责

文件：`internal/application/workset.go`

`ListWorksets` 当前流程：

1. 调用 `ResolveWorksetListIncludeSpec(args.Includes)` 得到 include 规格
2. 构建基础 `query options`：team_id 过滤、index 排序、分页
3. 若 `NeedTeam == true`，只在 options 数组追加 `query_option.WorksetQuery().IncludeTeamInfo()`
4. 统一调用一次 `worksetRepository.List(...)`
5. 在返回组装阶段按 `NeedTeam` 决定是否填充 `workset.team`

app 层只做编排，不直接拼装两种 repository 结果，也不走分叉的提前返回路径。

## 4. repo 层聚合实现

文件：`internal/infrastructure/repository/query_option/workset.go`

新增 query option：

- `IncludeTeamInfo()`

功能：

- `JOIN team_table`
- `SELECT` workset 主表字段 + team 别名字段

文件：`internal/infrastructure/repository/workset.go`

聚合与普通列表都通过统一的 `List(executor, options...)` 执行。

文件：`internal/infrastructure/repository/entity/workset.go`

新增 team 聚合列映射与转换函数：

- `ToWorksetWithTeamInfo`

## 5. 返回值注入

文件：`internal/value/workset.go`

仅保留单一转换函数：

- `NewWorksetInfoFromModel`

用途：

- 在包含 team 的查询路径中，将 `TeamInfo` 注入 `WorksetInfo.Team`
- 在不包含 include 的路径中，保持 `Team` 为空，避免过度返回
