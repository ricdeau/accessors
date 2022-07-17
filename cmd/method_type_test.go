package cmd

import (
	"encoding/json"
	"testing"

	"github.com/ricdeau/accessors/consts"
	"github.com/stretchr/testify/assert"
)

func Test_MethodType_UnmarshalJSON(t *testing.T) {
	type data struct {
		Type MethodType `json:"type"`
	}
	tests := []struct {
		name string
		text []byte
		want data
	}{
		{
			name: "single",
			text: []byte(`{"type":"getters"}`),
			want: data{
				Type: MethodType{
					types: map[string]struct{}{
						consts.MethodTypeGetters: {},
					},
				},
			},
		},
		{
			name: "multiple",
			text: []byte(`{"type":"getters,setters"}`),
			want: data{
				Type: MethodType{
					types: map[string]struct{}{
						consts.MethodTypeGetters: {},
						consts.MethodTypeSetters: {},
					},
				},
			},
		},
		{
			name: "multiple space",
			text: []byte(`{"type":"getters, setters"}`),
			want: data{
				Type: MethodType{
					types: map[string]struct{}{
						consts.MethodTypeGetters: {},
						consts.MethodTypeSetters: {},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got data
			err := json.Unmarshal(tt.text, &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_MethodType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		kind    MethodType
		wantErr bool
	}{
		{
			name: "success two",
			kind: MethodType{
				types: map[string]struct{}{
					consts.MethodTypeGetters: {},
					consts.MethodTypeSetters: {},
				},
			},
			wantErr: false,
		},
		{
			name: "success one",
			kind: MethodType{
				types: map[string]struct{}{
					consts.MethodTypeGetters: {},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid none",
			kind:    MethodType{},
			wantErr: true,
		},
		{
			name: "invalid length",
			kind: MethodType{
				types: map[string]struct{}{
					consts.MethodTypeGetters: {},
					consts.MethodTypeSetters: {},
					"invalid":                {},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid value",
			kind: MethodType{
				types: map[string]struct{}{
					"invalid": {},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.kind.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
