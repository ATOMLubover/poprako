# RFC 001：引入 OSSManagerSvc 与本地消息表

## 背景

当前系统对 OSS 采用 storage aside 方式：主数据在本地数据库，文件直接由客户端与第三方 OSS 交互，主服务器只负责签发预签名 URL、记录本地状态、以及在部分删除链路中直接调用 OSS delete。

这导致一个核心问题：主服务器无法可靠控制远程对象的最终状态。

典型失败场景：

1. 客户端拿到预签名 URL 后上传成功，但没有成功调用 confirm，本地仍显示未上传，远程却已经存在对象，形成资源泄漏。
2. chapter/page/comic/user/team 删除时，主服务器同步调用 OSS delete；若第三方 OSS 短暂失败，则本地删除失败或行为不一致，且删除链路本身高度不可靠。
3. 现有 page OSS key 为 `page_{index}` 这类全局扁平命名，不同 chapter 之间会发生对象 key 冲突。

本 RFC 的目标，是将所有 OSS 资源的创建确认与远程清理统一纳入一个本地可靠机制中。

## 目标

1. 所有 OSS 资源统一走一套可靠架构，不再按业务类型零散实现。
2. 所有远程 OSS delete 都只能由后台消费者执行，主服务器请求路径不再直接删除 OSS 对象。
3. 业务接口的本地状态变更必须同步提交，用户可见结果只依赖本地数据库，而不依赖远程 OSS 即时成功。
4. 上传未 confirm 导致的远程残留对象，必须可被后台任务最终清理。
5. 删除 chapter/page/comic/user/team 等业务实体后，其关联 OSS 对象必须可被后台任务最终清理。
6. 所有 OSS key 使用目录式命名，避免跨实体冲突。

## 非目标

1. 本阶段不引入对象版本号，也不为每次上传生成全新随机 key。
2. 本阶段不试图完全消除“旧消息误删新对象”的风险。该风险被接受，优先级低于防泄漏。
3. 本阶段不引入外部 MQ；仅使用本地数据库消息表与本地消费者。

## 总体方案

新增一个统一的 `OSSManagerSvc`，并引入一张私有本地消息表 `oss_message_table`。

`OSSManagerSvc` 负责：

1. 为所有 OSS 资源生成规范 object key。
2. 在 reserve/confirm/delete 等业务动作中，协调业务表与本地消息表的事务写入。
3. 为上传链路生成预签名 PUT URL。
4. 由后台 worker 消费本地消息表，并调用 OSS delete / delete batch 做最终清理。

设计原则：

1. 主请求链路中的“业务结果”只由本地数据库决定。
2. 远程 OSS 的最终一致性由后台任务兜底，不阻塞业务请求。
3. 消费者只能依赖消息表 payload，不能依赖原业务记录仍然存在。

## OSS 资源范围

以下资源全部纳入统一管理：

1. 用户头像
2. 汉化组头像
3. 漫画封面
4. 章节页面图片
5. 后续新增 OSS 资源

## OSS Key 命名规范

所有 object key 必须使用目录式命名，禁止继续使用扁平全局 key。

规范示例：

1. `user_{user_id}/avatar`
2. `team_{team_id}/avatar`
3. `comic_{comic_id}/cover`
4. `chapter_{chapter_id}/page_{page_index}`

约束：

1. key 必须包含所属实体 ID，避免跨实体冲突。
2. key 必须表达资源槽位，而不是仅表达文件类型。
3. page key 必须至少包含 `chapter_id + page_index`，不能再使用 `page_{index}`。
4. 业务代码生成 key 时，必须通过统一的 key builder，不允许手写字符串拼接。

备注：

1. 本阶段默认 key 为稳定槽位 key，不引入版本号。
2. 若未来需要彻底规避旧消息误删新对象，可在后续 RFC 中演进为版本化 key。

## 本地消息表

统一使用单表：`oss_message_table`。

建议字段：

1. `id`
2. `resource_type`
3. `resource_id`
4. `operation`
5. `status`
6. `object_key`
7. `payload_json`
8. `visible_at`
9. `expire_at`
10. `processing_started_at`
11. `attempt_count`
12. `last_error`
13. `created_at`
14. `updated_at`

