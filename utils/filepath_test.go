package utils

import (
	"testing"
)

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		suffixes []string
		want     bool
	}{
		{"HasSuffix(\"test.toml\",\"\") == false", "test.toml", []string{""}, false},
		{"HasSuffix(\"test.toml\",\".toml\") == true", "test.toml", []string{".toml"}, true},
		{"HasSuffix(\"test.exe\",\".toml\") == false", "test.exe", []string{".toml"}, false},
		{"HasSuffix(\"toml\",\".toml\") == false", "toml", []string{".toml"}, false},
		{"HasSuffix(\"toml.txt\",\".toml\") == false", "toml.txt", []string{".toml"}, false},
		{"HasSuffix(\"toml.toml\",\".toml\") == true", "toml.toml", []string{".toml"}, true},
		{"HasSuffix(\"./test.toml\",\".toml\") == true", "./test.toml", []string{".toml"}, true},
		{"HasSuffix(\"./test.txt\",\".toml\") == false", "./test.txt", []string{".toml"}, false},
		{"HasSuffix(\"/home/test/test.toml\",\".toml\") == true", "/home/test/test.toml", []string{".toml"}, true},
		{"HasSuffix(\"/home/test/test.txt\",\".toml\") == false", "/home/test/test.txt", []string{".toml"}, false},
		{"HasSuffix(\"../test.toml\",\".toml\") == true", "../test.toml", []string{".toml"}, true},
		{"HasSuffix(\"../test.txt\",\".toml\") == false", "../test.txt", []string{".toml"}, false},
		{"HasSuffix(\"/home/test/test.txt\",\".toml\", \".txt\") == true", "/home/test/test.txt", []string{".toml", ".txt"}, true},
		{"HasSuffix(\"/home/test/test.txt\",\"\") == false", "/home/test/test.txt", []string{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := HasSuffix(tt.input, tt.suffixes...)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestIsToml(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"IsToml(\"test.toml\") == true", "test.toml", true},
		{"IsToml(\"test.exe\") == false", "test.exe", false},
		{"IsToml(\"toml\") == false", "toml", false},
		{"IsToml(\"toml.txt\") == false", "toml.txt", false},
		{"IsToml(\"toml.toml\") == true", "toml.toml", true},
		{"IsToml(\"./test.toml\") == true", "./test.toml", true},
		{"IsToml(\"./test.txt\") == false", "./test.txt", false},
		{"IsToml(\"/home/test/test.toml\") == true", "/home/test/test.toml", true},
		{"IsToml(\"/home/test/test.txt\") == false", "/home/test/test.txt", false},
		{"IsToml(\"../test.toml\") == true", "../test.toml", true},
		{"IsToml(\"../test.txt\") == false", "../test.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := IsToml(tt.input)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestIsYaml(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"IsYaml(\"test.yaml\") == true", "test.yaml", true},
		{"IsYaml(\"test.yml\") == true", "test.yml", true},
		{"IsYaml(\"test.exe\") == false", "test.exe", false},
		{"IsYaml(\"yaml\") == false", "yaml", false},
		{"IsYaml(\"yml\") == false", "yml", false},
		{"IsYaml(\"yaml.txt\") == false", "yaml.txt", false},
		{"IsYaml(\"yml.txt\") == false", "yml.txt", false},
		{"IsYaml(\"yaml.yaml\") == true", "yaml.yaml", true},
		{"IsYaml(\"yml.yml\") == true", "yml.yml", true},
		{"IsYaml(\"./test.yaml\") == true", "./test.yaml", true},
		{"IsYaml(\"./test.yml\") == true", "./test.yml", true},
		{"IsYaml(\"./test.txt\") == false", "./test.txt", false},
		{"IsYaml(\"/home/test/test.yaml\") == true", "/home/test/test.yaml", true},
		{"IsYaml(\"/home/test/test.yml\") == true", "/home/test/test.yml", true},
		{"IsYaml(\"/home/test/test.txt\") == false", "/home/test/test.txt", false},
		{"IsYaml(\"../test.yaml\") == true", "../test.yaml", true},
		{"IsYaml(\"../test.yml\") == true", "../test.yml", true},
		{"IsYaml(\"../test.txt\") == false", "../test.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := IsYaml(tt.input)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestGetFileNamePart(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"GetFileNamePart(\"name\") == \"name\"", "name", "name"},
		{"GetFileNamePart(\"name.\") == \"name\"", "name.", "name"},
		{"GetFileNamePart(\"name.yaml\") == \"name\"", "name.yaml", "name"},
		{"GetFileNamePart(\"name.exe\") == \"name\"", "name.exe", "name"},
		{"GetFileNamePart(\"./name\") == \"name\"", "./name", "name"},
		{"GetFileNamePart(\"./name.yaml\") == \"name\"", "./name.yaml", "name"},
		{"GetFileNamePart(\"./name.\") == \"name\"", "./name.", "name"},
		{"GetFileNamePart(\"./name.txt\") == \"name\"", "./name.txt", "name"},
		{"GetFileNamePart(\"/home/test/name\") == \"name\"", "/home/test/name", "name"},
		{"GetFileNamePart(\"/home/test/name.yaml\") == \"name\"", "/home/test/name.yaml", "name"},
		{"GetFileNamePart(\"/home/test/name.\") == \"name\"", "/home/test/name.", "name"},
		{"GetFileNamePart(\"../name.yaml\") == \"name\"", "../name.yaml", "name"},
		{"GetFileNamePart(\"../name.\") == \"name\"", "../name.", "name"},
		{"GetFileNamePart(\"../name\") == \"name\"", "../name", "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := GetFileNamePart(tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
