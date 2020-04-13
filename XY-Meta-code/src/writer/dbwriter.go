//The data output module
//This modele is used to output calculated data from Xy-Meta such as Target-Decoy database and annotation result

package writer

import (
	"dataframe"
	"fmt"
	"math"
	"os"
	"reader"
	"strconv"
)

func DBwriter(path_database string, spectral_database []reader.Spectrum) {
	fmt.Println("Writting database!")
	f_db_writer, err := os.OpenFile(path_database, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := 0; i < len(spectral_database); i++ {
		precursor_mass := strconv.FormatFloat(spectral_database[i].Precusor_mass, 'f', -1, 64)
		charge := strconv.FormatFloat(spectral_database[i].Charge, 'f', -1, 64) + "+"
		retiontime := strconv.FormatFloat(spectral_database[i].Retention_time, 'f', -1, 64)
		f_db_writer.WriteString("BEGIN IONS\r\n")
		f_db_writer.WriteString("TITLE=" + spectral_database[i].Spectrl_number + "\r\n")
		f_db_writer.WriteString("PEPMASS=" + precursor_mass + "\r\n")
		f_db_writer.WriteString("CHARGE=" + charge + "\r\n")
		f_db_writer.WriteString("RTINSECONDS=" + retiontime + "\r\n")
		var Signal_content string
		for j := 0; j < len(spectral_database[i].Peaks); j++ {
			mz_content := strconv.FormatFloat(spectral_database[i].Peaks[j].Peak_mass, 'f', -1, 64)
			intensity_content := strconv.FormatFloat(spectral_database[i].Peaks[j].Peak_intensity, 'f', -1, 64)
			Signal_content = Signal_content + mz_content + " " + intensity_content + "\r\n"
		}
		Signal_content = Signal_content + "END IONS\r\n"
		f_db_writer.WriteString(Signal_content)
		//}
	}
	f_db_writer.Close()
}

func Mappingwriter(Match_result_list []dataframe.Grade, writepath string, Query, Database []reader.Spectrum, Adduct_isotope_list dataframe.Adduct_isotope) {
	f_o, _ := os.OpenFile(writepath, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	f_o.WriteString("ID" + "\t" + "Score" + "\t" + "FDR" + "\t" + "Reference_spectrum" + "\t" + "Match_Score" + "\t" + "Match_Cosine" + "\t" + "TSNR" + "\t" + "ESNR" + "\t" + "Query_precursor_retention_time" + "\t" + "Query_precursor_mass" + "\t" + "Reference_precursor_mass" + "\t" + "Diviation_mass" + "\t" + "Adduct" + "\t" + "Adduct_mass" + "\t" + "Query_peaks number" + "\t" + "Reference_peaks number" + "\t" + "Match_label" + "\r\n")
	for i := len(Match_result_list) - 1; i >= 0; i-- {
		Score := strconv.FormatFloat(Match_result_list[i].Score, 'f', -1, 64)
		FDR := strconv.FormatFloat(Match_result_list[i].FDR, 'f', -1, 64)
		Potproduct := strconv.FormatFloat(Match_result_list[i].Dot_product, 'f', -1, 64)
		TSNR := strconv.FormatFloat(Match_result_list[i].TSNR, 'f', -1, 64)
		ESNR := strconv.FormatFloat(Match_result_list[i].ESNR, 'f', -1, 64)
		Query_precursor_mass := strconv.FormatFloat(Query[Match_result_list[i].Queryindex].Precusor_mass, 'f', -1, 64)
		Reference_precursor_mass := strconv.FormatFloat(Database[Match_result_list[i].Index].Precusor_mass, 'f', -1, 64)
		Diviation_precursor_mass := strconv.FormatFloat(Query[Match_result_list[i].Queryindex].Precusor_mass-Database[Match_result_list[i].Index].Precusor_mass, 'f', -1, 64)
		Cosine := strconv.FormatFloat(Match_result_list[i].Cosine_similar, 'f', -1, 64)
		Rtime := strconv.FormatFloat(Query[Match_result_list[i].Queryindex].Retention_time, 'f', -1, 64)
		Adduct := ""
		Adduct_mass := ""
		Query_peaks_number := strconv.Itoa(len(Query[Match_result_list[i].Queryindex].Peaks))
		Reference_peaks_number := strconv.Itoa(len(Database[Match_result_list[i].Index].Peaks))
		lab := "NA"
		if Match_result_list[i].Score != 0.0 {
			lab = Match_result_list[i].Match_label
		}
		for j := 0; j < len(Adduct_isotope_list.Isotope_type_list); j++ {
			if math.Abs(Adduct_isotope_list.Isotope_mass_list[j]-(Query[Match_result_list[i].Queryindex].Precusor_mass-Database[Match_result_list[i].Index].Precusor_mass)) <= 0.1 {
				Adduct = Adduct_isotope_list.Isotope_type_list[j]
				Adduct_mass = strconv.FormatFloat(Adduct_isotope_list.Isotope_mass_list[j], 'f', -1, 64)
				break
			}
		}
		f_o.WriteString(Query[Match_result_list[i].Queryindex].Spectrl_number + "\t" + Score + "\t" + FDR + "\t" + Database[Match_result_list[i].Index].Spectrl_number + "\t" + Potproduct + "\t" + Cosine + "\t" + TSNR + "\t" + ESNR + "\t" + Rtime + "\t" + Query_precursor_mass + "\t" + Reference_precursor_mass + "\t" + Diviation_precursor_mass + "\t" + Adduct + "\t" + Adduct_mass + "\t" + Query_peaks_number + "\t" + Reference_peaks_number + "\t" + lab + "\r\n")
	}
}
