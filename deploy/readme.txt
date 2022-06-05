api和cms线上部署分别以api.conf和cms.conf替代默认的nginx.conf，配置值可根据实际需要调整优化

开发环境可将dev.conf复制到nginx的include包含目录下，并在hosts文件添加
127.0.0.1 api-dev.domain.cn
127.0.0.1 cms-dev.domain.cn

注意：如在nginx添加了跨域头，则路由不能使用Cors中间件，否则会出现两个相同header的错误。
反之，如未在nginx添加跨域头，则路由需使用Cors中间件，否则前端浏览器会拦截请求。
