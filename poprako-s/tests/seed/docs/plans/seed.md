# 一、先把你的“真实规则”抽象清楚（不然一定写错）

你系统的核心约束其实是：

---

## 1️⃣ Assignment 是权限激活器（不是附属物）

```text
没有 assignment
→ 用户对 chapter 没有操作权限
→ 包括：units / workflow 子阶段
```

---

## 2️⃣ Workflow 是硬 gating（不是记录）

```text
upload_complete 之前 → page 不允许写 unit
translate_start 之前 → translator 不允许写
proofread_start 之前 → proofreader 不允许写
```

---

## 3️⃣ Unit ID 完全由服务端控制

```text
客户端 insert → server 生成 id
客户端 patch → 必须用 server 返回的 id
```

👉 所以：

> ❗任何“写完马上 patch”的 seed 都是错的

---

## 4️⃣ 权限 = Team ∩ Chapter ∩ Role ∩ Workflow

不是简单 RBAC，而是：

```text
can_write_unit =
  in_team
  AND has_assignment(role)
  AND workflow_phase_active
```

---

# 二、正确的 Seed 链路（完全重写）

我按**真实执行顺序**给你，不跳步骤。

---

# Phase 0：SuperAdmin → Team bootstrap（正确）

（这部分你是对的，我不赘述逻辑，直接给代码）

```ts
// login
const admin = await api("POST", "/auth/login", undefined, SUPER_ADMIN);

// create team
const team = await api("POST", "/teams", admin.token, {...});

// inject self as admin
await api("POST", "/members", admin.token, {
  team_id: team.id,
  user_id: admin.user_id,
  roles: ROLE_ADMIN | ROLE_REVIEWER
});
```

---

# Phase 1：创建业务用户（必须先有 assignment 才能干活）

---

## 1.1 创建 invitation

```ts
const invTranslator = await api("POST", "/invitations", admin.token, {
  team_id: team.id,
  invitee_qq: "10002",
  roles: ROLE_TRANSLATOR,
});

const invProofreader = await api("POST", "/invitations", admin.token, {
  team_id: team.id,
  invitee_qq: "10003",
  roles: ROLE_PROOFREADER,
});
```

---

## 1.2 注册用户

```ts
const translator = await api("POST", "/auth/register", undefined, {
  qq: "10002",
  password: "test123",
  invitation_code: invTranslator.code,
});

const proofreader = await api("POST", "/auth/register", undefined, {
  qq: "10003",
  password: "test123",
  invitation_code: invProofreader.code,
});
```

---

# Phase 2：创建内容（此时没人能操作 unit）

```ts
workset → comic → chapter
pages（reserve + uploaded=true）
```

⚠️ 注意：

> 到这里为止，没有任何人可以写 unit

---

# Phase 3：创建 Assignment（**必须在 unit 之前**）

---

## 3.1 translator assignment

```ts
await api("POST", "/assignments", admin.token, {
  chapter_id: chapter.id,
  user_id: translator.user_id,
  roles: ROLE_TRANSLATOR,
});
```

---

## 3.2 proofreader assignment

```ts
await api("POST", "/assignments", admin.token, {
  chapter_id: chapter.id,
  user_id: proofreader.user_id,
  roles: ROLE_PROOFREADER,
});
```

---

# Phase 4：Workflow 推进（关键 gating）

---

## 4.1 上传完成（允许进入翻译）

```ts
await api("PATCH", `/chapters/${chapter.id}`, admin.token, {
  workflow_transition: "upload_complete",
});
```

---

## 4.2 开始翻译（解锁 translator 权限）

```ts
await api("PATCH", `/chapters/${chapter.id}`, admin.token, {
  workflow_transition: "translate_start",
});
```

---

# Phase 5：Translator 写 Unit（第一次出现 unit）

---

## 5.1 insert（⚠️ 不能假设 id）

```ts
await api("PUT", "/units", translator.token, {
  page_id: pageId,
  unit_diff: {
    insert: [
      {
        index: 0,
        translated_text: "hello",
      },
    ],
  },
});
```

---

## 5.2 重新拉取 unit（必须）

```ts
const units = await api("GET", `/units?page_id=${pageId}`, translator.token);

const unitId = units[0].id;
```

---

👉 这是你刚才骂的点，完全正确：

> ❗客户端 id 不可信 → 必须 round-trip

---

# Phase 6：完成翻译 → 开启校对

---

## 6.1 translate_complete

```ts
await api("PATCH", `/chapters/${chapter.id}`, admin.token, {
  workflow_transition: "translate_complete",
});
```

---

## 6.2 proofread_start

```ts
await api("PATCH", `/chapters/${chapter.id}`, admin.token, {
  workflow_transition: "proofread_start",
});
```

---

# Phase 7：Proofreader 行为（必须重新 fetch）

---

## ❗重新获取 unit（不能复用 translator 的结果）

```ts
const unitsForProof = await api(
  "GET",
  `/units?page_id=${pageId}`,
  proofreader.token,
);

const unitId2 = unitsForProof[0].id;
```

---

## 7.2 patch

```ts
await api("PUT", "/units", proofreader.token, {
  page_id: pageId,
  unit_diff: {
    patch: [
      {
        id: unitId2,
        proofread_text: "HELLO",
        is_proofread: true,
      },
    ],
  },
});
```

---

# Phase 8：后续 workflow（略，但顺序必须对）

```ts
proofread_complete;
typeset_start;
typeset_complete;
review_complete;
publish_complete;
```

---

# 三、关键断言（你这个系统不测这些一定炸）

---

## ✅ 1. 未 assignment 写 unit → 必须失败

```ts
expect(403);
```

---

## ✅ 2. 未 translate_start 写 unit → 必须失败

---

## ✅ 3. proofreader 在 translate 阶段写 → 必须失败

---

## ✅ 4. 使用 fake unit id → 必须失败

---

## ✅ 5. 不重新 fetch 就 patch → 必须失败

---

# 四、你刚刚骂的点，我帮你总结成“系统设计原则”

你其实已经踩到了几个非常对的点：

---

## 1️⃣ Unit 是 server-owned state

不是 DTO，而是：

```text
CRDT-like authoritative resource
```

---

## 2️⃣ Assignment 是 capability token

不是关系表，而是：

```text
权限开关
```

---

## 3️⃣ Workflow 是 state machine guard

不是日志，而是：

```text
状态门禁
```

---

# 五、为什么你之前的 seed 全是错的

因为默认假设：

```text
“只要登录了就能操作资源”
```

但你这个系统是：

```text
必须满足 4 条件：
team + assignment + role + workflow
```

---

# 六、最终总结（这版是正确的）

这条链路：

✔ 先 assignment，再 unit
✔ workflow 控制写入窗口
✔ unit id round-trip
✔ proofreader 强制重新 fetch
✔ 权限严格按 chapter scope

---

# 如果你下一步要继续（建议）

我可以帮你做三件真正关键的东西：

---

## 1️⃣ 权限校验矩阵（你现在是隐式的）

直接生成：

```ts
can(user, action, resource);
```

---

## 2️⃣ workflow transition validator（防止非法跳转）

---

## 3️⃣ unit diff 冲突模型（你这个迟早要做）
