package export

import (
	"crowsnest/internal/files"
	"crowsnest/internal/sqlite"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func WriteIStringToFile(iString sqlite.IString, outputFile string, fileType files.FileType) error {
	var data []byte
	var err error

	switch fileType {
	case files.JSON:
		data, err = json.MarshalIndent(iString, "", "  ")
	case files.XML:
		data, err = xml.MarshalIndent(iString, "", "  ")
	case files.YAML:
		data, err = yaml.Marshal(iString)
	case files.TEXT:
		data = []byte(iString.String())
	default:
		return err
	}

	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s.%s", outputFile, fileType.String())
	return os.WriteFile(filePath, data, 0644)
}
