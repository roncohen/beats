// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package server

import (
	"time"

	serverhelper "github.com/elastic/beats/metricbeat/helper/server"
	"github.com/elastic/beats/metricbeat/helper/server/udp"
	"github.com/elastic/beats/metricbeat/mb"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	mb.Registry.MustAddMetricSet("statsd", "server", New, mb.DefaultMetricSet())
}

// Config for the statsd server metricset.
type Config struct {
	ReservoirSize int `config:"reservoir_size"`
}

func defaultConfig() Config {
	return Config{
		ReservoirSize: 1000,
	}
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	server    serverhelper.Server
	processor *metricProcessor
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	config := defaultConfig()
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	svc, err := udp.NewUdpServer(base)
	if err != nil {
		return nil, err
	}

	processor := newMetricProcessor(config.ReservoirSize)
	return &MetricSet{
		BaseMetricSet: base,
		server:        svc,
		processor:     processor,
	}, nil
}

// Run method provides the module with a reporter with which events can be reported.
func (m *MetricSet) Run(reporter mb.PushReporterV2) {
	period := m.Module().Config().Period

	// Start event watcher
	m.server.Start()

	for {
		select {
		case <-reporter.Done():
			m.server.Stop()
			return
		case <-time.After(period):
			event := mb.Event{
				MetricSetFields: m.processor.GetAll(),
			}

			event.Namespace = "statsd"
			reporter.Event(event)
		case msg := <-m.server.GetEvents():
			err := m.processor.Process(msg)
			if err != nil {
				reporter.Error(err)
			}
		}
	}
}
