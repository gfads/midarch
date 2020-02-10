package messages

type Invocation struct {
	Host string
	Port string
	Op string
	Args [] interface{}
}
