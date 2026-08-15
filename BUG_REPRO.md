# BUG_REPRO

## Bug 是什么
worker 生产者/消费者不检查 ctx、service 用 Background 丢 ctx、store 去掉 ctx 检查，且 SortReservations 排序方向反了，导致取消后仍继续处理、排序错乱。

## 如何触发
`go test ./...`

## 错误信息
- TestRunCancel 失败（Reserved != 0）。
- TestSortReservations 失败（顺序反）。
