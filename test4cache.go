package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Customize struct {
	httpClient_ *http.Client
	username_   string
	passwd_     string
}

func doInit() *Customize {
	u := new(Customize)
	u.username_ = "yang"
	u.passwd_ = "yang@4258770"
	timeout := time.Duration(10 * time.Second)
	u.httpClient_ = &http.Client{
		Timeout: timeout,
	}
	return u
}

//substr like "2014-12-17"
func (u *Customize) doPut(suburi string) {

	realurl := "http://v0.api.upyun.com/test4cache/test/" + suburi + ".dat"

	request, _ := http.NewRequest("PUT", realurl, strings.NewReader("test"))
	request.SetBasicAuth(u.username_, u.passwd_)

	resp, err := u.httpClient_.Do(request)

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

	realurl := "http://test4cache.b0.upaiyun.com/test/" + suburi + ".dat"
	//fmt.Println(realurl)
	request, err := http.NewRequest("GET", realurl, nil)
	trans := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.Dial(netw, ipaddr+":80")
			return conn, err
		},
	}

	u.httpClient_.Transport = trans
	resp, err := u.httpClient_.Do(request)
	if err != nil {
		fmt.Printf("response docache %s:%s fail, err: %s\n", ipaddr, suburi, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		fmt.Printf("docache %s:%s fail, %v\n", ipaddr, suburi, resp.StatusCode)
	} else {
		fmt.Printf("docache %s:%s success\n", ipaddr, suburi)
	}
}

func (u *Customize) testCache(suburi string, ipaddr string) {

	realurl := "http://test4cache.b0.upaiyun.com/test/" + suburi + ".dat"
	//fmt.Println(realurl)
	request, err := http.NewRequest("GET", realurl, nil)
	trans := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.Dial(netw, ipaddr+":80")
			return conn, err
		},
	}

	u.httpClient_.Transport = trans
	resp, err := u.httpClient_.Do(request)
	if err != nil {
		fmt.Printf("response docache %s:%s fail, err: %s\n", ipaddr, suburi, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		fmt.Printf("testcache %s:%s fail, %v\n", ipaddr, suburi, resp.StatusCode)
	} else {
		fmt.Printf("testcache %s:%s success\n", ipaddr, suburi)
	}
}

var u = doInit()
var loc, _ = time.LoadLocation("Asia/Shanghai")
var iplist = map[string]string{}

func uploadTimer() {
	datestr := time.Now().Format("2006-01-02")
	if u != nil {
		u.doPut(datestr)
		for _, ip := range iplist {
			go u.doCache(datestr, ip)
		}
	}
}

func testTimer() {
	var datestr string
	for _, ip := range iplist {
		for i := 0; i < 7; i++ {
			datestr = time.Now().AddDate(0, 0, i*(-1)).Format("2006-01-02")
			if u != nil {
				go u.testCache(datestr, ip)
			}
		}
	}
}

func timerUpload() {
	nextdate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	nexttime, _ := time.ParseInLocation("2006-01-02", nextdate, loc)
	timer1 := time.NewTicker(nexttime.Sub(time.Now()))
	//timer1 := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timer1.C:
			uploadTimer()
			timer1 = time.NewTicker(24 * 3600 * time.Second)
		}
	}

}

func timerTest() {
	nextdate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	nexttime, _ := time.ParseInLocation("2006-01-02", nextdate, loc)
	timer2 := time.NewTicker(nexttime.Sub(time.Now()) + 23*3600*time.Second)
	//timer2 := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer2.C:
			testTimer()
			timer2 = time.NewTicker(24 * 3600 * time.Second)
		}
	}
}

func generateIplist(filename string) {
	f, err := os.Open(filename)
	defer f.Close()

	if nil == err {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" && line[0] != '#' {
				str_map := strings.Split(line, " ")
				if str_map != nil && len(str_map) == 2 {
					key := str_map[0]
					v_map := strings.Split(str_map[1], "=")
					if v_map != nil && len(v_map) == 2 {
						value := v_map[1]
						//fmt.Printf("key: %s, value: %s", key, value)
						iplist[key] = value
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("reading from %s, error: %s", filename, err)
		}

	}
}

func main() {
	generateIplist("iplist.txt")
	go timerUpload()
	timerTest()
}
