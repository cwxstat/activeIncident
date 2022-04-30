package db

import (
	"testing"
)

func TestPopulateWeather(t *testing.T) {
	tests := []struct {
		name    string
		want    *WeatherEntry
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			want:    &WeatherEntry{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PopulateWeather()
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulateWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
		})
	}
}

func TestRunInGoRoutine(t *testing.T) {
	type args struct {
		countlimit []int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				countlimit: []int64{3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunInGoRoutine(tt.args.countlimit...); (err != nil) != tt.wantErr {
				t.Errorf("RunInGoRoutine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
