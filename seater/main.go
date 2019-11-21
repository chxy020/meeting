package main

import (
	"seater/cmd"
	"seater/cron"
	_ "seater/routers"

	"github.com/astaxie/beego"
)

func runCronJobs() (manager *cron.JobManager) {
	manager = cron.GetJobManager()
	err := manager.Init()
	if err != nil {
		cmd.Exit()
	}
	go manager.Start()
	return manager
}

func main() {
	cmd.InitBeego()
	cmd.InitDb(false)

	manager := runCronJobs()
	defer manager.Stop()
	beego.Run()
}
