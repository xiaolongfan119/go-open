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

	ecode "go-open/library/ecode"

	breaker "go-open/library/net/netutil/breaker"
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

	bkConfig *breaker.BreakerConfig
}

type Client struct {
	conf    *ClientConfig
	client  *http.Client
	dialer  *net.Dialer
	bkgroup *breaker.BreakerGroup
}

func NewClient(c *ClientConfig) *Client {
	client := &Client{conf: c}
	client.setConfig(c)

	client.dialer = &net.Dialer{
		Timeout:   time.Duration(c.TimeOut),
		KeepAlive: time.Duration(c.KeepAlive),
	}

	transport := &http.Transport{DialContext: client.dialer.DialContext}
	client.client = &http.Client{Transport: transport}

	client.bkgroup = breaker.NewBreakerGroup(c.bkConfig)

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

func (client *Client) Do(c context.Context, req *http.Request, resp interface{}) error {

	bs, err := client.Raw(c, req)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, resp); err != nil {
		return err
	}

	return nil
}

func (client *Client) Raw(c context.Context, req *http.Request) (bs []byte, err error) {

	var (
		timeout time.Duration
		cancel  func()
	)

	timeout = time.Duration(client.conf.TimeOut)
	if deadline, ok := c.Deadline(); ok {
		if _timeout := time.Until(deadline); _timeout < timeout {
			timeout = _timeout
			c, cancel = context.WithTimeout(c, timeout)
			defer cancel()
		}
	}

	req = req.WithContext(c)

	bk := client.bkgroup.GetBreaker(req.URL.Path)
	_resp, _err := bk.Execute(func() (*http.Response, error) {
		resp, err := client.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return nil, ecode.BadRequest
		}

		return resp, nil
	})

	if _err != nil {
		return nil, _err
	}

	return ioutil.ReadAll(_resp.Body)
}

func (client *Client) Get(c context.Context, url string, resp interface{}) error {

	req, err := client.NewRequest(http.MethodGet, url, make(map[string]interface{}))
	if err != nil {
		return err
	}

	return client.Do(c, req, resp)
}
