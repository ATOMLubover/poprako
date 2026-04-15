# ComicDetailModal 可操作按钮分析

## 1. 背景与布局

**ComicDetailModal** 的大致结构：

```
┌──────────────────────────────────────────────────────┐
│  上半部分：当前章节详情                                 │
│   - 基本信息（章节序号、副标题、页数、unit 翻译进度）      │
│   - 工作流时间线（9 个阶段的状态）                       │
│   - 分配成员列表                                       │
│   - 操作按钮区（因身份不同而变化）                       │
├──────────────────────────────────────────────────────┤
│  底部页脚：章节横向滑动列表                              │
│   - 默认选中 pinned chapter                            │
│   - 可切换查看其他章节                                  │
└──────────────────────────────────────────────────────┘
```

---

## 2. 角色体系

### 2.1 两层权限

| 层级       | 来源                   | 作用                         |
| ---------- | ---------------------- | ---------------------------- |
| **团队层** | `MemberInfo.Roles`     | 控制漫画/章节的增删改权限    |
| **章节层** | `AssignmentInfo.Roles` | 控制工作流推进和人员分配权限 |

> **关键区别**：`Admin` 只存在于团队层（MemberInfo）；`RoleReviewer` 只存在于章节层（AssignmentInfo）。两者互不依赖。

### 2.2 角色清单

| 角色常量          | 中文名 | 所在层 | 说明                                                                         |
| ----------------- | ------ | ------ | ---------------------------------------------------------------------------- |
| `RoleAdmin`       | 管理员 | 团队层 | 可对漫画/章节进行 CRUD 操作                                                  |
| `RoleReviewer`    | 监修   | 章节层 | 章节创建时自动分配给创建者；可代替所有职位推进工作流；可管理该章节的所有分配 |
| `RoleRawProvider` | 图源   | 章节层 | 上传页面原图                                                                 |
| `RoleTranslator`  | 翻译   | 章节层 | 负责翻译各 unit                                                              |
| `RoleProofreader` | 校对   | 章节层 | 负责校对各 unit                                                              |
| `RoleTypesetter`  | 嵌字   | 章节层 | 负责嵌字排版                                                                 |
| `RolePublisher`   | 发布   | 章节层 | 负责最终发布                                                                 |

### 2.3 前端需要获取的上下文数据

```
myMember      : MemberInfo         // 当前用户的团队成员信息
myAssignment  : AssignmentInfo?    // 当前用户在本章节的分配，可能为 nil
chapter       : ChapterInfo        // 当前章节的完整信息（含工作流时间戳）
```

---

## 3. 工作流阶段与事件

章节的工作流是**非线性**的——各阶段可以并行推进，互不阻塞。

### 3.1 工作流时间戳字段

| 字段              | 含义             | 阶段状态                                       |
| ----------------- | ---------------- | ---------------------------------------------- |
| `uploaded_at`     | 原图上传完成时间 | `nil` = Pending / 非nil = Completed            |
| `translating_at`  | 翻译开始时间     | `nil` = Pending / 非nil = Ongoing              |
| `translated_at`   | 翻译完成时间     | `nil` = Pending or Ongoing / 非nil = Completed |
| `proofreading_at` | 校对开始时间     | `nil` = Pending / 非nil = Ongoing              |
| `proofread_at`    | 校对完成时间     | `nil` = Pending or Ongoing / 非nil = Completed |
| `typesetting_at`  | 嵌字开始时间     | `nil` = Pending / 非nil = Ongoing              |
| `typeset_at`      | 嵌字完成时间     | `nil` = Pending or Ongoing / 非nil = Completed |
| `reviewed_at`     | 监修完成时间     | `nil` = Pending / 非nil = Completed            |
| `published_at`    | 发布完成时间     | `nil` = Pending / 非nil = Completed            |

### 3.2 WorkflowTransition 事件与对应按钮

