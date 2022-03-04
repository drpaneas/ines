package ines

import (
	"encoding/binary"
	"fmt"
	"math"
)

// nolint: gomnd
func parseINES2(b []byte) Rom {
	headerless := b[16:] // without header
	header := b[:16]     // header 16 bytes
	trainer := getTrainer2(b, header)
	prgrom := getPrgRom2(header, headerless, trainer)
	chrrom := getChrRom2(header, headerless, trainer, prgrom)
	miscrom := getMiscRom(header, trainer, prgrom, chrrom, headerless)
	mapper, subMapper := getMappers(header)
	mirroring := getMirroring2(header)
	hasBattery, prgnvram := getPrgNVRamIfHasBattery(header)
	consoleType := getConsoleType(header)
	vsSystemPPU, vsSystemType, consoleType := getPPUSystemAndConsoleTypes(header, consoleType)
	programRAM := getProgramRAM(header)
	_, chrram := getChrRAMAndShiftCount(header)
	chrnvram := getChrNVRam(header)
	cpuPPUTiming := byteToInt(header[12] & 0b0000011)
	tvSystem, cpuppuTiming := getTvSystemAndCPUPpuTiming(cpuPPUTiming)
	expansionDevice := getDefaultExpansionDevice(header[15] & 0b00111111)

	return Rom{
		HeaderType:      "iNES 2.0",
		Headerless:      headerless,
		Header:          header,
		Trainer:         trainer,
		ProgramRom:      prgrom,
		CharacterRom:    chrrom,
		MiscRom:         miscrom,
		HasBattery:      hasBattery,
		programRAM:      programRAM,
		CharacterRAM:    chrram,
		ProgramNVRam:    prgnvram,
		CharacterNVRam:  chrnvram,
		Mapper:          mapper,
		SubMapper:       subMapper,
		ConsoleType:     consoleType,
		TVSystem:        tvSystem,
		Mirroring:       mirroring,
		VsSystemPPU:     vsSystemPPU,
		VsSystemType:    vsSystemType,
		CPUPPUTiming:    cpuppuTiming,
		ExpansionDevice: expansionDevice,
	}
}

// nolint: gomnd
func getTvSystemAndCPUPpuTiming(cpuPPUTiming int) (string, string) {
	var msgTV, msgCPU string
	switch cpuPPUTiming {
	case 0:
		msgCPU = "RP2C02 (\"NTSC NES\")"
		msgTV = "North America, Japan, South Korea, Taiwan"
	case 1:
		msgCPU = "RP2C07 (\"Licensed PAL NES\")"
		msgTV = "Western Europe, Australia"
	case 2:
		msgCPU = "Multiple-region"
		msgTV = "Identical ROM content in both NTSC and PAL countries"
	case 3:
		msgCPU = "UMC 6527P (\"Dendy\")"
		msgTV = "Eastern Europe, Russia, Mainland China, India, Africa"
	default:
		msgCPU = "Unknown"
		msgTV = "Unknown"
	}

	return msgTV, msgCPU
}

// nolint: gomnd
func getChrNVRam(header []byte) []byte {
	var chrnvramSize int
	shiftCount := int(readHighNibbleByte(header[11]))
	if shiftCount != 0 {
		chrnvramSize = 64 << shiftCount
	}
	var chrnvram = make([]byte, chrnvramSize)
	return chrnvram
}

// getChrRAMAndShiftCount
// In NES 1.0 an emulator assumes that a ROM image without CHR-ROM, automatically has 8 KiB of CHR-RAM;
// But in NES 2.0 all CHR-RAM must instead be explicitly specified in Header byte 11.
// nolint: gomnd
func getChrRAMAndShiftCount(header []byte) (int, []byte) {
	var chrramSize int // If the shift count is zero, there is no CHR-(NV)RAM
	shiftCount := int(readLowNibbleByte(header[11]))
	if shiftCount != 0 {
		chrramSize = 64 << shiftCount // i.e. that is 8192 bytes for a shift count of 7.
	}
	var chrram = make([]byte, chrramSize)
	return shiftCount, chrram
}

