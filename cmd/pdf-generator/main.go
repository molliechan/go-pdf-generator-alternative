package main

import (
	"log"

	"github.com/molliechan/go-pdf-generator-alternative/internal/user"
)

func main() {
	pdfService := NewPDFService()

	err := pdfService.generatePDF(
		"../../template/sample.gohtml",
		user.GetUser(),
		"sample.pdf",
	)

	if err != nil {
		log.Println(err)
	}

}
