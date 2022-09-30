package main

import (
	"fmt"

	"github.com/molliechan/go-pdf-generator-alternative/internal/gotenberg"
	"github.com/molliechan/go-pdf-generator-alternative/internal/template"
	"github.com/molliechan/go-pdf-generator-alternative/internal/user"
)

type PDFService struct {
	logger Logger
}

type Logger interface {
	Error(...interface{})
	Info(...interface{})
}

const (
	hostName   = "http://localhost:3000"
	requestURL = "/forms/chromium/convert/html"
	timeout    = 10
)

func NewPDFService(logger Logger) *PDFService {
	return &PDFService{logger: logger}
}

func (p PDFService) generatePDF(templ string, data interface{}, dest string) bool {

	p.logger.Info("started to parse template...")

	// parse template
	htmlBytes, err := template.ParseTemplate(templ, user.GetUser())
	if err != nil {
		p.logger.Error(
			"Fail to parse template",
			err,
		)
		return false
	}

	p.logger.Info("started to convert pdf...")
	
	// create gotenberg client
	client := &gotenberg.Client{Hostname: hostName}

	// define client request
	req := gotenberg.NewRequest(requestURL)
	req.SetMargins([4]float64{1, 1, 1, 1})
	req.SetFormIndexFile(htmlBytes)

	err = client.Store(req, dest)
	if err != nil {
		p.logger.Error(
			"Fail to convert and store PDF",
			err,
		)
		return false
	}

	p.logger.Info(
		fmt.Sprintf("Successfully convert to %s", dest),
	)
	return true
}
