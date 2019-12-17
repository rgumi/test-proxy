package ctrl

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/process"
)

var (
	// working: "C:\Program Files\Cygwin\bin\mintty.exe" -h never -e /bin/bash -l -c java.sh de.telekom.irrp.StartTestcase -t <test-case>
	test = []string{
		filepath.Clean("C:/Program Files/Cygwin/bin/mintty.exe"),
		"-h",
		"never",
		"-e",
		"/bin/bash",
		"-l",
		"-c",
		"java.sh de.telekom.irrp.StartTestcase -t",
	}

	// working: "C:\Program Files\Cygwin\bin\mintty.exe" -h never -e /bin/bash -l -c java.sh de.telekom.irrp.StartRegression -f <regression-file>.conf
	regression = []string{
		filepath.Clean("C:/Program Files/Cygwin/bin/mintty.exe"),
		"-h",
		"never",
		"-e",
		"/bin/bash",
		"-l",
		"-c",
		"java.sh de.telekom.irrp.StartRegression -f", //TODO: Test this!
	}

	//ProcessImage is the Imagename to search the tasklist for
	ProcessImage = "mintty.exe"
	//WindowTitle is the Name of the Bash window that runs the task
	WindowTitle string
	MaxTries    = 10
)

//GetPidOfProcess searched for a running imagename-Process and returns its PID
//note that it only searches for the first occurrence
func GetPidOfProcess() (pid int32, err error) {
	reg := regexp.MustCompile("[0-9]+")
	var tmpPid int64

	log.Printf("Please make sure and only one instance of %s is running and it is the one that you want to monitor!", ProcessImage)

	// Loop until a process matching the filters is found or the MaxTries are reached!
	for i := 1; i < MaxTries; i++ {

		// debug purpose => TODO: delete
		log.Print("Checking for: tasklist", " /FI ", fmt.Sprintf("IMAGENAME eq %s", ProcessImage), " /FI ", fmt.Sprintf("WindowTitle eq %s", WindowTitle))

		// tasklist returns a list of all running processes
		// FI sets a filter for tasklist
		//	IMAGENAME is the name of the executable that was started
		//	WindowTitle is the name of the bash windows of the process
		// Command().Output() pipes stdout of process to var out => check if a PID is returned
		out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", ProcessImage), "/FI", fmt.Sprintf("WindowTitle eq %s", WindowTitle)).Output()
		if err != nil {
			return 0, err
		}

		// if a process matching the filter is found, get its PID
		if tmp := fmt.Sprintf("%s", reg.Find(out)); tmp != "" {

			// parse PID from string to int64 value
			tmpPid, _ = strconv.ParseInt(tmp, 10, 32)

			// process.NewProcess in CheckIfActive needs a int32 value
			// => cast int64 to int32
			return int32(tmpPid), nil
		}
		// actively wait until a process matching the filter is found
		time.Sleep(500 * time.Millisecond)
	}

	// Could not find process matching filter in time
	return 0, fmt.Errorf("Failed to get pid with regex")
}

//CheckIfActive is actively waiting for the given process by PID to stop running
// returns if process has exited
func CheckIfActive(pid int32) (err error) {
	// create an process-oberserver for the PID
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}

	// Loop until the process exits
	for {

		// can this run into an error?
		if b, err := p.IsRunning(); err == nil {
			if b == false {
				return nil
			}
			// actively wait ...
			time.Sleep(500 * time.Millisecond)
		} else {
			return err
		}
	}
}

//GetCmd returns a exec.Command depending on the taskType of the event
func GetCmd(taskType string, task string) (cmd *exec.Cmd) {

	// array to copy test or regression cmds into (see global variables)
	var cpy []string

	// check task type and set cmd accordingly
	if taskType == "test" {

		// WindowsTitle is needed to identify the process when it has started to get its PID
		WindowTitle = ""

		// loop over array and build string of its elemenets
		for _, val := range test[4:] {
			WindowTitle += val + " "
		}

		// complete the WindowTitle by adding the task (" " before task is set by loop above)
		WindowTitle += task

		// copy command template and add task to it
		cpy = make([]string, len(test))
		copy(cpy, test)

		// adding " " is important otherwise the WindowTitle will not match
		cpy[len(cpy)-1] += " "
		cpy[len(cpy)-1] += task

		// same as above but for the regression type
	} else if taskType == "regression" {

		WindowTitle = ""

		for _, val := range regression[4:] {
			WindowTitle += val + " "
		}

		WindowTitle += task

		// copy command template and add task to it
		cpy = make([]string, len(regression))
		copy(cpy, regression)
		cpy[len(cpy)-1] += " "
		cpy[len(cpy)-1] += task
	}

	// finally build the Command and return it
	cmd = exec.Command(cpy[0], cpy[1:]...)
	return
}
