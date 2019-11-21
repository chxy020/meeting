package main

import (
	"os"

	"github.com/urfave/cli"

	"seater/cmd"
	"seater/database"
)

func syncDb() {
	cmd.InitDb(false)

	database.CheckCreateMigrationsTable()
	database.UpgradeDB()
}

func main() {
	app := cli.NewApp()
	app.Name = "seater-manage"
	app.Usage = "Manage tool for seater"
	app.Commands = []cli.Command{
		{
			Name:  "db",
			Usage: "Database manage",
			Subcommands: []cli.Command{
				{
					Name:  "sync",
					Usage: "Migrate database",
					Action: func(c *cli.Context) {
						syncDb()
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
