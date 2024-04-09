# Restful Server

一个 go server，基于 Restful 的设计风格，用于简单模拟 task 调度。

项目介绍地址：[Go 中的 REST 服务器：第 1 部分 - 标准库 - Eli Bendersky 的网站 (thegreenplace.net)](https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/)

## 1. taskstore

一个简单的用map实现的task数据库，以及一些提供的接口，例如创建task,根据id,tags,due获取和删除task等。以及通过使用sync.Mutex来保证并发安全。

提供的接口如下所示：

```
POST   /task/              :  create a task, returns ID
GET    /task/<taskid>      :  returns a single task by ID
GET    /task/              :  returns all tasks
DELETE /task/<taskid>      :  delete a task by ID
GET    /tag/<tagname>      :  returns list of tasks with this tag
GET    /due/<yy>/<mm>/<dd> :  returns list of tasks due by this date
```

后续可以考虑更新为用mysql,redis等来实现。

## 2. taskServer

使用taskstore的相关接口，来实现server，使得能够正常使用store的功能。

使用最简单，传统的net/http提供的路由实现。

golang更新到 1.22 之后，官方的 net/http 的 pattern 就支持为如下格式了:

> [METHOD ][HOST]/[PATH]

其他具体的改变，可以查看文档:[http package - net/http - Go Packages](https://pkg.go.dev/net/http@master#ServeMux)获取更详细的信息。

## 3. 路由聚合

关于路由聚合的方法介绍，见[Go 中 HTTP 路由的不同方法 (benhoyt.com)](https://benhoyt.com/writings/go-routing/)。

这里介绍了相当多的路由匹配的方式，例如正则表达式，拆分路径，使用第三方 package 等。

也尝试了几种，例如使用gorilla/mux，以及使用gin来重写路由。

## 4. 使用Swagger

我们最开始规范的定义的store的接口并不符合rest规范，也没有合适的doc来描述这些接口。

因此，使用Swagger/OpenAPI来完善我们的接口设计。





