package loader

import (
	"bufio"
	"fmt"
	"gmidarch/development/artefacts/madl"
	"os"
	"shared"
)

type Loader struct{}

func (l Loader) Load(file string) madl.MADL {

	// read file and create go information
	m := l.read(file)

	// Syntatic Check of configuration
	err := m.Check()
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(1)
	}

	// Configure Channels and Maps
	m.ConfigureChannelsAndMaps()

	// Configure Components
	m.ConfigureComponents()

	// Configure Connectors
	m.ConfigureConnectors()

	return m
}

func (Loader) read(file string) madl.MADL {

	m := madl.MADL{}

	// Check file name
	err := shared.CheckFileName(file)
	if err != nil {
		fmt.Println("MADL:: " + err.Error())
		os.Exit(0)
	}

	// Configure File & Path
	m.File = file
	m.Path = shared.DIR_MADL
	fullPathAdlFileName := m.Path + "/" + m.File

	// Read file
	fileContent := []string{}
	fileTemp, err := os.Open(fullPathAdlFileName)
	if err != nil {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	defer fileTemp.Close()

	scanner := bufio.NewScanner(fileTemp)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Identify Configuration
	configurationName, err := m.IdentifyConfigurationName(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Configuration = configurationName

	// Identify Components
	comps, err := m.IdentifyComponents(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Components = comps

	// Identify Connectors
	connectors, err := m.IdentifyConnectors(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Connectors = connectors

	// Identify Attachments
	attachments, err := m.IdentifyAttachments(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Attachments = attachments
	//m.SetAttachmentTypes()

	// Identify adaptability
	adaptability, err := m.IdentifyAdaptability(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Adaptability = adaptability

	return m
}