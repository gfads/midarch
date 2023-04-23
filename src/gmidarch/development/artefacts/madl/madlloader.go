package madl

import (
	"bufio"
	"github.com/gfads/midarch/src/gmidarch/development/components/component"
	"github.com/gfads/midarch/src/gmidarch/development/connectors"
	"github.com/gfads/midarch/src/shared"
	"os"
	"strings"
)

type MADLLoader interface {
	Load(string) MADL
}

type MADLLoaderImpl struct{}

func NewMADLLoader() MADLLoader {
	var mp MADLLoader

	mp = MADLLoaderImpl{}

	return mp
}

func (m MADLLoaderImpl) Load(fileName string) MADL {
	r := MADL{}

	// Check file name
	shared.CheckFileName(fileName, shared.MADL_EXTENSION)

	// Configure File & Path
	r.FileName = fileName
	r.Path = shared.DIR_MADL // TODO dcruzb : remove shared.DIR_MADL, should be passed the entire path to the file (this is user's responsibility)
	fullPathAdlFileName := r.Path + "/" + r.FileName

	// Read file
	fileContent := []string{}
	fileTemp, err := os.Open(fullPathAdlFileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer fileTemp.Close()

	scanner := bufio.NewScanner(fileTemp)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Check begin and EndConf
	// TODO

	// Get Configuration name
	r.Configuration = getConfigurationName(fileContent)

	// Get Components
	r.Components = getComponents(fileContent)

	// Get Connectors
	r.Connectors = getConnectors(fileContent)

	// Get Attachments
	r.Attachments = getAttachments(fileContent)

	// Get adaptability
	r.Adaptability = getAdaptability(fileContent)

	return r
}

// Get Configuration Name
func getConfigurationName(content []string) string {
	r := ""
	found := false

	for l := range content {
		tempContent := content[l]
		if strings.Contains(strings.ToUpper(tempContent), "CONFIGURATION") {
			temp := strings.Split(tempContent, " ")
			r = strings.TrimSpace(temp[1])
			found = true
			break
		}
	}
	if !found {
		shared.ErrorHandler(shared.GetFunction(), "Command 'Configuration' not found!!")
	}
	if r == "" {
		shared.ErrorHandler(shared.GetFunction(), "Configuration name not defined.")
	}
	return r
}

// Get Components
func getComponents(content []string) []component.Component {
	foundComponents := false
	r := []component.Component{}

	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "COMPONENTS") {
			foundComponents = true
		} else {
			if foundComponents && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				compType := *new(string)
				temp := strings.Split(tempLine, ":")
				compId := strings.TrimSpace(temp[0])
				compType = strings.TrimSpace(temp[1])
				r = append(r, component.Component{Id: compId, TypeName: compType})
			} else {
				if foundComponents && !shared.SkipLine(tempLine) && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if !foundComponents {
		shared.ErrorHandler(shared.GetFunction(), "Command 'Components' not found.")
	}
	if len(r) == 0 {
		shared.ErrorHandler(shared.GetFunction(), "'Components' not well formed.")
	}

	return r
}

// Get Connectors
func getConnectors(content []string) []connectors.Connector {
	foundConnectors := false
	r := []connectors.Connector{}

	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "CONNECTORS") {
			foundConnectors = true
		} else {
			if foundConnectors && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				connId := strings.TrimSpace(temp[0])
				connType := strings.TrimSpace(temp[1])
				connTypeName := connType
				newConn := connectors.NewConnector(connId, connTypeName, "", 0, 0)
				r = append(r, newConn)
			} else {
				if foundConnectors && tempLine != "" && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(r) == 0 {
		shared.ErrorHandler(shared.GetFunction(), "'Connectors' not well formed.")
	}

	return r
}

// Get Attachments
func getAttachments(content []string) []Attachment {
	r := []Attachment{}

	// Identify Attachments
	foundAttachments := false
	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ATTACHMENTS") {
			foundAttachments = true
		} else {
			if foundAttachments && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ",") {
				atts := strings.Split(strings.TrimSpace(tempLine), ",")
				c1Temp := strings.TrimSpace(atts[0])
				tTemp := strings.TrimSpace(atts[1])
				c2Temp := strings.TrimSpace(atts[2])

				c1 := component.Component{Id: c1Temp}
				t := connectors.Connector{Id: tTemp}
				c2 := component.Component{Id: c2Temp}

				att := Attachment{c1, t, c2}
				r = append(r, att)
			} else {
				if foundAttachments && tempLine != "" && !strings.Contains(tempLine, ",") {
					break
				}
			}
		}
	}

	if len(r) == 0 {
		shared.ErrorHandler(shared.GetFunction(), "'Attachments' not well formed.")
	}

	return r
}

// Get Adaptability
func getAdaptability(content []string) []string {
	r := []string{}

	foundAdaptability := false
	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ADAPTABILITY") {
			foundAdaptability = true
		} else {
			if foundAdaptability && !shared.SkipLine(tempLine) && isAdaptationType(tempLine) {
				r = append(r, strings.ToUpper(strings.TrimSpace(tempLine)))
			} else {
				if foundAdaptability && !shared.SkipLine(tempLine) && !isAdaptationType(tempLine) {
					break
				}
			}
		}
	}

	if !foundAdaptability || len(r) == 0 {
		shared.ErrorHandler(shared.GetFunction(), "'Adaptability' NOT well defined!")
	}

	return r
}

// Check whether the adaptation type is one supported by gMidArch
func isAdaptationType(t string) bool {

	_, r := shared.AdaptationTypes[strings.ToUpper(strings.TrimSpace(t))]

	return r
}
