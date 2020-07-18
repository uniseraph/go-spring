package kafka


import (
	"context"
	"errors"
	"fmt"
)
import "github.com/confluentinc/confluent-kafka-go/kafka"

type Handler interface {
	Process(ctx context.Context, msg []*kafka.Message) error
	ProcessHeaders(ctx context.Context , headers []kafka.Header)  error
	Start(ctx context.Context) error
}

type HandlerCtor func ( context.Context)(Handler,error)

var type2ctor map[string]HandlerCtor

func init()  {
	type2ctor = make(map[string]HandlerCtor)
}

func Register(name string , ctor HandlerCtor) {
	_ , ok :=  type2ctor[name]

	if ok {
		panic(fmt.Sprintf(  " dup handler type %s",name))
	}
	//logrus.Infof("add a %s handler",name)
	type2ctor[name] = ctor

}
func NewHandler(ctx context.Context) (Handler, error) {

	schema, ok := ctx.Value("schema").(string)
	if !ok  {
		schema = "compute"
	}

	ctor , ok :=  type2ctor[schema]
	if !ok{
		return nil , errors.New( fmt.Sprintf(" No such a handler type %s" , schema))
	}
	return ctor(ctx)

}



