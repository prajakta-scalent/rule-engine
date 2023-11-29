package ruleengine_test

import (
	"reflect"
	"testing"

	ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine"
	mock_ruleengine "github.com/prajakta-scalent/rule-engine/pkg/rule-engine/mock"
	"github.com/prajakta-scalent/rule-engine/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	type args struct {
		ruleExporter ruleengine.ResultExporter
	}

	mockCtrl := gomock.NewController(t)
	mockResultExporter := mock_ruleengine.NewMockResultExporter(mockCtrl)

	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "rule exporter is nil",
			args:    args{nil},
			want:    ruleengine.ErrRuleExporterEmpty,
			wantErr: true,
		},
		{
			name:    "should pass",
			args:    args{ruleExporter: mockResultExporter},
			want:    ruleengine.ErrRuleExporterEmpty,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ruleengine.New(tt.args.ruleExporter)
			if tt.wantErr && assert.Error(t, err) {
				assert.Equal(t, tt.want, err)
			}
		})
	}
}

func TestRuleEngineImpl_RegisterGroup(t *testing.T) {
	type fields struct {
		resultExporter ruleengine.ResultExporter
	}
	type args struct {
		ruleGroup ruleengine.RuleGroup
	}

	mockCtrl := gomock.NewController(t)
	mockResultExporter := mock_ruleengine.NewMockResultExporter(mockCtrl)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "empty rule name",
			fields: fields{
				resultExporter: mockResultExporter,
			},
			args: args{
				ruleGroup: ruleengine.RuleGroup{
					Name:  "",
					Rules: []ruleengine.Rule{},
				},
			},
			want:    ruleengine.ErrRuleGroupNameEmpty,
			wantErr: true,
		},
		{
			name: "empty rules in rule group",
			fields: fields{
				resultExporter: mockResultExporter,
			},
			args: args{
				ruleGroup: ruleengine.RuleGroup{
					Name:  "empty rules",
					Rules: nil,
				},
			},
			want:    ruleengine.ErrRulesEmpty,
			wantErr: true,
		},
		{
			name: "should pass",
			fields: fields{
				resultExporter: mockResultExporter,
			},
			args: args{
				ruleGroup: ruleengine.RuleGroup{
					Name: "empty rules",
					Rules: []ruleengine.Rule{
						{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ruleengine.New(tt.fields.resultExporter)
			assert.NoError(t, err)
			err = r.RegisterGroup(tt.args.ruleGroup)
			if tt.wantErr && assert.Error(t, err) {
				assert.Equal(t, tt.want, err)
			}
		})
	}
}

func TestRuleEngineImpl_Execute(t *testing.T) {
	type fields struct {
		resultExporter ruleengine.ResultExporter
		ruleGroup      ruleengine.RuleGroup
	}
	type args struct {
		data map[string][]ruleengine.Input
	}

	mockCtrl := gomock.NewController(t)
	mockResultExporter := mock_ruleengine.NewMockResultExporter(mockCtrl)
	ruleGroup := ruleengine.RuleGroup{
		Name: "test rules",
		Rules: []ruleengine.Rule{
			{},
		},
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecutionID string
		wantResult      []ruleengine.RuleGroupResult
		wantErr         bool
		want            error
	}{
		{
			name: "empty rule input map",
			fields: fields{
				resultExporter: mockResultExporter,
				ruleGroup:      ruleGroup,
			},
			args:    args{nil},
			want:    ruleengine.ErrRuleInputDataEmpty,
			wantErr: true,
		},
		{
			name: "provided name with empty rule input",
			fields: fields{
				resultExporter: mockResultExporter,
				ruleGroup:      ruleGroup,
			},
			args: args{
				data: map[string][]ruleengine.Input{
					"EmptyValue": nil,
				},
			},
			want:    ruleengine.ErrRuleInputValueEmpty,
			wantErr: true,
		},
		{
			name: "length of rule and rule input not matching",
			fields: fields{
				resultExporter: mockResultExporter,
				ruleGroup:      ruleGroup,
			},
			args: args{
				data: map[string][]ruleengine.Input{
					"test rules": {
						{},
						{},
					},
				},
			},
			want:    ruleengine.ErrRuleInputLengthNotEqual,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ruleengine.New(tt.fields.resultExporter)
			if assert.NoError(t, err) {
				err = r.RegisterGroup(tt.fields.ruleGroup)
				if !assert.NoError(t, err) {
					return
				}
				_, gotResult, err := r.Execute(tt.args.data)
				// if gotExecutionID != tt.wantExecutionID {
				// 	t.Errorf("RuleEngineImpl.Execute() gotExecutionID = %v, want %v", gotExecutionID, tt.wantExecutionID)
				// }
				if !reflect.DeepEqual(gotResult, tt.wantResult) {
					t.Errorf("RuleEngineImpl.Execute() gotResult = %v, want %v", gotResult, tt.wantResult)
				}
				if tt.wantErr && assert.Error(t, err) {
					assert.Equal(t, tt.want, err)
				}
			}
		})
	}
}

