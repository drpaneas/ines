package ines

import (
	"encoding/binary"
	"fmt"
	"math"
)

func parseINES2(b []byte) Rom {
	headerless := b[16:] // without header
	header := b[:16]     // header 16 bytes

	/*	Trainer Area
		------------
		Trainer exists if bit 2 of Header byte 6 is set.
		It contains data to be loaded into CPU memory at 0x7000
		It is only used by some games that were modified to run on different hardware from the original cartridges,
		such as early RAM cartridges and emulators, adding some compatibility code into those address ranges.
		Trainer is placed between header and PRG ROM data, so PRG ROM should start in the next avail address
	*/
	var trainer []byte
	if hasBit(header[6], 2) {
		low := 16         // the Trainer Area follows the 16-byte Header and precedes the PRG-ROM area
		high := low + 512 // trainer has always fixed 512 bytes size
		trainer = b[low:high]
	}

	/*	PRG-ROM Area
		------------
		The PRG-ROM Area follows the 16-byte Header and the Trainer Area and precedes the CHR-ROM Area.
		Header byte 4 (LSB) and bits 0-3 of Header byte 9 (MSB) together specify its size.
		If the MSB nibble is $0-E, LSB and MSB together simply specify the PRG-ROM size in 16 KiB units:
	*/

	// Size of PRG ROM in 16 KB units
	// The PRG-ROM Area follows the 16-byte Header and the Trainer Area (if exists) and precedes the CHR-ROM Area.
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

	/*	CHR-ROM Area
		------------
		The CHR-ROM Area, if present, follows the Trainer and PRG-ROM Areas and precedes the Miscellaneous ROM Area.
		Header byte 5 (LSB) and bits 4-7 of Header byte 9 (MSB) specify its size.
		If the MSB nibble is $0-E, LSB and MSB together simply specify the CHR-ROM size in 8 KiB units:
	*/

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

	/* 	Miscellaneous ROM Area
	----------------------
	The Miscellaneous ROM Area, if present, follows the CHR-ROM area and occupies the remainder of the file.
	Its size is not explicitly denoted in the header, and can be deduced by subtracting
	the 16-byte Header, Trainer, PRG-ROM and CHR-ROM Area sizes from the total file size.
	The meaning of this data depends on the console type and mapper type;

	Header byte 14 is used to denote the presence of the Miscellaneous ROM Area and
	the number of ROM chips in case any disambiguation is needed.
	*/
	var miscrom []byte
	if (header[14] & 0b00000011) != 0 {
		start := len(trainer) + len(prgrom) + len(chrrom)
		miscrom = headerless[start:]
	}

	mapper1 := readHighNibbleByte(header[6]) // Lower bits of mapper
	mapper2 := readHighNibbleByte(header[7]) // Upper bits of mapper
	mapper3 := readLowNibbleByte(header[8])
	mapper := int(binary.LittleEndian.Uint16([]byte{mergeNibbles(mapper2, mapper1), mapper3}))
	// SubMapper number
	subMapper := int(readHighNibbleByte(header[8]))

	// Mirroring
	// TODO: Check the mapper here
	// Header Byte 6 bit 0 is relevant only if the mapper does not allow the mirroring type to be switched.
	// Otherwise, it must be ignored and should be set to zero.
	mirroring := "Ignored"
	if hasBit(header[6], 3) {
		mirroring = "Four-screen VRAM" //  Ignore mirroring control and the mirroring bit
	} else {
		mirroring = "Horizontal or mapper-controlled"
		if hasBit(header[6], 0) {
			mirroring = "Vertical"
		}
	}

	hasBattery := false
	var sizeEepromPrgnvram int
	if hasBit(header[6], 1) {
		hasBattery = true
		sizeEepromPrgnvram = int(readHighNibbleByte(header[10]))
	}

	// Console Type
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

	var vsSystemPPU, vsSystemType string
	if hasBit(header[7], 0) && hasBit(header[7], 1) {
		consoleType = fmt.Sprintf("Extended: %s", getExtendedConsoleType(header[7]&0b00000011)) // take bit 0 and 1
		// If it's an extended console then the Vs. System Type has the following PPU and Hardware Type
		vsSystemPPU = getVsPPUType(readLowNibbleByte(header[13]))
		vsSystemType = getVsSystemType(readHighNibbleByte(header[13]))
	}

	var sizePrgram int
	if readLowNibbleByte(header[10]) != 0 {
		sizePrgram = 64 << readLowNibbleByte(header[10])
	}

	// In NES 1.0 an emulator assumes that a ROM image without CHR-ROM, automatically has 8 KiB of CHR-RAM;
	// But in NES 2.0 all CHR-RAM must instead be explicitly specified in Header byte 11.
	chrramSize := int(readLowNibbleByte(header[11]))
	chrnvramSize := int(readHighNibbleByte(header[11]))

	cpuPPUTiming := byteToInt(header[12] & 0b0000011)
	var tvSystem, cpuppuTiming string
	if cpuPPUTiming == 0 {
		cpuppuTiming = "RP2C02 (\"NTSC NES\")"
		tvSystem = "North America, Japan, South Korea, Taiwan"
	} else if cpuPPUTiming == 1 {
		cpuppuTiming = "RP2C07 (\"Licensed PAL NES\")"
		tvSystem = "Western Europe, Australia"
	} else if cpuPPUTiming == 2 {
		cpuppuTiming = "Multiple-region"
		tvSystem = "Identical ROM content in both NTSC and PAL countries."
	} else if cpuPPUTiming == 3 {
		cpuppuTiming = "UMC 6527P (\"Dendy\")"
		tvSystem = "Eastern Europe, Russia, Mainland China, India, Africa"
	} else {
		cpuppuTiming = "Unknown"
	}

	// Default Expansion Device
	expansionDevice := getDefaultExpansionDevice(header[15] & 0b00111111)

	return Rom{
		Headerless:         headerless,
		Header:             header,
		Trainer:            trainer,
		ProgramRom:         prgrom,
		CharacterRom:       chrrom,
		HasBattery:         hasBattery,
		SizePRGRAM:         sizePrgram,
		CharacterRamSize:   sizeChrrom,
		MiscRom:            miscrom,
		Mapper:             mapper,
		SubMapper:          subMapper,
		ConsoleType:        consoleType,
		Title:              nil,
		TVSystem:           tvSystem,
		Mirroring:          mirroring,
		VsSystemPPU:        vsSystemPPU,
		VsSystemType:       vsSystemType,
		CPUPPUTiming:       cpuppuTiming,
		ExpansionDevice:    expansionDevice,
		EepromPrgnvramSize: sizeEepromPrgnvram,
		ChrramSize:         chrramSize,
		ChrnvramSize:       chrnvramSize,
	}
}
