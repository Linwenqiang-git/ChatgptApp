package convert

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(n))
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
	// b := make([]byte, 4)
	// binary.BigEndian.PutUint32(b, uint32(n))
	// return b
}
func BytesToInt(b []byte) int32 {
	var i int32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}
