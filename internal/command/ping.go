package command

import "github.com/ruhollahh/gopher-cache/internal/resp"

func handlePing(_ []resp.RESP) resp.RESP {
	return resp.RESP{RespType: resp.SimpleString, Value: "PONG"}
}
