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
go run -tags sonic,debug . cronjob #定时任务
go run -tags sonic,debug . avatar:to:cdn #临时头像链接转存CDN
```

