package cmd

import (
	gmp "github.com/firmanmm/go-mod-private"
	"github.com/urfave/cli"
)

type AddCredentialCmd struct {
	setting *gmp.Setting
}

func (a *AddCredentialCmd) Run(ctx *cli.Context) error {
	gompConf := ctx.GlobalString("gomp")
	if err := a.setting.LoadFromFile(gompConf); err != nil {
		return err
	}
	host := ctx.String("host")
	user := ctx.String("user")
	base := ctx.String("base")

	if err := a.setting.AddCredential(host, user, base); err != nil {
		return err
	}

	return a.setting.SaveToFile(gompConf)
}

func (a *AddCredentialCmd) Init() cli.Command {
	return cli.Command{
		Name:   "add_credential",
		Usage:  "Add credential to be used for [get] command",
		Action: a.Run,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "host",
				Usage:    "Host to be used for connection to ssh, could be Domain or IP",
				Required: true,
			},
			cli.StringFlag{
				Name:     "user",
				Usage:    "Username to be used for connection to ssh",
				Required: true,
			},
			cli.StringFlag{
				Name:  "base",
				Usage: "Base path for traversal, user@host:[base]",
				Value: "",
			},
		},
	}
}

func NewAddCredentialCmd(setting *gmp.Setting) *AddCredentialCmd {
	instance := new(AddCredentialCmd)
	instance.setting = setting
	return instance
}
