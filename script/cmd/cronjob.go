package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"project/script/internal/service"
	"syscall"
)

/*
spec格式:
	"Seconds Minutes Hours Day-of-month Month Day-of-week"

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

默认标准parser只有5位：分时日月周，如需秒开始的6位，创建*cron.Cron需使用cron.New(cron.WithSeconds())
*/

var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "定时任务",
	Long:  "精确定时执行的任务",
	Run: func(cmd *cobra.Command, args []string) {
		s := service.NewCronJob()
		c := cron.New()
		var err error

		_, err = c.AddFunc("5 0 * * *", s.GetDailySummary) // 每天0点5分拉取昨日访问概况
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.AddFunc("10 0 * * *", s.GetDailyVisitTrend) // 每天0点10分拉取昨日访问趋势
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.AddFunc("15 0 * * 1", s.GetWeeklyVisitTrend) // 每周一0点15分拉取上周访问趋势
		if err != nil {
			log.Fatal(err)
		}

		c.Start()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit
		ctx := c.Stop()
		<-ctx.Done()
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}
