package sinapay

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	userURL   = "https://gate.pay.sina.com.cn/mgs/gateway.do"
	orderURL  = "https://gate.pay.sina.com.cn/mas/gateway.do"
	debugMode = false
)

const (
	// UserMode 用户模式
	UserMode = iota
	// OrderMode 用户模式
	OrderMode
)

var (
	partnerid        string
	publicPEM        string
	privatePEM       string
	sinaPayPublicPEM string
)

// InitSinaPay 初始化证书
func InitSinaPay(pid, pubPEM, privPEM, sinapayPubPEM string) {
	partnerid = pid
	publicPEM = pubPEM
	privatePEM = privPEM
	sinaPayPublicPEM = sinapayPubPEM
}

// InitBaseParam 初始化请求参数
func initBaseParam() map[string]string {
	baseParam := make(map[string]string)
	baseParam["version"] = "1.2"
	baseParam["partner_id"] = partnerid
	baseParam["_input_charset"] = "UTF-8"
	return baseParam
}

// Request 发送请求
func Request(data *map[string]string, mode int) (string, error) {
	err := encryptPostData(data)
	if err != nil {
		log.Panicln(err)
	}
	(*data)["request_time"] = time.Now().Local().Format("20060102150405")
	sign, err := getSign(*data)
	if err != nil {
		log.Panicln(err)
	}

	(*data)["sign"] = sign
	(*data)["sign_type"] = "RSA"
	(*data)["sign_version"] = "1.0"

	postData, err := handlePostData(data)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	var rs []byte
	if mode == UserMode {
		rs, err = postForm(userURL, postData)
	} else {
		rs, err = postForm(orderURL, postData)
	}
	if err != nil {
		log.Println(err)
	}
	rsStr, _ := url.QueryUnescape(string(rs))
	debugPrint("respon json:", rsStr)
	return rsStr, err
}

// 获得排序后的KEY列表
func getKeylist(data *map[string]string) []string {
	var keylist []string
	for k := range *data {
		keylist = append(keylist, k)
	}
	sort.Strings(keylist)
	return keylist
}

// 获得验签的排序后的KEY列表
func getVerifyKeylist(data *map[string]interface{}) []string {
	var keylist []string
	for k := range *data {
		keylist = append(keylist, k)
	}
	sort.Strings(keylist)
	return keylist
}

// 加密请求参数
func encryptPostData(data *map[string]string) error {
	newrsa := GetRSA()
	for needKey := range needRSAKey {
		if v, ok := (*data)[needKey]; ok {
			tmp, err := newrsa.PublicEncrypt(v)
			if err != nil {
				log.Println(err)
				return err
			}
			(*data)[needKey] = tmp
		}
	}
	// log.Println(*data)
	return nil
}

//拼接待加密字符串
func getSign(data map[string]string) (string, error) {
	delete(data, "sign")
	delete(data, "sign_type")
	delete(data, "sign_version")
	keylist := getKeylist(&data)
	unsign := ""
	for _, v := range keylist {
		value, _ := data[v]
		if value == "" {
			continue
		}
		if unsign == "" {
			unsign = fmt.Sprintf("%s=%s", v, value)
		} else {
			unsign += fmt.Sprintf("&%s=%s", v, value)
		}
	}
	if unsign == "" {
		err := "cannot get param"
		log.Println(err)
		return "", errors.New(err)
	}
	debugPrint("unsign string:", unsign)
	newrsa := GetRSA()
	return newrsa.Sign(unsign)
}

// 验签
func getVerify(data map[string]interface{}) error {
	var sign string
	if v, ok := data["sign"]; ok {
		sign = v.(string)
	} else {
		if v2, ok := data["response_message"]; ok {
			return errors.New(v2.(string))
		}
		return errors.New("response is error, pls use debug mode for debug")
	}
	delete(data, "sign")
	delete(data, "sign_type")
	delete(data, "sign_version")
	keylist := getVerifyKeylist(&data)
	unsign := ""
	for _, v := range keylist {
		value, _ := data[v]
		if value == "" {
			continue
		}
		if unsign == "" {
			unsign = fmt.Sprintf("%s=%s", v, value)
		} else {
			unsign += fmt.Sprintf("&%s=%s", v, value)
		}
	}
	if unsign == "" {
		err := "cannot get param"
		log.Println(err)
		return errors.New(err)
	}
	newrsa := GetRSA()
	return newrsa.Verify(unsign, sign)
}

//将字段添加到发送数据
func handlePostData(data *map[string]string) (url.Values, error) {
	keylist := getKeylist(data)
	postData := url.Values{}
	// log.Println(data)
	for _, v := range keylist {
		value, _ := (*data)[v]
		postData.Set(v, value)
	}
	if postData == nil {
		err := "cannot get param"
		log.Println(err)
		return nil, errors.New(err)
	}
	return postData, nil
}

// 判断是否响应成功
func checkResponseCode(rs string) (map[string]interface{}, error) {
	rsMap := json2Map([]byte(rs))
	err := getVerify(rsMap)
	if err != nil {
		log.Println("sign verify fail:", err)
		return rsMap, err
	}
	responseCode := rsMap["response_code"].(string)
	responseMessage := rsMap["response_message"].(string)
	log.Println("request results:", responseCode, ":", responseMessage)
	if responseCode != "APPLY_SUCCESS" {
		return rsMap, errors.New(responseMessage)
	}
	return rsMap, nil
}

