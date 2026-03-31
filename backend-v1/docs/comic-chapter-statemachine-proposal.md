# Comic & Chapter 状态机架构重构方案

> 本方案将 Chapter 的工作流管理重构为显式状态机，并引入领域事件机制，使 Comic 的副本同步和各种副作用通过事件驱动自动完成。

---

## 1. 设计目标

1. **Chapter 工作流以状态机建模**：每个工作流拥有明确的状态集合、合法转换和前置条件
2. **引入领域事件**：状态转换产生事件，副作用由事件处理器消费，与核心逻辑解耦
3. **Comic 最新章节副本改为事件驱动**：不再需要手动调用 `SyncLatestChapterReplica`
4. **API 层改为面向事件的命令**：UpdateChapter 拆分为语义化的工作流操作命令
5. **保持数据库 schema 不变**：重构限于 Go 代码层，不改变表结构

---

## 2. Chapter 工作流状态机定义

### 2.1 状态定义

Chapter 包含 6 个独立的工作流状态机，每个对应一条流水线工位：

```
Upload:     Pending ──────────────────────────→ Completed
Translate:  Pending → InProgress → Completed
Proofread:  Pending → InProgress → Completed
Typeset:    Pending → InProgress → Completed
Review:     Pending ──────────────────────────→ Completed
Publish:    Pending ──────────────────────────→ Completed
```

### 2.2 跨工作流前置条件（可选的严格模式）

```
Translate.Start   requires  Upload == Completed
Proofread.Start   requires  Translate == Completed
Typeset.Start     requires  Upload == Completed
Review.Start      requires  Translate == Completed AND Proofread == Completed AND Typeset == Completed
Publish.Start     requires  Review == Completed
```

> 注意：当前实现不强制这些前置条件。引入状态机后可以轻松开启或关闭。

### 2.3 领域模型: ChapterWorkflowMachine

```go
// domain/model/chapter_workflow.go

// WorkflowPhase 表示单个工作流的当前阶段
type WorkflowPhase int

const (
    PhasePending    WorkflowPhase = iota // 待处理（*At 全 nil）
    PhaseInProgress                       // 进行中（startAt != nil, doneAt == nil；仅三态工作流）
    PhaseCompleted                        // 已完成（doneAt != nil）
)

// WorkflowTransition 表示一次合法的状态转换
type WorkflowTransition struct {
    Workflow Workflow
    From     WorkflowPhase
    To       WorkflowPhase
}

// ChapterWorkflowState 封装一个章节的完整工作流快照
type ChapterWorkflowState struct {
    Upload    WorkflowPhase
    Translate WorkflowPhase
    Proofread WorkflowPhase
    Typeset   WorkflowPhase
    Review    WorkflowPhase
    Publish   WorkflowPhase
}

// DeriveFromChapterInfo 从当前 ChapterInfo 的时间戳推导出各工作流阶段
func DeriveWorkflowState(chapter ChapterInfo) ChapterWorkflowState {
    return ChapterWorkflowState{
        Upload:    derivePhase(nil, chapter.UploadedAt),          // 二态
        Translate: derivePhase(chapter.TransalatingAt, chapter.TranslatedAt),
        Proofread: derivePhase(chapter.ProofreadingAt, chapter.ProofreadAt),
        Typeset:   derivePhase(chapter.TypesettingAt, chapter.TypesetAt),
        Review:    derivePhase(nil, chapter.ReviewedAt),          // 二态
        Publish:   derivePhase(nil, chapter.PublishedAt),         // 二态
    }
}

func derivePhase(startAt *time.Time, doneAt *time.Time) WorkflowPhase {
    if doneAt != nil {
        return PhaseCompleted
    }
    if startAt != nil {
        return PhaseInProgress
    }
    return PhasePending
}
```

### 2.4 转换合法性校验

