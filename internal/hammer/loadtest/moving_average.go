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

import "sync"

// MovingAverage calculates the moving average of a stream of floats.
// It is not thread-safe.
type MovingAverage struct {
	size  int
	slots []float64
	head  int
	count int
	sum   float64
}

// NewMovingAverage creates a new MovingAverage.
func NewMovingAverage(size int) *MovingAverage {
	return &MovingAverage{
		size:  size,
		slots: make([]float64, size),
	}
}

// Add adds a new value to the moving average.
func (ma *MovingAverage) Add(val float64) {
	if ma.count < ma.size {
		ma.count++
	} else {
		ma.sum -= ma.slots[ma.head]
	}
	ma.slots[ma.head] = val
	ma.sum += val
	ma.head = (ma.head + 1) % ma.size
}

// Avg returns the average of the values.
func (ma *MovingAverage) Avg() float64 {
	if ma.count == 0 {
		return 0
	}
	return ma.sum / float64(ma.count)
}

// Min returns the minimum value in the current window.
func (ma *MovingAverage) Min() (float64, bool) {
	if ma.count == 0 {
		return 0, false
	}
	m := ma.slots[0]
	for i := 1; i < ma.count; i++ {
		if ma.slots[i] < m {
			m = ma.slots[i]
		}
	}
	return m, true
}

// Max returns the maximum value in the current window.
func (ma *MovingAverage) Max() (float64, bool) {
	if ma.count == 0 {
		return 0, false
	}
	m := ma.slots[0]
	for i := 1; i < ma.count; i++ {
		if ma.slots[i] > m {
			m = ma.slots[i]
		}
	}
	return m, true
}

// ConcurrentMovingAverage is a thread-safe wrapper around MovingAverage.
type ConcurrentMovingAverage struct {
	mu sync.RWMutex
	ma *MovingAverage
}

// NewConcurrentMovingAverage creates a new ConcurrentMovingAverage.
func NewConcurrentMovingAverage(size int) *ConcurrentMovingAverage {
	return &ConcurrentMovingAverage{
		ma: NewMovingAverage(size),
	}
}

// Add adds a new value to the moving average.
func (c *ConcurrentMovingAverage) Add(val float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ma.Add(val)
}

// Avg returns the average of the values.
func (c *ConcurrentMovingAverage) Avg() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ma.Avg()
}

// Min returns the minimum value in the current window.
func (c *ConcurrentMovingAverage) Min() (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ma.Min()
}

// Max returns the maximum value in the current window.
func (c *ConcurrentMovingAverage) Max() (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ma.Max()
}
