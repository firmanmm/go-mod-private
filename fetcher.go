package gomodprivate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type iPackageFetcher interface {
	Fetch() error
}

type GoFetcher struct {
	name string
}

func (g *GoFetcher) Fetch() error {
	eCmd := exec.Command("go", []string{
		"get",
		g.name,
	}...)
	eCmd.Stdout = os.Stdout
	eCmd.Stderr = os.Stderr
	return eCmd.Run()
}

func NewGoFetcher(name string) *GoFetcher {
	instance := new(GoFetcher)
	instance.name = name
	return instance
}

type SshFetcher struct {
	name       string
	connString string
}

func (s *SshFetcher) Fetch() error {
	targetDir := fmt.Sprintf("./.vendor.gomp/%s", s.name)

	if _, err := os.Lstat(targetDir + "/.git"); err == nil {
		return s.update(targetDir)
	}

	cleanDir, err := filepath.Abs(targetDir + "/..")
	if err != nil {
		return err
	}

	packageName, tag, err := _ExtractTag(s.name)
	if err != nil {
		return err
	}

	if err := s._Fetch(packageName, tag, cleanDir); err != nil {
		return err
	}

	return nil
}

func (s *SshFetcher) _Fetch(name, tag, cleanDir string) error {

	cmdParam := make([]string, 0, 6)
	cmdParam = append(cmdParam,
		"clone",
		"--depth",
		"1")
	if len(tag) > 0 {
		cmdParam = append(cmdParam,
			"--branch", tag)
	}
	cmdParam = append(cmdParam, s.connString+name, cleanDir+"@"+tag)
	eCmd := exec.Command("git", cmdParam...)
	eCmd.Stdout = os.Stdout
	eCmd.Stderr = os.Stderr
	return eCmd.Run()
}

func (s *SshFetcher) update(dir string) error {
	eCmd := exec.Command("git", []string{
		"pull",
		"--rebase",
	}...)
	eCmd.Stdout = os.Stdout
	eCmd.Stderr = os.Stderr
	eCmd.Dir = dir
	return eCmd.Run()
}

func NewSshFetcher(name, username, host, basePath string) *SshFetcher {
	instance := new(SshFetcher)
	instance.name = name
	instance.connString = fmt.Sprintf("%s@%s:%s/", username, host, basePath)
	return instance
}