```go
// domain/model/chapter_workflow.go

// ValidateTransition 验证转换是否合法，并可校验跨工作流前置条件
func (state ChapterWorkflowState) ValidateTransition(
    transition WorkflowTransition,
    enforcePrerequisites bool,
) error {
    currentPhase := state.phaseOf(transition.Workflow)

    // 1. 校验 From 是否匹配当前状态
    if currentPhase != transition.From {
        return fmt.Errorf(
            "工作流 %s 当前阶段为 %s，无法执行 %s → %s 转换",
            transition.Workflow, currentPhase, transition.From, transition.To,
        )
    }

    // 2. 校验转换方向合法性（不允许回退）
    if transition.To <= transition.From {
        return fmt.Errorf("不允许回退: %s → %s", transition.From, transition.To)
    }

    // 3. 校验跨工作流前置条件
    if enforcePrerequisites {
        return state.checkPrerequisites(transition)
    }

    return nil
}

func (state ChapterWorkflowState) checkPrerequisites(transition WorkflowTransition) error {
    switch transition.Workflow {
    case WorkflowTranslating:
        if transition.To >= PhaseInProgress && state.Upload != PhaseCompleted {
            return errors.New("开始翻译前须完成上传")
        }
    case WorkflowProofreading:
        if transition.To >= PhaseInProgress && state.Translate != PhaseCompleted {
            return errors.New("开始校对前须完成翻译")
        }
    case WorkflowTypesetting:
        if transition.To >= PhaseInProgress && state.Upload != PhaseCompleted {
            return errors.New("开始嵌字前须完成上传")
        }
    case WorkflowReviewing:
        if transition.To >= PhaseCompleted &&
            (state.Translate != PhaseCompleted ||
             state.Proofread != PhaseCompleted ||
             state.Typeset != PhaseCompleted) {
            return errors.New("审核前须完成翻译、校对和嵌字")
        }
    case WorkflowPublishing:
        if transition.To >= PhaseCompleted && state.Review != PhaseCompleted {
            return errors.New("发布前须完成审核")
        }
    }
    return nil
}
```

---

## 3. 领域事件体系

### 3.1 事件定义

```go
// domain/model/event.go

// DomainEvent 是所有领域事件的基接口
type DomainEvent interface {
    EventName() string
}

// ChapterWorkflowTransitioned 章节工作流发生了一次状态转换
type ChapterWorkflowTransitioned struct {
    ChapterID string
    ComicID   string
    Workflow   Workflow
    FromPhase  WorkflowPhase
    ToPhase    WorkflowPhase
    ActorID    string    // 操作者
    OccurredAt time.Time
}

func (e ChapterWorkflowTransitioned) EventName() string {
    return "chapter.workflow.transitioned"
}

// ChapterCreated 章节被创建
type ChapterCreated struct {
    ChapterID string
    ComicID   string
    CreatorID string
    OccurredAt time.Time
}

func (e ChapterCreated) EventName() string {
    return "chapter.created"
}

// ChapterDeleted 章节被删除
type ChapterDeleted struct {
    ChapterID string
    ComicID   string
    ActorID   string
    OccurredAt time.Time
}

func (e ChapterDeleted) EventName() string {
    return "chapter.deleted"
}

// ComicCreated 漫画被创建
type ComicCreated struct {
    ComicID   string
    WorksetID string
    CreatorID string
    OccurredAt time.Time
}

func (e ComicCreated) EventName() string {
    return "comic.created"
}

// ComicDeleted 漫画被删除
type ComicDeleted struct {
    ComicID   string
    ActorID   string
    OccurredAt time.Time
}

func (e ComicDeleted) EventName() string {
    return "comic.deleted"
}
```

### 3.2 事件总线

