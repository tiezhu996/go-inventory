# BUG_REPRO

## Bug 是什么
New 未初始化 resvs map、ValidQty 判断反向、service 吞掉 RecordReservation 错误、worker 不检查 ctx，导致记录预约时 nil map panic 且负数可过校验。

## 如何触发
`go test ./...`

## 错误信息
- `panic: assignment to entry in nil map`（TestReservationRecordOrderFresh）。
- TestValidQty 失败。
