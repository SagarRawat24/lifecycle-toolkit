package analysis

import (
	"fmt"
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/stretchr/testify/require"
)

func TestObjectiveEvaluator_Evaluate(t *testing.T) {
	tests := []struct {
		name            string
		values          map[string]v1alpha3.ProviderResult
		o               v1alpha3.Objective
		want            types.ObjectiveResult
		mockedEvaluator ITargetEvaluator
	}{
		{
			name:   "no value in results map",
			values: map[string]v1alpha3.ProviderResult{},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{},
			want: types.ObjectiveResult{
				Score: 0.0,
				Error: fmt.Errorf("required value 'name' not available"),
				Objective: v1alpha3.Objective{
					AnalysisValueTemplateRef: v1alpha3.ObjectReference{
						Name: "name",
					},
				},
			},
		},
		{
			name: "evaluation passed",
			values: map[string]v1alpha3.ProviderResult{
				"name": {Value: "20", Query: "qqqqqqqq"},
			},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				Weight: 2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target *v1alpha3.Target) types.TargetResult {
					return types.TargetResult{
						Pass: true,
					}
				},
			},
			want: types.ObjectiveResult{
				Score: 2.0,
				Error: nil,
				Value: 20.0,
				Query: "qqqqqqqq",
				Result: types.TargetResult{
					Pass: true,
				},
				Objective: v1alpha3.Objective{
					AnalysisValueTemplateRef: v1alpha3.ObjectReference{
						Name: "name",
					},
					Weight: 2,
				},
			},
		},
		{
			name: "evaluation finished with warning",
			values: map[string]v1alpha3.ProviderResult{
				"name": {Value: "20", Query: "qqqqqqqq"},
			},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				Weight: 2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target *v1alpha3.Target) types.TargetResult {
					return types.TargetResult{
						Warning: true,
						Pass:    false,
					}
				},
			},
			want: types.ObjectiveResult{
				Score: 1.0,
				Error: nil,
				Value: 20.0,
				Query: "qqqqqqqq",
				Result: types.TargetResult{
					Pass:    false,
					Warning: true,
				},
				Objective: v1alpha3.Objective{
					AnalysisValueTemplateRef: v1alpha3.ObjectReference{
						Name: "name",
					},
					Weight: 2,
				},
			},
		},
		{
			name: "evaluation failed",
			values: map[string]v1alpha3.ProviderResult{
				"name": {Value: "20", Query: "qqqqqqqq"},
			},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				Weight: 2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target *v1alpha3.Target) types.TargetResult {
					return types.TargetResult{
						Warning: false,
						Pass:    false,
					}
				},
			},
			want: types.ObjectiveResult{
				Score: 0.0,
				Error: nil,
				Value: 20.0,
				Query: "qqqqqqqq",
				Result: types.TargetResult{
					Pass:    false,
					Warning: false,
				},
				Objective: v1alpha3.Objective{
					AnalysisValueTemplateRef: v1alpha3.ObjectReference{
						Name: "name",
					},
					Weight: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := NewObjectiveEvaluator(tt.mockedEvaluator)
			require.Equal(t, tt.want, oe.Evaluate(tt.values, &tt.o))
		})
	}
}

func TestGetValueFromMap(t *testing.T) {
	tests := []struct {
		name    string
		values  map[string]v1alpha3.ProviderResult
		in      string
		val     float64
		query   string
		wantErr bool
	}{
		{
			name: "happy path",
			values: map[string]v1alpha3.ProviderResult{
				"key1": {Value: "7", Query: "qqqqqqqq"},
			},
			in:      "key1",
			val:     7.0,
			query:   "qqqqqqqq",
			wantErr: false,
		},
		{
			name: "key not found",
			values: map[string]v1alpha3.ProviderResult{
				"key1": {Value: "7", Query: "qqqqqqqq"},
			},
			in:      "key",
			val:     0.0,
			query:   "",
			wantErr: true,
		},
		{
			name: "value not float",
			values: map[string]v1alpha3.ProviderResult{
				"key1": {Value: "", Query: "qqqqqqqq"},
			},
			in:      "key1",
			val:     0.0,
			query:   "qqqqqqqq",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, query, err := getResultFromMap(tt.values, tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.val, res)
			require.Equal(t, tt.query, query)
		})
	}
}

func TestComputeKey(t *testing.T) {
	obj := v1alpha3.ObjectReference{
		Name: "key",
	}

	require.Equal(t, "key", ComputeKey(obj))

	obj.Namespace = "namespace"

	require.Equal(t, "key-namespace", ComputeKey(obj))
}
