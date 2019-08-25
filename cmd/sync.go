package cmd

import (
	gmp "github.com/firmanmm/go-mod-private"
	"github.com/urfave/cli"
)

type SyncCmd struct {
	setting *gmp.Setting
	syncMgr *gmp.SyncManager
}

func (s *SyncCmd) Run(ctx *cli.Context) error {
	gompConf := ctx.GlobalString("gomp")
	if err := s.setting.LoadFromFile(gompConf); err != nil {
		return err
	}
	return s.syncMgr.Sync()
}

func (s *SyncCmd) Init() cli.Command {
	return cli.Command{
		Name:   "sync",
		Usage:  "Synchronize vendor.gomp and go.mod to mod.gomp",
		Action: s.Run,
	}
}

func NewSyncCmd(syncMgr *gmp.SyncManager) *SyncCmd {
	instance := new(SyncCmd)
	instance.syncMgr = syncMgr
	return instance
}
