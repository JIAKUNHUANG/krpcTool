package stub

import (
	"encoding/json"
	"github.com/JIAKUNHUANG/krpc/client"
	"github.com/JIAKUNHUANG/krpc/server"
)

type FindingRequest struct {
	ReqType string `json:"reqType"`
	Addr string `json:"addr"`
	ServiceName string `json:"serviceName"`
}

type FindingResponse struct {
	ErrMsg string `json:"errMsg"`
	Status string `json:"status"`
	Addr string `json:"addr"`
	ServiceName string `json:"serviceName"`
}

func RegisterServiceTestService(s *server.Service) {
	s.ServiceName = "ServiceFinding"
	s.AddMethod("Connect", ConnectFunc)
	addr := server.GetAddr(s.Config.ServiceAddr.Ip, s.Config.ServiceAddr.Port)
	err := s.RegisterService(addr)
	if err != nil {
		panic(err)
	}
}

type Proxy struct {
	client *client.Client
}

func NewProxy() *Proxy {
	p := &Proxy{}
	p.client = client.NewClient()
	return p
}

func (p *Proxy) RegisterProxy() error {
	serviceAddr, err := client.ConnectServiceFinding("127.0.0.1:8000", "ServiceFinding")
	if err != nil {
		return err
	}

	err = p.client.ConnectService(serviceAddr)
	if err != nil {
		return err
	}
	return nil
}
func (p *Proxy) Connect(clientReq FindingRequest) (clientRsp FindingResponse, err error) {
	req := client.Request{
		Method: "Connect",
		Params: clientReq,
	}
	rsp, err := p.client.Call(req)
	if err != nil {
		return clientRsp, err
	}
	rspResultBuf, _ := json.Marshal(rsp.Result)
	err = json.Unmarshal(rspResultBuf, &clientRsp)
	if err != nil {
		return clientRsp, err
	}
	return clientRsp, nil
}

