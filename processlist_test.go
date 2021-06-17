package pl

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	pl := ProcessList{}
	//the first
	err := pl.add("id1", "la -la")
	if err != nil {
		t.Errorf("The error should be nill, %v given", err)
	}
	exist, _ := pl.getProcess("id1")
	if exist == false {
		t.Error("The process id1 doesn't exist in the list")
	}
	err = pl.add("id1", "la -la")
	if err == nil {
		t.Error("The error shouldn't be nill beacuse there is already a process using the id 'id1'")
	}
}

func TestGetProcess(t *testing.T) {
	pl := ProcessList{}
	pl.add("id1", "la -la")

	exist, process := pl.getProcess("id1")
	if exist == false {
		t.Error("The process id1 should exist")
	}
	if process.ID != "id1" {
		t.Errorf("The process ID should be 'id1', given %s", process.ID)
	}

	exist, _ = pl.getProcess("id_dont_exist")
	if exist != false {
		t.Error("The process id1 shouldn't exist")
	}
}

func TestExecute(t *testing.T) {
	dir, err := os.Getwd()

	if err != nil {
		t.Error("Error to get the current path")
	}

	pl := ProcessList{}
	commandSuccess := dir + "/program.sh"

	pl.add("teste_success1", commandSuccess)
	pl.execute("teste_success1")
	pl.wait("teste_success1")
	_, process := pl.getProcess("teste_success1")

	if process.Status != PROCESS_STATUS_SUCCESS {
		t.Errorf("The process should finish with the status %s", PROCESS_STATUS_SUCCESS)
	}
	if process.Output != "It was well executed" {
		t.Errorf("The process should finish with the output 'It was well executed', '%s' given", process.Output)
	}

	commandError := dir + "/program.sh error"
	pl.add("teste_error1", commandError)
	pl.execute("teste_error1")
	pl.wait("teste_error1")
	_, processErr := pl.getProcess("teste_error1")

	if processErr.Status != PROCESS_STATUS_ERROR {
		t.Errorf("The process should finish with the status %s", PROCESS_STATUS_ERROR)
	}
	if processErr.Output != "exit status 1: It was finished with a error" {
		t.Errorf("The process should finish with the output 'exit status 1: It was finished with a error', '%s' given", processErr.Output)
	}

}
