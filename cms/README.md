### 示例接口
- GET/ping 连通测试
- GET/modules 所有模块列表(仅非生产环境)
- GET/captcha 获取登录验证码
- POST/user/login 登入
- DELETE/user/logout 登出
- POST/user/password 修改密码
- GET/admin/list 管理员列表
- POST/admin/account 创建管理员账号
- PUT/admin/password 重置账号密码
- PUT/admin/authority 更新账号权限
- PUT/admin/status 切换账号状态

### 目录结构
- /docs/log/ 日志文件输出
- /handler/ 请求控制层
- /internal/acl/ 权限控制模型
- /internal/proto/ 自定义结构，用于handler和service层交互
- /internal/service/ 逻辑处理层

### 权限管理设计
> - 以模块为单位，给每个账号指定各个模块的权限（0无权限，1只读，2读写）
> - 前端根据登录返回的账号权限模块加载菜单，根据是否有写权限显示操作按钮。
> - 后端根据账号模块权限判断接口权限，无权限的返回403。

