package prompt

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

// PromptEngine 管理所有的 Prompt 模板
type PromptEngine struct {
	templates *template.Template
	funcMap   template.FuncMap
}

func New() *PromptEngine {
	funcMap := template.FuncMap{
		"json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
		"indent": func(s string, spaces int) string {
			prefix := strings.Repeat(" ", spaces)
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				lines[i] = prefix + line
			}
			return strings.Join(lines, "\n")
		},
		"codeblock": func(lang, code string) string {
			return "```" + lang + "\n" + code + "\n```"
		},
		"add": func(a, b int) int { return a + b },
	}

	pe := &PromptEngine{funcMap: funcMap}
	pe.templates = template.New("prompt").Funcs(funcMap)
	return pe
}

// AddTemplate 注册一个模板
func (pe *PromptEngine) AddTemplate(name, tmpl string) error {
	_, err := pe.templates.New(name).Parse(tmpl)
	return err
}

// Render 渲染模板
func (pe *PromptEngine) Render(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := pe.templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
