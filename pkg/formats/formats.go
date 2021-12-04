package formats

type Config struct {
	ConvertibleToPDF []string `json:"DirectPrintable"`
}

var Available Config

func InConvertibleToPDF(ext string) bool {
	for _, val := range Available.ConvertibleToPDF {
		if (val == ext) && (val != "pdf") {
			return true
		}
	}
	return false
}
