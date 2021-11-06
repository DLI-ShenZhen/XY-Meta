package main

import (
	"bufio"
	"dataframe"
	"dbmodify"
	"flag"
	"fmt"
	"io"
	"log"
	"match"
	"os"
	"path"
	"path/filepath"
	"plug"
	"ranker"
	"reader"
	"strconv"
	"strings"
	"time"
	"writer"
)

var logger *log.Logger

func main() {
	fmt.Println("XY-Meta is running!")
	Start := time.Now()
	Current_path := GetCurrentDirectory()
	parameter_path_default := Current_path + "/config/parameter.default"
	parameter_path := flag.String("S", parameter_path_default, "Input a parameter file path")
	query_temp_path := flag.String("D", "null", "Set a query input file path")
	reference_temp_path := flag.String("R", "null", "Set a reference input file path")
	flag.Parse()
	Input_param := inputparam(*parameter_path, *query_temp_path, *reference_temp_path)
	Input_param.Current_path = Current_path
	query_path := path.Base(Input_param.Input)
	query_path_dir := strings.TrimSuffix(Input_param.Input, query_path)
	log_file := query_path_dir + "/log.txt"
	loging, err := os.OpenFile(log_file, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_RDWR, 666)
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	defer loging.Close()
	logger = log.New(loging, "", log.LstdFlags|log.Lshortfile)
	mw := io.MultiWriter(os.Stdout, loging)
	log.SetOutput(mw)
	defer func() {
		if r := recover(); r != nil {
			log.Println("WARN:", r)
		}
	}()
	var Adduct_isotope_list dataframe.Adduct_isotope
	Adduct_isotope_list.Isotope_mass_list, Adduct_isotope_list.Isotope_precusor_nums, Adduct_isotope_list.Isotope_type_list = reader.Isotope(Input_param)
	if Input_param.Search_pattern == 1 {
		log.Println("Search pipeline 1!")
		if Input_param.Decoy_pattern == 1 {
			log.Println("Reading experiment spectral")
			Query := reader.Spectrumreader(Input_param.Input)
			log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
			plug.SpectralNormalize(Query)
			log.Println("Reading reference spectral")
			Target, Decoty := database.Inputdatabase(Input_param.Reference, Input_param)
			log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
			database_file_path := path.Base(Input_param.Reference)
			database_dir_path := strings.TrimSuffix(Input_param.Reference, database_file_path)
			database_sufix := path.Ext(database_file_path)
			database_preix := strings.TrimSuffix(database_file_path, database_sufix)
			decoy_path := database_dir_path + database_preix + "_Decoy.mgf"
			writer.DBwriter(decoy_path, Decoty)
			Match_result_list, Reference := match.Alignment(Query, Target, Decoty, Adduct_isotope_list, Input_param)
			log.Println("Running Rank!")
			ranker.TDA(len(Target)+len(Decoty), Match_result_list, Input_param)
			log.Println("Writting Now！")
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			}
			log.Println("Cycle Done!", "cost:", time.Since(Start))
		} else if Input_param.Decoy_pattern == 2 {
			log.Println("Reading experiment spectral")
			Query := reader.Spectrumreader(Input_param.Input)
			log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
			plug.SpectralNormalize(Query)
			log.Println("Reading reference spectral")
			Target := reader.Spectrumreader(Input_param.Reference)
			log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
			plug.SpectralNormalize(Target)
			plug.Spectralclear(Target, Input_param)
			log.Println("Reading decoy spectral")
			Decoty := reader.Spectrumreader(Input_param.Decoyinput)
			plug.SpectralNormalize(Decoty)
			plug.Spectralclear(Decoty, Input_param)
			Match_result_list, Reference := match.Alignment(Query, Target, Decoty, Adduct_isotope_list, Input_param)
			log.Println("Running Rank!")
			ranker.TDA(len(Target)+len(Decoty), Match_result_list, Input_param)
			log.Println("Writting Now！")
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			}
			log.Println("Cycle Done!", "cost:", time.Since(Start))
		}
	} else if Input_param.Search_pattern == 2 {
		log.Println("Search pipeline 2!")
		log.Println("Reading experiment spectral")
		Query := reader.Spectrumreader(Input_param.Input)
		log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
		plug.SpectralNormalize(Query)
		log.Println("Reading reference spectral")
		Target := reader.Spectrumreader(Input_param.Reference)
		log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
		plug.SpectralNormalize(Target)
		plug.Spectralclear(Target, Input_param)
		var Decoty []reader.Spectrum
		Match_result_list, Reference := match.Alignment(Query, Target, Decoty, Adduct_isotope_list, Input_param)
		if Input_param.Output != "" {
			Input_param.Output = Input_param.Output + "_Result.meta"
			writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
		} else {
			Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
			Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
			filename := path.Base(Input_param.Input)
			dir_path := strings.TrimSuffix(Input_param.Input, filename)
			suffix_name := path.Ext(filename)
			prefix_name := strings.TrimSuffix(filename, suffix_name)
			filename_r := path.Base(Input_param.Reference)
			suffix_name_r := path.Ext(filename_r)
			prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
			Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Result.meta"
			writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
		}
		log.Println("Cycle Done!", "cost:", time.Since(Start))
	} else if Input_param.Search_pattern == 3 {
		log.Println("Search pipeline 3!")
		if Input_param.Decoy_pattern == 1 {
			log.Println("Reading experiment spectral")
			Query := reader.Spectrumreader(Input_param.Input)
			log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
			plug.SpectralNormalize(Query)
			log.Println("Reading reference spectral")
			Target, Decoty := database.Inputdatabase(Input_param.Reference, Input_param)
			log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
			database_file_path := path.Base(Input_param.Reference)
			database_dir_path := strings.TrimSuffix(Input_param.Reference, database_file_path)
			decoy_path := database_dir_path + "Decoy_half.mgf"
			writer.DBwriter(decoy_path, Decoty)
			Match_result_list, Reference := match.Alignment(Query, Target, Decoty, Adduct_isotope_list, Input_param)
			log.Println("Writting Now！")
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Half_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Half_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			}
			var Target_N []reader.Spectrum
			var Decoty_N []reader.Spectrum
			index_map := make(map[int]int)
			for i := 0; i < len(Match_result_list); i++ {
				if Match_result_list[i].Index < len(Target) {
					index_map[Match_result_list[i].Index] = 0
				}
			}
			for index, _ := range index_map {
				Target_N = append(Target_N, Target[index])
				Decoty_N = append(Decoty_N, Decoty[index])
			}
			log.Println("Search Again!", "cost:", time.Since(Start))
			database_file_path = path.Base(Input_param.Reference)
			database_dir_path = strings.TrimSuffix(Input_param.Reference, database_file_path)
			database_sufix := path.Ext(database_file_path)
			database_preix := strings.TrimSuffix(database_file_path, database_sufix)
			decoy_path = database_dir_path + database_preix + "_Decoy.mgf"
			writer.DBwriter(decoy_path, Decoty_N)
			Match_result_list_N, Reference_N := match.Alignment(Query, Target_N, Decoty_N, Adduct_isotope_list, Input_param)
			log.Println("Running Rank!")
			ranker.TDA(len(Target_N)+len(Decoty_N), Match_result_list_N, Input_param)
			log.Println("Writting Now！")
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Result.meta"
				writer.Mappingwriter(Match_result_list_N, Input_param.Output, Query, Reference_N, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Result.meta"
				writer.Mappingwriter(Match_result_list_N, Input_param.Output, Query, Reference_N, Adduct_isotope_list)
			}
			log.Println("Cycle Done!", "cost:", time.Since(Start))
		} else if Input_param.Decoy_pattern == 2 {
			log.Println("Reading experiment spectral")
			Query := reader.Spectrumreader(Input_param.Input)
			log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
			plug.SpectralNormalize(Query)
			log.Println("Reading reference spectral")
			Target := reader.Spectrumreader(Input_param.Reference)
			log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
			plug.SpectralNormalize(Target)
			plug.Spectralclear(Target, Input_param)
			log.Println("Reading decoy spectral")
			Decoty := reader.Spectrumreader(Input_param.Decoyinput)
			plug.SpectralNormalize(Decoty)
			plug.Spectralclear(Decoty, Input_param)
			Match_result_list, Reference := match.Alignment(Query, Target, Decoty, Adduct_isotope_list, Input_param)
			log.Println("Writting Now！")
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Half_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Half_Result.meta"
				writer.Mappingwriter(Match_result_list, Input_param.Output, Query, Reference, Adduct_isotope_list)
			}
			var Target_N []reader.Spectrum
			var Decoty_N []reader.Spectrum
			index_map := make(map[int]int)
			for i := 0; i < len(Match_result_list); i++ {
				if Match_result_list[i].Index < len(Target) {
					index_map[Match_result_list[i].Index] = 0
				}
			}
			for index, _ := range index_map {
				Target_N = append(Target_N, Target[index])
				Decoty_N = append(Decoty_N, Decoty[index])
			}
			log.Println("Search Again!", "cost:", time.Since(Start))
			Match_result_list_N, Reference_N := match.Alignment(Query, Target_N, Decoty_N, Adduct_isotope_list, Input_param)
			log.Println("Running Rank!")
			ranker.TDA(len(Target_N)+len(Decoty_N), Match_result_list_N, Input_param)
			log.Println("Writting Now！")
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Result.meta"
				writer.Mappingwriter(Match_result_list_N, Input_param.Output, Query, Reference_N, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Result.meta"
				writer.Mappingwriter(Match_result_list_N, Input_param.Output, Query, Reference_N, Adduct_isotope_list)
			}
			log.Println("Cycle Done!", "cost:", time.Since(Start))
		}
	} else if Input_param.Search_pattern == 4 {
		log.Println("Search pipeline 4!")
		if Input_param.Decoy_pattern == 1 {
			log.Println("Reading experiment spectral")
			Query := reader.Spectrumreader(Input_param.Input)
			log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
			plug.SpectralNormalize(Query)
			log.Println("Reading reference spectral")
			Target, Decoty := database.Inputdatabase(Input_param.Reference, Input_param)
			log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
			database_file_path := path.Base(Input_param.Reference)
			database_dir_path := strings.TrimSuffix(Input_param.Reference, database_file_path)
			database_sufix := path.Ext(database_file_path)
			database_preix := strings.TrimSuffix(database_file_path, database_sufix)
			decoy_path := database_dir_path + database_preix + "_Decoy.mgf"
			writer.DBwriter(decoy_path, Decoty)
			var empy_db []reader.Spectrum
			Match_result_list_target, Reference_target := match.Alignment(Query, Target, empy_db, Adduct_isotope_list, Input_param)
			Match_result_list_decoy, Reference_decoy := match.Alignment(Query, Decoty, empy_db, Adduct_isotope_list, Input_param)
			ranker.STDA(len(Target), Match_result_list_target, Match_result_list_decoy, Input_param)
			Reference_Total := append(Reference_target, Reference_decoy...)
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Separated_Result.meta"
				writer.Mappingwriter(Match_result_list_target, Input_param.Output, Query, Reference_Total, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Separated_Result.meta"
				writer.Mappingwriter(Match_result_list_target, Input_param.Output, Query, Reference_Total, Adduct_isotope_list)
			}
		} else if Input_param.Decoy_pattern == 2 {
			log.Println("Reading experiment spectral")
			Query := reader.Spectrumreader(Input_param.Input)
			log.Println("Query spectral number ", len(Query), "cost:", time.Since(Start))
			plug.SpectralNormalize(Query)
			log.Println("Reading reference_target spectral")
			var empy_db []reader.Spectrum
			Target := reader.Spectrumreader(Input_param.Reference)
			log.Println("Reference spectral number ", len(Target), "cost:", time.Since(Start))
			plug.SpectralNormalize(Target)
			plug.Spectralclear(Target, Input_param)
			Match_result_list_target, Reference_target := match.Alignment(Query, Target, empy_db, Adduct_isotope_list, Input_param)
			log.Println("Reading reference_decoy spectral")
			Decoty := reader.Spectrumreader(Input_param.Decoyinput)
			log.Println("Reference spectral number ", len(Decoty), "cost:", time.Since(Start))
			plug.SpectralNormalize(Decoty)
			plug.Spectralclear(Decoty, Input_param)
			Match_result_list_decoy, Reference_decoy := match.Alignment(Query, Decoty, empy_db, Adduct_isotope_list, Input_param)
			ranker.STDA(len(Target), Match_result_list_target, Match_result_list_decoy, Input_param)
			Reference_Total := append(Reference_target, Reference_decoy...)
			if Input_param.Output != "" {
				Input_param.Output = Input_param.Output + "_Separated_Result.meta"
				writer.Mappingwriter(Match_result_list_target, Input_param.Output, Query, Reference_Total, Adduct_isotope_list)
			} else {
				Input_param.Input = strings.Replace(Input_param.Input, "\\", "/", -1)
				Input_param.Reference = strings.Replace(Input_param.Reference, "\\", "/", -1)
				filename := path.Base(Input_param.Input)
				dir_path := strings.TrimSuffix(Input_param.Input, filename)
				suffix_name := path.Ext(filename)
				prefix_name := strings.TrimSuffix(filename, suffix_name)
				filename_r := path.Base(Input_param.Reference)
				suffix_name_r := path.Ext(filename_r)
				prefix_name_r := strings.TrimSuffix(filename_r, suffix_name_r)
				Input_param.Output = dir_path + prefix_name + "_" + prefix_name_r + "_Separated_Result.meta"
				writer.Mappingwriter(Match_result_list_target, Input_param.Output, Query, Reference_Total, Adduct_isotope_list)
			}
		}
		log.Println("Cycle Done!", "cost:", time.Since(Start))
	}
}

