package main

import (
	"os/exec"
	"newsolution/shared/parameters"
	"newsolution/shared/error"
)

func main() {

	// createDot("confs.artefacts", "")

	createDot("confs.artefacts", "requestor")

}

func createDot(cspfile string, process string) {

	parser := parameters.DIR_CSPARSER + "/" + "CSParser.jar"

	inputFile := parameters.DIR_CSPARSER + "/" + cspfile

	exec.Command("/usr/bin/java", parameters.JAR_COMMAND, parser, inputFile, process)
	out, err := exec.Command("java").Output()
	if err != nil {
		myError := error.MyError{Source: "CSParser", Message: "Problem in creating .dot file"}
		myError.ERROR()
	}

	println(out)
}
