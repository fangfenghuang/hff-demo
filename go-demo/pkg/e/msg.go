package e

var MsgFlags = map[int]string{
	SUCCESS: "ok",
	ERROR:   "失败",
}

var EMsgFlags = map[int]string{
	SUCCESS: "Ok",
	ERROR:   "Fail",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

func GetEMsg(code int) string {
	msg, ok := EMsgFlags[code]
	if ok {
		return msg
	}

	return EMsgFlags[ERROR]
}
