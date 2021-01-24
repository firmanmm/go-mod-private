package gomodprivate

import (
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

type ModfileModUpdater struct {
	targetFile string
}

func (g *ModfileModUpdater) Update(repositories []string) error {
	repoMap := map[string]bool{}
	for _, repo := range repositories {
		repoMap[repo] = true
	}
	targetFile := g.targetFile
	if len(targetFile) == 0 {
		targetFile = "go.mod"
	}
	content, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return err
	}
	parsedFile, err := modfile.Parse(targetFile, content, nil)
	if err != nil {
		return err
	}
	for _, req := range parsedFile.Require {
		mod := req.Mod
		if _, ok := repoMap[mod.Path]; ok {
			repoMap[mod.Path] = false
		}
	}
	for _, path := range repositories {
		packageName, tag, err := _ExtractTag(path)
		if err != nil {
			return err
		}
		if len(tag) == 0 {
			tag = "v0.0.0"
		}
		if err := parsedFile.AddRequire(packageName, tag); err != nil {
			return err
		}
		if err := parsedFile.AddReplace(packageName, tag, "./.vendor.gomp/"+path, ""); err != nil {
			return err
		}

	}
	parsedFile.SortBlocks()
	result, err := parsedFile.Format()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(targetFile, result, 0666); err != nil {
		return err
	}
	return nil
}

func (g *ModfileModUpdater) SetTargetFile(targetFile string) {
	g.targetFile = targetFile
}

func NewModfileModUpdater() *ModfileModUpdater {
	return &ModfileModUpdater{}
}
