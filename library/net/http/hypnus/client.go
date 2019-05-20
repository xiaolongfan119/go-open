package hypnus

import (
	"context"
	"encoding/json"
	xtime "go-open/library/time"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

/*
	知识补充点: TCP keapAlive
	#. tcp_keepalive_time  单位是秒，表示TCP链接在多少秒之后没有数据报文传输启动探测报文,  默认值为7200s（2小时）
	#. tcp_keepalive_intvl 单位是也秒,表示前一个探测报文和后一个探测报文之间的时间间隔, 默认值为75s
	#. tcp_keepalive_probes 表示探测的次数, 默认值为9（次）

*/

type ClientConfig struct {

	// Timeout is the maximum amount of time a dial will wait for a connect to complete
	TimeOut xtime.Duration // 限制创建连接的时间

	// KeepAlive specifies the keep-alive period for an active network connection
	KeepAlive xtime.Duration // 对应TCP的 tcp_keepalive_time

}

type Client struct {
	conf   *ClientConfig
	client *http.Client
	dialer *net.Dialer
}

func NewClient(c *ClientConfig) *Client {
	client := &Client{}
	client.setConfig(c)

	client.dialer = &net.Dialer{
		Timeout:   time.Duration(c.TimeOut),
		KeepAlive: time.Duration(c.KeepAlive),
	}

	transport := &http.Transport{DialContext: client.dialer.DialContext}
	client.client = &http.Client{Transport: transport}

	return client
}

func (client *Client) setConfig(c *ClientConfig) {
	if c.TimeOut > 0 {
		client.conf.TimeOut = c.TimeOut
	}
	if c.KeepAlive > 0 {
		client.conf.KeepAlive = c.KeepAlive
	}
}

func (client *Client) NewRequest(method, url string, params map[string]interface{}) (req *http.Request, err error) {

	str, _ := json.Marshal(params)
	req, err = http.NewRequest(method, url, strings.NewReader(string(str)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return
}

func (client *Client) Get(c context.Context, url string, resp interface{}) error {

	req, err := client.NewRequest(http.MethodGet, url, make(map[string]interface{}))
	if err != nil {
		return err
	}

	response, _ := client.client.Do(req)

	body, _ := ioutil.ReadAll(response.Body)

	return nil
}
