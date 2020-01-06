package spring

import (
	"fmt"
	"gitlab.ustock.cc/core/go-spring/addr"
	"gitlab.ustock.cc/core/go-spring/info"
	"testing"
)

func Test_add(t *testing.T)  {

	addr.Extract("0.0.0.0:8282")

	fmt.Println(info.Target)
}