```go
// domain/model/event_bus.go

// EventHandler 处理一个特定类型的领域事件
type EventHandler func(executor repository.Executor, event DomainEvent) error

// EventBus 收集事务内事件和事务后任务
type EventBus struct {
    // transactionalHandlers 在事务内同步执行
    transactionalHandlers map[string][]EventHandler
    // postCommitHandlers 在事务提交后异步执行
    postCommitHandlers    map[string][]EventHandler
}

func NewEventBus() EventBus { ... }

func (bus EventBus) OnTransactional(eventName string, handler EventHandler) { ... }
func (bus EventBus) OnPostCommit(eventName string, handler EventHandler)    { ... }

// Dispatch 在事务内分发事件，收集 postCommit 任务
func (bus EventBus) Dispatch(
    executor repository.Executor,
    events []DomainEvent,
) (postCommitTasks []func(), err error) { ... }
```

> EventBus 是一个简单的进程内同步派发器。事务内 handler 在同一个事务中执行（共享 executor），保证一致性。postCommit handler 被收集后在事务提交成功后执行。

---

## 4. 重构后的 Application 层

### 4.1 ChapterApplication 新接口

当前的 `UpdateChapter` 是一个大杂烩接口，允许同时修改多个工作流状态。重构后拆分为语义化命令：

```go
type ChapterApplication interface {
    // --- 生命周期命令 ---
    CreateChapter(scope, currentUserID, args CreateChapterArgs) (CreateChapterResult, error)
    DeleteChapter(scope, currentUserID, chapterID string) error
    UpdateChapterSubtitle(scope, currentUserID, args UpdateChapterSubtitleArgs) error

    // --- 工作流转换命令 ---
    TransitionWorkflow(scope, currentUserID, args TransitionWorkflowArgs) error

    // --- 查询 ---
    ListChapters(scope, currentUserID, args ListChapterArgs) ([]ChapterInfo, error)
}
```

#### TransitionWorkflowArgs

```go
type TransitionWorkflowArgs struct {
    ChapterID string              `json:"chapter_id"`
    Workflow   model.Workflow      `json:"workflow"`        // e.g. "uploading"
    Action     model.WorkflowAction `json:"action"`         // e.g. "start" | "complete" | "reset"
}
```

**WorkflowAction 映射到状态转换**：

| Action     | 二态工作流(Upload/Review/Publish) | 三态工作流(Translate/Proofread/Typeset) |
| ---------- | --------------------------------- | --------------------------------------- |
| `start`    | — (不适用)                        | Pending → InProgress                    |
| `complete` | Pending → Completed               | InProgress → Completed                  |
| `reset`    | Completed → Pending               | \* → Pending                            |

### 4.2 TransitionWorkflow 实现流程

```
1. 验证参数
2. 开事务 → LockByID → Get 当前 ChapterInfo
3. DeriveWorkflowState(currentChapter)
4. 计算 WorkflowTransition{ Workflow, From: current, To: target }
5. state.ValidateTransition(transition, enforcePrerequisites)
6. 鉴权: PermChapterUpdate(userID, chapterID, []Workflow{args.Workflow}, ...)
7. 应用转换 → 生成 ChapterUpdate（仅修改一个工作流的时间戳）
8. ChapterRepository.Update
9. 生成事件: ChapterWorkflowTransitioned{...}
10. EventBus.Dispatch(executor, events)
    → 事务内 handler: SyncLatestChapterReplica
    → 事务内 handler: UpdateUserStats (if upload completed)
11. 提交事务
12. 执行 postCommit 任务
    → CleanupChapterPages (if upload completed)
    → 未来: 发送通知、触发 Webhook 等
```

### 4.3 事件处理器注册

