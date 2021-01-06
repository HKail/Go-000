# 学习笔记

## 课间记录

### Standard Go Project Layout

https://github.com/golang-standards/project-layout/blob/master/README_zh.md

### 依赖注入

https://blog.glang.org/wire

### Service Layout

https://github.com/go-kratos/service-layout

### gRPC

```shell
protoc --go out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
helloworld/helloworld.proto
```

### Protobuf v3（空基础类型封装）

https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto

### Google API Design Guide

FieldMask

https://www.bookstack.cn/books/API-design-guide

### 动态配置

https://pkg.go.dev/expvar

## 作业题目

1. 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

## 作业内容

