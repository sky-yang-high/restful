# Restful Server

一个 go server，基于 Restful 的设计风格，用于简单模拟 task 调度。

项目介绍地址：[Go 中的 REST 服务器：第 1 部分 - 标准库 - Eli Bendersky 的网站 (thegreenplace.net)](https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/)

## 1. taskstore

## 2. taskServer

go 更新到 1.22 之后，官方的 het/http 的 pattern 就支持为如下格式了:

> [METHOD ][HOST]/[PATH]

其他具体的改变，可以查看文档:[http package - net/http - Go Packages](https://pkg.go.dev/net/http@master#ServeMux)获取更详细的信息。

## 3. 路由聚合

关于路由聚合的方法介绍，见[Go 中 HTTP 路由的不同方法 (benhoyt.com)](https://benhoyt.com/writings/go-routing/)。

这里介绍了相当多的路由匹配的方式，例如正则表达式，拆分路径，使用第三方 package 等。
