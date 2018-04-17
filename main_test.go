package main

import "testing"

func Test_parseLine(t *testing.T) {

	tests := []struct {
		name  string
		line  string
		want  string
		want1 string
		want2 string
	}{
		{
			"one word",
			"help",
			"help",
			"",
			"",
		},
		{
			"2 words",
			"get /foo",
			"get",
			"/foo",
			"",
		},
		{
			"3 words",
			`post /foo '{"foo":"bar"}'`,
			"post",
			"/foo",
			`'{"foo":"bar"}'`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := parseLine(tt.line)
			if got != tt.want {
				t.Errorf("parseLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseLine() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseLine() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