说明：

1. `resource_type` 用于区分 `user_avatar`、`team_avatar`、`comic_cover`、`page_image` 等。
2. `resource_id` 表示业务资源 ID；对于 chapter 删除这类多对象清理，仍可使用 chapter ID 作为资源 ID。
3. `operation` 至少包括 `create_pending` 与 `delete_pending` 两类语义。
4. `object_key` 用于单对象场景；`payload_json` 用于多对象快照场景，例如 chapter 删除时记录所有 page object keys。
5. `visible_at` 表示消息何时可被 worker 消费，支持延迟重试。
6. `expire_at` 主要用于上传预留超时清理。

状态机：

1. `pending`
2. `processing`
3. `completed`

状态流转：

1. `pending -> processing -> completed`
2. `processing -> pending`，表示本次消费失败，回退等待下一次重试

本阶段不强制引入终态 `dead`，默认无限重试，但必须记录失败次数与最近错误。

## 消息唯一性与快照原则

### 上传类消息

上传消息语义上是“本地已预留，但尚未 confirm”。

要求：

1. 对同一逻辑资源，只允许存在一条活跃的上传待确认消息。
2. reserve 相同逻辑资源时，采用 upsert 语义覆盖旧的待确认消息，并重置状态为 `pending`。
3. confirm 时，将对应上传消息标记为 `completed`。

### 删除类消息

删除消息语义上是“本地已经不可见或已删除，但远程对象待清理”。

要求：

1. 删除消息必须携带完整快照，不能依赖原业务记录后续仍可查询。
2. 单对象删除可直接存 `object_key`。
3. 多对象删除必须把对象 key 列表写入 `payload_json`。

## 业务一致性原则

所有对用户可见的业务结果，以本地表状态为准。

具体要求：

1. 删除业务请求成功后，相关记录必须立即在本地不可见。
2. 资源一旦被本地标记为未上传、待删除或已删除，读路径不得继续暴露其 GET URL。
3. 远程 OSS 是否已经实际删除，不影响业务接口返回。

这意味着：

1. `page.ImageURL` 只能在 `page.is_uploaded = true` 时生成。
2. `comic.CoverURL` 只能在 `comic.is_cover_uploaded = true` 时生成。
3. `user.AvatarURL`、`team.AvatarURL` 同理。
4. chapter/comic/user/team 等记录若本地已被删除，则直接不可见，不等待 OSS delete 成功。

## 核心业务流程

### 1. 上传预留 reserve

适用于：用户头像、汉化组头像、漫画封面、页面图片。

流程：

1. 业务层通过 `OSSManagerSvc` 生成规范 object key。
2. `OSSManagerSvc` 生成预签名 PUT URL。
3. 开启事务。
4. 同步写入业务表，记录 object key，并将本地上传状态设为未完成。
5. 向 `oss_message_table` upsert 一条上传待确认消息，状态为 `pending`。
6. 提交事务。
7. 返回 PUT URL 给客户端。

约束：

1. reserve 成功后，本地必须有明确记录表明“该资源处于待确认上传状态”。
2. reserve 失败则不能留下孤立本地消息。

### 2. 上传确认 confirm

流程：

1. 开启事务。
2. 同步更新业务表，将本地上传状态设为已上传。
3. 将对应上传待确认消息标记为 `completed`。
4. 提交事务。

说明：

1. confirm 不调用远程 OSS。
2. 对 page 而言，当前外部接口若仍沿用 `UpdatePage(is_uploaded=true)`，内部也必须视为 confirm 语义并纳入 `OSSManagerSvc`，不能绕开消息表。

### 3. 上传超时清理

适用于：reserve 成功，但客户端未 confirm 的情况。

流程：

1. worker 周期性扫描所有 `create_pending` 且 `status = pending` 的消息。
2. 仅当 `expire_at <= now` 时才允许消费。
3. worker 使用 CAS 将消息从 `pending` 改为 `processing`。
4. 调用 OSS delete 删除对应 object key。
5. 删除成功后，将消息标记为 `completed`。
6. 删除失败则回写错误信息，增加 `attempt_count`，并将状态重置为 `pending`，同时更新 `visible_at` 以便下次重试。

