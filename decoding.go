package main

import "fmt"

type Register uint16

const (
	AL uint16 = iota
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

var registerStr = map[uint16]string{
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
	NO_DISP          = 0b00
	EIGHT_BIT_DISP   = 0b01
	SIXTEEN_BIT_DISP = 0b10
	REG_MODE         = 0b11
)

var effectiveAddrCalc = map[uint16]string{
	0b000: "bx + si",
	0b001: "bx + di",
	0b010: "bp + si",
	0b011: "bp + di",
	0b100: "si",
	0b101: "di",
	0b110: "bp", // I'm ignoring direct memory access
	0b111: "bx",
}

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
			reg := bits(secondByte, 3, 3)
			rm := bits(secondByte, 0, 3)

			if mod < REG_MODE {
				addrCalc := effectiveAddrCalc[rm]
				if mod == EIGHT_BIT_DISP {
					instrSize += 1
					addrCalc = fmt.Sprintf("[%s + %d]", addrCalc, uint16(stream[i+2]))
				} else if mod == SIXTEEN_BIT_DISP {
					instrSize += 2
					addrCalc = fmt.Sprintf("[%s + %d]", addrCalc, uint16(stream[i+2])|(uint16(stream[i+3])<<8))
				} else {
					addrCalc = fmt.Sprintf("[%s]", addrCalc)
				}

				reg = reg<<1 | w
				var dest, src string
				if d == REG_IS_SRC {
					src = registerStr[(reg)]
					dest = addrCalc
				} else {
					src = addrCalc
					dest = registerStr[reg]
				}

				fmt.Printf("mov %s, %s\n", dest, src)
			} else if mod == REG_MODE {
				var dest, src uint16
				if d == REG_IS_SRC {
					src = reg
					dest = rm
				} else {
					src = rm
					dest = reg
				}

				src = (src << 1) | w
				dest = (dest << 1) | w
				fmt.Printf("mov %s, %s\n", registerStr[dest], registerStr[src])
			} else {
				fatalError("should not happen..")
			}
		case matchOp(opcode, IMM_TO_REGMEM):
			fmt.Println("hi")
		case matchOp(opcode, IMM_TO_REG):
			firstByte := uint16(stream[i])
			wide := bits(firstByte, 3, 1) == 1
			reg := bits(firstByte, 0, 3) << 1
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
