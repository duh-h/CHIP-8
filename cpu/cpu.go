package main

import "fmt"

type CPU struct {

	//regitradores
	AF, BC, DE, HL, PC, SP uint16

	//flags
	Flags struct {
		Z, N, H, C bool
	}

	Memory *Memory

	/*
		O que cada registrador faz:
		AF: Acumulador (A, armazena resultados) + Flags (F, estado da CPU).

		BC, DE, HL: Registradores de propósito geral (podem ser usados como 16 bits ou 8 bits).

		PC: Aponta para a próxima instrução na memória.

		SP: Aponta para o topo da pilha (usado em chamadas de função/interrupções).


		Flags:
		Z (Zero): Ativado se o resultado da operação for zero.

		N (Subtract): Indica se a última operação foi uma subtração.

		H (Half-Carry): Indica "carry" nos primeiros 4 bits (ex: 0x0F + 0x01 = 0x10).

		C (Carry): Indica carry em operações de 8 bits (ex: 0xFF + 0x01 = 0x00 com carry).
	*/

}

func (c *CPU) Step() {

	//fetch: ler o opcode
	opcode := c.Memory.Read(c.PC)
	c.PC++

	switch opcode {

	case 0x00:
		// nao faz nada

	case 0x80: // soam b em a
		a := c.A()
		b := c.B()
		result := uint16(a) + uint16(b)

		// aualiza flags
		c.Flags.Z = (result & 0xFF) == 0
		c.Flags.N = false // ADD sempre reseta N
		c.Flags.H = (a&0x0F + b&0x0F) > 0x0F
		c.Flags.C = result > 0xFF

		c.SetA(byte(result & 0xFF))
		c.updateFlags()
	default:

		c.PC += 1
		panic(fmt.Sprintf("opcode não implementado: 0x%02X no endereço 0x%04X", opcode, c.PC-1))

	}

}

func (c *CPU) SetFlagsADD(result uint16, operand byte) {
	// Flag Zero resultado é zero
	c.Flags.Z = (result & 0xFF) == 0

	// Flag Half-Carry nos primeiros 4 bits
	c.Flags.H = (c.A()&0x0F)+(operand&0x0F) > 0x0F

	// Flag Carry resultado ultrapassou 8 bits
	c.Flags.C = result > 0xFF
}

type Memory struct {
	RAM []byte // 32 KB de RAM
	// ... (ROM, VRAM, I/O)
}

func (m *Memory) Read(addr uint16) byte {
	return m.RAM[addr]
}

func (m *Memory) Write(addr uint16, value byte) {
	m.RAM[addr] = value
}

func (c *CPU) updateFlags() {
	var f byte
	if c.Flags.Z {
		f |= 0x80
	} // Bit 7
	if c.Flags.N {
		f |= 0x40
	} // Bit 6
	if c.Flags.H {
		f |= 0x20
	} // Bit 5
	if c.Flags.C {
		f |= 0x10
	} // Bit 4
	c.AF = (c.AF & 0xFF00) | uint16(f) // Mantém A intacto
}

// Lê as flags do registrador F
func (c *CPU) readFlags() {
	f := byte(c.AF & 0xFF)
	c.Flags.Z = (f & 0x80) != 0
	c.Flags.N = (f & 0x40) != 0
	c.Flags.H = (f & 0x20) != 0
	c.Flags.C = (f & 0x10) != 0
}

// / Getters
func (c *CPU) A() byte { return byte(c.AF >> 8) }
func (c *CPU) F() byte { return byte(c.AF & 0xFF) }
func (c *CPU) B() byte { return byte(c.BC >> 8) }
func (c *CPU) C() byte { return byte(c.BC & 0xFF) }
func (c *CPU) D() byte { return byte(c.DE >> 8) }
func (c *CPU) E() byte { return byte(c.DE & 0xFF) }
func (c *CPU) H() byte { return byte(c.HL >> 8) }
func (c *CPU) L() byte { return byte(c.HL & 0xFF) }

// Setters
func (c *CPU) SetA(value byte) { c.AF = (uint16(value) << 8) | (c.AF & 0x00FF) }
func (c *CPU) SetF(value byte) { c.AF = (c.AF & 0xFF00) | uint16(value) }
func (c *CPU) SetB(value byte) { c.BC = (uint16(value) << 8) | (c.BC & 0x00FF) }
func (c *CPU) SetC(value byte) { c.BC = (c.BC & 0xFF00) | uint16(value) }
func (c *CPU) SetD(value byte) { c.DE = (uint16(value) << 8) | (c.DE & 0x00FF) }
func (c *CPU) SetE(value byte) { c.DE = (c.DE & 0xFF00) | uint16(value) }
func (c *CPU) SetH(value byte) { c.HL = (uint16(value) << 8) | (c.HL & 0x00FF) }
func (c *CPU) SetL(value byte) { c.HL = (c.HL & 0xFF00) | uint16(value) }

func main() {
	//teste

	mem := &Memory{
		RAM: make([]byte, 0xFFFF), //64KB de RAM
	}
	mem.RAM[0x0000] = 0x80

	cpu := CPU{Memory: mem}
	cpu.SetA(0x10)
	cpu.SetB(0x20)
	fmt.Printf("A antes: 0x%02X\n", cpu.A())

	// Executa ADD A, B (opcode 0x80)
	cpu.Step()

	fmt.Printf("A depois: 0x%02X\n", cpu.A())
	fmt.Printf("Flags: Z=%v, H=%v, C=%v\n", cpu.Flags.Z, cpu.Flags.H, cpu.Flags.C)
}
