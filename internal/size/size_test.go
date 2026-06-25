// Copyright 2026 The Tessera authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package size

import (
	"testing"
)

func TestParseBytes(t *testing.T) {
	tests := []struct {
		in      string
		want    uint64
		wantErr bool
	}{
		{"1", 1, false},
		{"0", 0, false},
		{"1B", 1, false},
		{"1 B", 1, false},
		{"1bytes", 1, false},
		{"1000", 1000, false},
		{"1k", 1000, false},
		{"1kb", 1000, false},
		{"1kib", 1024, false},
		{"1MB", 1000000, false},
		{"1MiB", 1048576, false},
		{"1GB", 1000000000, false},
		{"1GiB", 1073741824, false},
		{"768MB", 768000000, false},
		{"0MB", 0, false},
		{"1.5Mb", 1500000, false},
		{"1.5Mib", 1572864, false},
		{"", 0, true},
		{"blah", 0, true},
		{"1x", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseBytes(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseBytes(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseBytes(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestParseSI(t *testing.T) {
	tests := []struct {
		in       string
		wantVal  float64
		wantUnit string
		wantErr  bool
	}{
		{"256k", 256000, "", false},
		{"256kb", 256000, "b", false},
		{"256", 256, "", false},
		{"256.5M", 256500000, "", false},
		{"", 0, "", true},
		{"blah", 0, "", true},
	}

	for _, tt := range tests {
		gotVal, gotUnit, err := ParseSI(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseSI(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if gotVal != tt.wantVal || gotUnit != tt.wantUnit {
			t.Errorf("ParseSI(%q) = %f, %q, want %f, %q", tt.in, gotVal, gotUnit, tt.wantVal, tt.wantUnit)
		}
	}
}
