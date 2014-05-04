package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var url = "http://int.dpool.sina.com.cn/iplookup/iplookup.php"

type IpInfo struct {
	Country  string `json:country`
	Province string `json:province`
	City     string `json:city`
	District string `json:district`
	Isp      string `json:isp`
}

func GetIp(ip string) (*IpInfo, error) {
	getUrl := url + "?format=json&ip=" + ip
	resp, err := http.Get(getUrl)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ipInfo := &IpInfo{}
	err = json.Unmarshal(data, ipInfo)
	if err != nil {
		return nil, err
	}
	return ipInfo, nil
}
