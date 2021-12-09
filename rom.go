package ines

type Rom struct {
	Headerless         []byte // Romdump without the header
	Header             []byte // Added by a person, either iNES or iNES 2.0. Required by emulators.
	Trainer            []byte // Hacks and stuff
	ProgramRom         []byte // Memory chip connected to the CPU. Contains the code.
	CharacterRom       []byte // Memory chip connected to the PPU. Contains a fixed set of graphics tile data.
	HasBattery         bool   // Rare: There may be an additional chip like that to hold even more data.
	ProgramRam         []byte
	MiscRom            []byte
	Mapper             int
	SubMapper          int
	ConsoleType        string
	Title              []byte
	TVSystem           string
	Mirroring          string
	VsSystemPPU        string
	VsSystemType       string
	CPUPPUTiming       string
	ExpansionDevice    string
	CharacterRam       []byte
	CharacterNVRam     []byte
	ProgramNVRam       []byte	// EEPROM/Non-volatile Program RAM
}

func Decode(b []byte) Rom {
	rom := identifyFmt(b)
	return rom
}
