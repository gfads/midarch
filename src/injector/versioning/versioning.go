package versioning

import (
	"os/exec"
	"fmt"
	"os"
	"strings"
	"time"
	"newsolution/shared/parameters"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/framework/architecturallibrary"
)

type VersioningInjector struct {}

func (VersioningInjector) Start(conf madl.MADLGo, elementName string) {

	go start(conf,elementName)
}

func start(conf madl.MADLGo, elementName string) {
	pluginBase01 := elementName + "01"
	sourceCode01 := pluginBase01 + ".go"

	pluginBase02 := elementName + "02"
	sourceCode02 := pluginBase02 + ".go"

	dirPlugins := parameters.DIR_PLUGINS + "/"+ conf.ConfigurationName

	// remove old plugins
	outputLS, err := exec.Command("/bin/ls", dirPlugins).CombinedOutput()
	if err != nil {
		fmt.Println("Injector:: Something wrong in dir '"+dirPlugins)
		os.Stderr.WriteString(err.Error())
		os.Exit(0)
	}
	oldPlugins := strings.Split(string(outputLS), "\n")
	for plugin := range oldPlugins {
		if strings.Contains(oldPlugins[plugin], "_plugin_") {
			exec.Command("/bin/rm", "-r", dirPlugins+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
		}
	}
	
	// generate new plugin
	switch parameters.STRATEGY {
	case 1: // no change
	case 2: // change once
		pluginName := strings.TrimSpace(pluginBase01 + "_plugin_v1")
		_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", dirPlugins+"/"+pluginName, dirPlugins+"/"+pluginBase01+"/"+sourceCode01).CombinedOutput()
		if err != nil {
			fmt.Println("Injector:: Something wrong in generating plugin '"+pluginName+"'")
			os.Stderr.WriteString(err.Error())
			os.Exit(0)
		}
	case 3: // change same plugin
		for {
			pluginName := strings.TrimSpace(pluginBase01 + "_plugin_v1")
			_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", dirPlugins+"/"+pluginName, dirPlugins+"/"+pluginBase01+"/"+sourceCode01).CombinedOutput()
			if err != nil {
				fmt.Println("Injector:: Something wrong in generating plugin '"+pluginName+"'")
				os.Stderr.WriteString(err.Error())
				os.Exit(0)
			}
			time.Sleep(parameters.INJECTION_TIME * time.Second)
		}
	case 4: // alternate plugins
		currentPlugin := 1
		for {
			switch currentPlugin {
			case 1:
				currentPlugin = 2
				pluginName := strings.TrimSpace(pluginBase01 + "_plugin_v1")
				_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", dirPlugins+"/"+pluginName, dirPlugins+"/"+pluginBase01+"/"+sourceCode01).CombinedOutput()
				if err != nil {
					fmt.Println("Injector:: Something is wrong in generating plugin '"+pluginName+"'")
					os.Stderr.WriteString(err.Error())
					os.Exit(0)
				}
				fmt.Println("Injector:: [PLUGIN 01 GENERATED]")
			case 2:
				currentPlugin = 1
				pluginName := strings.TrimSpace(pluginBase02 + "_plugin_v1")
				_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", dirPlugins+"/"+pluginName, dirPlugins+"/"+pluginBase02+"/"+sourceCode02).CombinedOutput()
				if err != nil {
					fmt.Println("Injector:: Something is wrong in generating plugin '"+pluginName+"'")
					os.Stderr.WriteString(err.Error())
					os.Exit(0)
				}
				fmt.Println("Injector:: [PLUGIN 02 GENERATED]")
			}
			time.Sleep(parameters.INJECTION_TIME * time.Second)
		}
	default: // no change
	}
}

func confToGoType(tConf string) string {
	foundType := false
	tGo := ""

	repository := architecturallibrary.ArchitecturalLibrary{}
	repository.Load()

	for t := range repository.Lib {
		if strings.Contains(t, tConf) {
			tGo = t
			foundType = true
		}
	}

	if !foundType {
		fmt.Println("Type '" + tConf + "' NOT FOUND in Library")
		os.Exit(0)
	}
	return tGo
}

