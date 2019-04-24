package objects

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

type BlockID struct {
	Data [20]byte
}

func (b BlockID) BlockNumber() uint32 {
	return binary.BigEndian.Uint32(b.Data[:4])
}

func (b BlockID) ID() []byte {
	return b.Data[4:]
}

func (b BlockID) RefBlockNum() UInt16 {
	return UInt16(b.BlockNumber())
}

func (b BlockID) RefBlockPrefix() UInt32 {
	var prefix UInt32
	binary.Read(bytes.NewReader(b.Data[4:8]), binary.LittleEndian, &prefix)
	return prefix
}

func (b BlockID) String() string {
	return hex.EncodeToString(b.Data[:])
}

func (b BlockID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(hex.EncodeToString(b.Data[:]))), nil
}

func (b *BlockID) UnmarshalJSON(data []byte) error {
	unqoted, err := strconv.Unquote(string(data))
	if err != nil {
		return nil
	}
	bin, err := hex.DecodeString(unqoted)
	if err != nil {
		return err
	}
	if len(bin) != len(b.Data) {
		return fmt.Errorf("Invalid BlockID size")
	}

	copy(b.Data[:], bin)
	return err
}
