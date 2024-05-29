package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func Encode(r RESP) ([]byte, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	var err error
	switch r.RespType {
	case SimpleString:
		err = encodeSimpleString(writer, r.Value)
	case BulkString:
		err = encodeBulkString(writer, r.Value)
	case ErrorType:
		err = encodeError(writer, r.Value)
	default:
		return nil, fmt.Errorf("unknown RESP type: %v", r.RespType)
	}

	if err != nil {
		return nil, err
	}

	if err = writer.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func encodeSimpleString(writer *bufio.Writer, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid simple string Value")
	}

	_, err := writer.WriteString("+" + str + "\r\n")

	return err
}

func encodeBulkString(writer *bufio.Writer, value interface{}) error {
	if value == nil {
		_, err := writer.WriteString("$-1\r\n")
		return err
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid bulk string Value")
	}

	_, err := writer.WriteString("$" + strconv.Itoa(len(str)) + "\r\n" + str + "\r\n")

	return err
}

func encodeError(writer *bufio.Writer, value interface{}) error {
	errVal, ok := value.(error)
	if !ok {
		return errors.New("invalid error Value")
	}

	_, err := writer.WriteString("-" + errVal.Error() + "\r\n")
	return err
}
