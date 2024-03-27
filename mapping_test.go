package main

import (
	"testing"
)

func TestParsePackageMapping(t *testing.T) {
	tests := map[string]struct {
		spec    string
		want    PackageMapping
		wantErr bool
	}{
		"simple mapping":   {"/foo=bar", PackageMapping{"/foo", "git", "bar"}, false},
		"with hostname":    {"codello.dev/govanity=foobar", PackageMapping{"codello.dev/govanity", "git", "foobar"}, false},
		"with invalid vcs": {"/foo=xyz:bar", PackageMapping{"/foo", "git", "xyz:bar"}, false},
		"empty prefix":     {"=bar", PackageMapping{}, true},
		"empty repo":       {"/=", PackageMapping{}, true},
		"invalid spec":     {"hello world", PackageMapping{}, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParsePackageMapping(tt.spec)
			if err == nil && tt.wantErr {
				t.Errorf("ParsePackageMapping(%q) returned no error, but an error was expected", tt.spec)
				return
			} else if err != nil && !tt.wantErr {
				t.Errorf("ParsePackageMapping(%q) returned an unexpected error: %s", tt.spec, err)
				return
			}
			if err == nil && *got != tt.want {
				t.Errorf("ParsePackageMapping(%q) = %v, want %v", tt.spec, got, tt.want)
			}
		})
	}
}
