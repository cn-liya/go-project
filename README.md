## go-project
Go语言搭建的轻量级应用，包含以下三个服务，可分别编译后部署。
- api: 后端接口服务
- cms: 内容管理后台
- script: 常驻脚本任务
> api和cms服务基于web框架gin，script服务基于cli框架cobra。

### 目录结构
```
api/
    conf.yaml             # api服务本地配置文件
    docs/                 # 运行依赖目录
    internal/             # api服务私有包
        handler/          # 请求控制(参数校验、登录鉴权、数据组装)
            handler.go    # handler初始化和公共定义
            router.go     # api服务所有路由统一文件
            *.go          # 业务接口
        proto/            # 交互数据定义
        service/          # 数据处理(db、cache、mq等)
            service.go    # service初始化
            *.go          # 业务逻辑
    main.go               # api服务入口文件

cms/
    conf.yaml             # cms服务本地配置文件
    docs/                 # 运行依赖目录
    internal/             # cms服务私有包
        handler/          # 请求控制(参数校验、登录鉴权、数据组装)
            handler.go    # handler初始化和公共定义
            router.go     # cms服务所有路由统一文件
            *.go          # 业务接口
        proto/            # 交互数据定义
        service/          # 数据处理(db、cache、mq等)
            service.go    # service初始化
            *.go          # 业务逻辑
    main.go               # cms服务入口文件

script/
    conf.yaml             # script服务本地配置文件
    docs/                 # 运行依赖目录
    internal/             # script服务私有包
        handler/          # 逻辑控制
        proto/            # 交互数据定义
        service/          # 数据处理
    cmd/                  # 命令注册(任务的启动、轮次、停止控制)
        root.go           # 根命令
        *.go              # 一个文件代表一个命令
    main.go               # script服务入口文件

model/                    # 存储模型定义
    entity/               # 数据库实体
    cache/                # 缓存的key和结构
    queue/                # 消息队列名称和结构
pkg/                      # 公共方法包
    logger/               # 日志
    process/              # 进程相关方法
    core/                 # 核心上下文封装，主要为script使用
    cli/                  # cli服务扩展(script使用)
    web/                  # web服务扩展(api/cms使用)
        errcode/          # 响应错误码定义
        middleware.go     # 通用路由中间件方法
    db/                   # 数据库连接
    gredis/               # redis连接
    gnsq/                 # nsq连接
    rabbitmq/             # rabbitmq连接
    coss/                 # 云对象存储API
    wechat/               # 微信小程序API
design/                   # 设计文档
```

> + proto是【内部交互】的结构，用于服务内不同层级出入参。
> + model是【外部存储】的模型，用于各服务读写数据中间件。
> + entity是最底层结构，cache和queue可聚合entity。
> + proto定义可根据业务需要聚合model结构。
> + model可引入pkg，例如ORM钩子函数使用公共方法、中间件存储结构包含外部API出入参。
> + pkg不可引入model，模型的映射转换关系定义在model下，不能定义到pkg下。

> * 项目组织采用结合Go语言特性简化过的OOP设计模式，依赖注入关系：入口函数>handler层>service层。
> * handler和service目录下每个文件可视作一个"类"，通过receiver注入的Handler和Service视为"基类"，私有字段视为"基类”的受保护属性。
> * 更贴合其他语言OOP风格的做法是：对每个文件下的"类"定义一个结构体并嵌入作为"基类"的结构体实现继承，对外暴露一个"NewXXX"函数实现"构造函数"。
> * 这里做了简化：handler和service目录下不再拆分目录，直接在"基类"上扩展方法，省去为每个文件下的"类"单独定义结构体和实例化函数。

> Go有一种拒绝大型依赖关系树的文化："A little copying is better than a little dependency."

### 版本依赖
* Go(v1.22+)
* golangci-lint(v1.58+)
* MySQL(v8.0+)
* Redis(v5.0+)
* NSQ(v1.2+) / RabbitMQ(v3.08+)

### 初始化
- 需将api、cms、script目录下conf.dev.yaml复制为conf.yaml并修改相应配置。
- 如果依赖安装在本机，则可在hosts文件创建映射`127.0.0.1 host.docker.internal`，以使容器内外运行配置统一。
- MySQL需导入`design/sql`目录下的数据表。
- 消息队列若使用NSQ无需初始化，若改用RabbitMQ则需要执行`design/mq`目录下的shell命令或在管理UI界面创建队列。
- 升级和安装依赖包：
```shell
go get -u ./...
go mod tidy
go mod vendor
```

### 单元测试运行
```shell
go test ./... -cover
```

