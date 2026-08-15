# BUG_REPRO

## Bug 是什么
service 包装错误用 %v 丢错误链，store 返回普通错误而非哨兵，MergeSummary 丢 Failed，worker 派发失败不再累加 Failed，导致 errors.Is 失效且失败数漏统计。

## 如何触发
`go test ./...`

## 错误信息
- TestWrapping / TestMergeSummary 失败。
- TestRunFailed 失败（Failed=0 want 1）。
