//Database modification module
//The reference is read from local directory and is modified by our method in this package
package database

import (
	"dataframe"
	"fmt"
	"math"
	"math/rand"
	"plug"
	"reader"
)

func Inputdatabase(path_forward_database string, Input_param dataframe.Parameters) ([]reader.Spectrum, []reader.Spectrum) {
	fmt.Println("Loading Target Database")
	forward_database := reader.Spectrumreader(path_forward_database)
	fmt.Println("Generating Decoy Database!")
	Decoy_database := Decoygenerator(forward_database, Input_param.Min_mass, 2000*Input_param.Max_mass, Input_param.Decoy_similitude)
	plug.Spectralclear(forward_database, Input_param)
	plug.SpectralNormalize(forward_database)
	plug.Spectralclear(Decoy_database, Input_param)
	plug.SpectralNormalize(Decoy_database)
	return forward_database, Decoy_database
}

func Decoygenerator(forward_database []reader.Spectrum, lower_bound, upper_bound int, Decoy_similitude float64) []reader.Spectrum {
	approxication_position := make(map[float64]float64)
	peaks_map := make(map[float64][]int)
	var permutation []float64
	for i := lower_bound; i <= upper_bound; i++ {
		var err_mass float64
		err_mass = 0.0005 * float64(i)
		var spectrum_num_list []int
		permutation = append(permutation, err_mass)
		peaks_map[err_mass] = spectrum_num_list
	}
	for i := 0; i < len(forward_database); i++ {
		peaks_index := chop(forward_database[i].Precusor_mass, permutation)
		approxication_position[forward_database[i].Precusor_mass] = peaks_index
		peaks_map[peaks_index] = append(peaks_map[peaks_index], i)
	}
	var Decoy_database []reader.Spectrum
	r := rand.New(rand.NewSource(int64(len(forward_database))))
	for i := 0; i < len(forward_database); i++ {
		var decoy_spectrum reader.Spectrum
		var signal_random []reader.Signal
		if forward_database[i].Precusor_mass == 0 {
			var fragment_mass []float64
			for u := 0; u < len(forward_database[i].Peaks); u++ {
				fragment_mass = append(fragment_mass, forward_database[i].Peaks[u].Peak_mass)
			}
			forward_database[i].Precusor_mass = plug.Max(fragment_mass)
		}
		signal_spectral_mz := approxication_position[forward_database[i].Precusor_mass]
		signal_spectral := peaks_map[signal_spectral_mz]
		if len(signal_spectral) < 15 {
			for {
				if len(signal_spectral) == 15 {
					break
				}
				random_temp_index := r.Intn(len(forward_database))
				insert_index := true
				for _, test_index := range signal_spectral {
					if test_index == random_temp_index {
						insert_index = false
						break
					}
				}
				if insert_index == true {
					signal_spectral = append(signal_spectral, random_temp_index)
				}
			}
		}

		for _, spectrum_index := range signal_spectral {
			for j := 0; j < len(forward_database[spectrum_index].Peaks); j++ {
				if forward_database[spectrum_index].Peaks[j].Peak_mass <= forward_database[i].Precusor_mass {
					signal_random = append(signal_random, forward_database[spectrum_index].Peaks[j])
				}
			}
		}
		signal_random_counts := len(forward_database[i].Peaks)
		decoy_spectrum.Charge = forward_database[i].Charge
		decoy_spectrum.Peak_counts = signal_random_counts
		decoy_spectrum.Precusor_mass = forward_database[i].Precusor_mass - 1.00794
		decoy_spectrum.Retention_time = forward_database[i].Retention_time
		decoy_spectrum.Spectrl_number = forward_database[i].Spectrl_number
		decoy_spectrum.Peaks = append(decoy_spectrum.Peaks, forward_database[i].Peaks[forward_database[i].Peak_counts-1])
		fillnum := int(Decoy_similitude * float64(signal_random_counts))
		nums := generateRandomNumber(0, len(forward_database[i].Peaks)-1, fillnum)

		var sort_peaks_temp []float64
		for _, signal_unit := range signal_random {
			sort_peaks_temp = append(sort_peaks_temp, signal_unit.Peak_mass)
		}
		quickSort(sort_peaks_temp, 0, len(sort_peaks_temp))
		for j := 0; j < len(sort_peaks_temp); j++ {
			for h := j; h < len(signal_random); h++ {
				if signal_random[h].Peak_mass == sort_peaks_temp[j] {
					swapping := signal_random[j]
					signal_random[j] = signal_random[h]
					signal_random[h] = swapping
					break
				}
			}
		}

		for s := 0; s < len(nums); s++ {
			decoy_spectrum.Peaks = append(decoy_spectrum.Peaks, forward_database[i].Peaks[nums[s]])
		}

		signal_random_counts = forward_database[i].Peak_counts - fillnum
		random_initial := r.Intn(1)
		step_length := int(float64(1/2.0)*float64(len(signal_random))) / signal_random_counts
		search_index_front := int(float64(1/2.0) * float64(len(signal_random)))
		if len(signal_random) != 0 {
			for step := 0; step < signal_random_counts; step++ {
				search_index := search_index_front + step*step_length + step_length - random_initial
				if search_index >= len(signal_random) {
					decoy_spectrum.Peaks = append(decoy_spectrum.Peaks, signal_random[len(signal_random)-1])
					break
				}
				decoy_spectrum.Peaks = append(decoy_spectrum.Peaks, signal_random[search_index])
			}
		}

		decoy_spectrum = randomshift(decoy_spectrum, 0.7)
		decoy_spectrum.Spectrl_number = decoy_spectrum.Spectrl_number + "_REV"
		var sort_peaks []float64
		for _, signal_unit := range decoy_spectrum.Peaks {
			sort_peaks = append(sort_peaks, signal_unit.Peak_mass)
		}
		quickSort(sort_peaks, 0, len(sort_peaks)
		for j := 0; j < len(sort_peaks); j++ {
			for h := j; h < len(decoy_spectrum.Peaks); h++ {
				if decoy_spectrum.Peaks[h].Peak_mass == sort_peaks[j] {
					swapping := decoy_spectrum.Peaks[j]
					decoy_spectrum.Peaks[j] = decoy_spectrum.Peaks[h]
					decoy_spectrum.Peaks[h] = swapping
					break
				}
			}
		}

		Decoy_database = append(Decoy_database, decoy_spectrum)
	}
	fmt.Println("Generate Decoy Done!")
	return Decoy_database
}

func chop(search_dount float64, permutation []float64) float64 {
	var approximation float64
	start := 0
	end := len(permutation)
	middle := (end - start) / 2
	for {
		if permutation[middle]-search_dount > 0.0000 {
			if permutation[middle]-search_dount == 0.0005 {
				approximation = search_dount
				break
			} else if permutation[middle]-search_dount > 0.0005 {
				end = middle
				middle = (end-start)/2 + start
			} else {
				edge_left := math.Abs(permutation[middle-1] - search_dount)
				edge_right := math.Abs(permutation[middle] - search_dount)
				if edge_left <= edge_right {
					approximation = permutation[middle-1]
					break
				} else if edge_left > edge_right {
					approximation = permutation[middle]
					break
				}
			}
		} else if permutation[middle]-search_dount < 0.0000 {
			if permutation[middle]-search_dount == -0.0005 {
				approximation = search_dount
				break
			} else if permutation[middle]-search_dount < -0.0005 {
				start = middle
				middle = (end-start)/2 + start
			} else {
				edge_right := math.Abs(permutation[middle+1] - search_dount)
				edge_left := math.Abs(permutation[middle] - search_dount)
				if edge_left <= edge_right {
					approximation = permutation[middle]
					break
				} else if edge_left > edge_right {
					approximation = permutation[middle+1]
					break
				}
			}

		} else {
			approximation = search_dount
			break
		}
	}
	return approximation
}

func randomshift(decoy_spectrum reader.Spectrum, Decoy_similitude float64) reader.Spectrum {
	random_shift_counts := float64(len(decoy_spectrum.Peaks)) * (1 - Decoy_similitude)
	var step_length int
	if random_shift_counts < 1 {
		step_length = 1
	} else {
		step_length = len(decoy_spectrum.Peaks) / int(random_shift_counts)
	}
	r := rand.New(rand.NewSource(int64(len(decoy_spectrum.Peaks))))
	random_seed := r.Intn(3)
	for i := 0; i < int(random_shift_counts); i++ {
		random_search_index := step_length*i + random_seed
		if random_search_index >= int(random_shift_counts) {
			branch_random := r.Intn(1)
			if branch_random == 0 {
				decoy_spectrum.Peaks[len(decoy_spectrum.Peaks)-1].Peak_mass = decoy_spectrum.Peaks[len(decoy_spectrum.Peaks)-1].Peak_mass + decoy_spectrum.Precusor_mass*0.000003
			} else {
				decoy_spectrum.Peaks[len(decoy_spectrum.Peaks)-1].Peak_mass = decoy_spectrum.Peaks[len(decoy_spectrum.Peaks)-1].Peak_mass - decoy_spectrum.Precusor_mass*0.000005
			}
			break
		}
		branch_random := r.Intn(1)
		if branch_random == 0 {
			decoy_spectrum.Peaks[random_search_index].Peak_mass = decoy_spectrum.Peaks[random_search_index].Peak_mass + decoy_spectrum.Precusor_mass*0.000005
		} else {
			decoy_spectrum.Peaks[random_search_index].Peak_mass = decoy_spectrum.Peaks[random_search_index].Peak_mass - decoy_spectrum.Precusor_mass*0.000005
		}
	}

	return decoy_spectrum
}

//快速排序模块
func swap(a float64, b float64) (float64, float64) {
	return b, a
}

func partition(aris []float64, begin int, end int) int {
	pvalue := aris[begin]
	i := begin
	j := begin + 1
	for j < end {
		if aris[j] < pvalue {
			i++
			aris[i], aris[j] = swap(aris[i], aris[j])
		}
		j++
	}
	aris[i], aris[begin] = swap(aris[i], aris[begin])
	return i
}

func quickSort(aris []float64, begin int, end int) {
	if begin+1 < end {
		mid := partition(aris, begin, end)
		quickSort(aris, begin, mid)
		quickSort(aris, mid+1, end)
	}
}

func generateRandomNumber(start int, end int, count int) []int {
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	r := rand.New(rand.NewSource(int64(count)))
	for len(nums) < count {
		num := r.Intn((end - start)) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}
