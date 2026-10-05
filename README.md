# vfx-budget

按帧预算排特效播放计划，并渲染成对齐文本。只用 Go 标准库。

```
go test ./...
go run ./cmd/plan --sample=budget
go run ./cmd/plan --sample=preempt
go run ./cmd/plan --sample=cooldown
go run ./cmd/plan --sample=order
go run ./cmd/plan --sample=work
```

## 口径

- **每帧预算**：同一帧内最多播放 `Budget` 个特效，超出后仍可抢占。
- **抢占**：预算已满时，新请求的 `Priority` 严格高于当前该帧优先级最低的那个才顶替它，被顶替的计入 `Preempted`；不高于则丢弃。
- **冷却**：同一 `Key` 距它上一次**被播放**的帧不足 `CooldownFrames` 时不播放。
- **顺序**：同一帧内输出按请求 ID 升序排列。
- **幂等**：同一个 `requestID` 只处理一次。

## 不变量

- 每帧输出的条数不超过 `Budget`。
- 被顶替的特效不会出现在最终计划里。
- `Work` 记判定期间检查过的历史记录数，不得随请求总数放大。

## 输出契约

对齐表格，第一行表头 `frame id key`，最后一行是 `preempted` 计数：

```
frame  id  key
1      a   fire
-      preempted  1
```

`--sample=work` 只打印 `{"played": N, "work": N}`。
