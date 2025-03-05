# 项目结构说明

## 目录结构

```
.
├── configs/              # 配置文件目录
├── biz/             # 私有应用程序代码
│   ├── user/             # 用户服务
│   ├── chat/             # 聊天服务
│   ├── emotion/          # 情绪分析服务
│   └── pkg/              # 内部共享包
├── pkg/                  # 公共代码包
│   ├── auth/             # 认证相关
│   ├── database/         # 数据库操作
│   ├── mq/               # 消息队列
│   ├── websocket/        # WebSocket 实现
│   └── utils/            # 工具函数
├── scripts/              # 构建、部署脚本
├── test/                 # 测试文件
└── web/                  # Web 相关资源

```

## 服务说明

### 用户服务 (User Service)
负责用户管理，包括家长和儿童账户的注册、认证和信息管理。

### 聊天服务 (Chat Service)
处理语音聊天功能，包括语音识别、对话生成和语音合成。

### 情绪分析服务 (Emotion Service)
分析儿童情绪状态，生成情绪报告并负责定时推送。

### API 网关 (Gateway)
统一的 API 入口，处理请求路由、负载均衡和认证。

## 技术组件

### 数据库
- MySQL/PostgreSQL：存储用户数据、聊天记录和情绪分析结果
- Redis：缓存层，存储临时数据和会话信息

### 消息队列
- RocketMQ：处理异步消息、事件驱动和定时任务

### 监控和日志
- Prometheus：性能监控
- Grafana：可视化监控数据
- ELK Stack：日志收集和分析

### 部署
- Docker：容器化部署
- Kubernetes：容器编排和服务管理
- Nginx：API 网关和负载均衡