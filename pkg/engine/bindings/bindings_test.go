package bindings

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_checkBindingName(t *testing.T) {
	tests := []struct {
		name        string
		bindingName string
		wantErr     bool
	}{{
		name:        "empty",
		bindingName: "",
		wantErr:     true,
	}, {
		name:        "ok",
		bindingName: "foo",
		wantErr:     false,
	}, {
		name:        "not ok",
		bindingName: "foo bar",
		wantErr:     true,
	}, {
		name:        "not ok",
		bindingName: "$foo",
		wantErr:     true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkBindingName(tt.bindingName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegisterBinding(t *testing.T) {
	tests := []struct {
		name        string
		bindings    binding.Bindings
		bindingName string
		value       any
	}{{
		bindings:    binding.NewBindings(),
		bindingName: "foo",
		value:       "bar",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bindings := RegisterBinding(context.TODO(), tt.bindings, tt.bindingName, tt.value)
			assert.NotNil(t, bindings)
			got, err := bindings.Get("$" + tt.bindingName)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			value, err := got.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.value, value)
		})
	}
}

func TestResolveBinding(t *testing.T) {
	tests := []struct {
		name      string
		bindings  binding.Bindings
		input     any
		variable  v1alpha1.Binding
		wantName  string
		wantValue any
		wantErr   bool
	}{{
		name:     "ok",
		bindings: binding.NewBindings(),
		input:    nil,
		variable: v1alpha1.Binding{
			Name:  "foo",
			Value: v1alpha1.Any{Value: "bar"},
		},
		wantName:  "foo",
		wantValue: "bar",
		wantErr:   false,
	}, {
		name:     "error",
		bindings: binding.NewBindings(),
		input:    nil,
		variable: v1alpha1.Binding{
			Name:  "$foo",
			Value: v1alpha1.Any{Value: "bar"},
		},
		wantErr: true,
	}, {
		name:     "error",
		bindings: binding.NewBindings(),
		input:    nil,
		variable: v1alpha1.Binding{
			Name:  "foo",
			Value: v1alpha1.Any{Value: "($bar)"},
		},
		wantErr: true,
	}, {
		name:     "error",
		bindings: binding.NewBindings(),
		input:    nil,
		variable: v1alpha1.Binding{
			Name:  "($foo)",
			Value: v1alpha1.Any{Value: "bar"},
		},
		wantErr: true,
	}, {
		name:     "error",
		bindings: binding.NewBindings().Register("$foo", binding.NewBinding("abc")).Register("$bar", binding.NewBinding("def")),
		input:    nil,
		variable: v1alpha1.Binding{
			Name:  "($foo)",
			Value: v1alpha1.Any{Value: "($bar)"},
		},
		wantName:  "abc",
		wantValue: "def",
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, value, err := ResolveBinding(context.TODO(), tt.bindings, tt.input, tt.variable)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantName, name)
			assert.Equal(t, tt.wantValue, value)
		})
	}
}
