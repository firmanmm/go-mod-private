package gomodprivate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type ModUpdater struct {
	prefixMessage  string
	postfixMessage string
	templater      *template.Template
	remover        *regexp.Regexp
}

func (m *ModUpdater) Update(repositories []string) error {
	file, err := os.OpenFile("go.mod", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	replaced := m.remover.ReplaceAllString(string(body), "")
	replaced = strings.TrimRight(replaced, "\n")
	file.Truncate(0)
	file.Seek(0, 0)
	return m.templater.Execute(file, map[string]interface{}{
		"GoModBody":    replaced,
		"Prefix":       m.prefixMessage,
		"Postfix":      m.postfixMessage,
		"Repositories": repositories,
	})
}

func NewModUpdater() *ModUpdater {
	instance := new(ModUpdater)
	templateData := `{{.GoModBody}}

{{.Prefix}}
//This is an auto generated section made by Go Mod Private
//For more information visit https://github.com/firmanmm/go-mod-private
//Please add *.gomp to your .gitignore since any .gomp files is meant to be used locally

require (
	{{ range $idx, $repo := .Repositories }}
	{{ $repo }} v0.0.0{{ end }}
	
)

replace (
	{{ range $idx, $repo := .Repositories }}
	{{ $repo }} v0.0.0 => ./.vendor.gomp/{{ $repo }}{{ end }}
	
)
{{.Postfix}}
`
	tmpl, err := template.New("mod.go").Parse(templateData)
	if err != nil {
		log.Fatalln(err.Error())
	}
	instance.templater = tmpl
	instance.prefixMessage = "//GO_MOD_PRIVATE_START"
	instance.postfixMessage = "//GO_MOD_PRIVATE_END"
	removerPattern := fmt.Sprintf(`%s([\s\S]*)%s`, instance.prefixMessage, instance.postfixMessage)
	replacerPattern := regexp.MustCompile("/")
	removerPattern = replacerPattern.ReplaceAllString(removerPattern, "\\/")
	instance.remover = regexp.MustCompile(removerPattern)
	return instance
}
