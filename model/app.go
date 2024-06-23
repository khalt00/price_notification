package model

import (
	"bytes"
	"encoding/binary"
)

type AppStruct struct {
	Username string
	ChatID   int64
	Stocks   []string
}

// MarshalBinary serializes the AppStruct into a byte slice
func (a AppStruct) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	// Write Username
	if err := binary.Write(&b, binary.LittleEndian, int32(len(a.Username))); err != nil {
		return nil, err
	}
	if _, err := b.WriteString(a.Username); err != nil {
		return nil, err
	}

	// Write ChatID
	if err := binary.Write(&b, binary.LittleEndian, a.ChatID); err != nil {
		return nil, err
	}

	// Write Stocks
	if err := binary.Write(&b, binary.LittleEndian, int32(len(a.Stocks))); err != nil {
		return nil, err
	}
	for _, stock := range a.Stocks {
		if err := binary.Write(&b, binary.LittleEndian, int32(len(stock))); err != nil {
			return nil, err
		}
		if _, err := b.WriteString(stock); err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

// UnmarshalBinary deserializes the byte slice into an AppStruct
func (a *AppStruct) UnmarshalBinary(data []byte) error {
	b := bytes.NewReader(data)

	// Read Username
	var usernameLen int32
	if err := binary.Read(b, binary.LittleEndian, &usernameLen); err != nil {
		return err
	}
	username := make([]byte, usernameLen)
	if _, err := b.Read(username); err != nil {
		return err
	}
	a.Username = string(username)

	// Read ChatID
	if err := binary.Read(b, binary.LittleEndian, &a.ChatID); err != nil {
		return err
	}

	// Read Stocks
	var stocksLen int32
	if err := binary.Read(b, binary.LittleEndian, &stocksLen); err != nil {
		return err
	}
	a.Stocks = make([]string, stocksLen)
	for i := int32(0); i < stocksLen; i++ {
		var stockLen int32
		if err := binary.Read(b, binary.LittleEndian, &stockLen); err != nil {
			return err
		}
		stock := make([]byte, stockLen)
		if _, err := b.Read(stock); err != nil {
			return err
		}
		a.Stocks[i] = string(stock)
	}

	return nil
}
