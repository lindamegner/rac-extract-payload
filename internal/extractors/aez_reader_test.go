package extractors

import (
	"io"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

func TestDecodeAEZ(t *testing.T) {

	tests := []struct {
		name    string
		arg     common.DataRecord
		want    common.Exportable
		wantErr bool
	}{
		{
			"Package with error",
			common.DataRecord{Error: io.EOF},
			nil,
			true,
		},

		{
			"STAT package",
			common.DataRecord{
				SourceHeader: innosat.SourcePacketHeader{PacketID: 0x0864, PacketSequenceControl: 0xc89a, PacketLength: 0x31},
				TMHeader:     innosat.TMDataFieldHeader{PUS: 16, ServiceType: 3, ServiceSubType: 0x19, CUCTimeSeconds: 0, CUCTimeFraction: 0},
				Buffer: []byte{0x00, 0x01, 0x7f, 0x04, 0x02, 0x82,
					0x04, 0x02, 0x02, 0x06, 0x01, 0x1b, 0x12, 0x00, 0x00, 0xef, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0xa5, 0xd5,
				},
			},
			aez.STAT{SPID: 32516, SPREV: 2, FPID: 33284, FPREV: 2, SVNA: 2, SVNB: 6, SVNC: 1, TS: 454164480, TSS: 61344, MODE: 2, EDACE: 0, EDACCE: 0, EDACN: 16777216, SPWEOP: 1090519040, SPWEEP: 0, ANOMALY: 0},
			false,
		},
		{
			"Bad package",
			common.DataRecord{
				SourceHeader: innosat.SourcePacketHeader{PacketID: 0x0864, PacketSequenceControl: 0xc89a, PacketLength: 0x31},
				TMHeader:     innosat.TMDataFieldHeader{PUS: 16, ServiceType: 3, ServiceSubType: 0x19, CUCTimeSeconds: 0, CUCTimeFraction: 0},
				Buffer: []byte{0x00, 0x01, 0x7f, 0x04, 0x02, 0x82,
					0x04, 0x02, 0x02, 0x06, 0x01, 0x1b, 0x12, 0x00, 0x00, 0xef, 0xa0, 0x02, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			aez.STAT{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := make(chan common.DataRecord)
			target := make(chan common.DataRecord)
			go DecodeAEZ(target, source)
			source <- tt.arg
			close(source)
			got := <-target

			if (got.Error != nil) != tt.wantErr {
				t.Errorf("DataRecord.Error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want) {
				t.Errorf("DataRecord.Buffer = %v, want %v", got.Data, tt.want)
			}
		})
	}
}