// func TestRuleEngineImpl_ExecuteRulesConcurrently(t *testing.T) {
// 	type fields struct {
// 		ruleGroups     []RuleGroup
// 		resultExporter ResultExporter
// 	}
// 	type args struct {
// 		ruleGroup RuleGroup
// 		data      []Input
// 	}
// 	tests := []struct {
// 		name       string
// 		fields     fields
// 		args       args
// 		wantResult []RuleResult
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &RuleEngineImpl{
// 				ruleGroups:     tt.fields.ruleGroups,
// 				resultExporter: tt.fields.resultExporter,
// 			}
// 			if gotResult := r.ExecuteRulesConcurrently(tt.args.ruleGroup, tt.args.data); !reflect.DeepEqual(gotResult, tt.wantResult) {
// 				t.Errorf("RuleEngineImpl.ExecuteRulesConcurrently() = %v, want %v", gotResult, tt.wantResult)
// 			}
// 		})
// 	}
// }

func TestRuleEngineImpl_ExecuteRulesSequentially(t *testing.T) {
	type fields struct {
		ruleGroups     []ruleengine.RuleGroup
		resultExporter ruleengine.ResultExporter
	}
	type args struct {
		ruleGroup ruleengine.RuleGroup
		data      []ruleengine.Input
	}

	mockCtrl := gomock.NewController(t)
	mockResultExporter := mock_ruleengine.NewMockResultExporter(mockCtrl)

	rules := []ruleengine.Rule{
		{
			Name:        "AgeShouldBeMoreThan",
			Condition:   "gte",
			MatchValue:  18,
			IsMandatory: false,
		},
		{
			Name:        "APICallCheckAgeAllowed",
			Condition:   "callback",
			IsMandatory: true,
		},
	}

	ruleGroup := ruleengine.RuleGroup{
		Name:                "userRulesGroup",
		Rules:               rules,
		ExecuteConcurrently: true,
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []ruleengine.RuleResult
	}{
		{
			name: "sequential execution",
			fields: fields{
				ruleGroups: []ruleengine.RuleGroup{
					ruleGroup,
				},
				resultExporter: mockResultExporter,
			},
			args: args{
				ruleGroup: ruleGroup,
				data: []ruleengine.Input{
					{
						RuleName: "AgeShouldBeMoreThan",
						Value:    20,
					},
					{
						RuleName: "APICallCheckAgeAllowed",
						Value: services.User{
							Id:  1,
							Age: 18,
						},
					},
				},
			},
			wantResult: []ruleengine.RuleResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ruleengine.New(tt.fields.resultExporter)
			assert.NoError(t, err)
			err = r.RegisterGroup(tt.args.ruleGroup)
			if !assert.NoError(t, err) {
				return
			}
			result := r.ExecuteRulesSequentially(tt.args.ruleGroup, tt.args.data)
			if len(result) == 0 {
				assert.Fail(t, "result is empty")
			}
		})
	}
}

func TestRuleEngineImpl_ExecuteRulesConcurrently(t *testing.T) {
	type fields struct {
		ruleGroups     []ruleengine.RuleGroup
		resultExporter ruleengine.ResultExporter
	}
	type args struct {
		ruleGroup ruleengine.RuleGroup
		data      []ruleengine.Input
	}

	mockCtrl := gomock.NewController(t)
	mockResultExporter := mock_ruleengine.NewMockResultExporter(mockCtrl)

	rules := []ruleengine.Rule{
		{
			Name:        "AgeShouldBeMoreThan",
			Condition:   "gte",
			MatchValue:  18,
			IsMandatory: false,
		},
		{
			Name:        "APICallCheckAgeAllowed",
			Condition:   "callback",
			IsMandatory: true,
		},
	}

	ruleGroup := ruleengine.RuleGroup{
		Name:                "userRulesGroup",
		Rules:               rules,
		ExecuteConcurrently: true,
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []ruleengine.RuleResult
	}{
		{
			name: "sequential execution",
			fields: fields{
				ruleGroups: []ruleengine.RuleGroup{
					ruleGroup,
				},
				resultExporter: mockResultExporter,
			},
			args: args{
				ruleGroup: ruleGroup,
				data: []ruleengine.Input{
					{
						RuleName: "AgeShouldBeMoreThan",
						Value:    20,
					},
					{
						RuleName: "APICallCheckAgeAllowed",
						Value: services.User{
							Id:  1,
							Age: 18,
						},
					},
				},
			},
			wantResult: []ruleengine.RuleResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ruleengine.New(tt.fields.resultExporter)
			assert.NoError(t, err)
			err = r.RegisterGroup(tt.args.ruleGroup)
			if !assert.NoError(t, err) {
				return
			}
			result := r.ExecuteRulesConcurrently(tt.args.ruleGroup, tt.args.data)
			if len(result) == 0 {
				assert.Fail(t, "result is empty")
			}
		})
	}
}

// func Test_getRuleData(t *testing.T) {
// 	type args struct {
// 		name string
// 		data []Input
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    interface{}
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := getRuleData(tt.args.name, tt.args.data)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getRuleData() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("getRuleData() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
