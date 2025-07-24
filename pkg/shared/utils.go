package shared

import (
	"os"
	"strings"
)

func Contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func EnvironmentVariableValueWithDefault(variable, defaultValue string) string {
	value := EnvironmentVariableValue(variable)
	if value == "" {
		return defaultValue
	}
	return value
}

func EnvironmentVariableValue(variable string) (value string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == variable {
			value = pair[1]
			break
		}
	}
	return value
}

type FileModel struct {
	modelPath       string
	destinationPath string
	applyValues     map[string]string
}

func (f FileModel) Generate() (err error) {
	file, err := os.Open(f.modelPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	content := string(buffer)

	for key, value := range f.applyValues {
		content = strings.ReplaceAll(content, key, value)
	}

	//err = os.MkdirAll(f.destinationPath, os.ModePerm)
	//if err != nil {
	//	return err
	//}

	file, err = os.Create(f.destinationPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func GenerateFromModel(modelPath, destinationPath string, applyValues map[string]string) (err error) {
	fileModel := FileModel{modelPath, destinationPath, applyValues}
	return fileModel.Generate()
}
