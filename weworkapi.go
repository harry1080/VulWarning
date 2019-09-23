package main

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

const (
	atPath    = "/tmp/wework_access_token"
	weworkurl = "https://qyapi.weixin.qq.com/cgi-bin/"
	// menuData  = `{"button":[{"type":"click","name":"最新预警","key":"news"},{"type":"click","name":"接收推送","key":"add"},{"type":"click","name":"取消推送","key":"delete"}]}`
	menuData = `{"button":[{"type":"click","name":"最新预警","key":"news"}]}`
)

// MsgContent 微信消息体
type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
}

// AccessToken AccessToken
type AccessToken struct {
	Time  int64  `json:"time"`
	Token string `json:"token"`
}

// RespMsg access_token resp msg
type RespMsg struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires_in"`
}

// TextContent 文本消息内容
type TextContent struct {
	Content string `json:"content"`
}

// TextMessage 文本消息
type TextMessage struct {
	ToUser  string      `json:"touser"`
	ToParty string      `json:"toparty"`
	MsgType string      `json:"msgtype"`
	Agentid uint32      `json:"agentid"`
	Text    TextContent `json:"text"`
}

var (
	wxcpt *wxbizmsgcrypt.WXBizMsgCrypt
	at    *AccessToken
)

func init() {
	wxcpt = wxbizmsgcrypt.NewWXBizMsgCrypt(Token, EncodingAESKey, CorpID, wxbizmsgcrypt.XmlType)
	at = &AccessToken{
		Time:  0,
		Token: "",
	}
	// Read access token
	data, err := ioutil.ReadFile(atPath)
	if err == nil {
		json.Unmarshal(data, &at)
	}
	// getAccessToken()
}

func httpReq(url, method string, data map[string]string) ([]byte, error) {
	var err error
	var req *http.Request
	body := []byte("")
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	query := req.URL.Query()
	if at.Time > 0 {
		query.Add("access_token", at.Token)
		query.Add("agentid", string(AgentID))
	}
	if method == "POST" {
		var dataString string
		if stream, ok := data["_stream"]; ok {
			dataString = stream
		} else {
			q := req.URL.Query()
			for k, v := range data {
				q.Add(k, v)
			}
			dataString = q.Encode()
		}
		req, err = http.NewRequest("POST", url, strings.NewReader(dataString))
	} else {
		req, err = http.NewRequest("GET", url, nil)
		for k, v := range data {
			query.Add(k, v)
		}
	}
	req.URL.RawQuery = query.Encode()
	if err != nil {
		return body, err
	}
	resp, err := client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	return body, nil
}

// SendMessage 发送应用消息
func SendMessage(content string) error {
	err := getAccessToken()
	if err != nil {
		log.Println("[-] Access Token Error")
		return err
	}
	// _content :=
	textMsg := &TextMessage{
		MsgType: "text",
		Text: TextContent{
			Content: content,
		},
		Agentid: AgentID,
		ToUser:  "@all",
	}
	stream, err := json.Marshal(textMsg)
	if err != nil {
		log.Println("[-] TextMessage Marshal Error")
		return err
	}
	data := make(map[string]string, 1)
	data["_stream"] = string(stream)
	_, err = httpReq(weworkurl+"message/send", "POST", data)
	return nil
}

// CreateMenu 创建自定义菜单
func CreateMenu() error {
	err := getAccessToken()
	if err != nil {
		log.Println("[-] Access Token Error")
		return err
	}
	data := make(map[string]string, 1)
	data["_stream"] = menuData
	_, err = httpReq(weworkurl+"menu/create", "GET", data)
	if err != nil {
		log.Println("[-] Create Menu Error : ", err.Error())
		return err
	}
	return nil
}

func httpGetAccessToken() error {
	data := make(map[string]string, 2)
	data["corpid"] = CorpID
	data["corpsecret"] = Secret
	body, err := httpReq(weworkurl+"gettoken", "GET", data)
	if err != nil {
		at.Time = 0
		return err
	}
	var respMsg RespMsg
	json.Unmarshal(body, &respMsg)
	if respMsg.ErrCode != 0 {
		at.Time = 0
		return errors.New(respMsg.ErrMsg)
	}
	at.Token = respMsg.AccessToken
	at.Time = time.Now().Unix() + respMsg.Expires - 10
	return nil
}

func getAccessToken() error {
	if at.Time < time.Now().Unix() {
		err := httpGetAccessToken()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// Marshal access token
		sData, err := json.Marshal(at)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		// Save access token
		err = ioutil.WriteFile(atPath, sData, 0644)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

// verifyURL 验证回调URL
func verifyURL(c *gin.Context) {
	msgSignature := c.Query("msg_signature")
	timeStamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echoStr := c.Query("echostr")
	sEchoStr, cryptErr := wxcpt.VerifyURL(msgSignature, timeStamp, nonce, echoStr)
	if nil != cryptErr {
		log.Println("[-] VerifyURL fail", cryptErr)
	}
	log.Println("[+] VerifyURL success echoStr", string(sEchoStr))
	c.String(http.StatusOK, string(sEchoStr))
}

// receiveMsg 对用户回复的消息解密
func receiveMsg(c *gin.Context) {
	msgSignature := c.Query("msg_signature")
	timeStamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	reqData := []byte("")
	msg, cryptErr := wxcpt.DecryptMsg(msgSignature, timeStamp, nonce, reqData)
	if nil != cryptErr {
		log.Println("[-] Decrypt msg fail", cryptErr)
		return
	}
	log.Println("[+] Decrypt msg success: ", string(msg))
	var msgContent MsgContent
	err := xml.Unmarshal(msg, &msgContent)
	if nil != err {
		log.Println("[-] Unmarshal fail")
		return
	}
	respContent := ""
	// Action
	if msgContent.MsgType == "event" && msgContent.Event == "click" {
		switch msgContent.EventKey {
		case "news":
			// TODO: get news from database
			respContent = "news"
			break
		}
	}
	// Action end
	log.Println("[+] Struct", msgContent)
	respMsgData := &MsgContent{
		FromUsername: msgContent.ToUsername,
		ToUsername:   msgContent.FromUsername,
		CreateTime:   uint32(time.Now().Unix()),
		MsgType:      "text",
		Msgid:        msgContent.Msgid,
		Agentid:      msgContent.Agentid,
		Content:      respContent,
	}
	respData, err := xml.Marshal(respMsgData)
	if err != nil {
		log.Println("[-] Marshal fail")
		return
	}
	encryptMsg, cryptErr := wxcpt.EncryptMsg(string(respData), timeStamp, nonce)
	if nil != cryptErr {
		log.Println("DecryptMsg fail", cryptErr)
		return
	}
	sEncryptMsg := string(encryptMsg)
	log.Println("[+] Receive msg :", msgContent)
	c.String(http.StatusOK, sEncryptMsg)
}

// RunAPIServer RunAPIServer
func RunAPIServer() {
	router := gin.Default()
	// Test Ping
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/verifyurl", verifyURL)
	router.POST("/receive", receiveMsg)

	router.Run(fmt.Sprintf(":%s", PORT)) // listen and serve on 0.0.0.0:8080
}
