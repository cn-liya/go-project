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
go run main.go cronjob
go run main.go refresh:token
go run main.go message
```

### 示例任务
- cronjob 定时拉取微信analysis数据导入到db并通过机器人发送消息到钉钉、企业微信
- refresh:token 刷新小程序服务端access_token并保存到redis
- message 消费NSQ消息

