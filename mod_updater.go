package gomodprivate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
	templateData := `
{{.GoModBody}}

{{.Prefix}}
//This is an auto generated section made by Go Mod Private
//For more information visit https://github.com/firmanmm/go-mod-private
//Please add vendor.gomp to your .gitignore

replace (
	{{ range $idx, $repo := .Repositories }}
	{{ $repo }} => ./vendor.gomp/{{ $repo }}
	{{ end }}
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
	instance.remover = regexp.MustCompile(fmt.Sprintf("%s(.*)%s", instance.prefixMessage, instance.postfixMessage))
	return instance
}
