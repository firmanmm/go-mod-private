package gomodprivate

type GetterManager struct {
	setting *Setting
}

func (g *GetterManager) Get(name string) error {
	sshCredential := g.setting.GetMatchingCredential(name)
	var fetcher iPackageFetcher
	if sshCredential == nil {
		fetcher = NewGoFetcher(name)
	} else {
		fetcher = NewSshFetcher(
			name,
			sshCredential.Username,
			sshCredential.Host,
			sshCredential.BasePath)
	}
	if err := fetcher.Fetch(); err != nil {
		return err
	}
	if sshCredential != nil {
		if err := g.setting.AddRepository(name); err != nil {
			return err
		}
	}
	return nil
}

func NewGetterManager(setting *Setting) *GetterManager {
	instance := new(GetterManager)
	instance.setting = setting
	return instance
}
