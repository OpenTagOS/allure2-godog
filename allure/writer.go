package allure

import (
	"encoding/json"
	uuid5 "github.com/satori/go.uuid"
	"os"
)

type ReportWriter struct {
	dir      string
	tmpDir   string
	archiver *Archiver
}

func NewReportWriter(dir string) *ReportWriter {
	tmpDir := os.TempDir() + "/allure_godog/" + uuid5.NewV4().String() + "/"

	archiver := NewArchiver(tmpDir)

	return &ReportWriter{
		archiver: archiver,
		dir:      dir,
		tmpDir:   tmpDir,
	}
}

func (w *ReportWriter) Init() error {
	if err := os.MkdirAll(w.dir, 0775); err != nil {
		return err
	}

	return os.MkdirAll(w.tmpDir, 0775)
}

func (w *ReportWriter) WriteTestCaseResults(testCase *TestCase) error {
	fileName := testCase.UUID + "-result.json"

	return w.writeFile(testCase, fileName)
}

func (w *ReportWriter) WriteContainerResults(container *Container) error {
	fileName := container.UUID + "-container.json"

	if err := w.writeFile(container, fileName); err != nil {
		return err
	}

	archivePath := w.dir + "report.zip"

	return w.archiver.Zip(archivePath)
	//return os.RemoveAll(w.tmpDir)
}

func (w *ReportWriter) writeFile(data interface{}, fileName string) error {
	filePath := w.dir + fileName

	serialized, err := json.Marshal(data)
	if err != nil {
		return err
	}

	output, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer output.Close()

	if _, err := output.Write(serialized); err != nil {
		return err
	}

	return nil
}
