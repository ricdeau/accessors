package main

import (
	"testing"
)

func TestFieldType_Includes(t *testing.T) {
	tests := []struct {
		name  string
		this  FieldType
		other FieldType
		want  bool
	}{
		{
			name:  "All Includes Public",
			this:  All,
			other: Public,
			want:  true,
		},
		{
			name:  "All Includes Private",
			this:  All,
			other: Private,
			want:  true,
		},
		{
			name:  "Public not Includes Private",
			this:  Public,
			other: Private,
			want:  false,
		},
		{
			name:  "Private not Includes Public",
			this:  Private,
			other: Public,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Includes(tt.other); got != tt.want {
				t.Errorf("this = %s, other = %s: Includes() = %v, want %v", tt.this, tt.other, got, tt.want)
			}
		})
	}
}
