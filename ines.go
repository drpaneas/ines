package ines

// An iNES file consists of the following sections, in order:

// 1. Header (16 bytes)
// 2. Trainer, if present (0 or 512 bytes)
// 3. PRG ROM data (16384 * x bytes)
// 4. CHR ROM data, if present (8192 * y bytes)
// 5. PlayChoice INST-ROM, if present (0 or 8192 bytes)
// 6. PlayChoice PROM, if present (16 bytes Data, 16 bytes CounterOut) (this is often missing, see PC10 ROM-Images for details)
// 7. Some ROM-Images additionally contain a 128-byte (or sometimes 127-byte) title at the end of the file.

func parseINES(b []byte) Rom {
	headerless := b[16:] // rom without header. It's useful for calculating checksums.
	header := b[:16]     // header (16 bytes)
	trainer := getTrainer(b, header)
	prgrom := getPrgRom(header, headerless, trainer)
	chrrom, sizeChrrom := getChrRomAndSize(header, headerless, trainer, prgrom)
	chrram := getChrRam(sizeChrrom)
	consoleType, playChoiceInstRom, playChoicePROMData, playChoiceRomCounterOut := getConsoleTypes(header, headerless, trainer, prgrom, chrrom)
	title := getTitle(headerless, trainer, prgrom, chrrom, playChoiceInstRom, playChoicePROMData, playChoiceRomCounterOut)
	hasBatteryPrgRam, prgram := getPrgRamIfHasBattery(header)
	mapper := getMapper(header)
	mirroring := getMirroring(header)
	tvSystem := getTvSystem(header)

	return Rom{
		HeaderType:      "iNES 1.0",
		Headerless:      headerless,
		Header:          header,
		Trainer:         trainer,
		ProgramRom:      prgrom,
		CharacterRom:    chrrom,
		HasBattery:      hasBatteryPrgRam,
		ProgramRam:      prgram,
		MiscRom:         []byte{},
		Mapper:          mapper,
		SubMapper:       0,
		ConsoleType:     consoleType,
		Title:           title,
		TVSystem:        tvSystem,
		Mirroring:       mirroring,
		VsSystemPPU:     "Unknown",
		VsSystemType:    "Unknown",
		CPUPPUTiming:    "Unknown",
		ExpansionDevice: "Unknown",
		CharacterRam:    chrram,
		CharacterNVRam:  []byte{},
		ProgramNVRam:    []byte{},
	}
}

// getTitles fetches the game's name
// Some ROM-Images additionally contain a 128-byte (or sometimes 127-byte) title at the end of the file.
// This is an additional block of data appended to the end of the rom that's not read by the emulator and has the title
// of the game in ascii/half-width katakana and so would be ideal to categorize ROMs.
// Sadly this idea never really was applied for most rom dumps but if one really wanted they could go through their collection of roms and add that title block so that future tools parsing such info could find something inside the rom to recognize.

func getTitle(headerless []byte, trainer []byte, prgrom []byte, chrrom []byte, playChoiceInstRom []byte, playChoicePROMData []byte, playChoiceRomCounterOut []byte) []byte {
	var title []byte
	if leftover := len(headerless) - len(trainer) + len(prgrom) + len(chrrom) + len(playChoiceInstRom) + len(playChoicePROMData) + len(playChoiceRomCounterOut); leftover != 0 {
		title = headerless[len(trainer)+len(prgrom)+len(chrrom)+len(playChoiceInstRom)+len(playChoicePROMData)+len(playChoiceRomCounterOut):]
	}
	return title
}

// getConsoleTypes fetches several console types
// PlayChoice-10 INST-ROM Area is not part of the official specification, and most emulators simply ignores it.
// PlayChoice games are designed to look good with the 2C03 RGB PPU
// which handles color emphasis differently from a standard NES PPU.
// The detection of which palette a particular game uses is left unspecified.
// If present, it's 8192 bytes.
func getConsoleTypes(header []byte, headerless []byte, trainer []byte, prgrom []byte, chrrom []byte) (string, []byte, []byte, []byte) {
	var consoleType = nes
	var playChoiceInstRom, playChoicePROMData, playChoiceRomCounterOut []byte
	if hasBit(header[7], 1) {
		consoleType = playchoice
		playChoiceInstRomSize := 8192
		// 8 KB INST ROM (containing data and Z80 code for instruction screens)
		playChoiceInstRom = headerless[len(trainer)+len(prgrom)+len(chrrom) : len(trainer)+len(prgrom)+len(chrrom)+playChoiceInstRomSize]

		// PlayChoice PROM , if present (16 bytes Data, 16 bytes CounterOut)
		// 16 bytes RP5H01 PROM Data output (needed to decrypt the INST ROM)
		playChoicePROMDataSize := 8192 * 2
		playChoicePROMData = headerless[len(trainer)+len(prgrom)+len(chrrom)+len(playChoiceInstRom) : len(trainer)+len(prgrom)+len(chrrom)+len(playChoiceInstRom)+playChoicePROMDataSize]

		// 16 bytes RP5H01 PROM CounterOut output (needed to decrypt the INST ROM) (usually constant: 00,00,00,00,FF,FF,FF,FF,00,00,00,00,FF,FF,FF,FF)
		playChoiceRomCounterOutSize := 8192 * 2
		playChoiceRomCounterOut = headerless[len(trainer)+len(prgrom)+len(chrrom)+len(playChoiceInstRom)+len(playChoicePROMData) : len(trainer)+len(prgrom)+len(chrrom)+len(playChoiceInstRom)+len(playChoicePROMData)+playChoiceRomCounterOutSize]
	}
	if hasBit(header[7], 0) {
		consoleType = vs
	}
	return consoleType, playChoiceInstRom, playChoicePROMData, playChoiceRomCounterOut
}

