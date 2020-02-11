package browsermobproxy

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Params map[string]string
var ports = make(portMap,10)
type portMap map[int]bool
type Client struct {
	Host string 	`json:"host"`
	Port int		`json:"port"`
	Proxy string	`json:"proxy"`

}
func init(){
	for i:=1;i<=10; i++{
		ports[8080+i]=false
	}
}
func NewClient(urlStr string,param Params,options Params) *Client{

	host :="http://"+urlStr
	var port int
	q := new(url.URL).Query()
	if param!=nil {

		for key, val := range param {
			q.Add(key, val)
		}

	}
	if v,ok :=options["existing_proxy_port_to_use"]; ok{
		port,_=strconv.Atoi(v)
	}else{

		for k,v :=range ports{
			if v==false{
				port=k
				ports[k]=true
				q.Add("port",strconv.Itoa(port))
				break
			}
		}
		resp,err :=http.PostForm(fmt.Sprintf("%s/proxy",host),q)
		if err != nil {
			fmt.Printf("Fail to connect, %s\n", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		res := make(map[string]interface{})
		_ =json.Unmarshal(body,&res)
		port = int(res["port"].(float64))
	}
	client :=&Client{Host:host,Port:port}
	urlParts :=strings.Split(host,":")
	client.Proxy=urlParts[1][2:]+":"+strconv.Itoa(port)

	return client
}
//关闭客户端
func (c *Client)Clost() int{
	req,_ :=http.NewRequest("DELETE",fmt.Sprintf("%s/proxy/%d",c.Host,c.Port),nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	//释放port
	ports[c.Port]=false
	return resp.StatusCode
}
//get the ports
func (c *Client)proxyPorts() *simplejson.Json{
	r,err :=http.Get(fmt.Sprintf("%s/proxy",c.Host))
	if err!=nil{
		fmt.Printf("Fail to get proxy, %s\n", err)
	}
	body, err := ioutil.ReadAll(r.Body)
	dataJson, _ := simplejson.NewJson(body)
	return dataJson
}

func(c *Client)Har() *simplejson.Json{
	res,err :=http.Get(fmt.Sprintf("%s/proxy/%d/har",c.Host,c.Port))
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	dataJson, _ := simplejson.NewJson(body)
	return dataJson
}

func (c *Client)NewHar(ref string,options map[string]string){
	options["initialPageRef"] = ref
	req,_ :=http.NewRequest("PUT",fmt.Sprintf("%s/proxy/%d/har",c.Host,c.Port),nil)
	q := req.URL.Query()
	if options != nil {
		for key, val := range options {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Fail to set new har, %s\n", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	res := make(map[string]interface{})
	_ =json.Unmarshal(body,&res)

}
//加入黑名单
func(c *Client)Blacklist(regexp string, statusCode int){
	req,_ :=http.NewRequest("PUT",fmt.Sprintf("%s/proxy/%d/blacklist",c.Host,c.Port),nil)
	q := req.URL.Query()
	q.Add("regex", regexp)
	q.Add("status", strconv.Itoa(statusCode))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	status := resp.StatusCode
	fmt.Printf("blacklist statusCode, %d\n", status)
}
//拦截请求内容
func (c *Client) ResponseInterceptor(js string) int{
	resp,err :=http.Post(fmt.Sprintf("%s/proxy/%d/filter/response",c.Host,c.Port),"text/plain",strings.NewReader(js))
	if err!=nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
//请求拦截
func (c *Client) RequestInterceptor(js string) int{
	resp,err :=http.Post(fmt.Sprintf("%s/proxy/%d/filter/request",c.Host,c.Port),"text/plain",strings.NewReader(js))
	if err!=nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}