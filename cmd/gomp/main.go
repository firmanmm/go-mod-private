package main

import (
	"log"
	"os"

	"github.com/firmanmm/go-mod-private/cmd"

	gmp "github.com/firmanmm/go-mod-private"

	"github.com/urfave/cli"
)

func main() {

	setting := gmp.NewSetting()
	getter := gmp.NewGetterManager(setting)
	syncMgr := gmp.NewSyncManager(setting, getter)

	getCmd := cmd.NewGetCmd(setting, getter)
	addCredentialCmd := cmd.NewAddCredentialCmd(setting)
	syncCmd := cmd.NewSyncCmd(setting, syncMgr)

	cliApp := cli.NewApp()
	cliApp.Name = "GoModPrivate"
	cliApp.Usage = "Go Module for Private Git Server Repository"
	cliApp.Description = "Allow you to use Git Private Server as your source dependency while using Go Module"
	cliApp.Author = `Firman "Rendoru" Maulana`
	cliApp.Commands = []cli.Command{
		getCmd.Init(),
		addCredentialCmd.Init(),
		syncCmd.Init(),
	}
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "gomp",
			Usage: "Read from given gomp file",
			Value: "mod.gomp",
		},
	}
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalln("Failed to start GoModPrivate CLI, ", err.Error())
	}
}
