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
	prefixMessage         string
	postfixMessage        string
	requireMessage        string
	templater             *template.Template
	remover               *regexp.Regexp
	requireMatcher        *regexp.Regexp
	requireRemover        *regexp.Regexp
	requirePackageRemover *regexp.Regexp
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
	requireText := m.requireMatcher.FindString(replaced)
	replaced = m.requireMatcher.ReplaceAllString(replaced, "")
	requireText = m.requireRemover.ReplaceAllString(requireText, "")
	requireText = m.requirePackageRemover.ReplaceAllString(requireText, "")
	file.Truncate(0)
	file.Seek(0, 0)
	return m.templater.Execute(file, map[string]interface{}{
		"GoGetRepository": requireText,
		"GoModBody":       replaced,
		"Prefix":          m.prefixMessage,
		"Postfix":         m.postfixMessage,
		"Requirefix":      m.requireMessage,
		"Repositories":    repositories,
	})
}

func NewModUpdater() *ModUpdater {
	instance := new(ModUpdater)
	templateData := `{{.GoModBody}}
require (
	{{.GoGetRepository}}
	{{ range $idx, $repo := .Repositories }}
	{{ $repo }} v0.0.0 {{ .Requirefix }}{{ end }}
	
)
{{.Prefix}}
//This is an auto generated section made by Go Mod Private
//For more information visit https://github.com/firmanmm/go-mod-private
//Please add *.gomp to your .gitignore since any .gomp files is meant to be used locally

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
	instance.requireMessage = `//GO_MOD_PRIVATE_REQUIRE`
	removerPattern := fmt.Sprintf(`(%s([\s\S]*)%s)|([\n]{2,})`, instance.prefixMessage, instance.postfixMessage)
	replacerPattern := regexp.MustCompile("/")
	removerPattern = replacerPattern.ReplaceAllString(removerPattern, "\\/")
	instance.requireMatcher = regexp.MustCompile(`^require[ ]+\([\sa-zA-Z0-9\/\-.]+\)$`)
	instance.requireRemover = regexp.MustCompile(`(require[ ]+\(|\))|(\t+)`)
	requirePackageRemoverPattern := fmt.Sprintf(`(.*)%s[\s]+`, instance.requireMessage)
	requirePackageRemoverPattern = replacerPattern.ReplaceAllString(requirePackageRemoverPattern, "\\/")
	instance.requirePackageRemover = regexp.MustCompile(requirePackageRemoverPattern)
	instance.remover = regexp.MustCompile(removerPattern)
	return instance
}
