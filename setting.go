package gomodprivate

import (
	"encoding/json"
	"errors"
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

type Setting struct {
	modUpdater *ModUpdater
	data       *SettingData
}

func (s *Setting) GetMatchingCredential(name string) *SshCredential {
	return s.data.GetMatchingCredential(name)
}

func (s *Setting) AddCredential(matcher, host, username, basePath string) error {
	if idx := sort.Search(len(s.data.SshCredentials), func(i int) bool {
		return s.data.SshCredentials[i].Matcher >= matcher
	}); idx < len(s.data.SshCredentials) && s.data.SshCredentials[idx].Matcher == matcher {
		return errors.New("Duplicate credential with the same matcher, skipping")
	}
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
	if idx := sort.SearchStrings(s.data.PrivateRepositories, name); idx < len(s.data.PrivateRepositories) && s.data.PrivateRepositories[idx] == name {
		return errors.New("Repository already exist")
	}
	s.data.PrivateRepositories = append(s.data.PrivateRepositories, name)
	sort.Sort(stringSort(s.data.PrivateRepositories))
	return s.modUpdater.Update(s.data.PrivateRepositories)
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
	return instance
}

type SshCredential struct {
	Matcher  string
	Host     string
	Username string
	BasePath string
}
