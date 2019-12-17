package ctrl

import (
	"fmt"
	"regexp"
)

var (
	//TestcaseDir is the host-directory where all testcases are stored
	TestcaseDir string
)

func validateXML(basedir, task string) (fb []Feedback, err error) {
	filename := task + ".xml"
	// check if file exists and return instantly if not
	if err := FileExists(basedir, filename); err != nil {
		return []Feedback{}, err
	}

	// get the whole *first* row in file which contains "RunInformation"
	runInformation, err := getStringInFile(basedir, filename, "RunInformation")
	if err != nil {
		return []Feedback{}, err
	}

	// array of Feedbacks
	lastRunDate, runDurationInSec, runResult := getInfoFromXML(runInformation)
	fb = append(fb, Feedback{
		task,
		lastRunDate,
		runDurationInSec,
		runResult})

	return fb, nil
}

// runs regexp on a string to find the results of the corresponding tests
// => results for *single tests*
func getInfoFromXML(runInformation string) (lastRunDate, runDurationInSec, runResult string) {

	reLastRunDate := regexp.MustCompile("lastRunDate=\"([0-9]{1,2}.[0-9]{1,2}.[0-9]{4} [0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2})\"")
	lastRunDate = fmt.Sprintf("%s", reLastRunDate.FindAllStringSubmatch(runInformation, -1)[0][1])

	reRunDurationInSec := regexp.MustCompile("runDurationInSec=\"([0-9]+)\"")
	runDurationInSec = fmt.Sprintf("%s", reRunDurationInSec.FindAllStringSubmatch(runInformation, -1)[0][1])

	reRunResult := regexp.MustCompile("runResult=\"([A-Z]+)\"")
	runResult = fmt.Sprintf("%s", reRunResult.FindAllStringSubmatch(runInformation, -1)[0][1])

	return
}
