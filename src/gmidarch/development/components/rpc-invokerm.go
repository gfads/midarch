package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

//type Route struct {
//	Mapping		string
//	Function	interface{}
//	Parameters  []string
//}

type RPCInvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
	routes		[]Route
}

func NewRPCInvokerM() RPCInvokerM {
	r := new(RPCInvokerM)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	r.routes = getControllers()

	return *r
}

func (inv RPCInvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	inv.I_Process(msg, info)
}

func (inv RPCInvokerM) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	// unmarshall
	httpMessage := msg.Payload.(messages.HttpMessage)
	//request := messages.HttpRequest{}
	//request.Unmarshal(payload)

	//response := messages.HttpResponse{}

	//var result []reflect.Value
	//var err error
	//errorInCall := false
	//
	//found := false
	//for _, route := range inv.routes {
	//	if strings.HasPrefix(httpMessage.Request.RequestURI, route.Mapping) {
	//		//log.Println(httpMessage.Request.RequestURI,  "=>", route.Mapping)
	//		uriParameters := getURIParameters(httpMessage.Request.RequestURI)
	//		// Sort parameters based on function
	//		var parameters []interface{}
	//		for _, param := range route.Parameters {
	//			value := uriParameters[param]
	//			if value != nil {
	//				parameters = append(parameters, value)
	//			}
	//		}
	//		//log.Println("parameters:", parameters, "route.Parameters", route.Parameters)
	//
	//		result, err = Call(route.Function, parameters)
	//		if err != nil {
	//			// Every error must be returned to client, dont raise a fatal error
	//			errorInCall = true
	//		}
	//		found = true
	//		break
	//	}
	//}
	//
	//if !found {
	//	httpMessage.Response.WriteHeader(http.StatusNotFound)
	//	httpMessage.Response.Write([]byte("404 - Not found.\n\nThe requested URL " + httpMessage.Request.RequestURI + " was not found on this server!"))
	//	return
	//} else if errorInCall {
	//	httpMessage.Response.WriteHeader(http.StatusBadRequest)
	//	httpMessage.Response.Write([]byte("400 - Bad Request.\n\n" + err.Error()))
	//	return
	//}
	//
	//httpMessage.Response.Write([]byte(result[0].Interface().(string)))

	*msg = messages.SAMessage{Payload: httpMessage}
	//log.Println("Http2Invoker.I_Process")
}

//func getURIParameters(uri string) (parameters map[string]interface{}) {
//	parameters = make(map[string]interface{})
//	paramRegex, err := regexp.Compile("([?][\\w|=|%|(|)|+|-|.]+)|([&][\\w|=|%|(|)|+|-|.]+)")
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	params := paramRegex.FindAllString(uri, -1)
//	for _, param := range params {
//		param := param[1:]
//		keyValue := strings.Split(param, "=")
//		parameters[keyValue[0]] = keyValue[1]
//	}
//	return parameters
//}
//
//
//func getProjectFiles(projectRoot string) (files []string) {
//	files = []string{}
//	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
//		//log.Println("getProjectFiles.func: file->", path)
//		if info == nil || info.IsDir() {
//			return nil
//		}
//		if filepath.Ext(path) == ".go" {
//			//log.Println("getProjectFiles.appending",path,"where",files)
//			files = append(files, path)
//		}
//		return nil
//	})
//	if err!= nil {
//		panic(err)
//	}
//
//	return files
//}
//
//func getControllers() (routes []Route) {
//	files := getProjectFiles(shared.DIR_BASE + "/src/apps/")
//	routes = []Route{}
//
//	for _, file := range files {
//		routes = append(routes, getRoutesFromFile(file)...)
//	}
//
//	sort.Slice(routes, func(i, j int) bool {
//		return routes[i].Mapping < routes[j].Mapping
//	})
//
//	return routes
//}
//
//func getRoutesFromFile(filename string) (routes []Route) {
//	//log.Println("Filename:", filename)
//	routes = []Route{}
//
//	parsedAst, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	pkg := &ast.Package{
//		Name:    filename,
//		Files:   make(map[string]*ast.File),
//	}
//	pkg.Files[filename] = parsedAst
//
//
//	d := doc.New(pkg, filename, doc.AllDecls)
//
//	if !(strings.Contains(d.Doc, "@Controller")) {
//		return nil
//	}
//
//	rgxAnnotations, err := regexp.Compile("([@][\\w|=|\"|(|)|,| |/]+)")
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	rgxParameters, err := regexp.Compile("((?:\\()[\\w|=|\"|,| |/]+[)])")
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//
//	mapping := ""
//	requestMapping := ""
//	mappingIndex := strings.Index(d.Doc, "@RequestMapping(")
//	if mappingIndex >= 0 {
//		requestMapping = d.Doc[mappingIndex+17 : len(d.Doc)-3]
//	}
//
//	for _, theFunc := range d.Funcs {
//		//log.Println("*********************", theFunc.Doc)
//		var function interface{}
//		function = http2.GetFunction(theFunc.Name)
//		var parameters []string
//
//		annotations := rgxAnnotations.FindAllString(theFunc.Doc, -1)
//		//log.Println("annotations:", annotations, "theFunc.Doc", theFunc.Doc)
//		for _, annotation := range annotations {
//			param := rgxParameters.FindString(annotation)
//			//log.Println("param:", param)
//			if param != "" {
//				param = param[1:len(param)-1]
//			}
//			params := strings.Split(param, ",")
//			//log.Println("param:", param, "params:", params)
//
//			if strings.Contains(annotation, "@GetMapping") {
//				for _, parameter := range params {
//					p := strings.Split(parameter, "=")
//					key := p[0]
//					value := strings.ReplaceAll(p[1], "\"", "")
//					//rmp := strings.Index(theFunc.Doc, "@GetMapping(value=")
//					//if rmp >= 0 {
//					if key == "value" {
//						mapping = requestMapping + value //theFunc.Doc[rmp+19 : len(theFunc.Doc)-3]
//					}
//				}
//			}
//			if strings.Contains(annotation, "@RequestParam") {
//				//log.Println("@RequestParam:", params)
//				for _, parameter := range params {
//					p := strings.Split(parameter, "=")
//					key := p[0]
//					value := strings.ReplaceAll(p[1], "\"", "")
//					if key == "key" {
//						parameters = append(parameters, value)
//					}
//				}
//			}
//		}
//
//		route := Route{
//			Mapping:  mapping,
//			Function: function,
//			Parameters: parameters,
//		}
//
//		routes = append(routes, route)
//		//log.Println("Filename:", filename, "Controllers:", route)
//	}
//
//	if len(routes) <= 0 {
//		return nil
//	}
//
//	return routes
//}
//
//func Call(function interface{}, params []interface{}) (result []reflect.Value, err error) {
//	f := reflect.ValueOf(function)
//	numIn := f.Type().NumIn()
//	if len(params) != numIn {
//		err = errors.New("Unexpected number of params.")
//		return nil, err
//	}
//
//	in := make([]reflect.Value, numIn)
//	for k, param := range params {
//		switch f.Type().In(k).Kind() { // TODO: dcruzb Implement support to other types too
//		case reflect.Int:
//			value, err := strconv.Atoi(param.(string))
//			if err != nil {
//				return nil, err
//			}
//			in[k] = reflect.ValueOf(value)
//		}
//	}
//
//	result = f.Call(in)
//	//log.Println("Call.Result:", result)
//	return result, nil
//}