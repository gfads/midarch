package exec

type Exec struct{}

/*
func (Exec) Create(id string, elem interface{}, typeName string, dot dot.DOTGraph, maps map[string]string, channels map[string]chan messages.SAMessage) (ExecGraph) {
	r1 := NewExecGraph(dot.NumNodes)

	// Check dot actions against elem's interface
	checkInterface(elem, id, dot)

	// initialisation of message and info of a given element
	msg := new(messages.SAMessage)
	*msg = messages.SAMessage{Payload: ""} // TODO
	info := make([]*interface{}, 3, 3)     // 3 can be set any value
	for i := 0; i < 3; i++ {
		info[i] = new(interface{})
		*info[i] = new(interface{})
	}

	for e1 := range dot.EdgesDot {
		for e2 := range dot.EdgesDot [e1] {
			eActions := ExecEdgeInfo{}
			edgeTemp := dot.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action
			actionNameExec := ""
			if strings.Contains(actionNameFDR, ".") {
				actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
			}
			if sharedadaptive.IsExternal(actionNameExec) { // External action
				actionNameTemp := strings.Split(actionNameFDR, ".")
				key1 := id + "." + actionNameTemp[1]
				key2 := id + "." + actionNameTemp[0] + "." + maps[key1]
				channel, _ := channels[key2]
				params := ExecEdgeInfo{}
				switch actionNameExec {
				case sharedadaptive.INVR:
					invr := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName: "InvR", IsInternal: false, Message: msg, ActionChannel: &invr}
				case sharedadaptive.TERR:
					terr := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionName: "TerR", IsInternal: false, Message: msg, ActionChannel: &terr}
				case sharedadaptive.INVP:
					invp := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionName: "InvP", IsInternal: false, Message: msg, ActionChannel: &invp}
				case sharedadaptive.TERP:
					terp := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionName: "TerP", IsInternal: false, Message: msg, ActionChannel: &terp}
				}
				mapType := ExecEdgeInfo{}
				mapType = params
				eActions = mapType
			}

			if sharedadaptive.IsInternal(actionNameFDR) {
				channel := make(chan messages.SAMessage,sharedadaptive.CHAN_BUFFER_SIZE)

				// Configure selector of each individual element
				s := components.ConfigureSelector(reflect.TypeOf(elem).Name())

				params := ExecEdgeInfo{InternalAction: s, ActionName: actionNameFDR, IsInternal: true, ActionChannel: &channel, Message: msg, Info: info}
				mapType := params
				eActions = mapType
			}
			r1.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
		}
	}

	return *r1
}

*/
