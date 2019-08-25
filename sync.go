package gomodprivate

import (
	"errors"
	"strings"
)

type SyncManager struct {
	getter  *GetterManager
	setting *Setting
}

func (s *SyncManager) Sync() error {
	errorList := make([]string, 0)
	repositories := s.setting.GetAllRepositories()
	for _, repository := range repositories {
		if err := s.getter.Get(repository); err != nil {
			errorList = append(errorList, err.Error())
		}
	}
	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	}
	return nil
}

func NewSyncManager(setting *Setting, getter *GetterManager) *SyncManager {
	instance := new(SyncManager)
	instance.getter = getter
	instance.setting = setting
	return instance
}
