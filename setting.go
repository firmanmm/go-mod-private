package gomodprivate

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
)

type SettingData struct {
	SshCredentials      []*SshCredential
	PrivateRepositories []string
}

func (s *SettingData) GetMatchingCredential(request string) *SshCredential {
	lastMatchIdx := -1
	lastLength := 0
	for idx, credential := range s.SshCredentials {
		if lastLength >= len(credential.Matcher) {
			continue
		}
		isMatch, err := regexp.MatchString(credential.Matcher, request)
		if err != nil {
			log.Fatalln(err.Error())
		}
		if isMatch {
			lastMatchIdx = idx
		}
	}
	if lastMatchIdx < 0 {
		return nil
	}
	return s.SshCredentials[lastMatchIdx]
}

func NewSettingData() *SettingData {
	instance := new(SettingData)
	instance.SshCredentials = make([]*SshCredential, 0)
	instance.PrivateRepositories = make([]string, 0)
	return instance
}

type IModUpdater interface {
	Update([]string) error
}

type Setting struct {
	modUpdater IModUpdater
	data       *SettingData
}

func (s *Setting) GetMatchingCredential(name string) *SshCredential {
	return s.data.GetMatchingCredential(name)
}

func (s *Setting) AddCredential(matcher, host, username, basePath string) error {
	s.data.SshCredentials = append(s.data.SshCredentials, &SshCredential{
		Matcher:  matcher,
		Host:     host,
		Username: username,
		BasePath: basePath,
	})
	sort.Sort(sshCredentialSort(s.data.SshCredentials))
	return nil
}

func (s *Setting) AddRepository(name string) error {
	for _, repo := range s.data.PrivateRepositories {
		if repo == name {
			return nil
		}
	}
	s.data.PrivateRepositories = append(s.data.PrivateRepositories, name)
	return nil
}

func (s *Setting) Sync() error {
	return s.modUpdater.Update(s.data.PrivateRepositories)
}

func (s *Setting) GetAllRepositories() []string {
	return s.data.PrivateRepositories
}

func (s *Setting) SaveToFile(fileName string) error {
	body, err := json.MarshalIndent(s.data, "", "    ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile("mod.gomp", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(body)
	return err
}

func (s *Setting) LoadFromFile(fileName string) error {
	_, err := os.Lstat(fileName)
	if err != nil {
		return s.SaveToFile(fileName)
	}
	file, err := os.OpenFile("mod.gomp", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	settingData := NewSettingData()
	s.data = settingData
	return json.Unmarshal(body, settingData)
}

func NewSetting() *Setting {
	instance := new(Setting)
	instance.data = NewSettingData()
	instance.modUpdater = NewModfileModUpdater()
	return instance
}

type SshCredential struct {
	Matcher  string
	Host     string
	Username string
	BasePath string
}
