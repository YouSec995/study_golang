package cmd

import "github.com/YouSec995/study_golang/conf"

func GetAgtInfo(read conf.Reader) conf.AgtInfo {
	r := string(read.Msg)
	agtinfo := conf.AgtInfo{}
	// TODO 将r中数据解析到对应的AgtInfo中并返回

	return agtinfo
}
