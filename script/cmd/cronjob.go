package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"project/pkg/cli"
	"project/pkg/process"
	"project/script/internal/handler"
	"project/script/internal/service"
)

/*
spec格式： Seconds Minutes Hours Day-of-month Month Day-of-week

字段说明：
Field name   | Allowed values  | Allowed special characters
----------   | --------------  | --------------------------
Seconds      | 0-59            | * / , -
Minutes      | 0-59            | * / , -
Hours        | 0-23            | * / , -
Day-of-month | 1-31            | * / , - ?
Month        | 1-12 or JAN-DEC | * / , -
Day-of-week  | 0-6 or SUN-SAT  | * / , - ?

匹配符：
* 匹配任意值
/ 范围的增量，如在Minutes字段用"3-59/15"表示一小时中第3分钟开始到第59分钟，每隔15分钟
, 用于分隔列表中的项目，如在Day-of-week字段用"MON,WED,FRI"表示星期一三五
- 用于定义范围，如在Hours字段用"9-18"表示9:00~18:00之间的每小时
? 可用于代替*，将Day-of-month或Day-of-week留空

预定义简写：
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *

@every <duration> 表达式，duration为time.ParseDuration支持的格式，例如："@every 1h30m10s"，不对齐整点时分。

默认标准parser只有5位：分时日月周，如需秒开始的6位，创建*cron.Cron使用的cron.New函数需传入选项cron.WithSeconds()
cron.WithChain选项可以指定一些链式方法，例如包装一层recover。
默认每次到达调度时间就会执行一次job，不论上一次job是否结束。
可在链式选项中指定，当上次job未结束时如何处理：
DelayIfStillRunning等待上一次job结束再执行本次job，或SkipIfStillRunning跳过执行本次job。
*/

var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "定时任务",
	Run: func(cmd *cobra.Command, args []string) {
		h := handler.New(&cfg.Handler, service.New(service.Mysql(&cfg.Service.Mysql), service.Redis(&cfg.Service.Redis)))
		c := cron.New()

		cli.CronAdd(c, "@every 3m", h.WechatServerToken)
		cli.CronAdd(c, "1 0 * * *", h.ReportNewUser)

		c.Start()
		process.Notify()
		<-c.Stop().Done()
	},
}
