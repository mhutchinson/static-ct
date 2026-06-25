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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseBytes parses a string representation of bytes into the number of bytes it represents.
// It supports SI suffixes (KB, MB, GB, TB) which are 1000-based,
// and binary suffixes (KiB, MiB, GiB, TiB) which are 1024-based.
// It also supports 'k', 'm', 'g', 't' as SI suffixes (1000-based).
func ParseBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty string")
	}

	i := 0
	hasDot := false
	for i < len(s) {
		r := rune(s[i])
		if unicode.IsDigit(r) {
			i++
		} else if r == '.' && !hasDot {
			hasDot = true
			i++
		} else {
			break
		}
	}

	numStr := s[:i]
	unitStr := strings.TrimSpace(s[i:])

	if numStr == "" {
		return 0, fmt.Errorf("invalid size %q", s)
	}

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	var multiplier float64 = 1
	switch strings.ToLower(unitStr) {
	case "":
		multiplier = 1
	case "b", "byte", "bytes":
		multiplier = 1
	case "k", "kb":
		multiplier = 1000
	case "kib":
		multiplier = 1024
	case "m", "mb":
		multiplier = 1000 * 1000
	case "mib":
		multiplier = 1024 * 1024
	case "g", "gb":
		multiplier = 1000 * 1000 * 1000
	case "gib":
		multiplier = 1024 * 1024 * 1024
	case "t", "tb":
		multiplier = 1000 * 1000 * 1000 * 1000
	case "tib":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown unit %q in %q", unitStr, s)
	}

	return uint64(val * multiplier), nil
}

// ParseSI parses a string with an SI prefix and returns the value and the remaining unit.
// Supported prefixes: k/K (10^3), m/M (10^6), g/G (10^9), t/T (10^12).
// If no prefix is present, multiplier is 1.
// It returns the parsed value, the remaining unit string, and an error if parsing failed.
func ParseSI(s string) (float64, string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, "", errors.New("empty string")
	}

	i := 0
	hasDot := false
	for i < len(s) {
		r := rune(s[i])
		if unicode.IsDigit(r) {
			i++
		} else if r == '.' && !hasDot {
			hasDot = true
			i++
		} else {
			break
		}
	}

	numStr := s[:i]
	rem := strings.TrimSpace(s[i:])

	if numStr == "" {
		return 0, "", fmt.Errorf("invalid SI value %q", s)
	}

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, "", err
	}

	if rem == "" {
		return val, "", nil
	}

	// Check if the first character of rem is an SI prefix
	prefix := rem[0]
	var multiplier float64 = 1
	hasPrefix := true
	switch prefix {
	case 'k', 'K':
		multiplier = 1e3
	case 'm', 'M':
		multiplier = 1e6
	case 'g', 'G':
		multiplier = 1e9
	case 't', 'T':
		multiplier = 1e12
	default:
		hasPrefix = false
	}

	if hasPrefix {
		val *= multiplier
		rem = strings.TrimSpace(rem[1:])
	}

	return val, rem, nil
}
