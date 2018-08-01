package main

import (
	"testing"
)

func Test_isKevArg(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"name and value", args{"NAME=VALUE"}, true},
		{"name only", args{"NAME="}, true},
		{"value only", args{"=VALUE"}, false},
		{"empty", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isKevArg(tt.args.s); got != tt.want {
				t.Errorf("isKevArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kevArg_parse(t *testing.T) {
	tests := []struct {
		name      string
		a         kevArg
		wantKey   string
		wantValue string
	}{
		{"name and value", kevArg("NAME=VALUE"), "NAME", "VALUE"},
		{"name only", kevArg("NAME="), "NAME", ""},
		{"value only", kevArg("=VALUE"), "", ""},
		{"empty", kevArg(""), "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue := tt.a.parse()
			if gotKey != tt.wantKey {
				t.Errorf("kevArg.parse() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("kevArg.parse() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func Test_fileArg_parse(t *testing.T) {
	tests := []struct {
		name         string
		a            fileArg
		wantFileType string
		wantFilePath string
	}{
		{"normal", fileArg("@testfile"), "", "testfile"},
		{"with type", fileArg("@yaml:testfile"), "YAML", "testfile"},
		{"only type", fileArg("@yaml:"), "YAML", ""},
		{"only path", fileArg("@:testfile"), "", "testfile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileType, gotFilePath := tt.a.parse()
			if gotFileType != tt.wantFileType {
				t.Errorf("fileArg.parse() gotFileType = %v, want %v", gotFileType, tt.wantFileType)
			}
			if gotFilePath != tt.wantFilePath {
				t.Errorf("fileArg.parse() gotFilePath = %v, want %v", gotFilePath, tt.wantFilePath)
			}
		})
	}
}

func Test_isFileArg(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"normal", args{"@testfile"}, true},
		{"with type", args{"@yaml:testfile"}, true},
		{"only type", args{"@yaml:"}, true},
		{"only path", args{"@:testfile"}, true},
		{"at at", args{"@@"}, false},
		{"empty", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFileArg(tt.args.s); got != tt.want {
				t.Errorf("isFileArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringArg_String(t *testing.T) {
	tests := []struct {
		name string
		a    stringArg
		want string
	}{
		{"@@", stringArg("@@"), "@"},
		{"!!", stringArg("!!"), "!"},
		{"@!", stringArg("@!"), "@!"},
		{"!@", stringArg("!@"), "!@"},
		{"empty string", stringArg(""), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("stringArg.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
