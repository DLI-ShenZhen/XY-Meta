//Data normalize module
//This package is used to normalize the intensity of spectral
//Some data calculation module is also redacted at this package
package plug

import (
	"fmt"
	"math"
	"reader"
)

func SpectralNormalize(spectral_original []reader.Spectrum) []reader.Spectrum {
	fmt.Println("Spectral Normalization!")
	for i := 0; i < len(spectral_original); i++ {
		var calculation_signal []float64
		for j := 0; j < len(spectral_original[i].Peaks); j++ {
			calculation_signal = append(calculation_signal, spectral_original[i].Peaks[j].Peak_intensity)
		}
		min_max(calculation_signal)
		for j := 0; j < len(spectral_original[i].Peaks); j++ {
			spectral_original[i].Peaks[j].Peak_intensity = calculation_signal[j]
		}
	}
	return spectral_original
}

func min_max(dataline []float64) []float64 {
	min := Min(dataline)
	max := Max(dataline)
	if min == max {
		for i := 0; i < len(dataline); i++ {
			dataline[i] = 1
		}
	} else {
		for i := 0; i < len(dataline); i++ {
			normalize_x := (1 / max) * dataline[i]
			dataline[i] = normalize_x
		}
	}

	return dataline
}

func Min(data []float64) float64 {
	initial := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] < initial {
			initial = data[i]
		}
	}
	return initial
}

func Max(data []float64) float64 {
	initial := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > initial {
			initial = data[i]
		}
	}
	return initial
}

func Mean(data []float64) float64 {
	var sum float64
	for i := 0; i < len(data); i++ {
		sum = sum + data[i]
	}
	return sum / float64(len(data))
}

func Standard_deviation(data []float64) float64 {
	mean := Mean(data)
	var variance float64
	for i := 0; i < len(data); i++ {
		variance += math.Pow((data[i] - mean), 2)
	}
	return math.Sqrt(variance / float64(len(data)))
}

func Grubbs(test, mean, standard_deviation float64) float64 {
	return (math.Abs(test - mean)) / standard_deviation
}

func PPm(actual, theoretical float64) float64 {
	return (math.Abs(actual-theoretical) / theoretical) * math.Pow10(6)
}

func Precursor_match_Da(query_mz float64, reference_mz float64, tolerance float64) bool {
	precursormapping := false
	if math.Abs(query_mz-reference_mz) <= tolerance {
		precursormapping = true
	}
	return precursormapping
}

func Precursor_match_PPm(query_mz float64, reference_mz float64, tolerance float64) bool {
	precursormapping := false
	if PPm(query_mz, reference_mz) <= tolerance {
		precursormapping = true
	}
	return precursormapping
}

func SelectKthMin(s []float64, k int) float64 {
	k--
	lo, hi := 0, len(s)-1
	for {
		j := partition(s, lo, hi)
		if j < k {
			lo = j + 1
		} else if j > k {
			hi = j - 1
		} else {
			return s[k]
		}
	}
}

func partition(s []float64, lo, hi int) int {
	i, j := lo, hi+1
	for {
		for {
			i++
			if i == hi || s[i] > s[lo] {
				break
			}
		}
		for {
			j--
			if j == lo || s[j] <= s[lo] {
				break
			}
		}
		if i >= j {
			break
		}
		swap(s, i, j)
	}
	swap(s, lo, j)
	return j
}

func swap(s []float64, i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func SelectMid(s []float64) float64 {
	return SelectKthMin(s, len(s)/2)
}
