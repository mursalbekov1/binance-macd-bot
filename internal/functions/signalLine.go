package functions

//// Функция для вычисления сигнальной линии MACD
//func CalculateSignalLine(macd []float64, signalPeriod int) []float64 {
//	signalLine := make([]float64, len(macd))
//	multiplier := 2.0 / float64(signalPeriod+1)
//
//	// Вычисляем начальное значение сигнальной линии как первое значение MACD
//	signalLine[signalPeriod-1] = macd[signalPeriod-1]
//
//	// Вычисляем сигнальную линию для остальных элементов
//	for i := signalPeriod; i < len(macd); i++ {
//		signalLine[i] = (macd[i]-signalLine[i-1])*multiplier + signalLine[i-1]
//	}
//
//	return signalLine
//}
