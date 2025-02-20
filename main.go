package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"unicode"
)

type Message struct {
	Message map[string]Field `json:"message"`
}

type Field struct {
	Type string `json:"type"`
}

type Service struct {
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
}

type Method struct {
	Name   string `json:"name"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type InterfaceDescription struct {
	Messages map[string]Message `json:"messages"`
	Service  Service            `json:"service"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <json_file>")
		return
	}

	jsonFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var InterfaceDescription InterfaceDescription
	json.Unmarshal(byteValue, &InterfaceDescription)
	stubCode := stubCodeGeneration(InterfaceDescription)
	logicCode := logicCodeGeneration(InterfaceDescription)
	mainCode := mainCodeGeneration(InterfaceDescription)

	//创建一个文件夹
	os.Mkdir("./"+InterfaceDescription.Service.Name, 0755)
	os.Mkdir("./"+InterfaceDescription.Service.Name+"/stub", 0755)

	mainFile, err := os.Create("./" + InterfaceDescription.Service.Name + "/main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mainFile.Close()
	mainFile.WriteString(mainCode)

	stubCodeFile, err := os.Create("./" + InterfaceDescription.Service.Name + "/stub/stubCode.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stubCodeFile.Close()
	stubCodeFile.WriteString(stubCode)

	logicFile, err := os.Create("./" + InterfaceDescription.Service.Name + "/stub/logicCode.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logicFile.Close()
	logicFile.WriteString(logicCode)

	err = configMod(InterfaceDescription)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done!")
}

func stubCodeGeneration(InterfaceDescription InterfaceDescription) (code string) {
	code = "package stub\n"
	code += "\n"
	code += "import (\n"
	code += "\t\"encoding/json\"\n"
	code += "\t\"github.com/JIAKUNHUANG/krpc/client\"\n"
	code += "\t\"github.com/JIAKUNHUANG/krpc/server\"\n"
	code += ")\n\n"

	// 生成结构体
	for name, message := range InterfaceDescription.Messages {
		code += "type " + name + " struct {\n"
		for fieldName, field := range message.Message {
			firstLetterUpperFieldName := capitalizeFirstLetter(fieldName)
			code += "\t" + firstLetterUpperFieldName + " " + field.Type + " `json:\"" + fieldName + "\"`\n"
		}
		code += "}\n\n"
	}

	// 生成服务注册函数
	code += "func RegisterServiceTestService(s *server.Service) {\n"
	for _, method := range InterfaceDescription.Service.Methods {
		code += "\ts.AddMethod(\"" + method.Name + "\", " + method.Name + "Func)\n"
	}
	code += "	err := s.RegisterService(\"127.0.0.1:8000\")\n"
	code += "	if err != nil {\n"
	code += "\t\tpanic(err)\n"
	code += "	}\n"
	code += "}\n\n"

	// 生成代理注册
	code += "type Proxy struct {\n"
	code += "\tclient *client.Client\n"
	code += "}\n\n"

	code += "func NewProxy() *Proxy {\n"
	code += "\tp := &Proxy{}\n"
	code += "\tp.client = client.NewClient()\n"
	code += "\treturn p\n"
	code += "}\n\n"

	code += "func (p *Proxy) RegisterProxy() error {\n"
	code += "\terr := p.client.ConnectService(\"127.0.0.1:8000\")\n"
	code += "\tif err != nil {\n"
	code += "\t\treturn err\n"
	code += "\t}\n"
	code += "\treturn nil\n"
	code += "}\n"

	// 生成服务调用代码
	for _, method := range InterfaceDescription.Service.Methods {
		code += "func (p *Proxy) " + method.Name + "(clientReq " + method.Input + ") (clientRsp " + method.Output + ", err error) {\n"
		code += "\treq := client.Request{\n"
		code += "\t\tMethod: \"" + method.Name + "\",\n"
		code += "\t\tParams: clientReq,\n"
		code += "\t}\n"

		code += "\trsp, err := p.client.Call(req)\n"
		code += "\tif err != nil {\n"
		code += "\t\treturn clientRsp, err\n"
		code += "\t}\n"

		code += "\trspResultBuf, _ := json.Marshal(rsp.Result)\n"
		code += "\terr = json.Unmarshal(rspResultBuf, &clientRsp)\n"
		code += "\tif err != nil {\n"
		code += "\t\treturn clientRsp, err\n"
		code += "\t}\n"

		code += "\treturn clientRsp, nil\n"
		code += "}\n"
		code += "\n"
	}
	return code
}

func logicCodeGeneration(InterfaceDescription InterfaceDescription) (code string) {
	code = "package stub\n"
	code += "\n"
	for _, method := range InterfaceDescription.Service.Methods {
		code += "var " + method.Name + "Func = func(input " + method.Input + ") (output " + method.Output + ") {\n"
		code += "\treturn\n"
		code += "}\n"
		code += "\n"
	}
	return
}

func mainCodeGeneration(InterfaceDescription InterfaceDescription) (code string) {
	code = "package main\n"
	code += "\n"
	code += "import (\n"
	code += "\t\"" + InterfaceDescription.Service.Name + "/stub\"\n"
	code += "\n"
	code += "\t\"github.com/JIAKUNHUANG/krpc/server\"\n"
	code += ")\n\n"

	code += "func main() {\n"
	code += "\ts := server.CreateService()\n"
	code += "\n"
	code += "\tstub.RegisterServiceTestService(s)\n"
	code += "\tdefer s.Listener.Close()\n"
	code += "\n"
	code += "\tgo s.Service()\n"
	code += "\n"
	code += "\tselect {}\n"
	code += "}\n"
	return
}

func configMod(InterfaceDescription InterfaceDescription) error {
	targetDir := "./" + InterfaceDescription.Service.Name + "/"
	bagDir := "github.com/JIAKUNHUANG/krpc"
	err := os.Chdir(targetDir)
	if err != nil {
		return err
	}
	
	cmd := exec.Command("go", "mod", "init", InterfaceDescription.Service.Name)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "get", bagDir)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil

}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}
