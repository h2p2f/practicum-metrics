// Code generated by gowrap. DO NOT EDIT.
// template: ../../../../../templates/gowrap/zap
// gowrap: http://github.com/hexdigest/gowrap

package updatesmetrics

//go:generate gowrap gen -p github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatesmetrics -i Updater -t ../../../../../templates/gowrap/zap -o updatesmetrics_with_logging.go -l ""

import (
	"go.uber.org/zap"
)

// UpdaterWithZap implements Updater that is instrumented with zap logger
type UpdaterWithZap struct {
	_log  *zap.Logger
	_base Updater
}

// NewUpdaterWithZap instruments an implementation of the Updater with simple logging
func NewUpdaterWithZap(base Updater, log *zap.Logger) UpdaterWithZap {
	return UpdaterWithZap{
		_base: base,
		_log:  log,
	}
}

// SetCounter implements Updater
func (_d UpdaterWithZap) SetCounter(name string, value int64) {
	_d._log.Debug("UpdaterWithZap: calling SetCounter", zap.Reflect("params", map[string]interface{}{
		"name":  name,
		"value": value}))
	defer func() {
		_d._log.Debug("UpdaterWithZap: SetCounter finished")
	}()
	_d._base.SetCounter(name, value)
	return
}

// SetGauge implements Updater
func (_d UpdaterWithZap) SetGauge(name string, value float64) {
	_d._log.Debug("UpdaterWithZap: calling SetGauge", zap.Reflect("params", map[string]interface{}{
		"name":  name,
		"value": value}))
	defer func() {
		_d._log.Debug("UpdaterWithZap: SetGauge finished")
	}()
	_d._base.SetGauge(name, value)
	return
}
