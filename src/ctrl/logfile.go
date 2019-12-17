package ctrl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	//RegressionLogDir is the directory to find the logfiles of the regressions in
	RegressionLogDir string
)

func validateLog(basedir, filename string) (fb []Feedback, err error) {

	f, err := os.Open(filepath.Join(basedir, filename))
	if err != nil {
		return []Feedback{}, err
	}
	log.Printf("Opened %s", filename)
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		testCase, lastRunDate, runDurationInSec, runResult, err := getInfoFromLog(scanner.Text())
		if err == nil {
			fb = append(fb, Feedback{
				testCase,
				lastRunDate,
				runDurationInSec,
				runResult})
		}
	}

	// array of Feedbacks
	return fb, nil
}

// runs regexp on a string to find the results of the corresponding tests
// => results for *regression*
func getInfoFromLog(row string) (testCase, lastRunDate, runDurationInSec, runResult string, err error) {

	reTestCase := regexp.MustCompile("Test case: \"([A-Za-z0-9]+)\";")
	if tmpTestCase := reTestCase.FindAllStringSubmatch(row, -1); len(tmpTestCase) > 0 {
		testCase = fmt.Sprintf("%s", tmpTestCase[0][1])
	} else {
		return "", "", "", "", fmt.Errorf("Could not find element in logfile")
	}

	reLastRunDate := regexp.MustCompile("Run date: \"([0-9]{1,2}.[0-9]{1,2}.[0-9]{4} [0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2})\";")
	if tmpLastRunDate := reLastRunDate.FindAllStringSubmatch(row, -1); len(tmpLastRunDate) > 0 {
		lastRunDate = fmt.Sprintf("%s", tmpLastRunDate[0][1])
	}
	// double escape \\(
	reRunDurationInSec := regexp.MustCompile("Run duration \\(sec\\): \"([0-9]+)\";")
	if tmpRunDurationInSec := reRunDurationInSec.FindAllStringSubmatch(row, -1); len(tmpRunDurationInSec) > 0 {
		runDurationInSec = fmt.Sprintf("%s", tmpRunDurationInSec[0][1])
	}

	reRunResult := regexp.MustCompile("Run result: \"([A-Z]+)\"")
	if tmpRunResult := reRunResult.FindAllStringSubmatch(row, -1); len(tmpRunResult) > 0 {
		runResult = fmt.Sprintf("%s", tmpRunResult[0][1])
	}

	return
}