| 事件（`workflow_transition` 值） | 按钮文案         | 触发条件                                        | 有权操作的角色            |
| -------------------------------- | ---------------- | ----------------------------------------------- | ------------------------- |
| `upload_complete`                | **标记上传完成** | `uploaded_at == nil`                            | `RawProvider`、`Reviewer` |
| `translate_start`                | **开始翻译**     | `translating_at == nil && translated_at == nil` | `Translator`、`Reviewer`  |
| `translate_complete`             | **完成翻译**     | `translating_at != nil && translated_at == nil` | `Translator`、`Reviewer`  |
| `proofread_start`                | **开始校对**     | `proofreading_at == nil && proofread_at == nil` | `Proofreader`、`Reviewer` |
| `proofread_complete`             | **完成校对**     | `proofreading_at != nil && proofread_at == nil` | `Proofreader`、`Reviewer` |
| `typeset_start`                  | **开始嵌字**     | `typesetting_at == nil && typeset_at == nil`    | `Typesetter`、`Reviewer`  |
| `typeset_complete`               | **完成嵌字**     | `typesetting_at != nil && typeset_at == nil`    | `Typesetter`、`Reviewer`  |
| `review_complete`                | **完成监修**     | `reviewed_at == nil`                            | `Reviewer`（仅此角色）    |
| `publish_complete`               | **发布**         | `published_at == nil`                           | `Publisher`、`Reviewer`   |

> **注意**：`Reviewer` 是"超级权限"角色，可代替任意角色推进工作流中的任意阶段。

---

## 4. 按职位分类的可操作按钮

### 4.1 无分配的普通成员（仅有 MemberInfo，无 AssignmentInfo）

仅能**查看**章节信息，没有任何可操作按钮。

_唯一例外_：若该成员同时是团队管理员（Admin），则可见管理类按钮（见 4.6）。

---

### 4.2 图源（RawProvider）

**前置条件**：`myAssignment.HasAnyRole(RoleRawProvider) == true`

| 按钮                               | API 调用                                                 | 可见条件                        |
| ---------------------------------- | -------------------------------------------------------- | ------------------------------- |
| **上传页面**（预留并获取上传地址） | `PageApp.Reserve`                                        | 随时可见（图源可随时追加页面）  |
| **标记页面已上传**                 | `PageApp.Update(is_uploaded=true)`                       | 对应页面 `is_uploaded == false` |
| **标记上传完成**                   | `ChapterApp.Update(workflow_transition=upload_complete)` | `uploaded_at == nil`            |

---

### 4.3 翻译（Translator）

**前置条件**：`myAssignment.HasAnyRole(RoleTranslator) == true`

| 按钮               | API 调用                                                    | 可见条件                                        |
| ------------------ | ----------------------------------------------------------- | ----------------------------------------------- |
| **进入翻译编辑器** | 打开页面编辑器，调用 `UnitApp.Save`                         | 随时可见                                        |
| **开始翻译**       | `ChapterApp.Update(workflow_transition=translate_start)`    | `translating_at == nil && translated_at == nil` |
| **完成翻译**       | `ChapterApp.Update(workflow_transition=translate_complete)` | `translating_at != nil && translated_at == nil` |

---

### 4.4 校对（Proofreader）

**前置条件**：`myAssignment.HasAnyRole(RoleProofreader) == true`

| 按钮               | API 调用                                                    | 可见条件                                        |
| ------------------ | ----------------------------------------------------------- | ----------------------------------------------- |
| **进入校对编辑器** | 打开页面编辑器（校对模式），调用 `UnitApp.Save`             | 随时可见                                        |
| **开始校对**       | `ChapterApp.Update(workflow_transition=proofread_start)`    | `proofreading_at == nil && proofread_at == nil` |
| **完成校对**       | `ChapterApp.Update(workflow_transition=proofread_complete)` | `proofreading_at != nil && proofread_at == nil` |

---

### 4.5 嵌字（Typesetter）

**前置条件**：`myAssignment.HasAnyRole(RoleTypesetter) == true`

| 按钮         | API 调用                                                  | 可见条件                                     |
| ------------ | --------------------------------------------------------- | -------------------------------------------- |
| **开始嵌字** | `ChapterApp.Update(workflow_transition=typeset_start)`    | `typesetting_at == nil && typeset_at == nil` |
| **完成嵌字** | `ChapterApp.Update(workflow_transition=typeset_complete)` | `typesetting_at != nil && typeset_at == nil` |

---

### 4.6 发布（Publisher）

**前置条件**：`myAssignment.HasAnyRole(RolePublisher) == true`

| 按钮     | API 调用                                                  | 可见条件              |
| -------- | --------------------------------------------------------- | --------------------- |
| **发布** | `ChapterApp.Update(workflow_transition=publish_complete)` | `published_at == nil` |

---

### 4.7 监修（Reviewer）

**前置条件**：`myAssignment.HasAnyRole(RoleReviewer) == true`

