package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "GoModPrivate"
	cliApp.Description = "Allow you to use Go's Git Private Server as your source dependency"
	cliApp.Email = "maulana.firman56@gmail.com"
	cliApp.Author = `Firman "Rendoru" Maulana`
	cliApp.Commands = []cli.Command{}
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalln("Failed to start GoModPrivate CLI, ", err.Error())
	}
}
