package ctrl

import (
	"fmt"
	"testing"
)

func TestGetInfoFromXML(t *testing.T) {
	testLastRunDate := "13.11.2019 9:15:45"
	testRunDurationInSec := "23"
	testRunResult := "OK"
	runInformation := fmt.Sprintf(`lastRunDate="%s" runDurationInSec="%s" runResult="%s"`, testLastRunDate, testRunDurationInSec, testRunResult)

	lastRunDate, runDurationInSec, runResult := getInfoFromXML(runInformation)
	if lastRunDate != testLastRunDate {
		t.Errorf("getInfoFromXML failed to get lastRunDate correct")
	}
	if runDurationInSec != testRunDurationInSec {
		t.Errorf("getInfoFromXML failed to get runDurationInSec correct")
	}
	if runResult != testRunResult {
		t.Errorf("getInfoFromXML failed to get runResult correct")
	}
}