// getPrgRom data (16384 * x bytes)
// The PRG-ROM Area follows the Header and the Trainer and precedes the CHR-ROM Area.
// Size of Program ROM (in 16 KB units).
func getPrgRom(header []byte, headerless []byte, trainer []byte) []byte {
	var prgrom []byte
	sizePrgrom := int(header[4]) * 16384
	prgrom = headerless[len(trainer) : len(trainer)+sizePrgrom] // if trainer is 0, this will still work
	return prgrom
}

// getTrainer exists if bit 2 of Header byte 6 is set.
// It contains data to be loaded into CPU memory at 0x7000
// It is only used by some games that were modified to run on different hardware from the original cartridges,
// such as early RAM cartridges and emulators, adding some compatibility code into those address ranges.
func getTrainer(b []byte, header []byte) []byte {
	var trainer []byte
	if hasBit(header[6], 2) {
		trainer = b[16:528] // starts from b[16] and has 512 bytes length, so it goes up to b[16+512]
	}
	return trainer
}

// getChrRomAndSize The CHR-ROM Area, if present, follows the Trainer and PRG-ROM Areas and precedes the PlayChoice INST-ROM Area.
// CHR ROM data, if present (8192 * y bytes).
// Size of Character ROM (in 8 KB units).
func getChrRomAndSize(header []byte, headerless []byte, trainer []byte, prgrom []byte) ([]byte, int) {
	var chrrom []byte
	sizeChrrom := int(header[5]) * 8192
	chrrom = headerless[len(trainer)+len(prgrom) : len(trainer)+len(prgrom)+sizeChrrom]
	return chrrom, sizeChrrom
}

// getChrRam If CHR ROM size is 0; it means the board uses 8 KB CHR RAM
// The ROM file doesn't contain RAM contents (since they'd be lost at power-off anyhow).
// CHR RAM is located at the normal place in the PPU's memory map.
// So, the CPU will set the PPU's VRAM pointer, and just start writing data to the input port.
// Games with CHR ROM need to ignore this "write". Games with CHR RAM need to accept it.
func getChrRam(sizeChrrom int) []byte {
	var sizeChrram int
	if sizeChrrom == 0 {
		sizeChrram = 8192
	}
	var chrram = make([]byte, sizeChrram)
	return chrram
}

// getTvSystem fetches the TV system
// According to the official specification very few emulators honor this bit
// virtually no ROM images in circulation make use of it.
func getTvSystem(header []byte) string {
	tvSystem := "NTSC"
	if hasBit(header[9], 0) {
		tvSystem = "PAL"
	}
	return tvSystem
}

// getPrgRamIfHasBattery fetches Battery or any other non-volatile memory (PRG RAM).
func getPrgRamIfHasBattery(header []byte) (bool, []byte) {
	hasBatteryPrgRam := false
	var prgRamBatterySize int
	if hasBit(header[6], 1) {
		hasBatteryPrgRam = true
		// The PRG RAM Size value (stored in byte 8) was recently added to the official specification;
		// as such, virtually no ROM images in circulation make use of it.
		prgRamBatterySize = 8192 // default is 8 KB
	}
	var prgram = make([]byte, prgRamBatterySize)
	return hasBatteryPrgRam, prgram
}

// getMapper fetches the mapper.
func getMapper(header []byte) int {
	lowerNibbleMapper := readHighNibbleByte(header[6])
	upperNibbleMapper := readHighNibbleByte(header[7])
	mapper := byteToInt(mergeNibbles(upperNibbleMapper, lowerNibbleMapper))
	return mapper
}

// getMirroring fetches the mirror value.
func getMirroring(header []byte) string {
	mirroring := "Ignored"
	if hasBit(header[6], 3) {
		mirroring = "Four-screen VRAM" //  Ignore mirroring control and the mirroring bit
	} else {
		mirroring = "Horizontal or mapper controlled"
		if hasBit(header[6], 0) {
			mirroring = "Vertical"
		}
	}
	return mirroring
}
