package vin

import "testing"

func TestApp_parseOwnerAndRepo(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "abc",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: "abc/def",
			},
			want:    "abc",
			want1:   "def",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "abc/def/ghi",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{}
			got, got1, err := a.parseOwnerAndRepo(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.parseOwnerAndRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.parseOwnerAndRepo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("App.parseOwnerAndRepo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s       string
		substrs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				s:       "abc def",
				substrs: []string{"abc"},
			},
			want: true,
		},
		{
			name: "",
			args: args{
				s:       "abc def",
				substrs: []string{"abc", "def"},
			},
			want: true,
		},
		{
			name: "",
			args: args{
				s:       "abc def",
				substrs: []string{"def", "efg"},
			},
			want: false,
		},
		{
			name: "",
			args: args{
				s:       "abc def",
				substrs: []string{"efg", "hij"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.substrs); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
