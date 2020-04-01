//The reader module
//This module is used to parase the parameter and adduct
package reader

import (
	"bufio"
	"dataframe"
	"os"
	"strconv"
	"strings"
)

func Isotope(Input_param dataframe.Parameters) ([]float64, []float64, []string) {
	var add_isotope_list []float64
	var add_isotope_pre_nums []float64
	var isotope_type []string
	if Input_param.HPLC_pattern != 0 {
		if Input_param.Electric_pattern == 0 {
			if Input_param.HPLC_pattern == 1 {
				POS_isotope_path := Input_param.Current_path + "/isotope/HILIC_POS.txt"
				NEG_isotope_path := Input_param.Current_path + "/isotope/HILIC_NEG.txt"
				POS, pos_pre, pos_isotope_type := dataTable(POS_isotope_path)
				NEG, neg_pre, neg_isotope_type := dataTable(NEG_isotope_path)
				add_isotope_list = append(add_isotope_list, POS...)
				add_isotope_list = append(add_isotope_list, NEG...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, pos_pre...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, neg_pre...)
				isotope_type = append(isotope_type, pos_isotope_type...)
				isotope_type = append(isotope_type, neg_isotope_type...)
			} else if Input_param.HPLC_pattern == 2 {
				POS_isotope_path := Input_param.Current_path + "/isotope/RP_POS.txt"
				NEG_isotope_path := Input_param.Current_path + "/isotope/RP_NEG.txt"
				POS, pos_pre, pos_isotope_type := dataTable(POS_isotope_path)
				NEG, neg_pre, neg_isotope_type := dataTable(NEG_isotope_path)
				add_isotope_list = append(add_isotope_list, POS...)
				add_isotope_list = append(add_isotope_list, NEG...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, pos_pre...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, neg_pre...)
				isotope_type = append(isotope_type, pos_isotope_type...)
				isotope_type = append(isotope_type, neg_isotope_type...)
			}

		} else if Input_param.Electric_pattern == 1 {
			if Input_param.HPLC_pattern == 1 {
				POS_isotope_path := Input_param.Current_path + "/isotope/HILIC_POS.txt"
				POS, pos_pre, pos_isotope_type := dataTable(POS_isotope_path)
				add_isotope_list = append(add_isotope_list, POS...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, pos_pre...)
				isotope_type = append(isotope_type, pos_isotope_type...)
			} else if Input_param.HPLC_pattern == 2 {
				POS_isotope_path := Input_param.Current_path + "/isotope/RP_POS.txt"
				POS, pos_pre, pos_isotope_type := dataTable(POS_isotope_path)
				add_isotope_list = append(add_isotope_list, POS...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, pos_pre...)
				isotope_type = append(isotope_type, pos_isotope_type...)
			}

		} else if Input_param.Electric_pattern == 2 {
			if Input_param.HPLC_pattern == 1 {
				NEG_isotope_path := Input_param.Current_path + "/isotope/HILIC_NEG.txt"
				NEG, neg_pre, neg_isotope_type := dataTable(NEG_isotope_path)
				add_isotope_list = append(add_isotope_list, NEG...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, neg_pre...)
				isotope_type = append(isotope_type, neg_isotope_type...)
			} else if Input_param.HPLC_pattern == 2 {
				NEG_isotope_path := Input_param.Current_path + "/isotope/RP_NEG.txt"
				NEG, neg_pre, neg_isotope_type := dataTable(NEG_isotope_path)
				add_isotope_list = append(add_isotope_list, NEG...)
				add_isotope_pre_nums = append(add_isotope_pre_nums, neg_pre...)
				isotope_type = append(isotope_type, neg_isotope_type...)
			}

		}
	} else if Input_param.HPLC_pattern == 0 {
		isotope_mass, iso_pre, isotope_type_list := dataTable(Input_param.Adduct_path)
		add_isotope_list = append(add_isotope_list, isotope_mass...)
		add_isotope_pre_nums = append(add_isotope_pre_nums, iso_pre...)
		isotope_type = append(isotope_type, isotope_type_list...)
	}
	return add_isotope_list, add_isotope_pre_nums, isotope_type
}

func dataTable(data_path string) ([]float64, []float64, []string) {
	var isotope_data []float64
	var isotope_pre_nums []float64
	var isotope_type []string
	f_data, _ := os.Open(data_path)
	bread_data := bufio.NewReader(f_data)
	var counts int
	for {
		a, err := bread_data.ReadString('\n')
		if err != nil {
			break
		}
		if counts == 0 {
			counts++
			continue
		} else {
			a = strings.Replace(a, "\r", "", -1)
			a = strings.TrimSuffix(a, "\n")
			if len(a) != 0 {
				iso_splt := strings.Split(a, "	")
				isotope_mass_type, _ := strconv.ParseFloat(iso_splt[1], 64)
				isotope_data = append(isotope_data, isotope_mass_type)
				isotope_type = append(isotope_type, iso_splt[0])
				iso_pre := strings.Split(a[1:], "M")
				if len(iso_pre[0]) == 0 {
					isotope_pre_nums = append(isotope_pre_nums, 1.0)
				} else {
					isotope_pre_num, _ := strconv.ParseFloat(iso_pre[0], 64)
					isotope_pre_nums = append(isotope_pre_nums, isotope_pre_num)
				}
			}
		}
	}
	f_data.Close()
	return isotope_data, isotope_pre_nums, isotope_type
}
