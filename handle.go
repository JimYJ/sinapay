package sinapay

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"time"
)

var (
	userURL  = "https://gate.pay.sina.com.cn/mgs/gateway.do"
	orderURL = "https://gate.pay.sina.com.cn/mas/gateway.do"
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
	// log.Println(unsign)
	newrsa := GetRSA()
	return newrsa.Sign(unsign)
}

// 验签
func getVerify(data map[string]interface{}) error {
	sign := data["sign"].(string)
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
	log.Println(responseCode, ":", responseMessage)
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
