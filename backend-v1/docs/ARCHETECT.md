# 架构思路文档

核心开发架构方式为部分借鉴 DDD 的三层架构。

## Domain 部分

### Domain Model 层

核心数据建模层。采用部分的充血模型，对与 model 强关联的操作进行了定义。

### Domain Service 层

对于不与单独某个 model 强关联的泛逻辑进行聚合，仅包含纯函数。行为上类似添补 Go 没有 static 函数的生态空缺。

注意，与常规 DDD 的范式不同，此项目的 domain service 不依赖 Domain Repository 层，而是将包括事务在内所有 repository 操作全部上提至 Application 层进行统一处理。

一些需要特定 repository 数据的，采用注入回调函数懒加载的方式获取。

### Domain Repository 层

与普通 DDD 范式基本一致，为 repository interface 层。

但是与理想 DDD 不同的是，它依赖实现细节 gorm.DB，因为在 Go 的框架下几乎无法在不使用任何泛型、any 和断言的情况下，规定一个 executor 事务上下文对象的接口。

## Application 层

## 其他部分

### Repository 模块

为简便起见，没有单独放置在 infrastructure 模块下。

### Value 模块

基本等价于 DTO 模块，提供对 model 层的封装。
