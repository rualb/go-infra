package controller

import (
	"embed"
	"html/template"
	"strings"
	"sync"
)

//go:embed template/*.html
var fsTemplate embed.FS

type templateEmailSecretCodeData struct {
	LangCode        string
	AppTitle        string
	LabelSecretCode string
	SecretCode      string
	Subject         string
}

var templateEmailSecretCode *template.Template

var mu sync.Mutex

// TemplateEmailSecretCode get template
func TemplateEmailSecretCode() *template.Template {

	tmpl := templateEmailSecretCode

	if tmpl != nil {
		return tmpl
	}
	// may be use mutex
	mu.Lock()
	defer mu.Unlock()

	tmpl = templateEmailSecretCode // recheck

	if tmpl == nil {

		{

			data, err := fsTemplate.ReadFile("template/email_secret_code.html")
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
			err := tmpl.Execute(&bu, templateEmailSecretCodeData{})
			if err != nil {
				panic(err)
			}
		}

		templateEmailSecretCode = tmpl

	}

	return tmpl
}
