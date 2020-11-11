package components

import (
	"errors"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/repositories/plugins/http2"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"shared"
	"sort"
	"strings"
)

type Route struct {
	Mapping		string
	Function	interface{}
}

type Http2InvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewHttp2InvokerM() Http2InvokerM {
	r := new(Http2InvokerM)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e Http2InvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg, info)
}

func (Http2InvokerM) I_Process(msg *messages.SAMessage, info [] *interface{}) { // TODO
	// unmarshall
	httpMessage := msg.Payload.(messages.HttpMessage)
	//request := messages.HttpRequest{}
	//request.Unmarshal(payload)

	//response := messages.HttpResponse{}

	routes := getControllers()
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Mapping < routes[j].Mapping
	})

	var result []reflect.Value
	var err error

	found := false
	for _, route := range routes {
		log.Println(httpMessage.Request.RequestURI,  "=>", route.Mapping)
		if strings.HasPrefix(httpMessage.Request.RequestURI, route.Mapping) {
			result, err = Call(route.Function)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Result:", len(result))
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintf(httpMessage.Response, "404 - Not found.\n\nThe requested URL " + httpMessage.Request.RequestURI + " was not found on this server!")
		return
	}

	//httpMessage.Response.Write([]byte(result.(string)))

	log.Println("Response:", httpMessage.Response)
	log.Println("Result:", len(result))

	fmt.Fprintf(httpMessage.Response, result[0].Interface().(string))

	//impl.Handler(httpMessage.Response, httpMessage.Request)

	//msgTemp := response.Marshal()
	*msg = messages.SAMessage{Payload: httpMessage}
	log.Println("Http2Invoker.I_Process")
}


func getProjectFiles(projectRoot string) (files []string) {
	files = []string{}
	log.Println("getProjectFiles")
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		//log.Println("getProjectFiles.func: file->", path)
		if info == nil || info.IsDir() {
			return nil
		}
//		log.Println("getProjectFiles.before append")
		if filepath.Ext(path) == ".go" {
			//log.Println("getProjectFiles.appending",path,"where",files)
			files = append(files, path)
		}
		//log.Println("getProjectFiles.after append")
		return nil
	})
	if err!= nil {
		panic(err)
	}

	return files
}

func getControllers() (routes []Route) {
	files := getProjectFiles(shared.DIR_BASE)  // TODO: get project base dir, not the "go" base dir
	routes = []Route{}

	for _, file := range files {
		if strings.Contains(file, "/home/dcruzb/go/src/") { ///components/") {
			//log.Println(file)
			routes = append(routes, FuncDescription(file)...)
		}
	}

	return routes
}

func FuncDescription(filename string) (routes []Route) {
	//log.Println("Filename:", filename)
	routes = []Route{}

	parsedAst, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}

	pkg := &ast.Package{
		Name:    filename,
		Files:   make(map[string]*ast.File),
	}
	pkg.Files[filename] = parsedAst


	d := doc.New(pkg, filename, doc.AllDecls)

	if !(strings.Contains(d.Doc, "@Controller")) {
		return nil
	}

	//log.Println(d)
	for _, theFunc := range d.Funcs {
		//log.Println("*********************", theFunc.Doc)
		if strings.Contains(theFunc.Doc, "@GetMapping") {
			mapping := ""
			var function interface{}

			mappingIndex := strings.Index(d.Doc, "@RequestMapping(")
			if mappingIndex >= 0 {
				mapping = d.Doc[mappingIndex + 17:len(d.Doc)-3]
			}

			rmp := strings.Index(theFunc.Doc, "@GetMapping(value=")
			if rmp >= 0 {
				mapping += theFunc.Doc[rmp + 19:len(theFunc.Doc)-3]
			}

			//if theFunc.Name == "GetHealth" {
			//	function = impl.GetHealth
			//}else{
			//	function = impl.GetItems
			//}
			function = http2.GetFunction(theFunc.Name)

			route := Route{
				Mapping:  mapping,
				Function: function,
			}

			routes = append(routes, route)
			log.Println("Filename:", filename, "Controllers:", route)
		}
	}

	if len(routes) <= 0 {
		return nil
	}

	return routes
}

func Call(function interface{}, params ... interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(function)
	if len(params) != f.Type().NumIn() {
		err = errors.New("Unexpected number of params.")
		return nil, err
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)
	log.Println("Call.Result:", result)
	return result, nil
}