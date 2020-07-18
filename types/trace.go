package types


type Op string

const OpTrace  = "trace"
const OpUntrace = "untrace"

type OpEvent struct {
	Op Op
	Body interface{}
}
