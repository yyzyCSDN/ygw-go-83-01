基于 Go 实现的风力发电机组变桨与偏航控制系统项目，一款新能源发电设备控制服务，完成叶片变桨、机舱偏航对风、转速限制与安全停机联动管理。

## 项目简介

本项目面向风力发电机组的控制侧，提供叶片变桨、机舱偏航对风、转子转速限制与安全停机联动能力。控制服务采集风速与转速，按风速调整桨距限制转速，偏航对风提高捕获，并在大风或故障时触发安全停机。

## 构建与运行

本机运行需要 Go 1.23 及以上版本，依赖已通过 vendor 目录离线提供。

```bash
go build -mod=vendor ./...
go vet -mod=vendor ./...
go test -mod=vendor ./...
```

启动控制服务：

```bash
go run -mod=vendor ./cmd/turbine
```

服务默认监听 8080 端口，访问根路径打开控制台页面，健康检查地址为 /healthz。

## 目录结构

- cmd/turbine：控制服务入口与 HTTP 接口
- internal/pitch：叶片变桨控制
- internal/yaw：机舱偏航对风
- internal/speed：转速限制与保护
- internal/safe：安全停机联动
- internal/sensor：风速与转速采集
- internal/alarm：告警与恢复
- internal/record：运行事件记录
- web/console.html：控制台前端页面