监修是该 Modal 权限最高的角色，可以代替任意职位推进工作流，并管理全部分配关系。

#### 工作流按钮（代理所有职位）

| 按钮             | 可见条件                                        |
| ---------------- | ----------------------------------------------- |
| **标记上传完成** | `uploaded_at == nil`                            |
| **开始翻译**     | `translating_at == nil && translated_at == nil` |
| **完成翻译**     | `translating_at != nil && translated_at == nil` |
| **开始校对**     | `proofreading_at == nil && proofread_at == nil` |
| **完成校对**     | `proofreading_at != nil && proofread_at == nil` |
| **开始嵌字**     | `typesetting_at == nil && typeset_at == nil`    |
| **完成嵌字**     | `typesetting_at != nil && typeset_at == nil`    |
| **完成监修**     | `reviewed_at == nil`                            |
| **发布**         | `published_at == nil`                           |

#### 分配管理按钮

| 按钮                     | API 调用               | 备注                |
| ------------------------ | ---------------------- | ------------------- |
| **指派成员**（新增分配） | `AssignmentApp.Create` | 选择团队成员 + 角色 |
| **修改成员角色**         | `AssignmentApp.Update` | 全量替换角色掩码    |
| **移除成员分配**         | `AssignmentApp.Remove` | 从章节分配中移出    |

---

### 4.8 管理员（Admin）

**前置条件**：`myMember.HasAnyRole(RoleAdmin) == true`（团队层，不依赖 AssignmentInfo）

| 按钮               | API 调用                            | 位置 / 可见条件                            |
| ------------------ | ----------------------------------- | ------------------------------------------ |
| **新建章节**       | `ChapterApp.Create`                 | 底部页脚，始终可见                         |
| **删除章节**       | `ChapterApp.Remove`                 | 章节详情区，始终可见                       |
| **修改章节副标题** | `ChapterApp.Update(subtitle=...)`   | 内联编辑，始终可见                         |
| **设为顶置章节**   | `ChapterApp.Update(is_pinned=true)` | 章节详情区，当 `is_pinned == false` 时可见 |

> **注意**：Admin 不能直接执行工作流转换，除非他们同时被分配了 Reviewer 角色。

---

## 5. 复合身份示例

一个用户可以同时拥有多个角色，例如：

- **Admin + 无分配**：只能做章节 CRUD，不能推进工作流
- **Admin + Reviewer**：既可以做 CRUD，又可以推进全部工作流并管理分配
- **Translator + Proofreader**：可以同时操作翻译和校对的工作流按钮
- **无 Admin + Reviewer**：可以管理分配和全部工作流，但不能新建/删除章节

---

## 6. 按钮可见性决策树

```
用户打开 ComicDetailModal
├── 是否为该章节有分配？（myAssignment != nil）
│   ├── 否 → 仅展示信息，无操作按钮
│   │         （若同时是 Admin → 展示章节管理按钮）
│   └── 是 → 根据 myAssignment.Roles 展示对应按钮：
│             ├── RawProvider  → 上传/上传完成 按钮
│             ├── Translator   → 翻译进度 + 编辑器入口
│             ├── Proofreader  → 校对进度 + 编辑器入口
│             ├── Typesetter   → 嵌字进度 按钮
│             ├── Publisher    → 发布 按钮
│             └── Reviewer     → 全部工作流按钮 + 分配管理按钮
│
└── 是否为团队 Admin？（myMember.HasAnyRole(RoleAdmin)）
    └── 是 → 额外展示章节管理按钮（新建/删除/改标题/设顶置）
```

---

## 7. 重要业务约束

1. **每个 `WorkflowTransition` 只能触发一次**（状态机不可逆），前端提交后应将对应按钮立即禁用或隐藏。
2. **开始/完成 配对约束**：`translate_complete` 要求先有 `translate_start`；`proofread_complete` 要求先有 `proofread_start`；`typeset_complete` 要求先有 `typeset_start`。
3. **Reviewer 由章节创建触发自动分配**（`ChapterCreatorAssignedEvent`），前端无需手动设置。
4. **分配管理操作（AssignmentApp）要求操作者在该章节有 Reviewer 分配**，与团队 Admin 角色无关。
5. **章节 CRUD（Admin）要求团队层 Admin**，与章节分配无关。
6. **`publish_complete` 是特殊事务**，会同步更新 `comic_table` 的 `pinned_*` 字段（comic replica 列），并更新 `user_stats`。
