package template

import (
	"bytes"
	templatePackage "html/template"

	"github.com/aymerick/raymond"
)

type ITemplateInterface interface {
	RenderTemplate(html string, data interface{}) (string, error)
	RenderMandrillTemplate(template string, data interface{}) (string, error)
}

type templateService struct {
}

func (t templateService) RenderTemplate(html string, data interface{}) (string, error) {
	var templateBuffer bytes.Buffer

	template, err := templatePackage.New("invoice").Parse(html)

	if err != nil {
		return "", err
	}

	err = template.Execute(&templateBuffer, data)
	if err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}
func (t templateService) RenderMandrillTemplate(template string, data interface{}) (string, error) {
	html, err := raymond.Render(template, data)
	if err != nil {
		return "", err
	}

	return html, nil
}

func NewTemplateService() ITemplateInterface {
	return &templateService{}
}
