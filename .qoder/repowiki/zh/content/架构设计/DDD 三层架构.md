# DDD 三层架构

<cite>
**本文引用的文件**
- [main.go](file://main.go)
- [workset.go](file://internal/application/workset.go)
- [workset.go](file://internal/infrastructure/repository/workset.go)
- [workset.go](file://internal/domain/model/workset.go)
- [workset.go](file://internal/api/http/workset.go)
- [types.go](file://internal/domain/repository/types.go)
- [loader.go](file://internal/application/adapter/loader.go)
- [permission.go](file://internal/domain/model/permission.go)
- [workset.go](file://internal/application/assembler/workset.go)
- [workset.go](file://internal/value/workset.go)
- [http.go](file://internal/api/http/http.go)
- [oss_factory.go](file://internal/infrastructure/external/oss_factory.go)
- [oss.go](file://internal/domain/external/oss.go)
- [workset_include.go](file://internal/domain/service/workset_include.go)
</cite>

## 目录
1. [引言](#引言)
2. [项目结构](#项目结构)
3. [核心组件](#核心组件)
4. [架构总览](#架构总览)
5. [详细组件分析](#详细组件分析)
6. [依赖分析](#依赖分析)
7. [性能考虑](#性能考虑)
8. [故障排查指南](#故障排查指南)
9. [结论](#结论)
10. [附录](#附录)

## 引言
本文件面向 poprako-web-ms 项目的 DDD 三层架构实践，系统性阐述 Domain Layer（领域层）、Application Layer（应用层）、Infrastructure Layer（基础设施层）的设计理念、职责边界与交互关系。特别地，本文聚焦项目在传统 DDD 上的创新实践：
- 领域服务不依赖仓储层：权限校验等业务规则通过“加载器”以函数注入的方式在应用层完成，避免领域模型直接感知仓储接口。
- 事务处理在应用层：应用层统一开启/提交/回滚事务，确保业务操作的原子性与一致性。

同时，本文提供多幅架构图与序列图，帮助读者从宏观到微观全面理解代码组织与数据流。

## 项目结构
项目采用按层次划分的目录结构，清晰分离三层职责：
- Domain Layer：领域模型、领域服务、外部接口定义
- Application Layer：应用服务、装配器、适配器、值对象
- Infrastructure Layer：基础设施实现（仓储、外部服务）

```mermaid
graph TB
subgraph "应用入口"
MAIN["main.go<br/>应用初始化与依赖注入"]
end
subgraph "表示层(HTTP)"
HTTP_API["internal/api/http/*.go<br/>路由与控制器"]
end
subgraph "应用层"
APP_WS["internal/application/workset.go<br/>应用服务"]
APP_ASSEMBLER["internal/application/assembler/*.go<br/>装配器"]
APP_ADAPTER["internal/application/adapter/*.go<br/>适配器"]
end
subgraph "领域层"
DOMAIN_MODEL["internal/domain/model/*.go<br/>领域模型/值对象"]
DOMAIN_REPO_IF["internal/domain/repository/types.go<br/>仓储接口类型"]
DOMAIN_SERVICE["internal/domain/service/*.go<br/>领域服务"]
DOMAIN_EXT_OSS["internal/domain/external/oss.go<br/>外部接口"]
end
subgraph "基础设施层"
INFRA_REPO["internal/infrastructure/repository/*.go<br/>仓储实现"]
INFRA_EXT_OSS["internal/infrastructure/external/oss_factory.go<br/>外部客户端工厂"]
end
MAIN --> HTTP_API
HTTP_API --> APP_WS
APP_WS --> APP_ASSEMBLER
APP_WS --> APP_ADAPTER
APP_WS --> DOMAIN_MODEL
APP_WS --> DOMAIN_SERVICE
APP_WS --> DOMAIN_EXT_OSS
APP_WS --> DOMAIN_REPO_IF
APP_WS --> INFRA_REPO
APP_WS --> INFRA_EXT_OSS
```

**图表来源**
- [main.go:25-145](file://main.go#L25-L145)
- [http.go:16-151](file://internal/api/http/http.go#L16-L151)
- [workset.go:1-310](file://internal/application/workset.go#L1-L310)
- [workset.go:1-139](file://internal/infrastructure/repository/workset.go#L1-L139)
- [workset.go:1-82](file://internal/domain/model/workset.go#L1-L82)
- [types.go:1-12](file://internal/domain/repository/types.go#L1-L12)
- [loader.go:1-71](file://internal/application/adapter/loader.go#L1-L71)
- [oss_factory.go:1-36](file://internal/infrastructure/external/oss_factory.go#L1-L36)
- [oss.go:1-9](file://internal/domain/external/oss.go#L1-L9)

**章节来源**
- [main.go:25-145](file://main.go#L25-L145)
- [http.go:16-151](file://internal/api/http/http.go#L16-L151)

## 核心组件
- 表示层（HTTP）：负责路由注册、请求解析、响应封装与鉴权中间件集成。
- 应用层（WorksetApplication）：编排业务流程、鉴权、事务控制、装配结果。
- 领域层（WorksetInfo/Creation/Update、权限模型）：承载业务不变量与规则。
- 基础设施层（仓储实现、OSS 客户端）：提供持久化与外部能力。

**章节来源**
- [workset.go:20-70](file://internal/application/workset.go#L20-L70)
- [workset.go:12-18](file://internal/infrastructure/repository/workset.go#L12-L18)
- [workset.go:5-82](file://internal/domain/model/workset.go#L5-L82)
- [permission.go:198-800](file://internal/domain/model/permission.go#L198-L800)
- [oss_factory.go:16-36](file://internal/infrastructure/external/oss_factory.go#L16-L36)

## 架构总览
下图展示了从 HTTP 请求到应用层再到仓储与外部服务的整体调用链路与依赖方向。

```mermaid
sequenceDiagram
participant C as "客户端"
participant H as "HTTP 控制器<br/>internal/api/http/workset.go"
participant A as "应用服务<br/>internal/application/workset.go"
participant AD as "适配器<br/>internal/application/adapter/loader.go"
participant S as "领域服务/权限<br/>internal/domain/model/permission.go"
participant R as "仓储实现<br/>internal/infrastructure/repository/workset.go"
participant E as "外部OSS客户端<br/>internal/infrastructure/external/oss_factory.go"
C->>H : "GET /api/v1/worksets"
H->>A : "ListWorksets(scope, userID, args)"
A->>AD : "HandleLoadMemberInfo(repo)"
AD->>R : "Get(...FilterByUserID/TeamID)"
R-->>AD : "MemberInfo"
AD-->>A : "OnLoadMemberInfo"
A->>S : "PermWorksetList().Check(userID, teamID, onLoadMemberInfo)"
S-->>A : "鉴权结果"
A->>R : "List(...IncludeTeamInfo, ...Paginate)"
R-->>A : "[]WorksetInfo"
A->>E : "GenerateGetPresignedURL(key)"
E-->>A : "预签名URL"
A-->>H : "[]WorksetInfo"
H-->>C : "200 OK"
```

**图表来源**
- [workset.go:25-52](file://internal/api/http/workset.go#L25-L52)
- [workset.go:72-126](file://internal/application/workset.go#L72-L126)
- [loader.go:16-24](file://internal/application/adapter/loader.go#L16-L24)
- [permission.go:212-219](file://internal/domain/model/permission.go#L212-L219)
- [workset.go:32-49](file://internal/infrastructure/repository/workset.go#L32-L49)
- [oss_factory.go:16-36](file://internal/infrastructure/external/oss_factory.go#L16-L36)

## 详细组件分析

### 应用层：WorksetApplication
- 职责边界
  - 参数校验与日志追踪
  - 鉴权：基于 PermWorkset* 权限模型
  - 事务管理：统一开启/提交/回滚
  - 结果装配：使用装配器将仓储模型转换为对外值对象
- 关键交互
  - 依赖 OSS 客户端生成预签名 URL
  - 依赖成员仓储加载成员信息用于鉴权
  - 依赖工作集仓储进行 CRUD 与统计
- 事务处理位置
  - 在创建工作集时，应用层显式开启事务并负责提交或回滚

```mermaid
classDiagram
class WorksetApplication {
+ListWorksets(scope, userID, args) []WorksetInfo
+CreateWorkset(scope, userID, args) CreateWorksetResult
+UpdateWorkset(scope, userID, args) error
+DeleteWorkset(scope, userID, id) error
}
class Adapter {
+HandleLoadMemberInfo(memberRepo) OnLoadMemberInfo
}
class Permission {
+PermWorksetList() Check(...)
+PermWorksetCreate() Check(...)
+PermWorksetUpdate() Check(...)
+PermWorksetDelete() Check(...)
}
class OSSClient {
+GenerateGetPresignedURL(key) string
}
class WorksetRepository {
+BeginTransaction() Executor
+List(executor, options) []WorksetInfo
+Create(executor, creation) string
+Update(executor, update) error
+Delete(executor, id) error
}
WorksetApplication --> Adapter : "加载成员信息"
WorksetApplication --> Permission : "权限校验"
WorksetApplication --> OSSClient : "生成预签名URL"
WorksetApplication --> WorksetRepository : "CRUD/统计"
```

**图表来源**
- [workset.go:20-70](file://internal/application/workset.go#L20-L70)
- [loader.go:16-24](file://internal/application/adapter/loader.go#L16-L24)
- [permission.go:212-219](file://internal/domain/model/permission.go#L212-L219)
- [workset.go:12-18](file://internal/infrastructure/repository/workset.go#L12-L18)
- [oss.go:3-8](file://internal/domain/external/oss.go#L3-L8)

**章节来源**
- [workset.go:72-126](file://internal/application/workset.go#L72-L126)
- [workset.go:128-213](file://internal/application/workset.go#L128-L213)
- [workset.go:215-262](file://internal/application/workset.go#L215-L262)
- [workset.go:264-309](file://internal/application/workset.go#L264-L309)

### 领域层：Workset 领域模型与权限
- 领域模型
  - WorksetInfo：对外暴露的只读信息载体
  - WorksetCreation/Update：创建/更新的输入载体
- 权限模型
  - 通过 PermWorksetList/Create/Update/Delete 的 Check 方法实现细粒度权限控制
  - Check 方法接收 OnLoadUserInfo/OnLoadMemberInfo 等加载器函数，避免领域模型直接依赖仓储

```mermaid
classDiagram
class WorksetInfo {
+ID : string
+TeamID : string
+Team : *TeamInfo
+Index : int
+Name : string
+Description : string
+ComicCount : int
+CreatedAt : time
+UpdatedAt : time
}
class WorksetCreation {
+TeamID : string
+Index : int
+Name : string
+Description : string
}
class WorksetUpdate {
+ID : string
+Name : string
+Description : *string
}
class Permission {
+PermWorksetList() Check(userID, teamID, onLoadMemberInfo) bool
+PermWorksetCreate() Check(userID, teamID, onLoadMemberInfo) bool
+PermWorksetUpdate() Check(userID, teamID, onLoadMemberInfo) bool
+PermWorksetDelete() Check(userID, teamID, onLoadMemberInfo) bool
}
WorksetInfo --> TeamInfo : "可选关联"
WorksetCreation --> WorksetInfo : "创建后映射"
WorksetUpdate --> WorksetInfo : "更新后映射"
Permission --> WorksetInfo : "基于团队/成员权限"
```

**图表来源**
- [workset.go:5-82](file://internal/domain/model/workset.go#L5-L82)
- [permission.go:797-800](file://internal/domain/model/permission.go#L797-L800)

**章节来源**
- [workset.go:20-82](file://internal/domain/model/workset.go#L20-L82)
- [permission.go:797-800](file://internal/domain/model/permission.go#L797-L800)

### 基础设施层：仓储与外部服务
- 仓储实现
  - WorksetRepository：封装 GORM 执行器，提供 List/Get/LockByTeamID/Count/Create/Update/Delete 等方法
  - Executor 接口抽象，Transactor 接口定义 BeginTransaction
- 外部服务
  - OSS 客户端工厂：根据环境变量选择 R2 或阿里云 OSS 实现
  - OSSClient 接口：提供预签名 URL 生成与批量删除能力

```mermaid
classDiagram
class Executor {
<<interface>>
}
class Transactor {
<<interface>>
+BeginTransaction() Executor
}
class WorksetRepository {
-executor : Executor
+BeginTransaction() Executor
+List(executor, options) []WorksetInfo
+Get(executor, options) WorksetInfo
+LockByTeamID(executor, teamID) error
+Count(executor, options) int64
+Create(executor, creation) string
+Update(executor, update) error
+Delete(executor, id) error
}
class OSSClient {
<<interface>>
+GeneratePutPresignedURL(key) string
+GenerateGetPresignedURL(key) string
+Delete(key) error
+DeleteBatch(keys) error
}
class OSSFactory {
+NewOSSClient() OSSClient
}
WorksetRepository --> Executor : "使用"
WorksetRepository --> Transactor : "开启事务"
OSSFactory --> OSSClient : "创建实例"
```

**图表来源**
- [workset.go:12-18](file://internal/infrastructure/repository/workset.go#L12-L18)
- [types.go:5-12](file://internal/domain/repository/types.go#L5-L12)
- [oss_factory.go:16-36](file://internal/infrastructure/external/oss_factory.go#L16-L36)
- [oss.go:3-8](file://internal/domain/external/oss.go#L3-L8)

**章节来源**
- [workset.go:28-30](file://internal/infrastructure/repository/workset.go#L28-L30)
- [workset.go:32-49](file://internal/infrastructure/repository/workset.go#L32-L49)
- [types.go:9-12](file://internal/domain/repository/types.go#L9-L12)
- [oss_factory.go:16-36](file://internal/infrastructure/external/oss_factory.go#L16-L36)

### HTTP 层：路由与控制器
- 路由注册集中在 HTTP 初始化函数中，按模块划分（认证、用户、团队、成员、邀请、漫画、工作集、章节、页面、分配、单元）
- 控制器负责参数解析、鉴权提取、调用应用层并封装响应

```mermaid
flowchart TD
Start(["请求进入"]) --> Parse["解析参数/鉴权"]
Parse --> Route{"路由匹配"}
Route --> |工作集列表| List["ListWorksets"]
Route --> |创建工作集| Create["CreateWorkset"]
Route --> |更新工作集| Update["UpdateWorkset"]
Route --> |删除工作集| Delete["DeleteWorkset"]
List --> CallApp["调用 WorksetApplication"]
Create --> CallApp
Update --> CallApp
Delete --> CallApp
CallApp --> Resp["封装响应/状态码"]
Resp --> End(["返回"])
```

**图表来源**
- [http.go:16-151](file://internal/api/http/http.go#L16-L151)
- [workset.go:25-52](file://internal/api/http/workset.go#L25-L52)
- [workset.go:67-95](file://internal/api/http/workset.go#L67-L95)
- [workset.go:111-148](file://internal/api/http/workset.go#L111-L148)
- [workset.go:162-189](file://internal/api/http/workset.go#L162-L189)

**章节来源**
- [http.go:16-151](file://internal/api/http/http.go#L16-L151)
- [workset.go:25-52](file://internal/api/http/workset.go#L25-L52)
- [workset.go:67-95](file://internal/api/http/workset.go#L67-L95)
- [workset.go:111-148](file://internal/api/http/workset.go#L111-L148)
- [workset.go:162-189](file://internal/api/http/workset.go#L162-L189)

## 依赖分析
- 层间依赖方向
  - HTTP 层依赖应用层
  - 应用层依赖领域层（模型/服务/外部接口）与基础设施层（仓储/OSS）
  - 领域层不依赖基础设施层；通过接口与适配器解耦
- 关键依赖点
  - 应用层通过适配器注入加载器函数，完成权限校验所需的成员信息加载
  - 事务在应用层统一管理，仓储实现仅负责执行 SQL

```mermaid
graph LR
HTTP["HTTP 层"] --> APP["应用层"]
APP --> DOMAIN["领域层"]
APP --> INFRA["基础设施层"]
DOMAIN -.-> INFRA
```

**图表来源**
- [main.go:56-122](file://main.go#L56-L122)
- [workset.go:43-47](file://internal/application/workset.go#L43-L47)
- [loader.go:16-24](file://internal/application/adapter/loader.go#L16-L24)

**章节来源**
- [main.go:56-122](file://main.go#L56-L122)
- [workset.go:43-47](file://internal/application/workset.go#L43-L47)

## 性能考虑
- 事务范围最小化：仅在创建工作集时开启事务，减少锁竞争与事务持有时间
- 查询优化：通过 IncludeTeamInfo 与分页选项减少不必要的数据传输
- 外部服务：OSS 预签名 URL 降低上传/下载时的带宽与延迟
- 日志与追踪：应用层使用 TraceScope 记录关键步骤，便于定位性能瓶颈

## 故障排查指南
- 参数校验失败
  - 现象：返回参数错误
  - 排查：确认 List/Create/Update 参数结构与必填字段
  - 参考路径：[workset.go:79-82](file://internal/application/workset.go#L79-L82), [workset.go:135-138](file://internal/application/workset.go#L135-L138), [workset.go:222-225](file://internal/application/workset.go#L222-L225)
- 权限不足
  - 现象：返回权限不足
  - 排查：确认当前用户在目标团队的角色与分工
  - 参考路径：[workset.go:98-100](file://internal/application/workset.go#L98-L100), [workset.go:154-156](file://internal/application/workset.go#L154-L156), [workset.go:250-252](file://internal/application/workset.go#L250-L252), [workset.go:300-301](file://internal/application/workset.go#L300-L301)
- 事务异常
  - 现象：创建失败、回滚日志
  - 排查：检查 BeginTransaction 返回与 Commit/Rollback 调用
  - 参考路径：[workset.go:158-172](file://internal/application/workset.go#L158-L172), [workset.go:207-210](file://internal/application/workset.go#L207-L210)
- 外部服务异常
  - 现象：OSS 预签名 URL 生成失败
  - 排查：确认 OSS_PROVIDER 环境变量与对应 SDK 配置
  - 参考路径：[oss_factory.go:16-36](file://internal/infrastructure/external/oss_factory.go#L16-L36)

**章节来源**
- [workset.go:79-82](file://internal/application/workset.go#L79-L82)
- [workset.go:135-138](file://internal/application/workset.go#L135-L138)
- [workset.go:222-225](file://internal/application/workset.go#L222-L225)
- [workset.go:98-100](file://internal/application/workset.go#L98-L100)
- [workset.go:154-156](file://internal/application/workset.go#L154-L156)
- [workset.go:250-252](file://internal/application/workset.go#L250-L252)
- [workset.go:300-301](file://internal/application/workset.go#L300-L301)
- [workset.go:158-172](file://internal/application/workset.go#L158-L172)
- [workset.go:207-210](file://internal/application/workset.go#L207-L210)
- [oss_factory.go:16-36](file://internal/infrastructure/external/oss_factory.go#L16-L36)

## 结论
本项目在 DDD 三层架构基础上，实现了“领域服务不依赖仓储层”与“事务在应用层”的创新实践。通过适配器注入加载器函数，权限校验逻辑完全位于应用层，既保持了领域模型的纯净，又提升了系统的可测试性与可维护性。应用层统一管理事务，确保业务操作的一致性；仓储实现专注于数据访问，外部服务通过工厂模式灵活切换。整体架构清晰、职责明确、扩展性强，适合复杂业务场景下的长期演进。

## 附录
- 值对象与参数校验
  - ListWorksetArgs/ CreateWorksetArgs/ UpdateWorksetArgs 的校验逻辑集中于应用层，保证输入合法性
  - 参考路径：[workset.go:29-43](file://internal/value/workset.go#L29-L43), [workset.go:51-65](file://internal/value/workset.go#L51-L65), [workset.go:78-88](file://internal/value/workset.go#L78-L88)
- 关联展开规格解析
  - ResolveWorksetListIncludeSpec 将 includes[] 解析为 NeedTeam 等规格，指导仓储查询
  - 参考路径：[workset_include.go:9-20](file://internal/domain/service/workset_include.go#L9-L20)

**章节来源**
- [workset.go:29-43](file://internal/value/workset.go#L29-L43)
- [workset.go:51-65](file://internal/value/workset.go#L51-L65)
- [workset.go:78-88](file://internal/value/workset.go#L78-L88)
- [workset_include.go:9-20](file://internal/domain/service/workset_include.go#L9-L20)