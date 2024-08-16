package unisms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const DEF_ENDPOINT = "https://uni.apistd.com"
const DEF_SIGNING_ALGORITHM = "hmac-sha256"
const REQUEST_ID_HEADER_KEY = "x-uni-request-id"
const VERSION = "0.0.2"

type UniSMS struct {
	AccessKey        string
	AccessSecret     string
	Endpoint         string
	SigningAlgorithm string
	Timeout          uint32
	Data             clJson.M
}

func New(accesskey, Secret string, timeout uint32) *UniSMS {
	return &UniSMS{
		Endpoint:         DEF_ENDPOINT,
		SigningAlgorithm: DEF_SIGNING_ALGORITHM,
		Timeout:          timeout,
		Data:             clJson.M{},
		AccessKey:        accesskey,
		AccessSecret:     Secret,
	}
}

func (m *UniSMS) SetTo(phoneNumbers ...string) *UniSMS {
	m.Data["to"] = phoneNumbers
	return m
}

func (m *UniSMS) SetSignature(signature string) *UniSMS {
	m.Data["signature"] = signature
	return m
}

func (m *UniSMS) SetTemplateId(templateId string) *UniSMS {
	m.Data["templateId"] = templateId
	return m
}

func (m *UniSMS) SetTemplateData(templateData map[string]string) *UniSMS {
	m.Data["templateData"] = templateData
	return m
}

func (m *UniSMS) SetContent(content string) *UniSMS {
	m.Data["content"] = content
	return m
}

func (m *UniSMS) SetText(text string) *UniSMS {
	m.Data["text"] = text
	return m
}

func (c *UniSMS) GenerateRandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
func (c *UniSMS) Sign(query url.Values) url.Values {
	if c.AccessKey != "" {
		query.Add("algorithm", c.SigningAlgorithm)
		query.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
		query.Add("nonce", c.GenerateRandomString(8))

		message := query.Encode()
		mac := hmac.New(sha256.New, []byte(c.AccessSecret))
		mac.Write([]byte(message))
		query.Add("signature", hex.EncodeToString(mac.Sum(nil)))
	}

	return query
}
func (m *UniSMS) Send(phone ...string) error {
	u := m.Endpoint
	m.SetTo(phone...)
	query := url.Values{}
	query.Add("action", "sms.message.send")
	query.Add("accessKeyId", m.AccessKey)
	query = m.Sign(query)
	querystr := query.Encode()
	jsonbytes, err := json.Marshal(m.Data)
	if err != nil {
		clLog.Error("短信发送失败错误:%v", err)
		return err
	}

	reader := bytes.NewReader(jsonbytes)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		//Timeout: time.Second * time.Duration(m.Timeout),
	}
	if m.Timeout > 0 {
		client.Timeout = time.Second * time.Duration(m.Timeout)
	}
	req, err := http.NewRequest("POST", u+"/?"+querystr, reader)

	if err != nil {
		clLog.Error("短信发送失败错误:%v", err)
		return err
	}

	req.Header.Set("User-Agent", "uni-go-sdk"+"/"+VERSION)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)

	if err != nil {
		clLog.Error("短信发送失败错误:%v", err)
		return err
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		clLog.Error("短信发送失败错误:%v", err)
		return err
	}
	resp := UniSMSResp{}
	json.Unmarshal(content, &resp)
	if resp.Code != "0" {
		requestId := res.Header.Get(REQUEST_ID_HEADER_KEY)
		clLog.Error("短信发送失败ID:[%v] MSG:[%v]", requestId, string(content))
		return fmt.Errorf(resp.Message)
	}
	fmt.Println(string(content))
	return nil
}

type UniSMSResp struct {
	Data struct {
		Currency      string `json:"currency"`
		Recipients    int    `json:"recipients"`
		MessageCount  int    `json:"messageCount"`
		TotalAmount   string `json:"totalAmount"`
		PayAmount     string `json:"payAmount"`
		VirtualAmount string `json:"virtualAmount"`
		Messages      []struct {
			Id           string `json:"id"`
			To           string `json:"to"`
			RegionCode   string `json:"regionCode"`
			CountryCode  string `json:"countryCode"`
			MessageCount int    `json:"messageCount"`
			Status       string `json:"status"`
			Upstream     string `json:"upstream"`
			Price        string `json:"price"`
		} `json:"messages"`
	} `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

//func NewClient(AccessKeyId, AccessKeySecret string) *UniClient {
//	return &UniClient{
//		AccessKeyId:      AccessKeyId,
//		AccessKeySecret:  AccessKeySecret,
//		Endpoint:         DEF_ENDPOINT,
//		SigningAlgorithm: DEF_SIGNING_ALGORITHM,
//	}
//}
//
//func NewResponse(res *http.Response) (*UniResponse, error) {
//	var data map[string]interface{}
//	var code string
//	var message string
//
//	status := res.StatusCode
//	requestId := res.Header.Get(REQUEST_ID_HEADER_KEY)
//	rawBody, err := ioutil.ReadAll(res.Body)
//
//	if err != nil {
//		return nil, err
//	}
//
//	if rawBody != nil {
//		body := make(map[string]interface{})
//		err := json.Unmarshal(rawBody, &body) //第二个参数要地址传递
//
//		if err != nil {
//			return nil, err
//		}
//
//		code = body["code"].(string)
//		message = body["message"].(string)
//
//		if code != "0" {
//			return nil, errors.New(fmt.Sprintf("[%s] %s, RequestId: %s", code, message, requestId))
//		}
//		data = body["data"].(map[string]interface{})
//	}
//
//	return &UniResponse{
//		Raw:       res,
//		Status:    status,
//		Code:      code,
//		Message:   message,
//		Data:      data,
//		RequestId: requestId,
//	}, nil
//}
