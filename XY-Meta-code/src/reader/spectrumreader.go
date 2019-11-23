//The spectrum reader
//This module is used to parase the query mgf file and reference mgf file
//Some main data struct which will be aplay to calculate the match score in spectra match modele
package reader

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Spectrum struct {
	Spectrl_number string
	Precusor_mass  float64
	Retention_time float64
	Ionmode        string
	Charge         float64
	Peaks          []Signal
	Peak_counts    int
}

type Signal struct {
	Peak_mass      float64
	Peak_intensity float64
}

func Spectrumreader(path_spectrum string) []Spectrum {
	var spectral_list []Spectrum
	f_spectrum, err := os.Open(path_spectrum)
	if err != nil {
		fmt.Println(err.Error())
	}
	mgf_unit := Spectrum{}
	buf_spectral := bufio.NewReader(f_spectrum)
	for {
		a, err_read := buf_spectral.ReadString('\n')
		if err_read != nil {
			break
		}
		a = strings.Replace(a, "\r", "", -1)
		a = strings.TrimSuffix(a, "\n")
		kemp := strings.Split(a, "\t")
		if len(kemp) == 1 {
			bemp := strings.Split(a, "=")
			if bemp[0] == "BEGIN IONS" {
				mgf_unit = Spectrum{}
				continue
			} else if bemp[0] == "END IONS" {
				if len(mgf_unit.Peaks) < 1 {
					continue
				} else {
					mgf_unit.Peak_counts = len(mgf_unit.Peaks)
					spectral_list = append(spectral_list, mgf_unit)
					continue
				}
			} else if bemp[0] == "TITLE" || bemp[0] == "NAME" {
				if len(bemp) > 2 {
					var spectrum_num string
					for o := 1; o < len(bemp); o++ {
						if o != len(bemp)-1 {
							spectrum_num = spectrum_num + bemp[o] + "="
						} else {
							spectrum_num = spectrum_num + bemp[o]
						}
					}
					mgf_unit.Spectrl_number = spectrum_num
				} else {
					mgf_unit.Spectrl_number = bemp[1]
				}
			} else if bemp[0] == "RTINSECONDS" {
				temp_inf, _ := strconv.ParseFloat(bemp[1], 64)
				mgf_unit.Retention_time = temp_inf
			} else if bemp[0] == "PEPMASS" {
				temp := strings.Split(bemp[1], " ")
				temp_inf, _ := strconv.ParseFloat(temp[0], 64)
				mgf_unit.Precusor_mass = temp_inf
			} else if bemp[0] == "CHARGE" {
				temp_charge := strings.Split(bemp[1], "+")
				if len(temp_charge) < 2 {
					continue
				} else {
					temp_inf, _ := strconv.ParseFloat(temp_charge[0], 64)
					mgf_unit.Charge = temp_inf
				}
			} else if len(bemp) == 1 {
				if len(a) == 0 {
					continue
				} else {
					signal_temp := Signal{}
					temp := strings.Split(a, " ")
					mz, _ := strconv.ParseFloat(temp[0], 64)
					intensity, _ := strconv.ParseFloat(temp[1], 64)
					signal_temp.Peak_mass = mz
					signal_temp.Peak_intensity = intensity
					mgf_unit.Peaks = append(mgf_unit.Peaks, signal_temp)
				}
			}
		} else {
			signal_temp := Signal{}
			temp := strings.Split(a, "\t")
			mz, _ := strconv.ParseFloat(temp[0], 64)
			intensity, _ := strconv.ParseFloat(temp[1], 64)
			signal_temp.Peak_mass = mz
			signal_temp.Peak_intensity = intensity
			mgf_unit.Peaks = append(mgf_unit.Peaks, signal_temp)
		}
	}
	return spectral_list
}
