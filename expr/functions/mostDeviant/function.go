package mostDeviant

import (
	"container/heap"
	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	"math"
)

type mostDeviant struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &mostDeviant{}
	functions := []string{"mostDeviant"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

// mostDeviant(seriesList, n) -or- mostDeviant(n, seriesList)
func (f *mostDeviant) Do(e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	var nArg int
	if !e.Args()[0].IsConst() {
		// mostDeviant(seriesList, n)
		nArg = 1
	}
	seriesArg := nArg ^ 1 // XOR to make seriesArg the opposite argument. ( 0^1 -> 1 ; 1^1 -> 0 )

	n, err := e.GetIntArg(nArg)
	if err != nil {
		return nil, err
	}

	args, err := helper.GetSeriesArg(e.Args()[seriesArg], from, until, values)
	if err != nil {
		return nil, err
	}

	var mh types.MetricHeap

	for index, arg := range args {
		variance := helper.VarianceValue(arg.Values)
		if math.IsNaN(variance) {
			continue
		}

		if len(mh) < n {
			heap.Push(&mh, types.MetricHeapElement{Idx: index, Val: variance})
			continue
		}

		if variance > mh[0].Val {
			mh[0].Idx = index
			mh[0].Val = variance
			heap.Fix(&mh, 0)
		}
	}

	results := make([]*types.MetricData, len(mh))

	for len(mh) > 0 {
		v := heap.Pop(&mh).(types.MetricHeapElement)
		results[len(mh)] = args[v.Idx]
	}

	return results, err
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *mostDeviant) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"mostDeviant": {
			Description: "Takes one metric or a wildcard seriesList followed by an integer N.\nDraws the N most deviant metrics.\nTo find the deviants, the standard deviation (sigma) of each series\nis taken and ranked. The top N standard deviations are returned.\n\n  Example:\n\n.. code-block:: none\n\n  &target=mostDeviant(server*.instance*.memory.free, 5)\n\nDraws the 5 instances furthest from the average memory free.",
			Function:    "mostDeviant(seriesList, n)",
			Group:       "Filter Series",
			Module:      "graphite.render.functions",
			Name:        "mostDeviant",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "n",
					Required: true,
					Type:     types.Integer,
				},
			},
		},
	}
}
