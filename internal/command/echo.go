package command

import (
	"fmt"
	"github.com/ruhollahh/gopher-cache/internal/resp"
)

func handleEcho(args []resp.RESP) resp.RESP {
	if len(args) != 1 || args[0].RespType != resp.BulkString {
		return resp.RESP{RespType: resp.ErrorType, Value: fmt.Errorf("ERR wrong number of arguments for 'echo' command")}
	}

	return resp.RESP{RespType: resp.BulkString, Value: args[0].Value}
}
