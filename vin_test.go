package vin

import (
	"reflect"
	"testing"
)

func TestVin_ReadToml(t *testing.T) {
	type fields struct {
		Apps   []App
		vinDir string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		after   *Vin
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				Apps:   []App{},
				vinDir: "",
			},
			args: args{
				path: "testdata/vin.toml",
			},
			after: &Vin{
				Apps: []App{
					{Repo: "cli/cli"},
				},
				vinDir: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vin{
				Apps:   tt.fields.Apps,
				vinDir: tt.fields.vinDir,
			}
			if err := v.ReadTOML(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Vin.ReadToml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(v, tt.after) {
				t.Errorf("Vin.ReadToml() = %v, want %v", v, tt.after)
			}
		})
	}
}