// 格式化开始结束时间
func handleStartEndTime(startTime, endTime string) (string, string) {
	var s, e string
	st, err := time.Parse("2006-01-02 15:04:05", startTime)
	et, err2 := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil || err2 != nil {
		//如果参数不正确或未填写，默认查最近7日内数据
		d, _ := time.ParseDuration("-24h")
		s = time.Now().Local().Add(d * 7).Format("20060102150405")
		e = time.Now().Local().Format("20060102150405")
	} else {
		//如果开始时间比结束时间晚，则依然查最近7日的数据
		if err == nil && st.Before(et) {
			s = st.Format("20060102150405")
			e = et.Format("20060102150405")
		} else {
			d, _ := time.ParseDuration("-24h")
			s = time.Now().Local().Add(d * 7).Format("20060102150405")
			e = time.Now().Local().Format("20060102150405")
		}
	}
	return s, e
}

// 解密响应参数
func decryptResponData(data *map[string]interface{}, decryptlist map[string]bool) error {
	newrsa := GetRSA()
	for needKey := range decryptlist {
		if v, ok := (*data)[needKey]; ok {
			strV := v.(string)
			tmp, err := newrsa.PrivateDecrypt(strV)
			if err != nil {
				log.Println(err)
				return err
			}
			(*data)[needKey] = tmp
		}
	}
	// log.Println(*data)
	return nil
}

// TestMode 测试模式,请求测试环境地址
func TestMode() {
	userURL = "https://testgate.pay.sina.com.cn/mgs/gateway.do"
	orderURL = "https://testgate.pay.sina.com.cn/mas/gateway.do"
}

// 拼接支付方式
func handlePayMethod(cardAttr, cardtype, amount, bankCode string) string {
	amount = strings.TrimSpace(amount)
	cardAttr = strings.TrimSpace(cardAttr)
	cardtype = strings.TrimSpace(cardtype)
	if bankCode == "" {
		bankCode = "SINAPAY"
	}
	return fmt.Sprintf("online_bank^%s^%s,%s,%s", amount, bankCode, cardAttr, cardtype)
}

// 拼接代收分账信息
func handleSplitList(splitList []map[string]string) string {
	var str string
	for _, item := range splitList {
		payerID := strings.TrimSpace(item["payerID"])
		payeeID := strings.TrimSpace(item["payeeID"])
		amount := strings.TrimSpace(item["amount"])
		remarks := strings.TrimSpace(item["remarks"])
		payeeIdentityType, err := strconv.Atoi(strings.TrimSpace(item["payeeIdentityType"]))
		payeeAccountType, err2 := strconv.Atoi(strings.TrimSpace(item["payeeAccountType"]))
		payerIdentityType, err3 := strconv.Atoi(strings.TrimSpace(item["payerIdentityType"]))
		payerAccountType, err4 := strconv.Atoi(strings.TrimSpace(item["payerAccountType"]))
		if err != nil || err2 != nil || err3 != nil || err4 != nil || payerID == "" || payeeID == "" || amount == "" {
			continue
		}
		if str == "" {
			str = fmt.Sprintf("%s^%s^%s^%s^%s^%s^%s^%s",
				payeeID, identityTypeList[payeeIdentityType], acountTypeList[payeeAccountType],
				payerID, identityTypeList[payerIdentityType], acountTypeList[payerAccountType],
				amount, remarks)
		} else {
			str = fmt.Sprintf("%s|%s^%s^%s^%s^%s^%s^%s^%s", str,
				payeeID, identityTypeList[payeeIdentityType], acountTypeList[payeeAccountType],
				payerID, identityTypeList[payerIdentityType], acountTypeList[payerAccountType],
				amount, remarks)
		}
	}
	return str
}

// DebugMode 调试模式,打印请求,响应报文
func DebugMode() {
	debugMode = true
}

// 打印debug日志
func debugPrint(logs ...interface{}) {
	if debugMode {
		log.Println(logs...)
	}
}

// 拼接收款方式
func handleCollectMethod(userID, cardID string, identityType int) string {
	cardID = strings.TrimSpace(cardID)
	userID = strings.TrimSpace(userID)
	IDType := identityTypeList[identityType]
	return fmt.Sprintf("binding_card^%s,%s,%s", userID, IDType, cardID)
}

// 拼接代收完成交易列表
func handleTradeList(TradeList []map[string]string) string {
	var str string
	for _, item := range TradeList {
		requestID := strings.TrimSpace(item["requestID"])
		tradeID := strings.TrimSpace(item["tradeID"])
		amount := strings.TrimSpace(item["amount"])
		summary := strings.TrimSpace(item["summary"])
		extendParam := strings.TrimSpace(item["extendParam"])
		if requestID == "" || tradeID == "" || amount == "" || extendParam == "" {
			continue
		}
		if str == "" {
			str = fmt.Sprintf("%s^%s^%s^%s^%s", requestID, tradeID, amount, summary, extendParam)
		} else {
			str = fmt.Sprintf("%s|%s^%s^%s^%s^%s", str, requestID, tradeID, amount, summary, extendParam)
		}
	}
	return str
}

// 拼接代收完成交易列表
func handleCancelTradeList(TradeList []map[string]string) string {
	var str string
	for _, item := range TradeList {
		requestID := strings.TrimSpace(item["requestID"])
		tradeID := strings.TrimSpace(item["tradeID"])
		summary := strings.TrimSpace(item["summary"])
		extendParam := strings.TrimSpace(item["extendParam"])
		if requestID == "" || tradeID == "" || extendParam == "" {
			continue
		}
		if str == "" {
			str = fmt.Sprintf("%s^%s^%s^%s", requestID, tradeID, summary, extendParam)
		} else {
			str = fmt.Sprintf("%s|%s^%s^%s^%s", str, requestID, tradeID, summary, extendParam)
		}
	}
	return str
}
