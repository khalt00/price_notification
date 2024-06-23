package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func ByteToInt64(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, fmt.Errorf("invalid byte slice length: %d", len(b))
	}

	var result int64
	// Use binary.Read to convert the byte slice to an int64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &result)
	if err != nil {
		return 0, fmt.Errorf("binary.Read failed: %v", err)
	}
	return result, nil
}

func Int64ToBytes(n int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		return nil, fmt.Errorf("binary.Write failed: %v", err)
	}
	return buf.Bytes(), nil
}
