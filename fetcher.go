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
	targetDir := fmt.Sprintf("./vendor.gomp/%s", s.name)
	if err := s.PrepareDir(targetDir); err != nil {
		return err
	}

	if _, err := os.Lstat(targetDir + "/.git"); err != nil {
		return s.update(targetDir)
	}

	eCmd := exec.Command("git", []string{
		"clone",
		s.connString,
	}...)
	eCmd.Stdout = os.Stdout
	eCmd.Stderr = os.Stderr
	cleanDir, err := filepath.Abs(targetDir + "/..")
	if err != nil {
		return err
	}
	eCmd.Dir = cleanDir
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

func (b *SshFetcher) PrepareDir(name string) error {
	return os.MkdirAll(fmt.Sprintf("./vendor.gomp/%s", name), os.ModeDir)
}

func NewSshFetcher(name, username, host, basePath string) *SshFetcher {
	instance := new(SshFetcher)
	instance.name = name
	instance.connString = fmt.Sprintf("%s@%s:%s/%s", username, host, basePath, name)
	return instance
}