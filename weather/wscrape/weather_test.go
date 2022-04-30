package wscrape

import (
	"fmt"
	"testing"
)

func TestZips(t *testing.T) {
	type args struct {
		zips []int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				zips: []int{
					19027,
					18041,
					18426},
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Zips(tt.args.zips)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zips() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
