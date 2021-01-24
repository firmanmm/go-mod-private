package cmd

import (
	"errors"
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
	if err := g.setting.LoadFromFile(gompConf); err != nil {
		return err
	}
	args := ctx.Args()
	if len(args) == 0 {
		return errors.New("Arguments required")
	}
	for _, arg := range args {
		if err := g.getter.Get(arg); err != nil {
			log.Printf("Error when trying to get from %s, error : %s", arg, err.Error())
		}
	}
	if err := g.setting.Sync(); err != nil {
		return err
	}
	return g.setting.SaveToFile(gompConf)
}

func (g *GetCmd) Init() cli.Command {
	return cli.Command{
		Name:   "get",
		Usage:  "Perform Go Get operation, will switch to git clone and git pull if a matching SSH Credential has been registered",
		Action: g.Run,
	}
}

func NewGetCmd(setting *gmp.Setting, getter *gmp.GetterManager) *GetCmd {
	instance := new(GetCmd)
	instance.setting = setting
	instance.getter = getter
	return instance
}