```go
// application/event_handler.go

func RegisterChapterEventHandlers(bus model.EventBus, deps Dependencies) {
    // --- 事务内 ---

    // 章节工作流变更 → 同步漫画最新章节副本
    bus.OnTransactional("chapter.workflow.transitioned", func(executor, event) error {
        e := event.(model.ChapterWorkflowTransitioned)
        return deps.ComicRepository.SyncLatestChapterReplica(executor, e.ComicID)
    })

    // 上传完成 → 更新用户统计
    bus.OnTransactional("chapter.workflow.transitioned", func(executor, event) error {
        e := event.(model.ChapterWorkflowTransitioned)
        if e.Workflow != model.WorkflowUploading || e.ToPhase != model.PhaseCompleted {
            return nil
        }
        return handleUploadCompletedStats(executor, deps, e.ChapterID)
    })

    // 章节创建 → 同步漫画章节计数
    bus.OnTransactional("chapter.created", func(executor, event) error {
        e := event.(model.ChapterCreated)
        return deps.ComicRepository.SyncLatestChapterReplica(executor, e.ComicID)
    })

    // 章节删除 → 同步漫画章节计数
    bus.OnTransactional("chapter.deleted", func(executor, event) error {
        e := event.(model.ChapterDeleted)
        return deps.ComicRepository.SyncLatestChapterReplica(executor, e.ComicID)
    })

    // --- 事务后 ---

    // 上传完成 → 异步清理页面 OSS 资源
    bus.OnPostCommit("chapter.workflow.transitioned", func(_, event) error {
        e := event.(model.ChapterWorkflowTransitioned)
        if e.Workflow != model.WorkflowUploading || e.ToPhase != model.PhaseCompleted {
            return nil
        }
        go cleanupChapterPagesAfterUploaded(deps, e.ChapterID)
        return nil
    })

    // 未来扩展示例：
    // bus.OnPostCommit("chapter.workflow.transitioned", notifyTeamMembers)
    // bus.OnPostCommit("chapter.workflow.transitioned", triggerWebhook)
    // bus.OnTransactional("chapter.workflow.transitioned", recordActivityLog)
}
```

---

## 5. 重构后的 API 层

### 5.1 路由变更

当前：

```
PATCH /chapters/{chapter_id}   — body 包含所有可修改字段
```

重构后：

```
PATCH  /chapters/{chapter_id}                            — 仅更新 subtitle
POST   /chapters/{chapter_id}/workflow/{workflow}/start    — 开始工作流
POST   /chapters/{chapter_id}/workflow/{workflow}/complete — 完成工作流
POST   /chapters/{chapter_id}/workflow/{workflow}/reset    — 重置工作流
```

或者保留单端点但改变语义：

```
POST /chapters/{chapter_id}/transitions   — body: { "workflow": "uploading", "action": "complete" }
```

> 推荐后者（单端点），更符合 CQRS 命令模式，且不需要改动太多路由。

### 5.2 Handler 示例

```go
// TransitionChapterWorkflow godoc
// @Summary     章节工作流状态转换
// @Description 对指定章节执行一次工作流状态转换（开始/完成/重置）
//
// @Tags        chapter
// @Security    ApiKeyAuth
// @Accept      json
// @Produce     json
// @Param       chapter_id path string true "章节 ID"
// @Param       body body value.TransitionWorkflowArgs true "工作流转换参数"
//
// @Success     200
//
// @Router      /chapters/{chapter_id}/transitions [post]
func TransitionChapterWorkflow(appState *state.AppState) iris.Handler {
    chapterApplication := appState.ChapterApplication

    return func(ctx iris.Context) {
        currentUserID, ok := extractCurrentUserID(ctx)
        if !ok {
            return
        }

        chapterID := ctx.Params().Get("chapter_id")
        if chapterID == "" {
            reject(ctx, iris.StatusBadRequest, "缺少 chapter_id 路径参数")
            return
        }

        var args value.TransitionWorkflowArgs
        if err := ctx.ReadJSON(&args); err != nil {
            reject(ctx, iris.StatusBadRequest, "请求体格式错误: "+err.Error())
            return
        }
        args.ChapterID = chapterID

        if err := chapterApplication.TransitionWorkflow(
            buildTraceScope(ctx),
            currentUserID,
            args,
        ); err != nil {
            reject(ctx, iris.StatusBadRequest, err.Error())
            return
        }

        accept(ctx, "工作流状态转换成功", nil)
    }
}
```

---

## 6. 新增文件清单

