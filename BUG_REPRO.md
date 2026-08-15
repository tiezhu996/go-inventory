# BUG_REPRO

## Bug 是什么
SortReservations 排序方向反了、BuildBatches 分页边界 off-by-one 每页漏最后一条、OrderIDs 返回内部切片、worker 每页丢第一条，导致分页顺序错、漏数据、列表串改。

## 如何触发
`go test ./...`

## 错误信息
- TestSortReservations / TestListPagesOrder / TestBuildBatchesFresh 失败。
- TestReservationRecordOrderFresh / TestRunSummary 失败。
