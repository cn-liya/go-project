## go-project
使用Golang搭建轻量级应用，仓库包含以下三个服务，可分别编译后部署。
- api: 后端接口服务(gin-gonic/gin)
- script: 常驻脚本任务(spf13/cobra)
- cms: 内容管理后台(gin-gonic/gin)

### 目录结构
```
api/
    conf.yaml             #api服务本地配置文件
    docs/                 #运行依赖目录
    handler/              #请求控制
        handler.go        #handler初始化和公共定义
        router.go         #api服务所有路由
        *.go              #业务接口
    internal/             #api服务私有包
        proto/            #内部交互数据定义
        service/          #业务处理
            service.go    #service初始化
            *.go          #业务逻辑
    main.go               #api服务入口文件

cms/
    conf.yaml             #cms服务本地配置文件
    docs/                 #运行依赖目录
    handler/              #请求控制
        handler.go        #handler初始化和公共定义
        router.go         #cms服务所有路由
        *.go              #业务接口
    internal/             #cms服务私有包
        acl/              #权限控制
            model.go      #模型
            modules.go    #模块
        proto/            #内部交互数据定义
        service/          #业务处理
            service.go    #service初始化
            *.go          #业务逻辑
    main.go               #cms服务入口文件           

script/           
    cmd/                  #命令注册
        root.go           #根命令
        *.go              #一个文件代表一个命令
    conf.yaml             #script服务本地配置文件  
    docs/                 #运行依赖目录
    internal/             #script服务私有包
        service/          #业务处理 
            service.go    #service初始化
            *.go          #任务的具体实现
    main.go               #script服务入口文件

model/                    #数据模型和常量，各服务共用
pkg/                      #公共方法包，各服务可按需引入
    logger/               #日志
    cache/                #缓存(eg:redis)
    db/                   #数据库(eg:mysql)
    mq/                   #消息队列(eg:nsq)
    wechat/               #微信小程序接口
    util/                 #其他公共方法
design/                   #设计相关文档
deploy/                   #部署相关配置

```

### 代码层级
> - handler(请求控制)：参数校验、登录鉴权、返回数据组装。
> - service(业务逻辑)：处理数据，db、cache、mq、openapi等。
> - cmd(命令注册)：任务的启动、轮次、停止控制。

### 安装和依赖
- Go Version >= v1.18
- 首次下载项目后需执行`go mod download`和`go mod vendor`
- 运行依赖mysql,redis,nsq，需将api、cms、script目录下conf.yaml相应配置修改为本机开发环境。
- mysql需导入 design/sql 目录下的数据表。

### 编译运行
> - 分别进入api、cms、script目录执行`go build`命令；再运行该目录下的二进制文件。
> - 编译后的二进制文件和各自的docs目录、conf.yaml配置文件应平行放置在同一目录层级下。
> - 本地也可直接在三个目录下运行`go run main.go`命令。

### 日志设计
- 基于官方log包封装，支持控制台标准输出、格式化输出、文件写入，通过配置app.logger指定，默认不输出。
- 容器部署可根据日志收集策略指定为`std`标准输出或`file`文件写入，本地调试可指定为`fmt`格式化输出。
- level表示日志等级，预设4个级别：
1. Fatal 内部程序错误(panic)
2. Error 外部程序错误(数据库、缓存、消息队列、第三方接口)
3. Warn 业务告警
4. Info 业务信息

调用日志方法示例：
```
//单条打印
logger.FromContext(c).Info("message","input","output")

//多条打印
l := logger.FromContext(ctx)
l.Error("message","input","output")
l.Warn("message","input","output")
l.Info("message","input","output")
```
- 使用了AccessLog中间件的接口会自动记录请求和响应，msg为`access`。
- 使用logger包的NewHttpClient或NewTransport初始化的client发起的http请求都会自动打印trace日志，msg为`request`。
  <br>如需串连上下文日志，封装第三方请求需使用http.NewRequestWithContext并传入context
- gorm会打印trace日志，input为sql语句，output为rows或error，msg为`gorm`。
> 日志内容的解析、脱敏、反序列化等，统一放到日志收集脚本处理，减少牺牲应用程序性能。

### 接口协议
- 使用json作为数据传输格式，文件流数据需base64编码后放到json对象中。
- 增删改查分别使用POST,DELETE,PUT,GET请求方法。
- 获取详情的唯一参数(如ID)放在path路径，获取为空返回404错误；条件查询参数放queryString，查询空返回空数组。
- 请求Header头需携带以下参数：
  + Authorization: (omitempty) 登录Token
  + X-Trace-Id: (required,min=8,max=40)随机生成的唯一请求id，用于追踪请求链路
- 使用http状态码表示执行状态，错误信息在响应body体用msg和detail(omitempty)返回，detail前端不显示。

#### 状态码列表
+ 200: 成功
+ 400: 参数错误
+ 401: 登录失效
+ 403: 禁止操作(无权限)
+ 404: 目标不存在
+ 409: 数据已存在
+ 413: 提交内容过大
+ 415: 错误的文件类型
+ 422: 数据格式错误或已过期
+ 423: 资源被锁定
+ 429: 请求频率限制
+ 500: 服务端通用错误
+ 502: 服务端响应错误
+ 503: 服务不可用(停机维护)
+ 504: 服务端请求错误
