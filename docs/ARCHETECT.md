# 架构思路文档

核心开发架构方式为部分借鉴 DDD 的三层架构。

## Domain Model 层

核心数据建模层。采用部分的充血模型，对与 model 强关联的操作进行了定义。

## Domain Service 层

对于不与单独某个 model 强关联的泛逻辑进行聚合，仅包含纯函数。行为上类似添补 Go 没有 static 函数的生态空缺。

注意，与常规 DDD 的范式不同，此项目的 domain service 不依赖 Domain Repository 层，而是将包括事务在内所有 repository 操作全部上提至 Application 层进行统一处理。

一些需要特定 repository 数据的，采用注入回调函数懒加载的方式获取。
