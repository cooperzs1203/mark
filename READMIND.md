编程规则:
1.函数错误皆不打印，由上层调用函数处理决定

前提条件:
1.单个客户端的消息按顺序发送，无并发无切片包

Mark_V_0_1:
1.MServer(interface) : 定义了Mark服务器框架的基本要求，启动和停止一个服务器
    1.1 启动
    1.2 停止

2.Server(struct) : 实现MServer
    2.1 包含字段有Name、NetType、Host、Port、Listener、CloseFlag
    