package main

// RSIAnalyzer represents an object for analyzing Relative Strength Index (RSI).
type RSIAnalyzer struct {
	window    int       // Length of the window
	values    []float64 // Array of float values
	sumGains  float64   // Sum of gains
	sumLosses float64   // Sum of losses
	lastValue float64   // Last value pushed
	first     bool
}

// NewRSIAnalyzer creates a new RSIAnalyzer object with the specified window length.
func NewRSIAnalyzer(window int) *RSIAnalyzer {
	return &RSIAnalyzer{
		window: window,
		values: []float64{},
		first:  true,
	}
}

// PushValue pushes a new float value to the analyzer.
func (r *RSIAnalyzer) PushValue(value float64) {
	if r.first == true {
		r.lastValue = value
		r.first = false
		return
	}
	diff := value - r.lastValue
	if diff < 0 {
		r.sumLosses -= diff
	} else {
		r.sumGains += diff
	}
	r.lastValue = value
	r.values = append(r.values, diff)
	if len(r.values) > r.window {
		oldestValue := r.values[0]
		if oldestValue > 0 {
			r.sumGains -= oldestValue
		} else {
			r.sumLosses += oldestValue
		}
		r.removeOldestValue()
	}
}

// removeOldestValue removes the oldest value from the values array.
func (r *RSIAnalyzer) removeOldestValue() {
	r.values = r.values[1:]
}

// SumGains returns the sum of all gains.
func (r *RSIAnalyzer) SumGains() float64 {
	return r.sumGains
}

// SumLosses returns the sum of all losses.
func (r *RSIAnalyzer) SumLosses() float64 {
	return r.sumLosses
}

// func main() {
// 	// Example usage
// 	rsiAnalyzer := NewRSIAnalyzer(5)

// 	// Push some values
// 	values := []float64{21456.65, 21349.4, 21731.4, 21710.8, 21894.55, 21622.4, 21352.6, 21853.8,
// 		21782.5, 22040.7, 22212.7, 22338.75, 22493.55, 22023.35, 22096.75, 22326.9,
// 		22513.7, 22519.4, 22147, 22419.95, 22648.2}
// 	for _, value := range values {
// 		rsiAnalyzer.PushValue(value)
// 	}

// 	// Print sum of gains and losses
// 	fmt.Printf("Sum of gains: %.2f\n", rsiAnalyzer.SumGains()/5)
// 	fmt.Printf("Sum of losses: %.2f\n", rsiAnalyzer.SumLosses()/5)
// 	rs := (rsiAnalyzer.SumGains() / rsiAnalyzer.SumLosses())
// 	rsi := 100 - (100 / (1 + rs))
// 	fmt.Printf("%.2f\n", rsi)
// }
