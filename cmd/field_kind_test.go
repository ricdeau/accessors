package cmd

import (
	"encoding/json"
	"testing"

	"github.com/ricdeau/accessors/consts"
	"github.com/stretchr/testify/assert"
)

func Test_fieldKind_UnmarshalText(t *testing.T) {
	type data struct {
		Kind FieldKind `json:"kind"`
	}
	tests := []struct {
		name string
		text []byte
		want data
	}{
		{
			name: "single",
			text: []byte(`{"kind":"public"}`),
			want: data{
				Kind: FieldKind{
					kinds: map[string]struct{}{
						consts.FieldKindPublic: {},
					},
				},
			},
		},
		{
			name: "multiple",
			text: []byte(`{"kind":"public,private"}`),
			want: data{
				Kind: FieldKind{
					kinds: map[string]struct{}{
						consts.FieldKindPublic:  {},
						consts.FieldKindPrivate: {},
					},
				},
			},
		},
		{
			name: "multiple space",
			text: []byte(`{"kind":"public, private"}`),
			want: data{
				Kind: FieldKind{
					kinds: map[string]struct{}{
						consts.FieldKindPublic:  {},
						consts.FieldKindPrivate: {},
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

func Test_fieldKind_Validate(t *testing.T) {
	tests := []struct {
		name    string
		kind    FieldKind
		wantErr bool
	}{
		{
			name: "success two",
			kind: FieldKind{
				kinds: map[string]struct{}{
					consts.FieldKindPrivate: {},
					consts.FieldKindPublic:  {},
				},
			},
			wantErr: false,
		},
		{
			name: "success one",
			kind: FieldKind{
				kinds: map[string]struct{}{
					consts.FieldKindPrivate: {},
				},
			},
			wantErr: false,
		},
		{
			name: "success none",
			kind: FieldKind{
				kinds: map[string]struct{}{},
			},
			wantErr: false,
		},
		{
			name: "invalid length",
			kind: FieldKind{
				kinds: map[string]struct{}{
					consts.FieldKindPrivate: {},
					consts.FieldKindPublic:  {},
					"invalid":               {},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid value",
			kind: FieldKind{
				kinds: map[string]struct{}{
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
