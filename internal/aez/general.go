package aez

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/ccsds"
)

// Specification describes what version the current implementation follows
var Specification string = "AEZICD002:E"

//STAT General status housekeeping report of the payload instrument.
type STAT struct { //(34 octets)
	SPID   uint16 // Software Part ID
	SPREV  uint8  // Software Part Revision
	FPID   uint16 // Firmware Part ID
	FPREV  uint8  // Firmware Part Revision
	TS     uint32 // Time, seconds (CUC time format)
	TSS    uint16 // Time, subseconds (CUC time format)
	MODE   uint8  // Payload mode 1..2
	EDACE  uint32 // EDAC detected single bit errors
	EDACCE uint32 // EDAC corrected single bit errors
	EDACN  uint32 // EDAC memory scrubber passes through memory
	SPWEOP uint32 // SpaceWire received EOPs
	SPWEEP uint32 // SpaceWire received EEPs
}

// Read STAT
func (stat *STAT) Read(buf io.Reader) error {
	return binary.Read(buf, binary.BigEndian, stat)
}

// Time returns the measurement time in UTC
func (stat *STAT) Time(epoch time.Time) time.Time {
	return ccsds.UnsegmentedTimeDate(stat.TS, stat.TSS, epoch)
}

// Nanoseconds returns the measurement time in nanoseconds since epoch
func (stat *STAT) Nanoseconds() int64 {
	return ccsds.UnsegmentedTimeNanoseconds(stat.TS, stat.TSS)
}
