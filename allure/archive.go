package allure

import (
	"archive/zip"
	"io/ioutil"
	"os"
)

type Archiver struct {
	path string
}

func NewArchiver(path string) *Archiver {
	return &Archiver{
		path: path,
	}
}

func (a *Archiver) Zip(archivePath string) error {
	outFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}

	defer outFile.Close()

	w := zip.NewWriter(outFile)

	if err := a.addFiles(w); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (a *Archiver) addFiles(w *zip.Writer) error {
	files, err := ioutil.ReadDir(a.path)
	if err != nil {
		return err
	}

	for _, file := range files {
		dat, err := ioutil.ReadFile(a.path + file.Name())
		if err != nil {
			return err
		}

		f, err := w.Create("" + file.Name())
		if err != nil {
			return err
		}

		_, err = f.Write(dat)
		if err != nil {
			return err
		}
	}

	return nil
}
