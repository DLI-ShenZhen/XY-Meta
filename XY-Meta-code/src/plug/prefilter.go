//The preprocessing modele
//This modele is mainly used to filter the noise signals in spectrum
package plug

import (
	"dataframe"
	"math"
	"reader"
)

type Filter struct {
	Error_type          int
	Tolerance           float64
	Activate            bool
	Threshold_intensity float64
}

func Spectralclear(spectral_experiment []reader.Spectrum, Input_param dataframe.Parameters) []reader.Spectrum {
	if Input_param.Clear == true {
		for i := 0; i < len(spectral_experiment); i++ {
			start_index := 0
			var peaks_filter []reader.Signal
			if Input_param.Merge_type == 0 {
				for start_index < len(spectral_experiment[i].Peaks) {
					if spectral_experiment[i].Peaks[start_index].Peak_intensity > Input_param.Threshold_peaks {{
						if len(peaks_filter) == 0 {
							peaks_filter = append(peaks_filter, spectral_experiment[i].Peaks[start_index])
						} else {
							if spectral_experiment[i].Peaks[start_index].Peak_mass != peaks_filter[len(peaks_filter)-1].Peak_mass {
								peaks_filter = append(peaks_filter, spectral_experiment[i].Peaks[start_index])
							}
						}
						for h := start_index + 1; h < len(spectral_experiment[i].Peaks); h++ {
							if math.Abs(spectral_experiment[i].Peaks[start_index].Peak_mass-spectral_experiment[i].Peaks[h].Peak_mass) >= Input_param.Merge_tolerance {
								peaks_filter = append(peaks_filter, spectral_experiment[i].Peaks[h])
								start_index = h
								break
							}

						}
						start_index++
					} else {
						start_index++
					}
				}
			} else if Input_param.Merge_type == 1 { 
				for start_index < len(spectral_experiment[i].Peaks) {
					if spectral_experiment[i].Peaks[start_index].Peak_intensity > Input_param.Threshold_peaks {
						if len(peaks_filter) == 0 {
							peaks_filter = append(peaks_filter, spectral_experiment[i].Peaks[start_index])
						} else {
							if spectral_experiment[i].Peaks[start_index].Peak_mass != peaks_filter[len(peaks_filter)-1].Peak_mass {
								peaks_filter = append(peaks_filter, spectral_experiment[i].Peaks[start_index])
							}
						}
						for h := start_index + 1; h < len(spectral_experiment[i].Peaks); h++ {
							if math.Abs(spectral_experiment[i].Peaks[start_index].Peak_mass-spectral_experiment[i].Peaks[h].Peak_mass) >= Input_param.Merge_tolerance {
								peaks_filter = append(peaks_filter, spectral_experiment[i].Peaks[h])
								start_index = h
								break
							}
						}
						start_index++
					} else {
						start_index++
					}
				}
			}
			spectral_experiment[i].Peaks = peaks_filter
		}
		return spectral_experiment
	} else {
		return spectral_experiment
	}
}
