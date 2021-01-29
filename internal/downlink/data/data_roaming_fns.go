package data

import (
	"context"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/brocaar/lorawan/backend"
	"github.com/ibrahimozekici/chirpstack-api/go/v4/gw"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/backend/gateway"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/band"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/helpers"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/logging"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/roaming"
)

// HandleRoamingFNS handles a downlink as fNS.
func HandleRoamingFNS(ctx context.Context, pl backend.XmitDataReqPayload) error {
	// Retrieve RXInfo from DLMetaData
	rxInfo, err := roaming.DLMetaDataToUplinkRXInfoSet(*pl.DLMetaData)
	if err != nil {
		return errors.Wrap(err, "get uplink rxinfo error")
	}

	if len(rxInfo) == 0 {
		return errors.New("GWInfo must not be empty")
	}

	sort.Sort(bySignal(rxInfo))

	var downID uuid.UUID
	if ctxID := ctx.Value(logging.ContextIDKey); ctxID != nil {
		if id, ok := ctxID.(uuid.UUID); ok {
			downID = id
		}
	}

	downlink := gw.DownlinkFrame{
		GatewayId:  rxInfo[0].GatewayId,
		DownlinkId: downID[:],
		Items:      []*gw.DownlinkFrameItem{},
	}

	if pl.DLMetaData.DLFreq1 != nil && pl.DLMetaData.DataRate1 != nil && pl.DLMetaData.RXDelay1 != nil {
		item := gw.DownlinkFrameItem{
			PhyPayload: pl.PHYPayload[:],
			TxInfo: &gw.DownlinkTXInfo{
				Frequency: uint32(*pl.DLMetaData.DLFreq1 * 1000000),
				Board:     rxInfo[0].Board,
				Antenna:   rxInfo[0].Antenna,
				Context:   rxInfo[0].Context,
				Timing:    gw.DownlinkTiming_DELAY,
				TimingInfo: &gw.DownlinkTXInfo_DelayTimingInfo{
					DelayTimingInfo: &gw.DelayTimingInfo{
						Delay: ptypes.DurationProto(time.Duration(*pl.DLMetaData.RXDelay1) * time.Second),
					},
				},
			},
		}

		item.TxInfo.Power = int32(band.Band().GetDownlinkTXPower(int(item.TxInfo.Frequency)))

		if err := helpers.SetDownlinkTXInfoDataRate(item.TxInfo, *pl.DLMetaData.DataRate1, band.Band()); err != nil {
			return errors.Wrap(err, "set downlink txinfo data-rate error")
		}

		downlink.Items = append(downlink.Items, &item)
	}

	if pl.DLMetaData.DLFreq2 != nil && pl.DLMetaData.DataRate2 != nil && pl.DLMetaData.RXDelay1 != nil {
		item := gw.DownlinkFrameItem{
			PhyPayload: pl.PHYPayload[:],
			TxInfo: &gw.DownlinkTXInfo{
				Frequency: uint32(*pl.DLMetaData.DLFreq2 * 1000000),
				Board:     rxInfo[0].Board,
				Antenna:   rxInfo[0].Antenna,
				Context:   rxInfo[0].Context,
				Timing:    gw.DownlinkTiming_DELAY,
				TimingInfo: &gw.DownlinkTXInfo_DelayTimingInfo{
					DelayTimingInfo: &gw.DelayTimingInfo{
						Delay: ptypes.DurationProto(time.Duration(*pl.DLMetaData.RXDelay1+1) * time.Second),
					},
				},
			},
		}

		item.TxInfo.Power = int32(band.Band().GetDownlinkTXPower(int(item.TxInfo.Frequency)))

		if err := helpers.SetDownlinkTXInfoDataRate(item.TxInfo, *pl.DLMetaData.DataRate2, band.Band()); err != nil {
			return errors.Wrap(err, "set downlink txinfo data-rate error")
		}

		downlink.Items = append(downlink.Items, &item)
	}

	if err := gateway.Backend().SendTXPacket(downlink); err != nil {
		return errors.Wrap(err, "send downlink-frame to gateway error")
	}

	return nil
}

type bySignal []*gw.UplinkRXInfo

func (s bySignal) Len() int {
	return len(s)
}

func (s bySignal) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySignal) Less(i, j int) bool {
	if s[i].LoraSnr == s[j].LoraSnr {
		return s[i].Rssi > s[j].Rssi
	}

	return s[i].LoraSnr > s[j].LoraSnr
}
