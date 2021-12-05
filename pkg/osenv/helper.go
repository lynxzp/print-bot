package osenv

import (
	"io"
	"net/http"
	"os"
	"path"
)

type Config struct {
	FormatsConvertibleToPDF []string `json:"DirectPrintableFormats"`
	TempPath                string
}

var Cfg Config

func InConvertibleToPDF(ext string) bool {
	for _, val := range Cfg.FormatsConvertibleToPDF {
		if (val == ext) && (val != "pdf") {
			return true
		}
	}
	return false
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func GenFilePath(docId string, prefix string, filename string) string {
	if (len(Cfg.TempPath) > 1) && (Cfg.TempPath[0] == '%') {
		cropped := Cfg.TempPath[1 : len(Cfg.TempPath)-1]
		Cfg.TempPath = os.Getenv(cropped)
	}
	if (len(Cfg.TempPath) > 1) && (Cfg.TempPath[0] == '$') {
		cropped := Cfg.TempPath[1:len(Cfg.TempPath)]
		Cfg.TempPath = os.Getenv(cropped)
	}
	return path.Join(Cfg.TempPath, docId+"_"+prefix+"_"+filename)
}
