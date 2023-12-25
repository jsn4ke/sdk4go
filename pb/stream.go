package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

// CheckStreamValid
// 二进制流 pb3
// 返回有效流的长度 从0开始
func CheckStreamValid(stream []byte, in proto.Message) int {
	var (
		state int

		index int
		wire  int

		length int

		shift int

		validEnd int

		skipEnd int
	)
	const (
		inHead = 0
		inWire = 1
	)
fail:
	for i, v := range stream {
		if i < skipEnd {
			continue
		}
		v := int(v)
		switch state {
		case inHead:
			if 0 != v&0x80 {
				index |= v & 0x7f
				shift += 7
			} else {
				index <<= shift
				index |= v >> 3
				wire = v & 0x7

				state = inWire
				shift = 0
			}
		case inWire:
			switch wire {
			case 0:
				if 0 != v&0x80 {
					length |= v & 0x7f
					shift += 7
				} else {
					length <<= shift
					length |= v

					shift = 0
					state = inHead

					validEnd = i + 1

					index = 0
					length = 0
					wire = 0
				}
			case 2:
				if 0 != v&0x80 {
					length |= v & 0x7f
					shift += 7
				} else {
					length <<= shift
					length |= v

					skipEnd = i + 1 + length
					if skipEnd > len(stream) {
						break fail
					} else {
						err := proto.Unmarshal(stream[:skipEnd], in)
						if nil != err {
							break fail
						}
					}
				}

				shift = 0
				state = inHead

				validEnd = skipEnd

				index = 0
				length = 0
				wire = 0
			default:
				break fail
			}
		}
	}
	fmt.Println(index, validEnd)
	err := proto.Unmarshal(stream[:validEnd], in)
	if nil != err {
		fmt.Errorf("func error")
	} else {
		fmt.Println(in.String())
	}
	return validEnd
}
