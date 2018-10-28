package work

import (
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"testing"
)

func TestPipe(t *testing.T) {
	c, err := winio.DialPipe(PipeName, nil)
	if err != nil {
		panic(err)
	}

	cli := jsonrpc2.NewClient(c)

	//var u settings.Config
	//err = cli.Call("SetsSvc.UserConfig", nil, &u)
	//if err != nil {
	//	panic(err)
	//}

	err = cli.Call("SetsSvc.SetPortName", []string{"COM1"}, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("close...")
	_ = c.Close()
	_ = cli.Close()

}
