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

## 用户查询

`GET /users` 按查询参数过滤内存中的示例用户。多个条件同时生效。不传参数时返回全部用户。

| 参数 | 含义 |
| --- | --- |
| `q` | 在 id、姓名、邮箱中做不区分大小写的包含匹配 |
| `id` | 按 id 精确匹配 |
| `name` | 姓名包含匹配 |
| `email` | 邮箱包含匹配 |

```bash
curl "http://127.0.0.1:8080/users?q=hopper"
curl "http://127.0.0.1:8080/users/u1"
```

查询响应示例：

```json
{"total":1,"users":[{"id":"u3","name":"Grace Hopper","email":"grace@example.com"}]}
```

`GET /users/{id}` 返回单个用户。id 不存在时返回 404。

## 订单查询

`GET /orders` 按查询参数过滤内存中的示例订单。多个条件同时生效。不传参数时返回全部订单。金额单位是分。

| 参数 | 含义 |
| --- | --- |
| `q` | 在订单 id、用户 id、商品、状态中做不区分大小写的包含匹配 |
| `id` | 按订单 id 精确匹配 |
| `userId` | 按用户 id 精确匹配 |
| `item` | 商品名称包含匹配 |
| `status` | 按状态精确匹配，例如 `pending`、`paid`、`shipped` |

```bash
curl "http://127.0.0.1:8080/orders?userId=u1&status=paid"
curl "http://127.0.0.1:8080/orders/o1"
```

查询响应示例：

```json
{"total":1,"orders":[{"id":"o2","userId":"u2","item":"Enigma notes","status":"pending","amount":4500}]}
```

`GET /orders/{id}` 返回单个订单。id 不存在时返回 404。

## 测试

```bash
go test ./...
```
