package helpers

import (
	"os"
	"path/filepath"
)

func GetStaticPath() (string, string) {
	pwd, _ := os.Getwd()
	staticPath := filepath.Join(pwd, "web", "static")
	htmlPath := filepath.Join(staticPath, "html", "*")
	return staticPath, htmlPath
}
