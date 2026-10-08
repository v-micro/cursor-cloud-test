# cursor-cloud-test

最小 Go ping 服务，用来验证仓库能编译、测试，以及服务能响应健康检查。

## 运行

```bash
go run .
```

默认监听 `:8080`。设置 `PORT` 可改端口。

```bash
curl http://127.0.0.1:8080/ping
```

响应示例：

```json
{"status":"ok","message":"pong","time":"2026-10-08T09:00:00Z"}
```

## 测试

```bash
go test ./...
```
