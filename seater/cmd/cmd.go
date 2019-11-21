package cmd

import (
	"os"
	"seater/database"

	"github.com/astaxie/beego"
)

// Exit exits with signal
func Exit() {
	os.Exit(1)
}

// InitDb initialize database connection
func InitDb(retry ...bool) {
	err := database.RegisterDb("127.0.0.1", retry...)
	if err != nil {
		Exit()
	}
}

// InitBeego initialize beego
func InitBeego() {
	beego.BConfig.WebConfig.DirectoryIndex = true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
}