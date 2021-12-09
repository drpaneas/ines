package ines

// https://wiki.nesdev.org/w/index.php/NES_2.0#Extended_Console_Type

const (
	playchoice = "Playchoice 10"
	vs = "Nintendo Vs. System"
	nes = "Regular NES/Famicom/Dendy"
)

func getExtendedConsoleType(consoleTypeByte uint8) string {
	switch consoleTypeByte {
	case 3:
		return "Regular Famiclone, but with CPU that supports Decimal Mode (e.g. Bit Corporation Creator)"
	case 4:
		return "V.R. Technology VT01 with monochrome palette"
	case 5:
		return "V.R. Technology VT01 with red/cyan STN palette"
	case 6:
		return "V.R. Technology VT02"
	case 7:
		return "V.R. Technology VT03"
	case 8:
		return "V.R. Technology VT09"
	case 9:
		return "V.R. Technology VT32"
	case 10:
		return "V.R. Technology VT369"
	case 11:
		return "UMC UM6578"
	default:
		return "Unknown/Undefined"
	}
}

func getVsSystemType(vsSystemTypeByte uint8) string {
	switch vsSystemTypeByte {
	case 0:
		return "Vs. Unisystem (normal)"
	case 1:
		return "Vs. Unisystem (RBI Baseball protection)"
	case 2:
		return "Vs. Unisystem (TKO Boxing protection)"
	case 3:
		return "Vs. Unisystem (Super Xevious protection)"
	case 4:
		return "Vs. Unisystem (Vs. Ice Climber Japan protection)"
	case 5:
		return "Vs. Dual System (normal)"
	case 6:
		return "Vs. Dual System (Raid on Bungeling Bay protection)"
	default:
		return "Unknown/Undefined"
	}
}

func getDefaultExpansionDevice(defaultExpansionDeviceByte uint8) string {
	switch defaultExpansionDeviceByte {
	case 0:
		return "Unspecified"
	case 1:
		return "Standard NES/Famicom controllers"
	case 2:
		return "NES Four Score/Satellite with two additional standard controllers"
	case 3:
		return "Famicom Four Players Adapter with two additional standard controllers"
	case 4:
		return "Vs. System"
	case 5:
		return "Vs. System with reversed inputs"
	case 6:
		return "Vs. Pinball (Japan)"
	case 7:
		return "Vs. Zapper"
	case 8:
		return "Zapper ($4017)"
	case 9:
		return "Two Zappers"
	case 10:
		return "Bandai Hyper Shot Lightgun"
	case 11:
		return "Power Pad Side A"
	case 12:
		return "Power Pad Side B"
	case 13:
		return "Family Trainer Side A"
	case 14:
		return "Family Trainer Side B"
	case 15:
		return "Arkanoid Vaus Controller (NES)"
	case 16:
		return "Arkanoid Vaus Controller (Famicom)"
	case 17:
		return "Two Vaus Controllers plus Famicom Data Recorder"
	case 18:
		return "Konami Hyper Shot Controller"
	case 19:
		return "Coconuts Pachinko Controller"
	case 20:
		return "Exciting Boxing Punching Bag (Blowup Doll)"
	case 21:
		return "Jissen Mahjong Controller"
	case 22:
		return "Party Tap"
	case 23:
		return "Oeka Kids Tablet"
	case 24:
		return "Sunsoft Barcode Battler"
	case 25:
		return "Miracle Piano Keyboard"
	case 26:
		return "Pokkun Moguraa (Whack-a-Mole Mat and Mallet)"
	case 27:
		return "Top Rider (Inflatable Bicycle)"
	case 28:
		return "Double-Fisted (Requires or allows use of two controllers by one player)"
	case 29:
		return "Famicom 3D System"
	case 30:
		return "Doremikko Keyboard"
	case 31:
		return "R.O.B. Gyro Set"
	case 32:
		return "Famicom Data Recorder (don't emulate keyboard)"
	case 33:
		return "ASCII Turbo File"
	case 34:
		return "IGS Storage Battle Box"
	case 35:
		return "Family BASIC Keyboard plus Famicom Data Recorder"
	case 36:
		return "Dongda PEC-586 Keyboard"
	case 37:
		return "Bit Corp. Bit-79 Keyboard"
	case 38:
		return "Subor Keyboard"
	case 39:
		return "Subor Keyboard plus mouse (3x8-bit protocol)"
	case 40:
		return "Subor Keyboard plus mouse (24-bit protocol)"
	case 41:
		return "SNES Mouse ($4017.d0)"
	case 42:
		return "Multicart"
	case 43:
		return "Two SNES controllers replacing the two standard NES controllers"
	case 44:
		return "RacerMate Bicycle"
	case 45:
		return "U-Force"
	case 46:
		return "R.O.B. Stack-Up"
	case 47:
		return "City Patrolman Lightgun"
	case 48:
		return "Sharp C1 Cassette Interface"
	case 49:
		return "Standard Controller with swapped Left-Right/Up-Down/B-A"
	case 50:
		return "Excalibor Sudoku Pad"
	case 51:
		return "ABL Pinball"
	case 52:
		return "Golden Nugget Casino extra buttons"
	default:
		return "Unknown/Undefined"
	}
}

func getVsPPUType(vsPPUTypeByte uint8) string {
	switch vsPPUTypeByte {
	case 0:
		return "RP2C03B"
	case 1:
		return "RP2C03G"
	case 2:
		return "RP2C04-0001"
	case 3:
		return "RP2C04-0002"
	case 4:
		return "RP2C04-0003"
	case 5:
		return "RP2C04-0004"
	case 6:
		return "RC2C03B"
	case 7:
		return "RC2C03C"
	case 8:
		return "RC2C05-01 ($2002 AND $?? =$1B)"
	case 9:
		return "RC2C05-02 ($2002 AND $3F =$3D)"
	case 10:
		return "RC2C05-03 ($2002 AND $1F =$1C)"
	case 11:
		return "RC2C05-04 ($2002 AND $1F =$1B)"
	case 12:
		return "RC2C05-05 ($2002 AND $1F =unknown)"
	default:
		return "Unknown/Undefined"
	}
}