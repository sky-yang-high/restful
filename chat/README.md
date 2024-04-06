# 说明

一个go的项目，来自project based learning中的推荐，[Building a Chat Application in Go with ReactJS](https://tutorialedge.net/projects/chat-system-in-go-and-react/)。

简单来说，就是在前端使用react实现ui和渲染，涉及js,html等，后端使用go来维护数据。

前面的几次改进过程忘记用git保存了，后悔莫及，不过也正好，后面可以尝试一下自己重写一次。

下面就简单说明一下每次改进的功能吧。

## 初始

把项目分为前端和后端两部分。后端部分先简单的写一个server即可。

前端部分，用npm安装并创建react实例。

## 简单交互

后端部分，重写server，改为使用socket来维护每次连接并收发数据。

前端部分，用一个最关键的代码：

```js
var socket = new WebSocket("ws://localhost:8080/ws")
```

来建立前端和后端的联系。或者是借助axios来使用HTTP进行交互也可以。

然后依次实现前端在不同事件下的动作即可。

这样一个简单的echo服务器即完成了。

## 前端优化

通过添加组件(Component)的方式来修饰前端，这里简单实现了两个组件-Hearder和chathisroty。chathistory实现有一定难度。

在前端的组件实现中，通常是这三个文件-c.jsx,c.scss,index.js。

## 后端优化

这部分最为关键，也略有难度，看的时候还懵懵懂懂的。

当前目标是实现一个服务器，当同时多个网页访问时，模拟他们在一起聊天的效果。

因此定义了一个Client变量描述客户，用一个全局的Pool来维护这些用户的登记，注销，广播等，怎么实现呢？这里就用到go的chan了。用户的不同动作通过chan就可以完美的实现阻塞/运行的过程。

需要注意，这个pool是一个全局的，即在用户之间传递的是它的指针。

(看websocket的时候，发现了一个chat的example，实现一模一样?乐，看来以后还得多看看不同包的example)













