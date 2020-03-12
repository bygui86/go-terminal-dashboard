package controller

import "container/ring"

// INFO: Package ring implements operations on circular lists.

// statRing, encapsulate container/ring,  adapt to chart data
type statRing struct {
	ring *ring.Ring
}

func newChartRing(n int) *statRing {
	return &statRing{
		ring: ring.New(n),
	}
}

func (r *statRing) push(n uint64) {
	r.ring.Value = n
	r.ring = r.ring.Next()
}

// convertData converts underline data to float64
func (r *statRing) convertData() []float64 {
	var l []float64
	r.ring.Do(func(x interface{}) {
		if v, ok := x.(uint64); ok {
			l = append(l, float64(v))
		} else {
			l = append(l, 0.0)
		}
	})
	return l
}

// normalizedData return normalized data between [0,1]
func (r *statRing) normalizedData() []float64 {
	max := r.findMax()
	if max == 0 {
		return make([]float64, r.ring.Len(), r.ring.Len())
	}

	var l []float64
	r.ring.Do(func(x interface{}) {
		var pct float64
		if v, ok := x.(uint64); ok {
			pct = float64(v) / float64(max)
		}
		l = append(l, pct)
	})
	return l
}

func (r *statRing) findMax() uint64 {
	var max uint64
	r.ring.Do(func(x interface{}) {
		if v, ok := x.(uint64); ok && v > max {
			max = v
		}
	})
	return max
}
