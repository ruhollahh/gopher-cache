package command

import (
	"fmt"
	"github.com/ruhollahh/gopher-cache/internal/resp"
	"strings"
)

type Handler struct {
	handlers map[string]func([]resp.RESP) resp.RESP
}

func NewHandler() Handler {
	return Handler{handlers: make(map[string]func([]resp.RESP) resp.RESP)}
}

func (h Handler) RegisterHandlers() {
	h.register("PING", handlePing)
	h.register("ECHO", handleEcho)
}

func (h Handler) register(command string, handler func([]resp.RESP) resp.RESP) {
	h.handlers[command] = handler
}

func (h Handler) Handle(cmd resp.RESP) resp.RESP {
	if cmd.RespType != resp.Array {
		return resp.RESP{RespType: resp.ErrorType, Value: fmt.Errorf("ERR invalid command")}
	}

	elements, ok := cmd.Value.([]resp.RESP)
	if !ok || len(elements) == 0 {
		return resp.RESP{RespType: resp.ErrorType, Value: fmt.Errorf("ERR invalid command")}
	}

	cmdName, ok := elements[0].Value.(string)
	if !ok {
		return resp.RESP{RespType: resp.ErrorType, Value: fmt.Errorf("ERR invalid command")}
	}

	handler, exists := h.handlers[strings.ToUpper(cmdName)]
	if !exists {
		return resp.RESP{RespType: resp.ErrorType, Value: fmt.Errorf("ERR unknown command")}
	}

	return handler(elements[1:])
}
