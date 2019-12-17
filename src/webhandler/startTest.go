package webhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"../ctrl"
)

func (e event) start() (jsonObj []byte, err error) {

	cmd := ctrl.GetCmd(e.Type, e.Task)
	// Get current timestamp to find logfiles that were created during the runtime of process
	timeStamp := time.Now()
	// don't wait for process
	err = cmd.Start()
	if err != nil {
		log.Print("cmd.Start: ", err)
		return jsonObj, fmt.Errorf("Error starting process for '%s'", e.Type)
	}
	log.Printf("Running process for '%s'", e.Type)

	pid, err := ctrl.GetPidOfProcess()
	if err != nil {
		log.Print("GetPidOfProcess: ", err)
		return jsonObj, fmt.Errorf("Error getting pid of process for '%s'", e.Type)
	}
	log.Printf("Found pid of process: %d", pid)

	// actively wait for process to exit
	err = ctrl.CheckIfActive(pid)
	if err != nil {
		log.Print("CheckIfActive: ", err)
		return jsonObj, fmt.Errorf("Error checking for process '%d'", pid)
	}
	log.Printf("Process %d finished in %fs", pid, (time.Now().Sub(timeStamp)).Seconds())

	// ---
	// Get Feedback for task

	fb, err := ctrl.Run(e.Type, e.Task, timeStamp)
	if err != nil {
		log.Print("ValidationError: ", err)
		//TODO: error handling
		return jsonObj, err
	}
	currentResponse := Response{
		e.Type,
		e.Title,
		e.Task,
		fb}

	jsonObj, err = json.Marshal(currentResponse)
	return jsonObj, err
}
