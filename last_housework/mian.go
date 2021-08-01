package main
import (
	"github.com/YouSec995/study_golang/cmd"
	"github.com/YouSec995/study_golang/job/kafka"
	"github.com/YouSec995/study_golang/conf"
)
func main(){
	read, myErr := kafka.MqGetAgtInfo()
	if myErr.ErrMsg != nil {
		// TODO 打印错误日志,并反回给消息队列kafka
		w := conf.Writer{}
		// TODO 将报错信息转换为Writer中去并推送到kafka中
		w.Msg = read.Msg
		kafka.MqPushAgtInfo(w)
		return
	}
	agtQueue := cmd.GetAgtInfo(read)
}