| 文件                               | 职责                                                              |
| ---------------------------------- | ----------------------------------------------------------------- |
| `domain/model/event.go`            | 领域事件类型定义                                                  |
| `domain/model/event_bus.go`        | 进程内事件总线                                                    |
| `domain/model/chapter_workflow.go` | ChapterWorkflowState, WorkflowPhase, ValidateTransition           |
| `domain/model/workflow_action.go`  | WorkflowAction 类型（start/complete/reset）及到 Transition 的映射 |
| `application/event_handler.go`     | 事件处理器注册                                                    |
| `value/transition.go`              | TransitionWorkflowArgs DTO                                        |

### 需修改的文件

| 文件                      | 变更                                                               |
| ------------------------- | ------------------------------------------------------------------ |
| `application/chapter.go`  | 拆分 UpdateChapter 为 UpdateChapterSubtitle + TransitionWorkflow   |
| `api/http/chapter.go`     | 新增 TransitionChapterWorkflow handler，修改 UpdateChapter handler |
| `domain/model/chapter.go` | NewChapterUpdate 简化为仅处理单个工作流转换（或保留兼容）          |
| `value/chapter.go`        | 新增 TransitionWorkflowArgs，UpdateChapterArgs 移除工作流字段      |
| `state/app_state.go`      | 注入 EventBus                                                      |

---

## 7. 迁移策略

### Phase 1: 引入状态机模型（纯重构，不改 API）

1. 新增 `domain/model/chapter_workflow.go`，实现 WorkflowPhase 和 ValidateTransition
2. 在现有 `NewChapterUpdate` 内部嵌入 ValidateTransition 校验
3. 新增 `domain/model/event.go` 和 `event_bus.go`
4. 将 `handleChapterUploaded` 和 `SyncLatestChapterReplica` 调用迁移为事件处理器
5. 保持现有 API 签名不变

### Phase 2: API 层重构

1. 新增 `POST /chapters/{chapter_id}/transitions` 端点
2. UpdateChapter (PATCH) 仅保留 subtitle 更新
3. 旧 API 内部调用新的 TransitionWorkflow 方法，保持向后兼容
4. 前端迁移完成后移除旧 API 中的工作流字段

### Phase 3: 启用严格前置条件（可选）

1. 通过配置开关 `enforceWorkflowPrerequisites` 控制
2. 逐步在生产环境验证后全量启用

---

## 8. 设计决策与权衡

### 为什么用进程内 EventBus 而非消息队列？

当前系统是单体架构，引入 MQ 会增加运维复杂度。进程内 EventBus 已足够解耦事件产生与消费。未来若拆分微服务，EventBus 接口可平滑替换为 MQ 适配器。

### 为什么不将 6 个工作流合并为单一状态机？

6 个工作流在业务上是可以并行推进的（如嵌字和翻译可以同时进行），它们不是线性流水线。合并为单一状态将丢失这种并行语义。每个工作流独立建模更准确地反映业务现实。

### reset 操作的必要性

当前实现中可以将状态从 completed 设回 pending（即回退）。这在实际工作中是必要的（例如审核打回需要重新翻译）。状态机中通过 `reset` 动作显式支持，并且 reset 可以携带额外语义（如记录回退原因）。

### subtitle 更新为什么独立出来？

subtitle 是纯属性更新，与工作流状态无关。将其与状态转换混在同一个接口会导致权限模型复杂化（更新 subtitle 的权限不应依赖工作流角色）。分离后 subtitle 更新可使用独立的权限规则。

---

## 9. 未来扩展方向

利用事件驱动架构，以下功能可通过添加事件处理器实现，无需修改核心业务逻辑：

- **活动日志**：监听所有事件，记录操作历史（谁在什么时候做了什么）
- **团队通知**：工作流完成事件触发下游角色通知（如翻译完成通知校对人员）
- **Webhook**：对外推送章节状态变更
- **数据统计**：工作流耗时统计（从 start 到 complete 的时间差）
- **自动化流水线**：上传完成后自动启动翻译流程