func inputparam(parameter_path, query_temp_path, reference_temp_path string) dataframe.Parameters {
	var Param dataframe.Parameters
	f_parameter, err_p := os.Open(parameter_path)
	if err_p != nil {
		fmt.Println(err_p.Error())
	}
	buf_read_p := bufio.NewReader(f_parameter)
	for {
		line, err := buf_read_p.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\r\n")
		if len(line) != 0 {
			if line[0:1] != "#" {
				slice2 := strings.Split(line, "=")
				if slice2[0] == "input" {
					if query_temp_path == "null" {
						slice2[1] = strings.Replace(slice2[1], "\\", "/", -1)
						Param.Input = slice2[1]
					} else if query_temp_path != "null" {
						Param.Input = strings.Replace(query_temp_path, "\\", "/", -1)
					}
				} else if slice2[0] == "output" {
					slice2[1] = strings.Replace(slice2[1], "\\", "/", -1)
					Param.Output = slice2[1]
				} else if slice2[0] == "clear" {
					if slice2[1] == "false" || slice2[1] == "False" {
						Param.Clear = false
					} else if slice2[1] == "true" || slice2[1] == "True" {
						Param.Clear = true
					} else {
						Param.Clear = false
					}
				} else if slice2[0] == "merge_tolerance" {
					temp, _ := strconv.ParseFloat(slice2[1], 64)
					Param.Merge_tolerance = temp
				} else if slice2[0] == "merge_type" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Merge_type = temp
				} else if slice2[0] == "threshold_peaks" {
					temp, _ := strconv.ParseFloat(slice2[1], 64)
					Param.Threshold_peaks = temp
				} else if slice2[0] == "precur" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Precur = temp
				} else if slice2[0] == "tolerance_precur" {
					temp, _ := strconv.ParseFloat(slice2[1], 64)
					Param.Tolerance_precur = temp
				} else if slice2[0] == "adduct" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Isotope = temp
				} else if slice2[0] == "tolerance_isotope" {
					temp, _ := strconv.ParseFloat(slice2[1], 64)
					Param.Tolerance_isotopic = temp
				} else if slice2[0] == "decoy_similitude" {
					temp, _ := strconv.ParseFloat(slice2[1], 64)
					Param.Decoy_similitude = temp
				} else if slice2[0] == "threads" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Threads = temp
				} else if slice2[0] == "Min_mass" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Min_mass = temp
				} else if slice2[0] == "Max_mass" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Max_mass = temp
				} else if slice2[0] == "search_pattern" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Search_pattern = temp
				} else if slice2[0] == "decoy_pattern" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Decoy_pattern = temp
				} else if slice2[0] == "decoyinput" {
					slice2[1] = strings.Replace(slice2[1], "\\", "/", -1)
					Param.Decoyinput = slice2[1]
				} else if slice2[0] == "reference" {
					if reference_temp_path == "null" {
						slice2[1] = strings.Replace(slice2[1], "\\", "/", -1)
						Param.Reference = slice2[1]
					} else if reference_temp_path != "null" {
						Param.Reference = reference_temp_path
					}
				} else if slice2[0] == "electric_pattern" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Electric_pattern = temp
				} else if slice2[0] == "hplc_pattern" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.HPLC_pattern = temp
				} else if slice2[0] == "adduct_path" {
					slice2[1] = strings.Replace(slice2[1], "\\", "/", -1)
					Param.Adduct_path = slice2[1]
				} else if slice2[0] == "database_filter" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Database_filter = temp
				} else if slice2[0] == "match_model" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.Match_model = temp
				} else if slice2[0] == "mmi" {
					temp, _ := strconv.Atoi(slice2[1])
					Param.MMI = temp
				}
			}
		}
	}
	f_parameter.Close()
	return Param
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
