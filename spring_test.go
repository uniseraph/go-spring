package spring

import (
	"fmt"
	"github.com/uniseraph/go-spring/addr"
	"github.com/uniseraph/go-spring/info"
	"testing"
)

func Test_add(t *testing.T)  {

	addr.Extract("0.0.0.0:8282")

	fmt.Println(info.Target)
}