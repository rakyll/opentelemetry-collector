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

package testdata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/internal"
	otlpcollectortrace "go.opentelemetry.io/collector/internal/data/protogen/collector/trace/v1"
)

type traceTestCase struct {
	name string
	td   pdata.Traces
	otlp *otlpcollectortrace.ExportTraceServiceRequest
}

func generateAllTraceTestCases() []traceTestCase {
	return []traceTestCase{
		{
			name: "one-empty-resource-spans",
			td:   GenerateTracesOneEmptyResourceSpans(),
			otlp: generateTracesOtlpOneEmptyResourceSpans(),
		},
		{
			name: "no-libraries",
			td:   GenerateTracesNoLibraries(),
			otlp: generateTracesOtlpNoLibraries(),
		},
		{
			name: "one-empty-instrumentation-library",
			td:   GenerateTracesOneEmptyInstrumentationLibrary(),
			otlp: generateTracesOtlpOneEmptyInstrumentationLibrary(),
		},
		{
			name: "one-span-no-resource",
			td:   GenerateTracesOneSpanNoResource(),
			otlp: generateTracesOtlpOneSpanNoResource(),
		},
		{
			name: "one-span",
			td:   GenerateTracesOneSpan(),
			otlp: generateTracesOtlpOneSpan(),
		},
		{
			name: "two-spans-same-resource",
			td:   GenerateTracesTwoSpansSameResource(),
			otlp: generateTracesOtlpSameResourceTwoSpans(),
		},
		{
			name: "two-spans-same-resource-one-different",
			td:   GenerateTracesTwoSpansSameResourceOneDifferent(),
			otlp: generateTracesOtlpTwoSpansSameResourceOneDifferent(),
		},
	}
}

func TestToFromOtlpTrace(t *testing.T) {
	allTestCases := generateAllTraceTestCases()
	// Ensure NumTraceTests gets updated.
	for i := range allTestCases {
		test := allTestCases[i]
		t.Run(test.name, func(t *testing.T) {
			td := pdata.TracesFromInternalRep(internal.TracesFromOtlp(test.otlp))
			assert.EqualValues(t, test.td, td)
			otlp := internal.TracesToOtlp(td.InternalRep())
			assert.EqualValues(t, test.otlp, otlp)
		})
	}
}
