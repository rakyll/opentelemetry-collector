// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loggingexporter

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/internal/otlptext"
)

type loggingExporter struct {
	logger *zap.Logger
	debug  bool
}

func (s *loggingExporter) pushTraceData(
	_ context.Context,
	td pdata.Traces,
) error {

	s.logger.Info("TracesExporter", zap.Int("#spans", td.SpanCount()))

	if !s.debug {
		return nil
	}

	s.logger.Debug(otlptext.Traces(td))

	return nil
}

func (s *loggingExporter) pushMetricsData(
	_ context.Context,
	md pdata.Metrics,
) error {
	s.logger.Info("MetricsExporter", zap.Int("#metrics", md.MetricCount()))

	if !s.debug {
		return nil
	}

	s.logger.Debug(otlptext.Metrics(md))

	return nil
}

// newTracesExporter creates an exporter.TracesExporter that just drops the
// received data and logs debugging messages.
func newTracesExporter(config config.Exporter, level string, logger *zap.Logger) (component.TracesExporter, error) {
	s := &loggingExporter{
		debug:  strings.ToLower(level) == "debug",
		logger: logger,
	}

	return exporterhelper.NewTracesExporter(
		config,
		logger,
		s.pushTraceData,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(loggerSync(logger)),
	)
}

// newMetricsExporter creates an exporter.MetricsExporter that just drops the
// received data and logs debugging messages.
func newMetricsExporter(config config.Exporter, level string, logger *zap.Logger) (component.MetricsExporter, error) {
	s := &loggingExporter{
		debug:  strings.ToLower(level) == "debug",
		logger: logger,
	}

	return exporterhelper.NewMetricsExporter(
		config,
		logger,
		s.pushMetricsData,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(loggerSync(logger)),
	)
}

// newLogsExporter creates an exporter.LogsExporter that just drops the
// received data and logs debugging messages.
func newLogsExporter(config config.Exporter, level string, logger *zap.Logger) (component.LogsExporter, error) {
	s := &loggingExporter{
		debug:  strings.ToLower(level) == "debug",
		logger: logger,
	}

	return exporterhelper.NewLogsExporter(
		config,
		logger,
		s.pushLogData,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
		exporterhelper.WithShutdown(loggerSync(logger)),
	)
}

func (s *loggingExporter) pushLogData(
	_ context.Context,
	ld pdata.Logs,
) error {
	s.logger.Info("LogsExporter", zap.Int("#logs", ld.LogRecordCount()))

	if !s.debug {
		return nil
	}

	s.logger.Debug(otlptext.Logs(ld))

	return nil
}

func loggerSync(logger *zap.Logger) func(context.Context) error {
	return func(context.Context) error {
		// Currently Sync() return a different error depending on the OS.
		// Since these are not actionable ignore them.
		err := logger.Sync()
		if osErr, ok := err.(*os.PathError); ok {
			wrappedErr := osErr.Unwrap()
			if knownSyncError(wrappedErr) {
				err = nil
			}
		}
		return err
	}
}
