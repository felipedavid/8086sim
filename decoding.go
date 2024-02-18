package main

import "fmt"

type Register uint16

const (
	AL Register = iota
	AX
	CL
	CX
	DL
	DX
	BL
	BX
	AH
	SP
	CH
	BP
	DH
	SI
	BH
	DI
)

var registerStr = map[Register]string{
	AL: "al",
	AX: "ax",
	CL: "cl",
	CX: "cx",
	DL: "dl",
	DX: "dx",
	BL: "bl",
	BX: "bx",
	AH: "ah",
	SP: "sp",
	CH: "ch",
	BP: "bp",
	DH: "dh",
	SI: "si",
	BH: "bh",
	DI: "di",
}

type Opcode struct {
	mask   uint16
	opcode uint16
}

var (
	REGMEM_TOFROM_REG = Opcode{mask: 0b1111_1100_0000_0000, opcode: 0b1000_1000_0000_0000}
	IMM_TO_REGMEM     = Opcode{mask: 0b1111_1110_0011_1000, opcode: 0b1100_0110_0000_0000}
	IMM_TO_REG        = Opcode{mask: 0b1111_0000_0000_0000, opcode: 0b1011_0000_0000_0000}
)

// d flag
const (
	REG_IS_SRC = 0
	REG_IS_DST = 1
)

// mod field
const (
	REG_MODE = 0b11
)

func disassemble(stream []byte) {
	for i := 0; i < len(stream); {
		instrSize := 2

		opcode := uint16(stream[i]) << 8
		opcode |= uint16(stream[i+1])

		switch {
		case matchOp(opcode, REGMEM_TOFROM_REG):
			firstByte := uint16(stream[i])
			secondByte := uint16(stream[i+1])
			d := bits(firstByte, 1, 1)
			w := bits(firstByte, 0, 1)
			mod := bits(secondByte, 6, 2)
			reg := (bits(secondByte, 3, 3) << 1) | w
			rm := (bits(secondByte, 0, 3) << 1) | w

			if mod != REG_MODE {
				fatalError("Unknown instruction: 0x%x", opcode)
			}

			var dest, src Register
			if d == REG_IS_SRC {
				src = Register(reg)
				dest = Register(rm)
			} else {
				src = Register(rm)
				dest = Register(reg)
			}

			fmt.Printf("mov %s, %s\n", registerStr[dest], registerStr[src])
		case matchOp(opcode, IMM_TO_REGMEM):
		case matchOp(opcode, IMM_TO_REG):
			firstByte := uint16(stream[i])
			wide := bits(firstByte, 3, 1) == 1
			reg := Register(bits(firstByte, 0, 3) << 1)
			immd := int16(stream[i+1])

			if wide {
				reg |= 1
				immd = immd | (int16(stream[i+2]) << 8)
				instrSize++
			} else {
				immd = (immd << 8) >> 8 // Trick to extend the sign bit
			}

			fmt.Printf("mov %s, %v\n", registerStr[reg], immd)
		default:
			fatalError("Unknown instruction: 0x%x", opcode)
		}

		i += instrSize
	}
}

func matchOp(encodedOp uint16, op Opcode) bool {
	return (encodedOp & op.mask) == op.opcode
}

func bits(n, start, len uint16) uint16 {
	return (n >> start) & ((1 << len) - 1)
}
