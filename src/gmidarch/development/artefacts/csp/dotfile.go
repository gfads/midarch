package csp

import (
	"newsolution/shared/parameters"
	"os"
	"bufio"
	"errors"
	"newsolution/shared/shared"
)

type DOTFile struct {
	FilePath string
	FileName string
	Content  []string
}

func (c *DOTFile) Read() error {
	r1 := *new(error)

	// Check file name
	err := c.CheckFileName(c.FileName)
	if err != nil {
		err = errors.New("DotFile"+err.Error())
		return r1
	}

	fullPathFileName := c.FilePath + "/" + c.FileName

	// read file
	fileContent := []string{}
	file, err := os.Open(fullPathFileName)
	shared.CheckError(err,"DOTFILE")

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// configure r
	c.Content = fileContent

	return r1
}

func (c DOTFile) LoadDotFiles(csp CSP)(map[string]DOTFile,error){
	r1 := map[string]DOTFile{}
	r2 := *new(error)

	// Components
	for i := range csp.CompProcesses{
		dotFile := DOTFile{}
		dotFile.FilePath = parameters.DIR_CSP + "/" + csp.CompositionName

		dotFile.FileName = i + parameters.DOT_EXTENSION
		err := c.CheckFileName(dotFile.FileName)
		if err != nil {
			r2 = errors.New("DOTFILE"+err.Error())
			return r1,r2
		}
		err = dotFile.Read()
		if err != nil {
			r2 = errors.New("DOTFILE"+err.Error())
			return r1,r2
		}

		r1[i] = dotFile
	}

	// Connectors
	for i := range csp.ConnProcesses{
		dotFile := DOTFile{}
		dotFile.FilePath = parameters.DIR_CSP + "/" + csp.CompositionName

		dotFile.FileName = i + parameters.DOT_EXTENSION
		err := c.CheckFileName(dotFile.FileName)
		if err != nil {
			r2 = errors.New("DOTFILE"+err.Error())
			return r1,r2
		}
		err = dotFile.Read()
		if err != nil {
			r2 = errors.New("DOTFILE"+err.Error())
			return r1,r2
		}

		r1[i] = dotFile
	}

	return r1,r2
}

func (DOTFile) CheckFileName(fileName string) error {
	r1 := *new(error)

	l := len(fileName)

	if l <= len(parameters.DOT_EXTENSION) {
		r1 = errors.New("File Name Invalid")
	} else {
		if fileName[l-len(parameters.DOT_EXTENSION):] != parameters.DOT_EXTENSION {
			r1 = errors.New("Invalid extension of '" + fileName + "'")
		} else {
			r1 = nil
		}
	}
	return r1
}