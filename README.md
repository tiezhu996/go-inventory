# inventory

一个用 Go 写的内存库存/预留服务，演示分层、并发预留、批量派发与上下文取消。

## 功能
- 创建商品、预留/释放库存、查询库存
- 预留记录分页查询
- 并发派发 worker 池，支持 context 取消

## 目录结构
```
cmd/inventory/      程序入口
internal/config/    环境配置
internal/model/     模型与纯工具函数
internal/store/     内存存储（商品 + 预留记录 + 锁）
internal/service/   业务逻辑
internal/worker/    派发 worker 池
```

## 运行与测试
```bash
go build ./...
go test ./...
go run ./cmd/inventory
```

## 环境变量
| 变量 | 说明 | 默认值 |
|------|------|--------|
| `INV_WORKERS` | worker 数量 | `2` |
| `INV_BATCH_SIZE` | 分页大小 | `2` |
