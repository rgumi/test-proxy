package ctrl

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//Feedback stores the results of the tests in
type Feedback struct {
	Task             string `json:"task"`
	LastRunDate      string `json:"lastRunDate"`
	RunDurationInSec string `json:"runDurationInSec"`
	Status           string `json:"status"`
}

// TODO: make this fit regression & test => return []Feedback
func Run(taskType, task string, modTime time.Time) (fb []Feedback, err error) {
	if taskType == "test" {
		currentBasedir := filepath.Join(TestcaseDir, task)

		fb, err = validateXML(currentBasedir, task)
		if err != nil {
			return []Feedback{}, err
		}
	} else if taskType == "regression" {
		reRegressionFile := regexp.MustCompile("regression_.*") // TODO: add \.log maybe?

		file, err := FindNewestFile(RegressionLogDir, reRegressionFile, modTime)
		if err != nil {
			return []Feedback{}, err
		}

		fb, err = validateLog(RegressionLogDir, file.Name())
		if err != nil {
			return []Feedback{}, err
		}
	}

	return
}

func FileExists(basedir, filename string) (err error) {

	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if fi.Name() == filename {
			return nil
		}
	}
	return fmt.Errorf("File not found")
}

// returns the first row that has a pattern matching token
func getStringInFile(basedir, file, token string) (line string, err error) {

	f, err := os.Open(filepath.Join(basedir, file))
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), token) {
			// return after first match!
			return scanner.Text(), nil
		}
	}

	// check for errors != EOF
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}

	return "", fmt.Errorf("Could not find string in file")
}

//FindNewestFile returns ordered list of newest files
func FindNewestFile(dir string, pattern *regexp.Regexp, modTime time.Time) (file os.FileInfo, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			// check if file matches the search pattern
			if pattern.MatchString(fi.Name()) {
				// check if file was created after modTime and is therefore the logfile we need
				if fi.ModTime().After(modTime) {
					return fi, nil
				}
			}
		}
	}
	return nil, os.ErrNotExist
}
