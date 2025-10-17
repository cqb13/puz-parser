package puz

type checksums struct {
	checksum           uint16
	cibChecksum        uint16
	maskedLowChecksum  [4]byte
	maskedHighChecksum [4]byte
}

func computeChecksums(bytes []byte, width int, height int, title string, author string, copyright string, clues []string, notes string) *checksums {
	//cib checksum
	computedCibChecksum := checksumRegion(bytes[44:52], 0)

	// primary checksum
	computedChecksum := computedCibChecksum
	offset := 52
	computedChecksum = checksumRegion(bytes[offset:offset+int(width*height)], computedChecksum)
	offset += int(width * height)
	computedChecksum = checksumRegion(bytes[offset:offset+int(width*height)], computedChecksum)

	computedChecksum = checksumStrings(title, author, copyright, clues, notes, computedChecksum)

	// masked checksum
	checksumCIB := checksumRegion(bytes[44:52], 0x0000)
	offset = 52
	solutionChecksum := checksumRegion(bytes[offset:offset+int(width*height)], 0x0000)
	offset += int(width * height)
	stateChecksum := checksumRegion(bytes[offset:offset+int(width*height)], 0x0000)

	stringsChecksum := checksumStrings(title, author, copyright, clues, notes, 0x0000)

	maskedLowCheck := make([]byte, 4)
	maskedLowCheck[0] = 0x49 ^ byte((checksumCIB & 0xFF))
	maskedLowCheck[1] = 0x43 ^ byte(solutionChecksum&0xFF)
	maskedLowCheck[2] = 0x48 ^ byte(stateChecksum&0xFF)
	maskedLowCheck[3] = 0x45 ^ byte(stringsChecksum&0xFF)

	maskedHighCheck := make([]byte, 4)
	maskedHighCheck[0] = 0x41 ^ byte((checksumCIB>>8)&0xFF)
	maskedHighCheck[1] = 0x54 ^ byte((solutionChecksum>>8)&0xFF)
	maskedHighCheck[2] = 0x45 ^ byte((stateChecksum>>8)&0xFF)
	maskedHighCheck[3] = 0x44 ^ byte((stringsChecksum>>8)&0xFF)

	return &checksums{
		computedChecksum,
		computedCibChecksum,
		[4]byte(maskedLowCheck),
		[4]byte(maskedHighCheck),
	}
}

func checksumStrings(title string, author string, copyright string, clues []string, notes string, checksum uint16) uint16 {
	if len(title) > 0 {
		checksum = checksumRegion(append([]byte(title), 0x00), checksum)
	}

	if len(author) > 0 {
		checksum = checksumRegion(append([]byte(author), 0x00), checksum)
	}

	if len(copyright) > 0 {
		checksum = checksumRegion(append([]byte(copyright), 0x00), checksum)
	}

	for _, clue := range clues {
		checksum = checksumRegion([]byte(clue), checksum)
	}

	if len(notes) > 0 {
		checksum = checksumRegion(append([]byte(notes), 0x00), checksum)
	}

	return checksum
}

func checksumRegion(buffer []byte, checksum uint16) uint16 {
	for i := range buffer {
		checksum = (checksum >> 1) | ((checksum & 1) << 15)
		checksum += uint16(buffer[i])
	}

	return checksum
}
