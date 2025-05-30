# dineflow

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.24%2B-blue?logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/frontend-react%20%7C%20vite-blue?logo=react)](https://react.dev/)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/FacundoChan/dineflow/.github%2Fworkflows%2Fgo-test.yml)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)

</div>

dineflow 是一个基于微服务架构的点餐系统，涵盖订单、库存、支付、厨房等核心业务，支持分布式事务、消息队列、链路追踪、服务发现等云原生特性。后端使用 Go 语言实现，前端采用 React + TypeScript + Vite。

## 功能概述

- [x] **订单服务（Order Service）**：负责订单的创建、查询、状态流转，支持 RESTful API 和 gRPC。
- [x] **库存服务（Stock Service）**：管理商品库存，校验下单时库存充足性。
- [x] **支付服务（Payment Service）**：对接 Stripe，生成支付链接，处理支付回调，更新订单支付状态。
- [x] **厨房服务（Kitchen Service）**：消费已支付订单，模拟出餐流程。
- [x] **消息队列（RabbitMQ）**：实现服务间事件驱动通信（如订单创建、支付完成等事件）。
- [x] **服务发现与配置（Consul）**：服务注册与发现。
- [x] **链路追踪（Jaeger）**：分布式链路追踪。
- [ ] **监控（Prometheus + Grafana）**：服务监控与可视化。
- [x] **数据库**：MySQL（库存）、MongoDB（订单）、Redis（分布式锁）。

## 技术栈

- **后端**：Go、gRPC、Gin、RabbitMQ、MySQL、MongoDB、Redis、Consul、Jaeger、~~Prometheus、Grafana~~
- **前端**：React、TypeScript、Vite、TailwindCSS
- **支付**：Stripe

## 快速开始

### 1. 克隆项目

```bash
git clone git@github.com:FacundoChan/dineflow.git
cd dineflow
```

### 2. 启动依赖服务

> [!NOTE]
> 确保本地已安装 Docker 和 Docker Compose。

```bash
docker compose -f ./docker-compose.yml -p dineflow up -d --build
```

### 3. 启动 Stripe Webhook（支付功能）

```bash
stripe listen --forward-to localhost:8284/api/webhook
```

### 4. 启动后端服务

每个微服务可单独运行，示例（以 order 服务为例）：

```bash
cd internal/order
go run .
```

其余服务（stock、payment、kitchen）同理。

### 5. 前端

访问 [http://localhost:3001](http://localhost:3001)

## 目录结构

```
.
├── internal/         # Go 微服务主目录（order, stock, payment, kitchen）
├── frontend/         # 前端项目（React + Vite）
├── api/              # OpenAPI/Proto 接口定义
├── docker-compose.yml
├── scripts/          # 辅助脚本
├── prometheus/       # Prometheus 配置
├── data/             # Redis/MongoDB 数据
├── mysql_data/       # MySQL 数据
└── ...
```

## API 文档

- OpenAPI 文档：`api/openapi/order.yml`
- 主要接口：
  - `POST /customer/{customer_id}/orders` 创建订单
  - `GET /customer/{customer_id}/orders/{order_id}` 查询订单
  - `GET /products` 获取商品列表

可用工具如 [Swagger Editor](https://editor.swagger.io/) 打开 `order.yml` 进行交互式测试。

## 其他命令

#### Docker & 服务编排

```sh
docker compose -f ./docker-compose.yml -p dineflow up -d --build
# 或启动部分服务
docker compose -f ./docker-compose.yml -p dineflow up -d --build consul mysql rabbit-mq jaeger order-mongo mongo-express
# 停止并移除容器
docker compose -f ./docker-compose.yml -p dineflow down
# 停止并移除容器及数据卷
docker compose -f ./docker-compose.yml -p dineflow down -v
```

#### Stripe 支付开发

```sh
stripe listen --forward-to localhost:8284/api/webhook
```

#### Makefile

```sh
make gen          # 生成所有代码（包括 Proto 和 OpenAPI）
make genproto     # 生成 Proto 代码
make genopenapi   # 生成 OpenAPI 代码
make lint         # 代码风格检查
```

单独执行 Lint：

```sh
lint: golangci-lint run --config ../../.golangci.yaml
```

#### 测试

测试文件示例：

- internal/order/tests/create_order_test.go
- internal/stock/adapters/stock_mysql_repository_test.go

运行部分测试（`stock`部分示例）：

```sh
cd internal/stock/adapters
go test -run 'OverSell'
go test -run 'Race'
```

运行全部测试：

```sh
go test ./...
```

#### 监控与链路追踪

- Consul: http://127.0.0.1:8500/ui/dc1/services
- Jaeger: http://localhost:16686/search
- MongoExpress: http://localhost:8082

## 贡献

- 遵循 DDD 分层架构（Entities, Use Cases, Interface Adapters, Infrastructure）
- 代码风格与规范见 `.golangci.yaml`

## License

MIT
