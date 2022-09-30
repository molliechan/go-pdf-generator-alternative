package main

import (
	"fmt"
	"os"
	"time"

	"github.com/molliechan/go-pdf-generator-alternative/internal/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := initLogger()
	pdfService := NewPDFService(logger)

	logger.Info("Start")

	fileSuffix := "0"
	if len(os.Args) > 1 {
		fileSuffix = os.Args[1]
	}

	output := fmt.Sprintf("output/sample-%s.pdf", fileSuffix)
	pdfService.generatePDF(
		"../../template/sample.gohtml",
		user.GetUser(),
		output,
	)

	logger.Info("Completed")
}

func initLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, _ := config.Build()
	return logger.Sugar()
}

/** result: 24 mins for 1000 files, asyn call for ~60 before crash
func stressTest(pdfService *PDFService) {
	count := 1001
	for i := 1; i < count; i++ {
		output := fmt.Sprintf("sample-%d.pdf", i)
		pdfService.generatePDF(
			"../../template/sample.gohtml",
			user.GetUser(),
			output,
		)

	}
}
*/

