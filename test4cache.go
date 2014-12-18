package main

import (
	"fmt"
	//"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Customize struct {
	httpClient_  *http.Client
	httpRequest_ *http.Request
	username_    string
	passwd_      string
}

func doInit() *Customize {
	u := new(Customize)
	u.username_ = "yang"
	u.passwd_ = "yang@4258770"
	timeout := time.Duration(10 * time.Second)
	u.httpClient_ = &http.Client{
		Timeout: timeout,
	}
	var err error
	u.httpRequest_, err = http.NewRequest("PUT", "http://v0.api.upyun.com", strings.NewReader("test"))
	if err != nil {
		fmt.Printf("create http request fail, err: %s\n", err)
		return nil
	}
	u.httpRequest_.SetBasicAuth(u.username_, u.passwd_)
	return u
}

//substr like "2014-12-17"
func (u *Customize) doPut(suburi string) {

	realurl := "http://v0.api.upyun.com/test4cache/test/" + suburi + ".dat"
	u.httpRequest_.URL, _ = url.ParseRequestURI(realurl)

	u.httpRequest_.Header.Add("Content-Length", "4")
	resp, err := u.httpClient_.Do(u.httpRequest_)

	if err != nil {
		fmt.Println("response fail, err: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		fmt.Printf("put %s fail, %v\n", suburi, resp.StatusCode)
	} else {
		fmt.Printf("put %s success\n", suburi)
	}
}

func (u *Customize) doCache(suburi string, ipaddr string) {

	realurl := "http://124.115.20.67/test/" + suburi + ".dat"
	//fmt.Println(realurl)
	request, err := http.NewRequest("GET", realurl, nil)
	request.Header.Add("Host", "test4cache.b0.upaiyun.com")

	resp, err := u.httpClient_.Do(request)
	if err != nil {
		fmt.Printf("response docache fail, err: %s\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		fmt.Printf("docache %s fail, %v\n", suburi, resp.StatusCode)
	} else {
		fmt.Printf("docache %s success\n", suburi)
	}
}

func (u *Customize) testCache(suburi string, ipaddr string) {

	realurl := "http://v0.api.upyun.com/test4cache/test/" + suburi + ".dat"
	request, _ := http.NewRequest("DELETE", realurl, nil)
	request.SetBasicAuth(u.username_, u.passwd_)
	resp, err := u.httpClient_.Do(request)

	if err != nil {
		fmt.Printf("response delete fail, err: %s\n", err)
		return
	}

	if resp.StatusCode/100 != 2 {
		fmt.Printf("delete %s fail, %v\n", suburi, resp.StatusCode)
	} else {
		fmt.Printf("delete %s success\n", suburi)
	}

	realurl = "http://test4cache.b0.upaiyun.com/test/" + suburi + ".dat"
	request, _ = http.NewRequest("GET", realurl, nil)
	request.Header.Add("Host", "test4cache.b0.upaiyun.com")

	resp, err = u.httpClient_.Do(request)
	if err != nil {
		fmt.Printf("response testcache fail, err: %s\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		fmt.Printf("testcache %s fail, %v\n", suburi, resp.StatusCode)
	} else {
		fmt.Printf("testcache %s success\n", suburi)
	}
}

func main() {
	u := doInit()

	if u != nil {
		u.doPut("2014-12-16")
		u.doCache("2014-12-16", "124.115.20.67")
		u.testCache("2014-12-16", "124.115.20.67")
	}
}
