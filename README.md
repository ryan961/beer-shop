# beer-shop

An experimental mono-repo microservices project based on Kratos.

本项目从上游 [`go-kratos/beer-shop`](https://github.com/go-kratos/beer-shop) fork 而来，目前作为一个实验性仓库持续演进，目标不是简单追随上游实现，而是围绕 **Kratos Project Layout Best Practices** 做迭代优化，并逐步引入一些优秀、现代的第三方库来改善工程组织与开发体验。

## 项目定位

这个仓库主要用于验证一套更适合中大型 Kratos mono-repo 的工程实践，重点关注：

- 在 mono-repo 下组织多个 service / job / interface 模块
- 统一 API、配置、依赖注入和代码生成的管理方式
- 演示多个微服务之间的依赖调用和基础设施集成方式
- 为后续项目布局、模块边界和工程工具链优化提供实验场

项目中的业务设计仍然以示例和实验为主，很多实现做了简化，不等同于完整生产级电商系统。

## 当前已调整的方向

目前已经开始落地的调整包括：

- 使用 `buf` 管理 Proto 定义与生成流程
- 使用 `samber/do` 替换 `google/wire` 依赖注入方案

### 1. 使用 `buf` 管理 Protos

仓库已经引入 `api/buf.yaml`、`api/buf.gen.yaml` 和 `api/buf.lock`，将 Proto 管理逐步收敛到 `buf` 工作流。

引入 `buf` 的目的主要是：

- 统一 Proto 组织、校验和代码生成配置
- 减少散落在各模块中的 `protoc` 调用脚本和参数维护成本
- 提升多服务 API 演进时的可维护性和一致性

简单理解，`buf` 是一套围绕 Protobuf 的现代化工具链，负责把“定义、检查、生成”几个环节标准化。对于 mono-repo 项目，它比直接手写 `protoc` 命令更容易维护。

### 2. 使用 `samber/do` 替换 `wire`

项目已不再继续依赖 `wire`，仓库中的 `wire` 目标也已经改为提示信息。当前依赖注入方案切换为 [`samber/do`](https://github.com/samber/do)。

这样调整的原因主要是：

- `wire` 已归档，不再适合作为后续持续演进的基础设施
- `samber/do` 更适合渐进式重构，不需要生成额外注入代码
- 对多模块 mono-repo 更友好，provider 组织与组合方式更直接

`samber/do` 是一个轻量级依赖注入容器，适合把构造逻辑、模块注册和应用装配放在运行时统一管理。当前仓库已经在多个 `app/*` 模块中逐步统一到这种 injector/provider 风格。

## 仓库结构

```text
.
├── api     // 对外 API Proto、错误定义、生成配置（buf）
├── app     // 各业务模块：service / job / interface / admin
├── pkg     // 公共库与通用组件
├── docs    // 设计与说明文档
├── deploy  // 部署相关文件
└── web     // 前端相关代码
```

更细化地看，项目仍然延续 Kratos mono-repo 的基本思路：

```text
.
├── api
│   ├── protos
│   │   └── <domain>/<kind>/v1/*.proto
│   └── _gen
│       └── go/...
├── app
│   └── <domain>
│       ├── service
│       ├── job
│       ├── interface
│       └── admin
├── pkg
└── docs
```

## 当前状态

- 项目仍处于实验和迭代阶段
- 代码可作为 Kratos mono-repo 组织与工程化演进的参考
- 目录结构、Proto 工作流、DI 方案仍会继续调整

如果你关注的是：

- Kratos 在 mono-repo 下的落地方式
- Proto 与代码生成工具链如何现代化
- `wire` 迁移到 `samber/do` 的组织方式

这个仓库会比上游更偏向“工程实验”和“持续优化”视角。

## 相关文档

- 架构与接口文档: [Docs](https://go-kratos.github.io/beer-shop/#/)
- 上游项目: [go-kratos/beer-shop](https://github.com/go-kratos/beer-shop)

## 注意

本项目当前仍以实验性验证为主，不承诺功能完整性或生产可用性。更适合作为：

- Kratos 项目布局参考
- 工程化改造样例
- Proto / DI / mono-repo 实践试验田
