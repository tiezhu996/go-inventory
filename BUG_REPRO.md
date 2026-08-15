# BUG_REPRO

## Bug 是什么
BuildBatches / OrderIDs / ListPages 返回共享底层数组子切片，worker 又用 `page[:len(page)-1]` 切掉最后一条，导致记录漏掉、列表串改。

## 如何触发
`go test ./...`

## 错误信息
- TestBuildBatchesFresh / TestReservationRecordOrderFresh 失败。
- TestRunSummary 失败（Reserved 计数不对）。
