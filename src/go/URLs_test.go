package main

import (
	// "fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"tomr/src/go/ratom"
)

type NoLog int

func (NoLog) Write([]byte) (int, error) {
	return 0, nil
}

// * START OF RESPONSIVENESS * \\

// Function to get responsiveness metric score
func TestGetResponsiveness(t *testing.T) {
	ratom.DebugLogger = log.New(new(NoLog), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	ratom.InfoLogger = log.New(new(NoLog), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// testUrl := "github.com/hugoday/resume"
	// cloneRepo(testUrl)
	// resp := ratom.GetResp(testUrl)
	// responsiveness := ratom.GetResponsiveness(resp)
	// fmt.Println(responsiveness)

}

func TestRemoveScores(t *testing.T) {
	ratom.RemoveScores()
}

// * END OF RESPONSIVENESS * \\

// * START OF RAMP-UP TIME * \\

// Function to get ramp-up time metric scor
func TestGetRampUpTime(t *testing.T) {

	ratom.GetRampUpTime("github.com/hugoday/resume")
}

// * END OF RAMP-UP TIME * \\

// * START OF BUS FACTOR * \\

// Function to get bus factor metric score
func TestGetBusFactor(t *testing.T) {

	// getBusFactor("github.com/hugoday/resume")

}

// * END OF BUS FACTOR * \\

// * START OF Correctness * \\

// Function to get Correctness metric score
func TestGetCorrectness(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	ratom.GetCorrectness("not/a/url")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "Did not find issues file from api, invalid url: not/a/url\n" {
		t.Errorf("Expected %s, got %s", "Did not find issues file from api, invalid url: not/a/url", out)
	}
}

func TestRunRestApi(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	exitStatus := ratom.RunRestApi("not/a/url")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if exitStatus != 1 {
		t.Errorf("Expected %s, got %s", "exit status 1", out)
	}
}

func TestTeardownRestApi(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	ratom.TeardownRestApi()

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "" {
		t.Errorf("Expected %s, got %s", "nothing", out)
	}
}

func TestCalc_score(t *testing.T) {

	out := ratom.Calc_score("2", "2")

	if out != 0.5 {
		t.Errorf("Expected %s, got %f", "0.5", out)
	}
}

// * END OF Correctness * \\

// * START OF LICENSE COMPATABILITY * \\

// Function to get license compatibility metric score
func TestGetLicenseCompatibility(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	ratom.GetLicenseCompatibility("testURL")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "" {
		t.Errorf("Expected %s, got %s", "nothing", out)
	}
}

func TestGetLicenseCompatibility2(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	ratom.GetLicenseCompatibility("https://github.com/hugoday/resume")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "" {
		t.Errorf("Expected %s, got %s", "nothing", out)
	}
}

func TestSearchForLicenses(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	ratom.SearchForLicenses("./src")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "" {
		t.Errorf("Expected %s, got %s", "nothing", out)
	}
}

func TestCheckFileForLicense1(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	ratom.CheckFileForLicense("./src/does/not/exist")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "" {
		t.Errorf("Expected %s, got %s", "nothing", out)
	}
}

func TestCheckFileForLicense2(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	// var r repo
	// r.URL = "testUrl"
	// r.netScore = 7.4
	ratom.CheckFileForLicense("main.go")

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	// t.Errorf(string(out))
	// t.Errorf("hi")
	// t.Fail()
	if string(out) != "" {
		t.Errorf("Expected %s, got %s", "nothing", out)
	}
}

// * END OF LICENSE COMPATABILITY * \\

// * START OF REPO CLONING/REMOVING  * \\

func TestCloneRepo(t *testing.T) {
	ratom.CloneRepo("not/a/repo")
}

func TestClearRepoFolder(t *testing.T) {
	ratom.ClearRepoFolder()
}

// * END OF REPO CLONING/REMOVING  * \\

// * START OF STDOUT * \\

func TestPrintRepo(t *testing.T) {
	r := ratom.Repo{URL: "testRepo"}
	r.BusFactor = -1
	r.Correctness = -1
	r.LicenseCompatibility = -1
	r.RampUpTime = -1
	r.NetScore = -1
	ratom.PrintRepo(&r)
}

func TestRepoOUT(t *testing.T) {
	rescueStdout := os.Stdout
	read, w, _ := os.Pipe()
	os.Stdout = w

	var r ratom.Repo
	r.URL = "testUrl"
	r.NetScore = 7.4
	ratom.RepoOUT(&r)

	w.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	if string(out) != "{\"URL\":\"testUrl\", \"NET_SCORE\":7.40, \"RAMP_UP_SCORE\":0.00, \"Correctness_SCORE\":0.00, \"BUS_FACTOR_SCORE\":0.00, \"RESPONSIVE_MAINTAINER_SCORE\":0.00, \"LICENSE_SCORE\":0.00} \n" {
		t.Errorf("Expected %s, got %s", "{\"URL\":\"testUrl\", \"NET_SCORE\":7.40, \"RAMP_UP_SCORE\":0.00, \"Correctness_SCORE\":0.00, \"BUS_FACTOR_SCORE\":0.00, \"RESPONSIVE_MAINTAINER_SCORE\":0.00, \"LICENSE_SCORE\":0.00}", out)
	}

}

// * END OF STDOUT * \\

// * START OF SORTING * \\

func TestAddRepo(t *testing.T) {

}
