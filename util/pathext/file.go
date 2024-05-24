package pathext

import (
	"os"
	"path/filepath"
)

// FileExists returns true if the specified file is exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// LoadTemplate gets template content by the specified file.
func LoadTemplate(category, file, builtin string) (string, error) {
	//dir, err := GetTemplateDir(category)
	//if err != nil {
	//	return "", err
	//}
	dir := "" // 模板所在目录

	file = filepath.Join(dir, file)
	if !FileExists(file) {
		return builtin, nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
