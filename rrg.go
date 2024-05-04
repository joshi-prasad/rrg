package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func ReadCsv() ([][]string, error) {
	// Open the CSV file
	file, err := os.Open("index_weekly.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	records := [][]string{}
	if record, err := reader.Read(); err != nil {
		fmt.Println("Error reading first header:", err)
		return nil, err
	} else {
		records = append(records, record)
	}

	// Skip the second header line
	if _, err := reader.Read(); err != nil {
		fmt.Println("Error reading second header:", err)
		return nil, err
	}

	// Read and process data records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading record:", err)
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func transposeRecords(records [][]string) map[string][]string {
	// Check if there's at least one record
	if len(records) == 0 {
		return map[string][]string{}
	}

	// Get the header row (assuming first row contains keys)
	headers := records[0]

	// Create a map to store transposed data
	transposed := make(map[string][]string)
	for _, header := range headers {
		transposed[header] = []string{}
	}

	// Iterate through remaining records (excluding header)
	for i := 1; i < len(records); i++ {
		record := records[i]

		// Create a new slice for the current key (header)
		// transposed[headers[i-1]] = make([]string, len(record))

		// Populate the transposed data with elements from the record
		for j, element := range record {
			transposed[headers[j]] = append(transposed[headers[j]], element)
		}
	}

	// for i, records := range transposed {
	// 	transposed[i] = reverseStringArray(records)
	// }

	return transposed
}

func convertStringToFloatMap(strMap map[string][]string) map[string][]float64 {
	floatMap := make(map[string][]float64)

	for key, stringSlice := range strMap {
		floatSlice := make([]float64, len(stringSlice))
		for i, strVal := range stringSlice {
			// Convert string to float using conversion function (adjust based on your needs)
			floatVal, err := strconv.ParseFloat(strVal, 64) // 64-bit float
			if err != nil {
				// Handle conversion error (e.g., log or skip element)
				fmt.Println("Error converting string to float:", err)
				continue
			}
			floatSlice[i] = floatVal
		}
		floatMap[key] = floatSlice
	}

	return floatMap
}

func validateData(data map[string][]string) error {
	array_length := -1
	for _, records := range data {
		if array_length == -1 {
			array_length = len(records)
			continue
		}
		if array_length != len(records) {
			msg := fmt.Sprintf("Array length does not match %d != %d", array_length, len(records))
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func calculateRateOfChange(prices []float64, benchmarkPrices []float64) ([]float64, error) {
	// Check for input validity
	if len(prices) == 0 || len(prices) != len(benchmarkPrices) {
		return nil, errors.New("prices and benchmarkPrices slices must have the same length")
	}

	rsRating := make([]float64, len(prices))

	for i := range prices {
		if i == 0 {
			// No previous data for the first week
			rsRating[i] = 0
			continue
		}
		priceChange := (prices[i] - prices[i-1]) / prices[i-1] * 100
		benchmarkChange := (benchmarkPrices[i] - benchmarkPrices[i-1]) / benchmarkPrices[i-1] * 100
		rsRating[i] = priceChange / benchmarkChange
	}

	return rsRating, nil
}

func calculateRRGRating(date []string, tickerMomentum map[string][]float64, tickerRSRating map[string][]float64) map[string]string {
	// Get the index of the latest week
	latestWeekIndex := len(date) - 2

	// Initialize RRG rating map
	rrgRating := make(map[string]string)

	// Iterate through each ticker symbol
	for ticker, momentum := range tickerMomentum {
		// Get the latest momentum and RS Rating values
		latestMomentum := momentum[latestWeekIndex]
		latestRSRating := tickerRSRating[ticker][latestWeekIndex]

		// Define quadrants based on signs of momentum and RS Rating (adjust ranges as needed)
		var quadrant string
		if latestMomentum >= 0 {
			if latestRSRating >= 70 {
				quadrant = "Leading"
			} else if latestRSRating >= 50 {
				quadrant = "Improving"
			} else {
				quadrant = "NoRS"
			}
		} else if latestMomentum < 0 {
			if latestRSRating >= 70 {
				quadrant = "Weakening"
			} else if latestRSRating >= 50 {
				quadrant = "Lagging"
			} else {
				quadrant = "NoMomentum"
			}
		}

		// Assign RRG rating to the map
		rrgRating[ticker] = quadrant
	}

	return rrgRating
}

// func getPerformanceForWeek(weekly_prices []float64, prev_week int) []float64 {
// 	skip := len(weekly_prices) - prev_week - 1
// 	m3w := skip - 12  // 12 weeks = 3 months before
// 	m6w := skip - 25  // 24 weeks = 6 months before
// 	m9w := skip - 38  // 36 weeks = 9 months before
// 	m12w := skip - 52 // 48 weeks = 12 months before

// 	ppw := weekly_prices[skip]
// 	p3m := weekly_prices[m3w]
// 	p6m := weekly_prices[m6w]
// 	p9m := weekly_prices[m9w]
// 	p12m := weekly_prices[m12w]

// 	gain_3m := (ppw - p3m) / p3m
// 	gain_6m := (ppw - p6m) / p6m
// 	gain_9m := (ppw - p9m) / p9m
// 	gain_12m := (ppw - p12m) / p12m

// 	// fmt.Println("Index ", len(weekly_prices), skip, m3w, m6w, m9w, m12w)
// 	// fmt.Println("Prices ", ppw, p3m, p6m, p9m, p12m)
// 	// fmt.Println("Gain ", gain_3m, gain_6m, gain_9m, gain_12m)

// 	return []float64{gain_3m, gain_6m, gain_9m, gain_12m}
// }

func percentileRank(values []float64, forValue float64) float64 {
	// Sort the slice in ascending order
	sort.Float64s(values)

	// Find the index of the forValue using binary search
	index := sort.SearchFloat64s(values, forValue)

	// Calculate percentile rank based on the index
	var rank float64
	if index == len(values) {
		rank = 100.0 // forValue is the highest value
	} else if index == 0 {
		rank = 0.0 // forValue is the lowest value
	} else {
		// Linear interpolation for non-exact matches
		rank = float64(index) / float64(len(values)) * 100.0
	}

	return rank
}

// func calculateRsRatingForWeek(
// 	ticker_records map[string][]float64,
// 	prev_week int) map[string]float64 {

// 	perf := map[string][]float64{}
// 	perf3m := []float64{}
// 	perf6m := []float64{}
// 	perf9m := []float64{}
// 	perf12m := []float64{}
// 	for ticker, weekly_prices := range ticker_records {
// 		// fmt.Println(weekly_prices)
// 		perf_values := getPerformanceForWeek(weekly_prices, prev_week)
// 		perf[ticker] = perf_values
// 		perf3m = append(perf3m, perf_values[0])
// 		perf6m = append(perf6m, perf_values[1])
// 		perf9m = append(perf9m, perf_values[2])
// 		perf12m = append(perf12m, perf_values[3])
// 		// fmt.Println("Perf ", ticker, perf_values)
// 	}

// 	rs_rating := map[string]float64{}
// 	for ticker, perf_values := range perf {
// 		rank3m := percentileRank(perf3m, perf_values[0])
// 		rank6m := percentileRank(perf6m, perf_values[1])
// 		rank9m := percentileRank(perf9m, perf_values[2])
// 		rank12m := percentileRank(perf12m, perf_values[3])
// 		rank := (rank3m * 0.40) + (rank6m * 0.20) + (rank9m * 0.20) + (rank12m * 0.20)
// 		rs_rating[ticker] = rank
// 	}
// 	return rs_rating
// }

// computeRSI calculates the Relative Strength Index (RSI) for a given set of values and window size.
func computeRSI(values []float64, window int) (float64, error) {
	var avgGain, avgLoss float64
	for i := 1; i < (len(values) - window); i += 1 {
		diff := values[i] - values[i-1]
		if diff >= 0 {
			avgGain += diff
		} else {
			avgLoss -= diff
		}
		// fmt.Println(i, i-1)
	}

	wf := float64(window)
	avgGain = avgGain / wf
	avgLoss = avgLoss / wf

	for i := len(values) - window; i < len(values); i += 1 {
		// fmt.Println("Inside ", i, i-1)

		diff := values[i] - values[i-1]
		if diff >= 0 {
			avgGain = ((avgGain * (wf - 1)) + diff) / wf
		} else {
			avgLoss = ((avgLoss * (wf - 1)) - diff) / wf
		}
	}
	rs := avgGain / (avgLoss + 1e-14) // Avoid division by zero
	rsi := 100 - (100 / (1 + rs))
	return rsi, nil
}

func calculateRSIRatingForWeek(
	ticker_records map[string][]float64,
	prev_week int,
	rsi_window int) map[string]float64 {

	rsi := map[string]float64{}
	for ticker, weekly_prices := range ticker_records {
		// if ticker != "INDEXNSE:NIFTY_ENERGY" {
		// 	continue
		// }
		start := len(weekly_prices) - prev_week - 1 - (rsi_window * 2)
		end := start + (rsi_window * 2) + 1
		slice := weekly_prices[start:end]
		rsi_for_ticker, _ := computeRSI(slice, rsi_window)
		rsi[ticker] = rsi_for_ticker
	}
	return rsi
}

func writeToFile(data map[string]map[string][]float64, fileName string) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	// Write JSON data to a file
	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return err
	}

	fmt.Println("Data successfully written to file.")
	return nil
}

func computeRsi(prices []float64, window int) []float64 {
	fmt.Println(prices)
	var avgGain, avgLoss float64
	for i := 1; i < window; i += 1 {
		diff := prices[i] - prices[i-1]
		if diff >= 0 {
			avgGain += diff
		} else {
			avgLoss -= diff
		}
	}

	allRsi := []float64{}

	wf := float64(window)
	avgGain = avgGain / wf
	avgLoss = avgLoss / wf
	fmt.Println("Initial ", avgGain, avgLoss)
	for i := window; i < len(prices); i += 1 {
		diff := prices[i] - prices[i-1]
		if diff >= 0 {
			avgGain = ((avgGain * (wf - 1)) + diff) / wf
		} else {
			avgLoss = ((avgLoss * (wf - 1)) - diff) / wf
		}
		rs := avgGain / (avgLoss + 1e-14) // Avoid division by zero
		rsi := 100 - (100 / (1 + rs))
		fmt.Println(prices[i], diff, avgGain, avgLoss, rs, rsi)
		allRsi = append(allRsi, rsi)
	}
	return allRsi
}

func computeEma(values []float64, emaLength int) []float64 {
	if emaLength <= 0 {
		return nil // Handle invalid window size
	}

	ema := make([]float64, len(values)) // Initialize slice to store EMA values
	for i := range ema {
		if i == 0 {
			ema[i] = values[i] // First element is the initial value
		} else {
			// Calculate EMA using smoothing factor
			smoothingFactor := 2.0 / (float64(emaLength) + 1.0)
			ema[i] = smoothingFactor*values[i] + (1.0-smoothingFactor)*ema[i-1]
		}
	}
	return ema
}

func getPerformanceForWeek(
	ticker_records map[string][]float64,
	week int,
	months int) map[string]float64 {

	perf := map[string]float64{}
	prev_week := week - (months * 4)
	for ticker, weekly_prices := range ticker_records {
		cur_price := weekly_prices[week]
		prev_price := weekly_prices[prev_week]

		gain := (cur_price - prev_price) / prev_price
		perf[ticker] = gain
	}
	return perf
}

func compteRsRating(
	perf_for_all_tickers map[string]float64) map[string]float64 {

	arr := []float64{}
	for _, perf := range perf_for_all_tickers {
		arr = append(arr, perf)
	}

	result := map[string]float64{}
	for ticker, perf := range perf_for_all_tickers {
		result[ticker] = percentileRank(arr, perf)
	}
	return result
}

func computeRsRating(
	ticker_records map[string][]float64,
	values_per_record int) map[string][]float64 {

	rs_ratings := map[string][]float64{}

	skip := 52 // weeks in an year
	for i := skip; i < values_per_record; i += 1 {
		ticker_3m_rs := compteRsRating(
			getPerformanceForWeek(ticker_records, i, 3))
		ticker_6m_rs := compteRsRating(
			getPerformanceForWeek(ticker_records, i, 6))
		ticker_9m_rs := compteRsRating(
			getPerformanceForWeek(ticker_records, i, 9))
		ticker_12m_rs := compteRsRating(
			getPerformanceForWeek(ticker_records, i, 12))

		for ticker, _ := range ticker_records {
			rs3m := ticker_3m_rs[ticker]
			rs6m := ticker_6m_rs[ticker]
			rs9m := ticker_9m_rs[ticker]
			rs12m := ticker_12m_rs[ticker]
			rank := (rs3m * 0.40) + (rs6m * 0.20) + (rs9m * 0.20) + (rs12m * 0.20)
			if _, ok := rs_ratings[ticker]; !ok {
				rs_ratings[ticker] = []float64{}
			}
			rs_ratings[ticker] = append(rs_ratings[ticker], rank)
		}
	}
	return rs_ratings
}

func main() {
	records, err := ReadCsv()
	if err != nil {
		return
	}
	ticker_records_map := transposeRecords(records)
	if err = validateData(ticker_records_map); err != nil {
		fmt.Println(err)
		return
	}

	dates := ticker_records_map["Date"]
	delete(ticker_records_map, "Date")
	ticker_records := convertStringToFloatMap(ticker_records_map)
	fmt.Println(len(dates))
	fmt.Println(len(ticker_records))

	rsi_ratings := map[string][]float64{}
	rsi_ema := map[string][]float64{}

	for ticker, weekly_prices := range ticker_records {
		// if ticker != "INDEXNSE:NIFTY_FMCG" {
		// 	continue
		// }
		rsi := computeRsi(weekly_prices, 10)
		ema := computeEma(rsi, 10)
		rsi_ratings[ticker] = rsi
		rsi_ema[ticker] = ema
	}

	rs_ratings := computeRsRating(ticker_records, len(dates))
	// fmt.Println(rs_ratings)
	data := map[string]map[string][]float64{}
	for ticker, rating := range rs_ratings {
		data[ticker] = map[string][]float64{}
		data[ticker]["rs"] = rating
	}
	for ticker, rating := range rsi_ratings {
		data[ticker]["rsi"] = rating
	}
	for ticker, rating := range rsi_ema {
		data[ticker]["rsi_ema"] = rating
	}
	// fmt.Println(data)
	err = writeToFile(data, "data.json")
	if err != nil {
		fmt.Println("writing to file failed ", err)
	}
}
