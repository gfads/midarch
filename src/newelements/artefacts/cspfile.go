package artefacts

import (
	"shared/parameters"
	"os"
	"bufio"
	"errors"
	"shared/shared"
	"strings"
	"framework/configuration/commands"
	"os/exec"
	"fmt"
)

type CSPFile struct {
	FilePath string
	FileName string
	Content  []string
}

func (c *CSPFile) Read(fileName string) {

	// Check file name
	err := c.CheckFileName(fileName)
	shared.CheckError(err, "CSPFile")

	// configure r
	c.FileName = fileName
	c.FilePath = parameters.DIR_CONF

	fullPathAdlFileName := c.FilePath + "/" + c.FileName

	// read file
	fileContent := []string{}
	file, err := os.Open(fullPathAdlFileName)
	shared.CheckError(err, "CSPFile")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// configure r
	c.Content = fileContent
}

func (c CSPFile) Save() (error) {
	r1 := *new(error)

	// create diretcory if it does not exist
	confDir := c.FilePath + "/" + strings.Replace(c.FileName, parameters.CSP_EXTENSION, "", 99)
	_, err := os.Stat(confDir);
	if  os.IsNotExist(err) {
		os.MkdirAll(confDir, os.ModePerm);
	}

	// create file if it does not exist && truncate otherwise
	fmt.Println(confDir + "/" + c.FileName)
	file, err := os.Create(confDir + "/" + c.FileName)
	if err != nil {
		r1 := errors.New("CSP File not created")
		return r1
	}
	defer file.Close()

	// save data
	for i := range c.Content {
		_, err = file.WriteString(c.Content[i])
		if err != nil {
			r1 := errors.New("CSP File not saved")
			return r1
		}
	}
	err = file.Sync()
	if err != nil {
		r1 := errors.New("CSP File not Synced")
		return r1
	}
	defer file.Close()

	return r1
}

func (CSPFile) CheckFileName(fileName string) error {
	r1 := *new(error)

	len := len(fileName)

	if len <= 5 {
		r1 = errors.New("File Name Invalid")
	} else {
		if fileName[len-5:] != parameters.CSP_EXTENSION {
			r1 = errors.New("Invalid extension of '" + fileName + "'")
		} else {
			r1 = nil
		}
	}
	return r1
}

func (c CSPFile) Check() error {
	r1 := *new(error)

	cmdExp := parameters.DIR_FDR + "/" + commands.FDR_COMMAND
	inputFile := c.FilePath + "/" + c.FileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		r1 = errors.New("File '" + inputFile + "' has a problem (e.g.,syntax error)")
		return r1
	}
	s := string(out[:])

	if !strings.Contains(s, "Passed") {
		r1 = errors.New("Deadlock detected")
		return r1
	}
	return r1
}
