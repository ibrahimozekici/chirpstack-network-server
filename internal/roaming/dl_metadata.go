package roaming

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/brocaar/lorawan/backend"
	"github.com/ibrahimozekici/chirpstack-api/go/v4/gw"
)

// DLMetaDataToRXInfoSet returns the UplinkRXInfo set from the given DLMetaData.
func DLMetaDataToUplinkRXInfoSet(dl backend.DLMetaData) ([]*gw.UplinkRXInfo, error) {
	var out []*gw.UplinkRXInfo

	for _, gwInfo := range dl.GWInfo {
		var rxInfo gw.UplinkRXInfo
		if err := proto.Unmarshal(gwInfo.ULToken[:], &rxInfo); err != nil {
			return nil, errors.Wrap(err, "protobuf unmarshal error")
		}

		out = append(out, &rxInfo)
	}

	return out, nil
}
