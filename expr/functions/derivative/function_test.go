package derivative

import (
	"math"
	"testing"
	"time"

	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	th "github.com/go-graphite/carbonapi/tests"
)

var (
	md []interfaces.FunctionMetadata = New("")
)

func init() {
	for _, m := range md {
		metadata.RegisterFunction(m.Name, m.F)
	}
}

func TestDerivative(t *testing.T) {
	now32 := int64(time.Now().Unix())

	tests := []th.EvalTestItem{
		{
			"derivative(metric1)",
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 0, Until: 1}: {types.MakeMetricData("metric1", []float64{2, 4, 6, 1, 4, math.NaN(), 8}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("derivative(metric1)",
				[]float64{math.NaN(), 2, 2, -5, 3, math.NaN(), 4}, 1, now32).SetTag("derivative", "1").SetNameTag("derivative(metric1)")},
		},
		{
			"derivative(metric1)",
			map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 0, Until: 1}: {types.MakeMetricData("metric1", []float64{math.NaN(), 4, 6, 1, 4, math.NaN(), 8}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("derivative(metric1)",
				[]float64{math.NaN(), math.NaN(), 2, -5, 3, math.NaN(), 4}, 1, now32).SetTag("derivative", "1").SetNameTag("derivative(metric1)")},
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			eval := th.EvaluatorFromFunc(md[0].F)
			th.TestEvalExpr(t, eval, &tt)
		})
	}

}
