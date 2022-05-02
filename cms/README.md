### 示例接口
- /captcha 获取登录验证码
- /user/login 登入
- /user/logout 登出
- /user/password 修改密码
- /system/action-menu/sync 一键同步接口和菜单
- /system/action-menu/list 获取接口列表和菜单树
- /system/action/update 更新接口描述(备注和排序)
- /system/role/list 角色分页列表
- /system/role/option 所有角色选项(仅包含id和name)
- /system/role/create 创建角色
- /system/role/update 更新角色
- /system/user/list 管理员分页列表
- /system/user/create 创建管理员账号
- /system/user/password 重置账号密码
- /system/user/role 分配账号角色
- /system/user/status 切换账号状态
- /upload/image 图片上传(500KB)
- /upload/audio 音频上传(2M)
- /upload/video 视频上传(5M)
- /upload/pdf PDF上传(1M)
- /crypto/encrypt 加密文本
- /crypto/decrypt 解密文本

### 权限管理设计
> - 每个角色分别勾选接口权限和菜单权限，数据库只保存接口和菜单的唯一标识。
> - 后端根据接口权限判断，无权限的返回403。
> - 前端根据菜单权限关键字渲染菜单，根据接口权限关键字显示隐藏操作按钮。
> - 接口和菜单采用一键同步操作，后端读取路由表保存到系统接口，前端读取路由表提交到系统菜单。
