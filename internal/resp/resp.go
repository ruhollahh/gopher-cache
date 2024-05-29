package resp

type RESPType int

const (
	SimpleString RESPType = iota
	BulkString
	Array
	ErrorType
)

type RESP struct {
	RespType RESPType
	Value    interface{}
}
