package template

import (
	"bytes"
	"text/template"
)

func ParseTemplate(fileName string, data interface{}) ([]byte, error) {
	// use Go's default HTML template generation tools to generate your HTML
	templ, err := template.ParseFiles(fileName)

	if err != nil {
		return nil, err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err := templ.Execute(&body, data); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}