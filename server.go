package browsermobproxy

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)
const (
	_LOG_FILE = "log/browsermob-proxy.out.log"
)
type Server struct {
	Path string `json:"path"`
	Host string `json:"host"`
	Port int	`json:"port"`
	Process *os.Process `json:"process"`
	Command string `json:"command"`
	Url string `json:"url"`
}
//initialises a server object
func NewServer(path string) *Server{
	server := new(Server)
	if runtime.GOOS=="windows" {
		if !strings.HasSuffix(path,".bat"){
			path +=".bat"
		}
	}
	server.Path=path
	server.Host="localhost"
	server.Port=8080
	server.Url = fmt.Sprintf("http://%s:%d",server.Host,server.Port)
	return server
}
//启动
func(s *Server) Start(){
	runPid ,_ :=CheckPidRunning("java")
	if runPid!="" {
		fmt.Println("process isexist now kill it...")
		pid,_ := strconv.Atoi(runPid)
		runProcess,_ :=os.FindProcess(pid)
		_ =runProcess.Kill()
	}
	stdOut, _ := os.OpenFile(_LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	//r, w, err := os.Pipe()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer r.Close()
	procAttr := &os.ProcAttr{
		Files: []*os.File{nil, stdOut, stdOut},
	}
	process, err := os.StartProcess(s.Path, nil, procAttr)  //运行脚本

	if err != nil{
		fmt.Println("look path error:", err)
		os.Exit(1)
	}

	//buf := make([]byte,1024)
	//for{
	//	n,err := r.Read(buf)
	//	if err != nil && err != io.EOF{panic(err)}
	//	if 0 ==n {break}
	//	fmt.Println(string(buf[:n]))
	//}
	//w.Close()
	s.Process = process
	time.Sleep(2 * time.Second)

}
//根据进程名判断进程是否运行
func CheckPidRunning(serverName string) (string, error) {
	a := `lsof -i:8080|sed -n '2p'|awk '{print $2}'`
	result, err := exec.Command("/bin/sh", "-c", a).Output()
	pid :=""
	if err != nil {
		return pid, err
	}
	pid =strings.TrimSpace(string(result))
	return pid, nil
}
//判断是否启动
func(s *Server)isListen() bool{
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d",s.Port))
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return false
	}
	defer conn.Close()
	return true
}
//停止
func(s *Server) Stop(){
	if s.Process.Pid==0{
		return
	}
	_ =s.Process.Kill()
	processStatus,_ :=s.Process.Wait()
	if processStatus.Exited(){
		return
	}
}
//创建代理
func(s *Server)CreateProxy(param Params) *Client{
	client :=NewClient(s.Url[7:],param,nil)
	return client
}
