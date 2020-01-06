package addr

import (
	"github.com/micro/go-micro/util/addr"
	"strconv"
	"strings"
)

func Extract( bind string) (string, int , error) {
	var advertiseIP string
	var advertisePort string

	parts := strings.Split(bind, ":")
	if len(parts) == 2 {
		advertiseIP = parts[0]
		advertisePort = parts[1]
	} else {
		advertiseIP = "0.0.0.0"
		advertisePort = parts[0]
	}

	port, _ := strconv.Atoi(advertisePort)
	advertiseIP, err := addr.Extract(advertiseIP)
	if err != nil {
		return "",0, err
	}
	return advertiseIP, port , nil
}
