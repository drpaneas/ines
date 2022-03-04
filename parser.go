package ines

// identifyFmt activates the appropriate section format.
// It returns nil if no format was identified.
// nolint: exhaustivestruct, varnamelen
func identifyFmt(b []byte) Rom {
	if hasHeader(b) {
		if isINES2(b) {
			return parseINES2(b)
		}

		return parseINES(b)
	}

	return Rom{}
}