### 编译运行
> - 本地调式可直接在api、cms、script三个目录下运行`go run .`命令临时编译并运行，`debug`标签会开启pprof性能分析
```shell
(cd api && go run -tags sonic,debug .)
(cd cms && go run -tags sonic,debug .)
(cd script && go run -tags sonic,debug . cronjob)
```
> - 本地构建docker镜像并运行
```shell
tag=v$(date "+%Y%m%d.%H%M%S")

docker build --build-arg srv=api -t api:$tag .
docker run -d -p 8000:8000 api:$tag

docker build --build-arg srv=cms -t cms:$tag .
docker run -d -p 6000:6000 cms:$tag

docker build --build-arg srv=script -t script:$tag .
docker run -d --name cronjob-$tag script:$tag cronjob
```
> - 生产部署需在CI发布流水fetch代码后，从配置中心（Zookeeper/Nacos/Apollo/etcd）拉取环境配置`conf.yaml`到对应的服务目录下，再构建镜像。

### 服务拆分方向
+ api面向C端用户
+ cms面向管理后台
+ script处理异步任务
> + 如需接入第三方异步通知，可根据服务规模考虑：在api服务中接入，或是拆分出一个notify服务用于接入。
> + script服务如果编译体积非常臃肿，也可以按一定规则拆分服务。
> + 同一个项目（业务相关且数据互通）的所有服务都放在同一个代码仓库和go模块下。
> + 不同项目（业务独立或数据隔离）的服务即使架构和逻辑都非常相似，也要剥离出去。

### 日志设计
- 基于标准库原生封装，支持控制台标准输出、格式化输出、文件输出。通过配置app.logger.output指定。
- 容器部署可根据日志收集策略指定为`std`标准输出或`file`文件输出，本地调试可指定为`fmt`格式化输出。
- 可设置字段加密，密钥为16/24/32位ASCII字符。
- 加密字段支持忽略大小写模糊匹配，字段名必须符合变量命名规范（仅包含字母、数字、下划线）。
- 若使用自建日志服务，可在应用程序中去掉加密处理，改由采集程序处理，以使应用程序获得更优的性能表现。
#### 字段含义
- v0 单次handler的链路追踪ID（http请求ID、队列消息ID、某次任务的标识……）
- v1 进入handler的路径标识（http请求path、cmd启动命令……）
- v2、v3 进程、用户、数据……其他标识
- level 日志等级，设置5个级别：
  * 1 Fatal 内部程序错误（panic）
  * 2 Error 外部程序错误（数据库、缓存、消息队列、外部api……）
  * 3 Warn 业务警告
  * 4 Info 普通信息
  * 5 Trace 追踪记录
- time 日志打印时间
- msg 简要说明
- input 入参、请求体……
- output 结果、响应体……
- elapsed （耗时环节）经过的毫秒数
#### 使用示例
```
//单次打印
logger.FromContext(c).Info("message","input","output")
//多次打印
l := logger.FromContext(ctx)
l.Error("message","input","output")
l.Warn("message","input","output")
l.Info("message","input","output")
```
- 使用了AccessLog中间件的接口会自动记录请求和响应，msg为`access`。
- 使用logger包的NewHttpClient或NewTransport初始化的client发起的http请求都会自动打印trace日志，msg为`request`。
  <br>如需串连上下文日志，封装第三方请求需使用http.NewRequestWithContext并传入context
- gorm配置开启tracelog会打印所有查询日志，input为sql语句，output为rows或error，msg为`gorm`。

### 接口协议
- 所有接口请求应使用https传输协议。
- 使用json作为数据传输格式，特殊接口(上传、下载)除外。
- 请求方法一般使用POST，特殊需要(点击链接激活、静态资源代理等)使用GET。
- 为便于网关监控对错误进行分类统计以及前端统一异常处理，使用http状态码对错误进行归类。
- 客户端关注4xx错误，服务端关注5xx错误。
- 状态码非200时，在响应body体用code和msg对错误进行描述。
- 主键查询获取为空返回404错误，条件查询为空返回200状态码和空数组。
- 请求Header头需携带以下参数：
  + Authorization: (omitempty) 登录Token
  + X-Request-Id: (required,min=8,max=40) 调用端随机生成的唯一请求id，用于追踪请求链路。
- 面向客户端的接口使用Token鉴权。面向服务端的接口使用账号密码鉴权，密码可以通过签名校验而无需网络传输。

### 安全策略
+ 密码使用argon2id单向哈希算法储存
+ 敏感信息落库使用AES算法CTR模式加密
+ 日志特定字段使用AES算法ECB模式加密
+ 网关层根据ip限制接口请求频率

### 参考文档
+ [gin-gonic](https://pkg.go.dev/github.com/gin-gonic/gin)
+ [validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
+ [cobra](https://pkg.go.dev/github.com/spf13/cobra)
+ [gorm](https://gorm.io)
+ [go-redis](https://redis.uptrace.dev)

