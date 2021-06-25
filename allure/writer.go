package allure

import (
	"encoding/json"
	"os"
)

type ReportWriter struct {
	dir string
}

func NewReportWriter(dir string) *ReportWriter {
	return &ReportWriter{
		dir: dir,
	}
}

func (w *ReportWriter) Init() error {
	if err := os.RemoveAll(w.dir); err != nil {
		return err
	}

	if err := os.MkdirAll(w.dir, 0775); err != nil {
		return err
	}

	return nil
}

func (w *ReportWriter) WriteTestCaseResults(testCase *TestCase) error {
	fileName := testCase.UUID + "-result.json"

	return w.writeFile(testCase, fileName)
}

func (w *ReportWriter) WriteContainerResults(container *Container) error {
	fileName := container.UUID + "-container.json"

	return w.writeFile(container, fileName)
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
