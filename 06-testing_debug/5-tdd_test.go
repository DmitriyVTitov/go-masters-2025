package testingdebug

import "testing"

func Test_stringrev(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single character",
			args: args{s: "a"},
			want: "a",
		},
		{
			name: "simple string",
			args: args{s: "hello world"},
			want: "dlrow olleh",
		},
		{
			name: "empty string",
			args: args{s: ""},
			want: "",
		},
		{
			name: "palindrome",
			args: args{s: "madam"},
			want: "madam",
		},
		{
			name: "special characters",
			args: args{s: "!@#$%^&*()"},
			want: ")(*&^%$#@!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringrev(tt.args.s); got != tt.want {
				t.Errorf("stringrev() = %v, want %v", got, tt.want)
			}
		})
	}
}
