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

package loadtest

import (
	"testing"
)

func TestMovingAverage(t *testing.T) {
	ma := NewMovingAverage(3)

	// Test empty
	if _, ok := ma.Min(); ok {
		t.Error("Min() on empty should return false")
	}
	if _, ok := ma.Max(); ok {
		t.Error("Max() on empty should return false")
	}
	if avg := ma.Avg(); avg != 0 {
		t.Errorf("Avg() on empty should be 0, got %f", avg)
	}

	// Add 1
	ma.Add(10)
	assertVal(t, ma, 10, 10, 10)

	// Add 2
	ma.Add(20)
	assertVal(t, ma, 15, 10, 20)

	// Add 3 (full)
	ma.Add(30)
	assertVal(t, ma, 20, 10, 30)

	// Add 4 (overflow, 10 should be removed)
	ma.Add(40)
	assertVal(t, ma, 30, 20, 40)

	// Add 5 (20 should be removed)
	ma.Add(15)
	assertVal(t, ma, 28.333333333333332, 15, 40)
}

func assertVal(t *testing.T, ma *MovingAverage, expectedAvg, expectedMin, expectedMax float64) {
	t.Helper()
	avg := ma.Avg()
	if !approxEqual(avg, expectedAvg) {
		t.Errorf("Expected Avg %f, got %f", expectedAvg, avg)
	}
	min, ok := ma.Min()
	if !ok || min != expectedMin {
		t.Errorf("Expected Min %f, got %f (ok=%t)", expectedMin, min, ok)
	}
	max, ok := ma.Max()
	if !ok || max != expectedMax {
		t.Errorf("Expected Max %f, got %f (ok=%t)", expectedMax, max, ok)
	}
}

func approxEqual(a, b float64) bool {
	const epsilon = 1e-9
	return (a-b) < epsilon && (b-a) < epsilon
}
