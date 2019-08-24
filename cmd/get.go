package main

import (
	"log"

	gmp "github.com/firmanmm/go-mod-private"
	"github.com/urfave/cli"
)

type GetCmd struct {
	setting *gmp.Setting
	getter  *gmp.GetterManager
}

func (g *GetCmd) Run(ctx *cli.Context) error {
	gompConf := ctx.GlobalString("gomp")
	if len(gompConf) == 0 {
		gompConf = "mod.gomp"
	}
	if err := g.setting.LoadFromFile(gompConf); err != nil {
		return err
	}
	args := ctx.Args()
	trailing := args.Tail()
	log.Println(trailing)
	//return g.setting.SaveToFile(gompConf)
	return nil
}

func (g *GetCmd) Init() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Perform Go Get operation and switch to ssh if a matching SSH Credential has been registered",
		Subcommands: []cli.Command{
			cli.Command{
				Action: g.Run,
			},
		},
	}
}

func NewGetCmd(setting *gmp.Setting, getter *gmp.GetterManager) *GetCmd {
	instance := new(GetCmd)
	instance.setting = setting
	instance.getter = getter
	return instance
}
