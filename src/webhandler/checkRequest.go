package webhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../ctrl"
)

func getHeaderInfo(r *http.Request) {
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
}

func checkRequestMethod(r *http.Request, method string) (err error) {
	if r.Method != method {
		return fmt.Errorf("Expected Method %v. Got %v", method, r.Method)
	}
	return nil
}

func checkType(currentType string) bool {
	var validInputs = []string{"test", "regression"}
	for _, possType := range validInputs {
		if possType == currentType {
			return true
		}
	}
	return false
}

func checkNotEmpty(str string) bool {
	if len(str) > 0 {
		return true
	}
	return false
}

// not supported right now!
func checkIfTaskExists(currentTask string) bool {
	err := ctrl.FileExists(ctrl.TestcaseDir, currentTask)
	if err != nil {
		return false
	}
	return true
}

func validateJSON(jsonObj []byte) (newEvent event, err error) {
	// 	jsonObj should contains all necessary field such as:
	// 		type, which can contain (test, regressin)
	//		title, which should contains a title for the task
	//		task, which should contain (a) name(s) of the tests
	err = json.Unmarshal(jsonObj, &newEvent)
	if err != nil {
		log.Print("Could not unmarshal JSON object")
		return event{}, fmt.Errorf("Could not unmarshal JSON object")
	}

	if !checkType(newEvent.Type) {
		return newEvent, fmt.Errorf("Wrong type: Got '%s'. Expected test|regression", newEvent.Type)
	}
	if !checkNotEmpty(newEvent.Title) {
		return newEvent, fmt.Errorf("Title is empty")
	}
	if !checkNotEmpty(newEvent.Task) {
		return newEvent, fmt.Errorf("Task is empty")
	}
	// checkIfTaskExists(newEvent.Task){}
	return
}
