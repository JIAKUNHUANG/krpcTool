package stub

var ServiceMap = make(map[string]string)

var ConnectFunc = func(input FindingRequest) (output FindingResponse) {
	if input.ReqType == "connect" {
		ServiceMap[input.ServiceName] = input.Addr
		output.ServiceName = input.ServiceName
		output.Status = "ok"
		output.ErrMsg = ""
		return
	}
	if input.ReqType == "disconnect" {
		delete(ServiceMap, input.ServiceName)
		output.ServiceName = input.ServiceName
		output.Status = "ok"
		output.ErrMsg = ""
		return
	}
	if input.ReqType == "finding" {
		if _, ok := ServiceMap[input.ServiceName]; ok {
			output.ServiceName = input.ServiceName
			output.Addr = ServiceMap[input.ServiceName]
			output.Status = "ok"
			output.ErrMsg = ""
			return
		}
		output.ServiceName = input.ServiceName
		output.Status = "fail"
		output.ErrMsg = "service not found"
		return
	}
	output.ServiceName = input.ServiceName
	output.Status = "ReqType not found"
	return
}
