package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Decode(reader *bufio.Reader) (RESP, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return RESP{}, err
	}

	switch prefix {
	case '+':
		return decodeSimpleString(reader)
	case '$':
		return decodeBulkString(reader)
	case '*':
		return decodeArray(reader)
	default:
		return RESP{}, fmt.Errorf("unknown prefix: %c", prefix)
	}
}

func decodeSimpleString(reader *bufio.Reader) (RESP, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESP{}, err
	}

	return RESP{RespType: SimpleString, Value: line}, nil
}

func decodeBulkString(reader *bufio.Reader) (RESP, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESP{}, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return RESP{}, err
	}
	if length == -1 {
		return RESP{RespType: BulkString, Value: nil}, nil
	}

	data := make([]byte, length+2) // +2 for \r\n
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return RESP{}, err
	}

	return RESP{RespType: BulkString, Value: string(data[:length])}, nil
}

func decodeArray(reader *bufio.Reader) (RESP, error) {
	line, err := readLine(reader)
	if err != nil {
		return RESP{}, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return RESP{}, err
	}
	if length == -1 {
		return RESP{RespType: Array, Value: nil}, nil
	}

	items := make([]RESP, length)
	for i := 0; i < length; i++ {
		item, err := Decode(reader)
		if err != nil {
			return RESP{}, err
		}
		items[i] = item
	}
	return RESP{RespType: Array, Value: items}, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(line, "\r\n"), nil
}
