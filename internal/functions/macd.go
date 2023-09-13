package functions

func CalculateMACD(data []float64, shortPeriod, longPeriod, signalPeriod int) ([]float64, []float64) {
	shortEMA := CalculateEMA(data, shortPeriod)
	longEMA := CalculateEMA(data, longPeriod)

	// Вычисляем MACD как разницу между short EMA и long EMA
	var macd []float64
	for i := 0; i < len(data); i++ {
		macd = append(macd, shortEMA[i]-longEMA[i])
	}

	// Вычисляем сигнальную линию MACD как EMA от MACD
	signalLine := CalculateEMA(macd, signalPeriod)

	return macd, signalLine
}
