package pl

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

//PROCESS_STATUS_ERRO - when the process is created. That is the first status
const PROCESS_STATUS_CREATED = "created"

//PROCESS_STATUS_ERRO - when the process ends well executed
const PROCESS_STATUS_SUCCESS = "success"

//PROCESS_STATUS_ERROR - when the process ends with a error
const PROCESS_STATUS_ERROR = "error"

//PROCESS_STATUS_ERRO - when the process is executing
const PROCESS_STATUS_EXECUTING = "executing"

//Process - The information about the process
type Process struct {
	ID      string
	Command string
	Status  string
	Output  string
	Pid     int
	Err     error
}

//ProcessList - list of process and functions to handle it
type ProcessList struct {
	list map[string]*Process
}

//add a process
func (pl *ProcessList) add(ID string, command string) error {

	if pl.list[ID] != nil {
		return fmt.Errorf("There is already a process using the ID: %s", ID)
	}

	if pl.list == nil {
		pl.list = make(map[string]*Process)
	}

	pl.list[ID] = &Process{
		ID:      ID,
		Command: command,
		Status:  PROCESS_STATUS_CREATED,
		Output:  "",
		Err:     nil,
	}
	return nil
}

//return a process
func (pl *ProcessList) getProcess(ID string) (exists bool, process Process) {
	if pl.list[ID] == nil {
		return false, Process{}
	}
	return true, *pl.list[ID]
}

//delete a process
func (pl *ProcessList) delete(ID string) (process Process, err error) {

	exist, process := pl.getProcess(ID)
	if exist == false {
		return Process{}, fmt.Errorf("There is no a process using the ID: %s", ID)
	}

	delete(pl.list, ID)
	return process, nil
}

//wait process execute
func (pl *ProcessList) wait(ID string) error {
	process := pl.list[ID]
	if process.Status == PROCESS_STATUS_CREATED {
		return fmt.Errorf("This process was never executed")
	}
	for {
		process = pl.list[ID]
		time.Sleep(time.Millisecond * 100)
		if process.Status == PROCESS_STATUS_ERROR || process.Status == PROCESS_STATUS_SUCCESS {
			return nil
		}
	}

}

//execute a process
func (pl *ProcessList) execute(ID string) error {

	if pl.list[ID] == nil {
		return fmt.Errorf("There is no a process using the ID: %s", ID)
	}

	process := pl.list[ID]
	if process.Status != PROCESS_STATUS_CREATED {
		return fmt.Errorf("The process is not in the status '%s'", PROCESS_STATUS_CREATED)
	}

	process.Status = PROCESS_STATUS_EXECUTING
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", process.Command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		process.Status = "error"
		process.Output = fmt.Sprintf("%v: %v", err, strings.TrimSpace(stderr.String()))
		return nil
	}

	process.Pid = cmd.Process.Pid

	go func() {
		err = cmd.Wait()
		if err != nil {
			process.Status = PROCESS_STATUS_ERROR
			process.Err = err
			process.Output = fmt.Sprintf("%v: %v", err, strings.TrimSpace(stderr.String()))
			return
		}

		process.Err = nil
		process.Status = PROCESS_STATUS_SUCCESS
		process.Output = strings.TrimSpace(stdout.String())
	}()

	return nil
}

//print a process
func (pl *ProcessList) printProcess() {
	for _, process := range pl.list {
		fmt.Printf("%v \n", process)
	}
}

// func main() {
// 	pl := ProcessList{}
// 	pl.add("1", "php /home/ronyldo12/projetos/golang/src/github.com/ronyldo12/proccess/teste.php 1")
// 	pl.add("2", "php /home/ronyldo12/projetos/golang/src/github.com/ronyldo12/proccess/teste.php 2")
// 	err := pl.add("2", "php /home/ronyldo12/projetos/golang/src/github.com/ronyldo12/proccess/teste.php")
// 	if err != nil {
// 		fmt.Printf("%v\n", err)
// 	}
// 	pl.execute("1")
// 	pl.execute("2")
// 	err = pl.execute("3")
// 	if err != nil {
// 		fmt.Printf("%v\n", err)
// 	}

// 	_, process := pl.getProcess("2")
// 	fmt.Printf("Get process: %v\n", process)

// 	for {
// 		pl.printProcess()
// 		time.Sleep(time.Second * 1)
// 		pl.delete("2")
// 		fmt.Print("\n========================\n")
// 	}
// }
