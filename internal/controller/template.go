package controller

import (
	"embed"
	"html/template"
	"strings"
	"sync"
)

//go:embed template/*.html
var fsTemplate embed.FS

type templateEmailPasscodeData struct {
	LangCode      string
	AppTitle      string
	LabelPasscode string
	Passcode      string
	Subject       string
}

var templateEmailPasscode *template.Template

var mu sync.Mutex

// TemplateEmailPasscode get template
func TemplateEmailPasscode() *template.Template {

	tmpl := templateEmailPasscode

	if tmpl != nil {
		return tmpl
	}
	// may be use mutex
	mu.Lock()
	defer mu.Unlock()

	tmpl = templateEmailPasscode // recheck

	if tmpl == nil {

		{

			data, err := fsTemplate.ReadFile("template/email_passcode.html")
			if err != nil {
				panic(err)
			}
			// create cached
			tmpl, err = template.New("main").Parse(string(data))
			if err != nil {
				panic(err)
			}

		}
		{
			// test cached
			bu := strings.Builder{}
			err := tmpl.Execute(&bu, templateEmailPasscodeData{})
			if err != nil {
				panic(err)
			}
		}

		templateEmailPasscode = tmpl

	}

	return tmpl
}
