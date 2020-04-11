//TDA method module
//This module is used to caculate the the FDR of annotation result after database search
package ranker

import (
	"dataframe"
)

func TDA(Compounds int, Match_result_list []dataframe.Grade, input_param dataframe.Parameters) []dataframe.Grade {
	if input_param.Match_model == 1 {
		QuickSort(Match_result_list, 0, len(Match_result_list)-1)
	} else if input_param.Match_model == 2 {
		QuickSort_C(Match_result_list, 0, len(Match_result_list)-1)
	}
	Qvalue(Compounds, Match_result_list)
	return Match_result_list
}

func BubbleBestSort(a []dataframe.Grade) []dataframe.Grade {
	lastSwap := len(a) - 1
	lastSwapTemp := len(a) - 1
	for j := 0; j < len(a)-1; j++ {
		lastSwap = lastSwapTemp
		for i := 0; i < lastSwap; i++ {
			if a[i].Score > a[i+1].Score {
				a[i], a[i+1] = a[i+1], a[i]
				lastSwapTemp = i
			}
		}
		if lastSwap == lastSwapTemp {
			break
		}
	}
	return a
}

func BubbleBestSort_C(a []dataframe.Grade) []dataframe.Grade {
	lastSwap := len(a) - 1
	lastSwapTemp := len(a) - 1
	for j := 0; j < len(a)-1; j++ {
		lastSwap = lastSwapTemp
		for i := 0; i < lastSwap; i++ {
			if a[i].Cosine_similar > a[i+1].Cosine_similar {
				a[i], a[i+1] = a[i+1], a[i]
				lastSwapTemp = i
			}
		}
		if lastSwap == lastSwapTemp {
			break
		}
	}
	return a
}

func Qvalue(Compounds int, Match_result_list []dataframe.Grade) []dataframe.Grade {
	Decoy_label := Compounds / 2
	var Target, Decoy int
	for i := len(Match_result_list) - 1; i >= 0; i-- {
		if Match_result_list[i].Index >= Decoy_label {
			Decoy++
		} else {
			Target++
		}
		if 2*float64(Decoy)/float64(Decoy+Target) >= 1.0 {
			Match_result_list[i].FDR = 1.0
		} else {
			Match_result_list[i].FDR = 2 * float64(Decoy) / float64(Decoy+Target)
		}
	}
	return Match_result_list
}

func QuickSort(arr []dataframe.Grade, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2].Score
		for i <= j {
			for arr[i].Score < key {
				i++
			}
			for arr[j].Score > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			QuickSort(arr, start, j)
		}
		if end > i {
			QuickSort(arr, i, end)
		}
	}
	return
}

func QuickSort_C(arr []dataframe.Grade, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2].Cosine_similar
		for i <= j {
			for arr[i].Cosine_similar < key {
				i++
			}
			for arr[j].Cosine_similar > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			QuickSort_C(arr, start, j)
		}
		if end > i {
			QuickSort_C(arr, i, end)
		}
	}
	return
}

func STDA(Compounds int, Match_result_list_target, Match_result_list_decoy []dataframe.Grade, input_param dataframe.Parameters) []dataframe.Grade { //分离搜索FDR评估
	if input_param.Match_model == 1 {
		QuickSort(Match_result_list_target, 0, len(Match_result_list_target)-1)
		QuickSort(Match_result_list_decoy, 0, len(Match_result_list_decoy)-1)
		SQvalue(Compounds, Match_result_list_target, Match_result_list_decoy)
	} else if input_param.Match_model == 2 {
		QuickSort_C(Match_result_list_target, 0, len(Match_result_list_target)-1)
		QuickSort_C(Match_result_list_decoy, 0, len(Match_result_list_decoy)-1)
		SQvalue_C(Compounds, Match_result_list_target, Match_result_list_decoy)
	}

	return Match_result_list_target
}

func SQvalue(Compounds int, Match_result_list_target, Match_result_list_decoy []dataframe.Grade) []dataframe.Grade {
	var Target, Decoy int
	decoy_index := len(Match_result_list_decoy) - 1
	for i := len(Match_result_list_target) - 1; i >= 0; i-- {
		if Match_result_list_target[i].Score >= Match_result_list_decoy[decoy_index].Score {
			Target++
			Match_result_list_target[i].FDR = float64(Decoy) / float64(Decoy+Target)
		} else {
			Decoy++
			Match_result_list_target[i].FDR = float64(Decoy) / float64(Decoy+Target)
			Match_result_list_target[i].Index = Compounds + Match_result_list_decoy[decoy_index].Index
			decoy_index--
		}
	}
	return Match_result_list_target
}

func SQvalue_C(Compounds int, Match_result_list_target, Match_result_list_decoy []dataframe.Grade) []dataframe.Grade {
	var Target, Decoy int
	decoy_index := len(Match_result_list_decoy) - 1
	for i := len(Match_result_list_target) - 1; i >= 0; i-- {
		if Match_result_list_target[i].Cosine_similar >= Match_result_list_decoy[decoy_index].Cosine_similar {
			Target++
			Match_result_list_target[i].FDR = float64(Decoy) / float64(Decoy+Target)
		} else {
			Decoy++
			Match_result_list_target[i].FDR = float64(Decoy) / float64(Decoy+Target)
			Match_result_list_target[i].Index = Compounds + Match_result_list_decoy[decoy_index].Index
			decoy_index--
		}
	}
	return Match_result_list_target
}
