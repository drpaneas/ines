package ines

type Rom struct {
	HeaderType      string
	Headerless      []byte // Romdump without the header
	Header          []byte // Added by a person, either iNES or iNES 2.0. Required by emulators.
	Trainer         []byte // Hacks and stuff
	ProgramRom      []byte // Memory chip connected to the CPU. Contains the code.
	CharacterRom    []byte // Memory chip connected to the PPU. Contains a fixed set of graphics tile data.
	HasBattery      bool   // Rare: There may be an additional chip like that to hold even more data.
	ProgramRam      []byte
	MiscRom         []byte
	Mapper          int
	SubMapper       int
	ConsoleType     string
	Title           []byte
	TVSystem        string
	Mirroring       string
	VsSystemPPU     string
	VsSystemType    string
	CPUPPUTiming    string
	ExpansionDevice string
	CharacterRam    []byte
	CharacterNVRam  []byte
	ProgramNVRam    []byte // EEPROM/Non-volatile Program RAM
}

func Decode(b []byte) Rom {
	// Default values
	var rom = Rom{
		HeaderType:      "Unknown",
		Headerless:      []byte{},
		Header:          []byte{},
		Trainer:         []byte{},
		ProgramRom:      []byte{},
		CharacterRom:    []byte{},
		HasBattery:      false,
		ProgramRam:      []byte{},
		MiscRom:         []byte{},
		Mapper:          0,
		SubMapper:       0,
		ConsoleType:     "Unknown",
		Title:           []byte{},
		TVSystem:        "Unknown",
		Mirroring:       "Unknown",
		VsSystemPPU:     "Unknown",
		VsSystemType:    "Unknown",
		CPUPPUTiming:    "Unknown",
		ExpansionDevice: "Unknown",
		CharacterRam:    []byte{},
		CharacterNVRam:  []byte{},
		ProgramNVRam:    []byte{},
	}

	// Parse it properly (iNES 1.0 or iNES 2.0)
	rom = identifyFmt(b)
	return rom
}
