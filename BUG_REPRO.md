# BUG_REPRO

## Bug 是什么
MergeSummary 丢掉 Released、store 重复记录静默返回 nil、service 吞掉记录错误、worker 派发失败不再累加 Failed，导致重复无报错、失败数漏统计。

## 如何触发
`go test ./...`

## 错误信息
- TestMergeSummary 失败（Released 丢失）。
- TestReservationRecordOrderFresh 失败（重复未报错）。
- TestRunFailed 失败（Failed=0 want 1）。
