package cmd

import (
	"fmt"

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
	matcher := ctx.String("pattern")
	if len(matcher) == 0 {
		matcher = fmt.Sprintf("%s(.*)", host)
	}

	if err := a.setting.AddCredential(matcher, host, user, base); err != nil {
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
			},
			cli.StringFlag{
				Name:  "pattern",
				Usage: "Pattern to be used to match dependency, Default will be [host](.*)",
			},
		},
	}
}

func NewAddCredentialCmd(setting *gmp.Setting) *AddCredentialCmd {
	instance := new(AddCredentialCmd)
	instance.setting = setting
	return instance
}
