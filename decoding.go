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

			if mod != 0b11 {
				fatalError("Unknown instruction: 0x%x", opcode)
			}

			var dest, src Register
			if d == 0 {
				dest = Register(reg)
				src = Register(rm)
			} else {
				dest = Register(rm)
				src = Register(reg)
			}

			fmt.Printf("mov %s, %s\n", registerStr[dest], registerStr[src])
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
