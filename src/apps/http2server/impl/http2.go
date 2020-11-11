//@Controller
//@RequestMapping("/api")
package impl

//@GetMapping(value="/health")
func GetHealth() string {
	return "{'status':'ok'}"
}

//@GetMapping(value="/items")
func GetItems() string {
	return "{'items': [{'name':'Item 1'},{'name':'Item 2'},{'name':'Item 3'}]}"
}

//func Handler(w http.ResponseWriter, r *http.Request) {
//	log.Println("Handler")
//	fmt.Fprintf(w, "<h1>Hi there %s!</h1>", r.URL.Path)
//}