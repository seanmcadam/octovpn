package settings

import (
	"os"
	"path"
	"testing"

	"github.com/seanmcadam/loggy"
)

func TestReadConfig_compile(t *testing.T) {

	testfiledir := "jsontestfile"

	errorFiles, err := ListFilesWithPattern(testfiledir, "e*.json")
	if err != nil {
		t.Fatalf("Cannot read directory %s", testfiledir)
	}
	for _, file := range errorFiles {
		if _, err := ReadConfig(testfiledir + "/" + file); err == nil {
			t.Errorf("ERROR Test:%s FAIL", file)
		} else {
			loggy.Infof("ERROR Test:%s: PASS:%s", file, err)
		}
	}

	goodFiles, err := ListFilesWithPattern(testfiledir, "g*.json")
	if err != nil {
		t.Fatalf("Cannot read directory %s", testfiledir)
	}
	for _, file := range goodFiles {
		if _, err := ReadConfig(testfiledir + "/" + file); err != nil {
			t.Errorf("GOOD Test:%s FAIL:%s", file, err)
		} else {
			loggy.Infof("GOOD Test:%s: PASS", file)
		}
	}
}

func ListFilesWithPattern(directory string, pattern string) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var matchingFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue // skip directories
		}

		match, err := path.Match(pattern, file.Name())
		if err != nil {
			return nil, err
		}

		if match {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}

	return matchingFiles, nil
}
