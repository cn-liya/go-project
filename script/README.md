### 安装cobra
gopath下执行
```bash
go install github.com/spf13/cobra-cli@latest
```

### 创建子命令
go-project/script$ 目录下执行
```bash
cobra-cli add [cmd] 
```

### 运行示例
go-project/script$ 目录下执行
```bash
go run main.go cronjob #定时任务
go run main.go refresh:token #刷新小程序服务端access_token
go run main.go avatar2cdn #临时头像链接转存CDN
```

