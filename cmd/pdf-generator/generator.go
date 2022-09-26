package main

import (
	"github.com/molliechan/go-pdf-generator-alternative/internal/gotenberg"
	"github.com/molliechan/go-pdf-generator-alternative/internal/template"
	"github.com/molliechan/go-pdf-generator-alternative/internal/user"
)

type PDFService struct{}

const (
	hostName      = "http://localhost:3000"
	requestURL    = "/forms/chromium/convert/html"
	timeout       = 10
)

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (p PDFService) generatePDF(templ string, data interface{}, dest string) error {
	// parse template
	htmlBytes, err := template.ParseTemplate(templ, user.GetUser())
	if err != nil {
		return err
	}

	// create gotenberg client
	client := &gotenberg.Client{Hostname: hostName}

	// define client request
	req := gotenberg.NewRequest(requestURL)
	req.SetMargins([4]float64{1, 1, 1, 1})
	req.SetFormIndexFile(htmlBytes)

	err = client.Store(req, dest)
	if err != nil {
		return err
	}

	return nil
}
