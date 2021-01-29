package framelog

import (
	"github.com/pkg/errors"

	"github.com/ibrahimozekici/chirpstack-api/go/v4/ns"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/models"
)

// CreateUplinkFrameLog creates a UplinkFrameLog.
func CreateUplinkFrameLog(rxPacket models.RXPacket) (ns.UplinkFrameLog, error) {
	b, err := rxPacket.PHYPayload.MarshalBinary()
	if err != nil {
		return ns.UplinkFrameLog{}, errors.Wrap(err, "marshal phypayload error")
	}

	return ns.UplinkFrameLog{
		PhyPayload: b,
		TxInfo:     rxPacket.TXInfo,
		RxInfo:     rxPacket.RXInfoSet,
	}, nil
}
