package chip8

import "io/ioutil"

type Memory struct {
	Memory [4096]byte
}

func (memory *Memory) LoadRom(romPath string) error {
	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(rom); i++ {
		memory.Memory[0x200+i] = rom[i]
	}
	return nil
}