// nolint: gomnd
func getProgramRAM(header []byte) []byte {
	var sizePrgram int
	if readLowNibbleByte(header[10]) != 0 {
		sizePrgram = 64 << readLowNibbleByte(header[10])
	}
	var programRAM = make([]byte, sizePrgram)
	return programRAM
}

// nolint: gomnd
func getPPUSystemAndConsoleTypes(header []byte, consoleType string) (string, string, string) {
	var vsSystemPPU, vsSystemType string
	if hasBit(header[7], 0) && hasBit(header[7], 1) {
		consoleType = fmt.Sprintf("%s", getExtendedConsoleType(header[7]&0b00000011)) // take bit 0 and 1
		// If it's an extended console then the Vs. System Type has the following PPU and Hardware Type
		vsSystemPPU = getVsPPUType(readLowNibbleByte(header[13]))
		vsSystemType = getVsSystemType(readHighNibbleByte(header[13]))
	}
	return vsSystemPPU, vsSystemType, consoleType
}

func getConsoleType(header []byte) string {
	var consoleType string
	if !hasBit(header[7], 0) && !hasBit(header[7], 1) {
		consoleType = nes
	}
	if hasBit(header[7], 0) && !hasBit(header[7], 1) {
		consoleType = vs
	}
	if !hasBit(header[7], 0) && hasBit(header[7], 1) {
		consoleType = playchoice
	}
	return consoleType
}

// nolint: gomnd
func getPrgNVRamIfHasBattery(header []byte) (bool, []byte) {
	hasBattery := false
	var sizeProgramNVRam int // If the shift count is zero, PRG-NVRAM or EEPROM (non-volatile) is zero
	if hasBit(header[6], 1) {
		hasBattery = true
		shiftCount := int(readHighNibbleByte(header[10]))
		if shiftCount != 0 {
			sizeProgramNVRam = 64 << shiftCount // i.e. that is 8192 bytes for a shift count of 7.
		}
	}
	var prgnvram = make([]byte, sizeProgramNVRam)
	return hasBattery, prgnvram
}

// getMirroring2
// Header Byte 6 bit 0 is relevant only if the mapper does not allow the mirroring type to be switched.
// Otherwise, it must be ignored and should be set to zero.
// nolint: gomnd
func getMirroring2(header []byte) string {
	mirroring := "Ignored"
	if hasBit(header[6], 3) {
		mirroring = "Four-screen" //  Ignore mirroring control and the mirroring bit
	} else {
		mirroring = "Horizontal or mapper-controlled"
		if hasBit(header[6], 0) {
			mirroring = "Vertical"
		}
	}
	return mirroring
}

func getMappers(header []byte) (int, int) {
	mapper1 := readHighNibbleByte(header[6]) // Lower bits of mapper
	mapper2 := readHighNibbleByte(header[7]) // Upper bits of mapper
	mapper3 := readLowNibbleByte(header[8])
	mapper := int(binary.LittleEndian.Uint16([]byte{mergeNibbles(mapper2, mapper1), mapper3}))
	// SubMapper number
	subMapper := int(readHighNibbleByte(header[8]))
	return mapper, subMapper
}

/* getMiscRom
Miscellaneous ROM Area
----------------------
The Miscellaneous ROM Area, if present, follows the CHR-ROM area and occupies the remainder of the file.
Its size is not explicitly denoted in the header, and can be deduced by subtracting
the 16-byte Header, Trainer, PRG-ROM and CHR-ROM Area sizes from the total file size.
The meaning of this data depends on the console type and mapper type;

Header byte 14 is used to denote the presence of the Miscellaneous ROM Area and
the number of ROM chips in case any disambiguation is needed.
*/
// nolint: gomnd
func getMiscRom(header []byte, trainer []byte, prgrom []byte, chrrom []byte, headerless []byte) []byte {
	var miscrom []byte
	if (header[14] & 0b00000011) != 0 {
		start := len(trainer) + len(prgrom) + len(chrrom)
		miscrom = headerless[start:]
	}
	return miscrom
}

