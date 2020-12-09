package vin

import (
	"reflect"
	"testing"
)

func TestVin_Filter(t *testing.T) {
	type fields struct {
		Apps []App
	}
	type args struct {
		filter func(app App) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vin
	}{
		{
			name: "",
			fields: fields{
				Apps: []App{{Repo: "a"}, {Repo: "b"}, {Repo: "c"}},
			},
			args: args{
				filter: func(app App) bool { return app.Repo == "a" },
			},
			want: &Vin{
				Apps: []App{{Repo: "a"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vin{
				Apps: tt.fields.Apps,
			}
			if got := v.Filter(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vin.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVin_FilterByPriority(t *testing.T) {
	type fields struct {
		Apps []App
	}
	type args struct {
		minPriority int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vin
	}{
		{
			name: "",
			fields: fields{
				Apps: []App{{Priority: 1}, {Priority: 2}, {Priority: 3}},
			},
			args: args{
				minPriority: 2,
			},
			want: &Vin{
				Apps: []App{{Priority: 2}, {Priority: 3}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vin{
				Apps: tt.fields.Apps,
			}
			if got := v.FilterByPriority(tt.args.minPriority); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vin.FilterByPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVin_FilterByRepo(t *testing.T) {
	type fields struct {
		Apps []App
	}
	type args struct {
		repos []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vin
	}{
		{
			name: "",
			fields: fields{
				Apps: []App{{Repo: "a"}, {Repo: "b"}, {Repo: "c"}},
			},
			args: args{
				repos: []string{"a", "b"},
			},
			want: &Vin{
				Apps: []App{{Repo: "a"}, {Repo: "b"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vin{
				Apps: tt.fields.Apps,
			}
			if got := v.FilterByRepo(tt.args.repos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vin.FilterByRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
