// Copyright 2021 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus-community/ipmi_exporter/freeipmi"
)

const (
	DCMICollectorName CollectorName = "dcmi"
)

var (
	powerConsumptionDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dcmi", "power_consumption_watts"),
		"Current power consumption in Watts.",
		[]string{},
		nil,
	)
)

type DCMICollector struct{}

func (c DCMICollector) Name() CollectorName {
	return DCMICollectorName
}

func (c DCMICollector) Cmd() string {
	return "ipmi-dcmi"
}

func (c DCMICollector) Args() []string {
	return []string{"--get-system-power-statistics"}
}

func (c DCMICollector) Collect(result freeipmi.Result, ch chan<- prometheus.Metric, target ipmiTarget) (int, error) {
	currentPowerConsumption, err := freeipmi.GetCurrentPowerConsumption(result)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to collect DCMI data", "target", targetName(target.host), "error", err)
		return 0, err
	}
	ch <- prometheus.MustNewConstMetric(
		powerConsumptionDesc,
		prometheus.GaugeValue,
		currentPowerConsumption,
	)
	return 1, nil
}
