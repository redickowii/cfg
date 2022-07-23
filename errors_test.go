package cfg

import (
	"testing"
)

func TestGetEnvError_Error(t *testing.T) {
	type fields struct {
		field string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test Error func",
			fields: fields{"string"},
			want:   "cant find env tag for field: string",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				e := GetEnvError{
					field: tt.fields.field,
				}
				if got := e.Error(); got != tt.want {
					t.Errorf("Error() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