### 4. 本地删除请求

适用于：page 删除、chapter 删除、comic 删除、user 删除、team 删除，以及未来所有 OSS 资源删除动作。

流程：

1. 开启事务。
2. 同步更新本地业务数据，使资源立即不可见。
3. 若需要物理删除业务记录，则在同一事务内完成。
4. 将消费所需的 object key 快照写入 `oss_message_table`，插入一条 `delete_pending` 消息，状态为 `pending`。
5. 提交事务。
6. 请求立即返回成功，不等待远程 OSS delete。

约束：

1. 主请求链路中禁止直接调用 OSS delete。
2. 所有同步统计字段必须在本地事务内完成更新，例如 chapter page_count、comic chapter_count 等。
3. 如果原记录会被立即删除，则消息表中必须提前保存完整清理信息。

### 5. 删除后台清理

流程：

1. worker 周期性扫描所有 `delete_pending` 且 `status = pending` 且 `visible_at <= now` 的消息。
2. 通过 CAS 抢占消息，将其置为 `processing`。
3. 根据 `object_key` 或 `payload_json.object_keys` 执行单删或批量删除。
4. 删除成功则标记为 `completed`。
5. 删除失败则记录错误并回退为 `pending`，等待重试。

### 6. chapter 删除特别说明

chapter 删除会级联影响 page。

要求：

1. 在事务开始前或事务内部写消息前，必须先收集该 chapter 的所有 page object keys。
2. 删除 chapter 时，本地 chapter/page/assignment 等业务数据仍按现有语义同步删除或级联删除。
3. worker 后续只依赖消息表中的快照执行远程清理，不再回查 page 表。

## Worker 设计

扫描间隔建议：`5 分钟`。

理由：

1. 上传超时清理不应等到 `1 小时` 级别才开始发现可消费消息。
2. 删除成功后的远程清理也不应过度延迟。

建议策略：

1. 统一一个 worker 循环，每 `5 分钟` 扫描一次。
2. 上传超时消息使用较长 `expire_at`，建议 `24 小时`。
3. 删除消息默认 `visible_at = now`，尽快清理。
4. 重试使用退避策略，例如固定 `5 分钟` 或逐步增长。
5. 每次按分页拉取有限批次，避免单次扫描过大。

## CAS 并发控制

消费者必须使用 CAS 抢占，不能使用“先 select 后处理”的无锁方式。

要求：

1. 抢占条件必须包含 `status = pending`。
2. 更新成功后才算获得处理权。
3. 未抢占成功的 worker 不得继续处理该消息。
4. `processing_started_at` 用于排查卡死消息。

这保证多实例部署时同一消息不会被重复消费为主要路径。

## OSS 接口约束

`oss.Client` 或其拆分后的删除能力接口，必须满足以下契约：

1. 删除不存在的对象必须视为成功。
2. 批量删除中，不存在对象不应导致整体失败。
3. provider 差异必须在 infra 层被屏蔽，不能泄露给 `OSSManagerSvc`。

## 架构与 DI 重构

这是一次跨多个 app 的统一重构。

### 新增组件

1. `repo.OSSMessageRepo`
2. `repo_infra.OSSMessageRepo` 实现
3. `app.OSSManagerSvc` 或独立 app service 实现
4. `app/worker` 或等价位置的本地消息消费者

### 依赖关系调整

业务 app 不应再直接负责 OSS 生命周期控制。

建议调整为：

1. `UserApp`、`TeamApp`、`ComicApp`、`PageApp`、`ChapterApp` 通过 `OSSManagerSvc` 完成 reserve/confirm/delete enqueue。
2. 业务 app 不再直接调用 `ossClient.Delete/DeleteBatch`。
3. 读路径若仍需生成 GET URL，可继续依赖只读签名能力。
4. 为降低职责耦合，建议把当前 `oss.Client` 拆为两组能力：
   1. URL 生成能力，用于 PUT/GET 预签名 URL
   2. 对象删除能力，仅供 `OSSManagerSvc` 的 worker 使用

### main.go 调整

