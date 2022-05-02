package cli

import (
	"github.com/robfig/cron/v3"
	"log"
)

func CronAdd(c *cron.Cron, spec string, handler HandlerFunc) {
	fn := wrap(handler)
	if _, err := c.AddFunc(spec, fn); err != nil {
		log.Fatal(err)
	}
}
