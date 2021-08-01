package conf

type MyError struct {
	ErrMsg	error
	ErrCode	int64
}

type Reader struct {
	Msg []byte
	Index int64
}

type Writer struct {
	Msg []byte
	Index int64
}

type AgtInfo struct {
	Type string
	AgtMsg AgentMsg
}

type AgentMsg struct {
	cnct uint8
	uuid int64
	phone string
	skills skill
	mainStatus int32
	subStatus int32
	sessionId string
	assData string
	// 等等
}

type skill struct {
	level int32
	msg string
	// 等等
}