启动阶段需要新增：

1. 初始化 `OSSMessageRepo`
2. 初始化 `OSSManagerSvc`
3. 初始化后台 worker
4. 在应用生命周期内启动并优雅关闭 worker

### 领域服务调整

各资源的 key 生成逻辑仍可由对应 domain service 提供，但必须统一收口到 `OSSManagerSvc` 调用。

例如：

1. UserService 生成 `user_{id}/avatar`
2. TeamService 生成 `team_{id}/avatar`
3. ComicService 生成 `comic_{id}/cover`
4. PageService 生成 `chapter_{chapter_id}/page_{page_index}`

其中 PageService 需要修改签名，不能再只靠 `index` 生成 key。

## 对现有代码的直接影响

### 1. 删除链路

以下同步删除 OSS 的实现都要移除：

1. `page.Remove`
2. `chapter.Remove`
3. `comic.Remove`
4. `user.Remove`
5. `team.Remove`

这些接口改为：同步更新本地业务表 + 写入 delete message。

### 2. 上传链路

以下 reserve/confirm 逻辑要统一收口：

1. `user reserve/confirm avatar`
2. `team reserve/confirm avatar`
3. `comic reserve/confirm cover`
4. `page reserve/confirm image`

### 3. 读路径

所有 GET URL 生成逻辑都必须以本地 uploaded 标志为前置条件。

### 4. key 生成逻辑

1. 所有旧的扁平命名生成函数需要改写。
2. Page key 生成签名必须纳入 `chapterID`。

## 数据迁移与兼容

1. 新增 `oss_message_table` 迁移。
2. 现有业务表无需为本 RFC 新增大量字段，优先复用已有 `is_uploaded`、`is_cover_uploaded`、`is_avatar_uploaded` 等字段。
3. 已存在的旧格式 object key 不强制立即迁移；老数据在下一次 reserve、delete 或被清理时自然过渡。
4. 新写入数据必须全部使用新目录式 key 规范。

## 可接受风险

本阶段接受如下风险：

1. 因为 object key 是稳定槽位 key，旧的超时清理消息在极端竞态下可能删除同一逻辑资源的新对象。
2. 该风险暂不通过版本号解决，后续若实际发生频率不可接受，再引入 versioned key 方案。

接受原因：

1. 当前首要目标是防止 OSS 资源泄漏。
2. 保持 key 简单、稳定、可推导，有利于先完成架构统一。

## 监控与可观测性

必须补充以下能力：

1. pending 消息总数
2. processing 消息总数
3. completed 消息吞吐
4. attempt_count 分布
5. 最老 pending 消息年龄
6. 每类资源的失败次数
7. 结构化日志，至少包含 message_id、resource_type、resource_id、operation、attempt_count

## 测试要求

至少覆盖：

1. reserve 成功时业务表与消息表同时写入
2. confirm 成功时业务表与消息表同时完成
3. 删除请求成功时，本地立即不可见，且消息表持久化成功
4. worker 成功消费 create timeout 消息
5. worker 成功消费 delete 消息
6. worker 并发抢占同一消息时，仅一方成功
7. provider 返回 not found 时删除视为成功
8. chapter 删除后消费者只依赖消息快照完成 page 对象清理
9. page key 生成为 `chapter_{chapter_id}/page_{page_index}`

## 落地顺序建议

1. 新增 `oss_message_table`、repo、worker、`OSSManagerSvc`
2. 重构 key 生成逻辑，先修正 page key 冲突问题
3. 重构上传 reserve/confirm 链路
4. 重构 page/chapter/comic/user/team 删除链路
5. 收敛读路径的 uploaded 判断
6. 加入指标、日志与测试

## 结论

本 RFC 之后，系统对 OSS 的控制方式将变为：

1. 主请求链路只负责本地状态与消息持久化。
2. 远程删除完全由后台 worker 异步执行。
3. 所有 OSS 资源统一由 `OSSManagerSvc` 管理。
4. 所有 object key 统一改为目录式命名。

这会带来一次较大的应用层与 DI 重构，但这是把 OSS 从“散落的第三方副作用”收敛为“受控本地一致性机制”的必要步骤。
