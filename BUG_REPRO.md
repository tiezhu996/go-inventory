# BUG_REPRO

## Bug 是什么
ValidQty 判断反向、Reserve/Release 去掉锁、service 去掉校验、worker 汇总去掉锁且 wg.Add 放进 goroutine，导致库存算错、负数可预留、统计错误并有数据竞争。

## 如何触发
`go test -race ./...`

## 错误信息
- TestValidQty / TestReserveFlow 失败（负数被接受、库存错）。
- `go test -race` 报 data race。
