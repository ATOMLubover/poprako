# UnitSvc 测试说明

**源文件**: `refactor/internal/domain/svc/test/unit_test.go`  
**目标**: `UnitSvc.ApplyOps`

## 测试范围

这组测试覆盖 `CandOrder` 协议校验、临时 id 映射、以及多人并发修改后的最终结果。

## 关键规则

`CandOrder` 必须满足以下条件：

1. 包含本次提交中所有 `UnitCre` 的 `LocalId`
2. 包含本次提交中所有 `UnitSave` 的 `Id`
3. 不包含本次提交中任何 `UnitDel` 的 `Id`
4. 不能出现重复 id

如果不满足，`ApplyOps` 返回 `BadRequest`。

## 测试辅助

`mustRejectBadRequest` 用于验证协议错误时不会修改仓库内容。

`mustApplyOpsAccept` 用于验证合法提交可以正常执行。

`mustPageOrder` 用于读取页面内 unit 的最终顺序。

`assertOrderEqual` 用于比较最终顺序是否完全一致。

`assertStringPtrEqual` 用于比较字段是否被正确覆盖，包含 `nil`。

## 测试用例

### 1. `TestUnitApplyOps_RejectsCandOrderMissingSubmittedSaveId`

验证 `CandOrder` 漏掉本次提交的 `UnitSave` id 时直接拒绝。

### 2. `TestUnitApplyOps_RejectsCandOrderContainingDeletedId`

验证 `CandOrder` 仍然包含本次提交的 `UnitDel` id 时直接拒绝。

### 3. `TestUnitApplyOps_RejectsCandOrderContainingDuplicateIds`

验证 `CandOrder` 出现重复 id 时直接拒绝。

### 4. `TestUnitApplyOps_ResolvesCreateLocalIdToRealId`

验证 `UnitCre.LocalId` 会被映射成服务端生成的真实 id，最终结果里不会留下临时 id。

### 5. `TestUnitApplyOps_UsesLaterSubmitAsFinalTruthUnderConcurrentEdits`

验证两个客户端先后提交时，后提交者覆盖前提交者的修改结果。

这个用例覆盖以下行为：

1. A 先创建新 unit 并修改已有 unit
2. B 后删除 A 修改过的 unit
3. B 把某个字段从非空改回 `nil`
4. 最终顺序以后提交者的 `CandOrder` 为准
5. A 创建但 B 不知道的 unit 仍然保留，并按当前索引聚簇到合适位置
