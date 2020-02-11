# browsermobproxy-go
browsermobproxy for go

How too use it?
## Create new server?
```
    //the path browsermob-proxy install
    Server =browsermobproxy.NewServer("/Users//bin/browsermob-proxy")
    Server.Start()
```

## create proxy 
```
proxy :=Server.CreateProxy(browsermobproxy.Params{"trustAllServers":"true"})
```
## 
```go
chromeCaps := chrome.Capabilities{
		Prefs:            imgCaps,
		Path:             "",
		Args: []string{
			//"--headless",
			"--start-maximized",
			"--window-size=1200x600",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
		},
	}

```
## add it to chrome
```go
var caps selenium.Capabilities
chromeCaps.Args = append(spider.ChromeCaps.Args, fmt.Sprintf("--proxy-server=%s",p.Proxy))
caps.AddChrome(chromeCaps)
wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
```

## use newhar
```go
options :=make(map[string]string)
	options["captureHeaders"] = "true"
	options["captureContent"] = "true"
	proxy.NewHar("loginform",options)
	wd.get("url")
	result :=proxy.Har()
```

v1.0 you can use
1、proxy
2、Har
3、NewHar
4、Blacklist
5、ResponseInterceptor
6、RequestInterceptor