package downlink

import (
	"time"

	"github.com/pkg/errors"

	"github.com/ibrahimozekici/chirpstack-network-server/internal/config"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/downlink/data"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/downlink/join"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/downlink/multicast"
	"github.com/ibrahimozekici/chirpstack-network-server/internal/downlink/proprietary"
)

var (
	schedulerBatchSize = 100
	schedulerInterval  time.Duration
)

// Setup sets up the downlink.
func Setup(conf config.Config) error {
	nsConfig := conf.NetworkServer
	schedulerInterval = nsConfig.Scheduler.SchedulerInterval

	if err := data.Setup(conf); err != nil {
		return errors.Wrap(err, "setup downlink/data error")
	}

	if err := join.Setup(conf); err != nil {
		return errors.Wrap(err, "setup downlink/join error")
	}

	if err := multicast.Setup(conf); err != nil {
		return errors.Wrap(err, "setup downlink/multicast error")
	}

	if err := proprietary.Setup(conf); err != nil {
		return errors.Wrap(err, "setup downlink/proprietary error")
	}

	return nil
}
