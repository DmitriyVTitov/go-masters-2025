package testingdebug

import (
	"testing"
)

// ***
// Unit Tests
// ***

func TestReverseString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{s: ""},
			want: "",
		},
		{
			name: "single character",
			args: args{s: "a"},
			want: "a",
		},
		{
			name: "palindrome",
			args: args{s: "madam"},
			want: "madam",
		},
		{
			name: "regular string",
			args: args{s: "hello"},
			want: "olleh",
		},
		{
			name: "string with spaces",
			args: args{s: "hello world"},
			want: "dlrow olleh",
		},
		{
			name: "string with special characters",
			args: args{s: "!@#$%^&*()"},
			want: ")(*&^%$#@!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.args.s); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseStringWithError(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{s: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "regular string",
			args:    args{s: "hello"},
			want:    "olleh",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReverseStringWithError(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReverseStringWithError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReverseStringWithError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lowercaseName(t *testing.T) {}

func TestT_Method(t *testing.T) {}

// ***
// Benchmarks
// ***

func BenchmarkReverseString(b *testing.B) {
	for b.Loop() {
		ReverseString("hello world!")
	}
}

func BenchmarkReverseStringSB(b *testing.B) {
	for b.Loop() {
		ReverseStringSB("hello world!")
	}
}