/*	getChrRom2
	CHR-ROM Area
	------------
	The CHR-ROM Area, if present, follows the Trainer and PRG-ROM Areas and precedes the Miscellaneous ROM Area.
	Header byte 5 (LSB) and bits 4-7 of Header byte 9 (MSB) specify its size.
	If the MSB nibble is $0-E, LSB and MSB together simply specify the CHR-ROM size in 8 KiB units:
*/
// nolint: gomnd
func getChrRom2(header []byte, headerless []byte, trainer []byte, prgrom []byte) []byte {
	var sizeChrrom int
	var chrrom []byte
	MSBNibbleByte9 := readHighNibbleByte(header[9])
	if byteToHex(MSBNibbleByte9) == "0F" {
		E := (header[5] & 0b11111100) >> 2
		MM := header[5] & 0b00000011
		sizeChrrom = int(math.Pow(2, float64(E))) * (byteToInt(MM)*2 + 1)
	} else {
		tmp := fmt.Sprintf("%v%v", byteToHex(MSBNibbleByte9), byteToHex(header[5]))
		sizeChrrom = hexToInt(tmp) * 8 * 1024
	}

	chrrom = headerless[len(trainer)+len(prgrom) : len(trainer)+len(prgrom)+sizeChrrom]
	return chrrom
}

/*	getPrgRom2
	PRG-ROM Area
	------------
	The PRG-ROM Area follows the 16-byte Header and the Trainer Area and precedes the CHR-ROM Area.
	Header byte 4 (LSB) and bits 0-3 of Header byte 9 (MSB) together specify its size.
	If the MSB nibble is $0-E, LSB and MSB together simply specify the PRG-ROM size in 16 KiB units:

	// Size of PRG ROM in 16 KB units
	// The PRG-ROM Area follows the 16-byte Header and the Trainer Area (if exists) and precedes the CHR-ROM Area.
*/
// nolint: gomnd
func getPrgRom2(header []byte, headerless []byte, trainer []byte) []byte {
	var prgrom []byte
	var sizeOfPrgRom int
	MSNibbleByte9 := readLowNibbleByte(header[9])
	if byteToHex(MSNibbleByte9) == "0F" {
		E := (header[4] & 0b11111100) >> 2
		MM := header[4] & 0b00000011
		sizeOfPrgRom = int(math.Pow(2, float64(E))) * (byteToInt(MM)*2 + 1)
	} else {
		tmp := fmt.Sprintf("%v%v", byteToHex(MSNibbleByte9), byteToHex(header[4]))
		sizeOfPrgRom = hexToInt(tmp) * 16 * 1024
	}

	prgrom = headerless[len(trainer) : len(trainer)+sizeOfPrgRom] // if trainer is 0, this will still work
	return prgrom
}

/*	getTrainer2
	Trainer Area
	------------
	Trainer exists if bit 2 of Header byte 6 is set.
	It contains data to be loaded into CPU memory at 0x7000
	It is only used by some games that were modified to run on different hardware from the original cartridges,
	such as early RAM cartridges and emulators, adding some compatibility code into those address ranges.
	Trainer is placed between header and PRG ROM data, so PRG ROM should start in the next avail address
*/
// nolint: gomnd
func getTrainer2(b []byte, header []byte) []byte {
	var trainer []byte
	if hasBit(header[6], 2) {
		low := 16         // the Trainer Area follows the 16-byte Header and precedes the PRG-ROM area
		high := low + 512 // trainer has always fixed 512 bytes size
		trainer = b[low:high]
	}
	return trainer
}
