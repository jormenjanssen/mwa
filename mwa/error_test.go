package main

import (
	"fmt"
	"testing"
)

func TestLastErrorFunc(t *testing.T) {

	kError := fmt.Errorf("Some expected error")
	//aError := fmt.Errorf("Alternate error")

	type args struct {
		f    func() error
		aerr error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		desErr  error
	}{
		{name: "Test without alternate or primary error", args: args{f: func() error { return nil }, aerr: nil}, wantErr: true},
		{name: "Test for alternate error", args: args{f: func() error { return nil }, aerr: kError}, wantErr: true},
		{name: "Test for primary errror ", args: args{f: func() error { return kError }, aerr: nil}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LastErrorFunc(tt.args.f, tt.args.aerr); (err != nil) != tt.wantErr {
				t.Errorf("LastErrorFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
