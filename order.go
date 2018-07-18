package sinapay

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	tradeTimeOut    = "1m"
	rechargeTimeOut = "15m"
	withdrawTimeOut = "5m"
	isRepayList     = []string{"Y", "N"}
	paytoTypeList   = []string{"FAST", "GENERAL"}
)

// isRepay 支付失败是否重付
const (
	Repay = iota
	Unrepay
)

// paytoType 提现到账类型(速度)
const (
	Fast = iota
	General
)

// CreateHostingCollectTrade 创建托管代收交易 weibopay服务名称：create_hosting_collect_trade
// param:交易订单号,摘要,标的号,付款用户ID,付款用户IP,卡属性,卡类型,金额,外部业务码,是否失败重付,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID,是否代付冻结
// return: 响应参数:交易订单号,交易状态,支付状态,ticket,转跳URL
func CreateHostingCollectTrade(tradeID, summary, goodsID, userID, userIP, cardAttr, cardType, amount string, outTradeCode, isRepay, mode, identityType int, isFreeze bool) (map[string]string, error) {
	data := initBaseParam()
	//业务参数
	data["service"] = "create_hosting_collect_trade"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["trade_close_time"] = tradeTimeOut
	data["can_repay_on_failed"] = isRepayList[isRepay]
	data["goods_id"] = strings.TrimSpace(goodsID)
	data["summary"] = strings.TrimSpace(summary)
	data["out_trade_code"] = outTradeCodeList[outTradeCode]
	data["trade_related_no"] = strconv.FormatInt(time.Now().Unix(), 10)[0:31]
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	//支付参数
	data["payer_id"] = strings.TrimSpace(userID)
	data["payer_identity_type"] = identityTypeList[identityType]
	data["payer_ip"] = strings.TrimSpace(userIP)
	data["pay_method"] = handlePayMethod(cardAttr, cardType, amount, "")
	if isFreeze {
		data["collect_trade_type"] = "pre_auth"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, err
	}
	rt := make(map[string]string)
	if v, ok := rsMap["out_trade_no"]; ok {
		rt["outTradeNo"] = v.(string)
	}
	if v, ok := rsMap["trade_status"]; ok {
		rt["tradeStatus"] = v.(string)
	}
	if v, ok := rsMap["pay_status"]; ok {
		rt["payStatus"] = v.(string)
	}
	if v, ok := rsMap["ticket"]; ok {
		rt["ticket"] = v.(string)
	}
	if v, ok := rsMap["redirect_url"]; ok {
		rt["redirectURL"] = v.(string)
	}
	return rt, nil
}

// CreateSingleHostingPayTrade 创建托管代付交易 weibopay服务名称：create_single_hosting_pay_trade
// param: 交易订单号,摘要,标的号,付款用户ID,收款用户ID,付款用户IP,金额,备注,付款用户标识类型，收款用户标识类型:UID,MemberID,Email,Mobile,付款账户类型,收款账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户，代收分账列表:分账不可超过10笔,分账信息中的付款人必须为收款信息中的收款人，或分账信息中的所有收款人
// return: 交易订单号,交易状态
func CreateSingleHostingPayTrade(tradeID, summary, goodsID, payerID, payeeID, userIP, amount, remarks string, payerIdentityType, payeeIdentityType, payerAccountType, payeeAccountType, outTradeCode int, splitList []map[string]string) (string, string, error) {
	data := initBaseParam()
	data["service"] = "create_single_hosting_pay_trade"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["trade_close_time"] = tradeTimeOut
	if splitList != nil {
		data["split_list"] = handleSplitList(splitList)
	}
	data["goods_id"] = strings.TrimSpace(goodsID)
	data["summary"] = strings.TrimSpace(summary)
	data["amount"] = strings.TrimSpace(amount)
	data["out_trade_code"] = outTradeCodeList[outTradeCode]
	data["trade_related_no"] = strconv.FormatInt(time.Now().Unix(), 10)[0:31]
	//支付参数
	data["payee_identity_id"] = strings.TrimSpace(payeeID)
	data["payee_identity_type"] = identityTypeList[payeeIdentityType]
	data["user_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	// data["creditor_info_list"] = strings.TrimSpace("")
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", "", err
	}
	return rsMap["out_trade_no"].(string), rsMap["trade_status"].(string), nil
}

// PayHostingTrade 托管交易支付 weibopay服务名称：pay_hosting_trade
// param: 支付订单号,付款用户IP,卡属性,卡类型,金额,交易订单号列表
// return: 响应参数:支付订单号,支付状态,ticket,转跳URL
func PayHostingTrade(outPayNo, userIP, cardAttr, cardType, amount string, list []string) (map[string]string, error) {
	data := initBaseParam()
	data["service"] = "pay_hosting_trade"
	data["out_pay_no"] = outPayNo
	// data["extend_param"] = ""
	data["payer_ip"] = userIP
	data["pay_method"] = handlePayMethod(cardAttr, cardType, amount, "")
	var tradeNoList string
	for i := 0; i < len(list); i++ {
		if tradeNoList == "" {
			tradeNoList = list[i]
		} else {
			tradeNoList = fmt.Sprintf("%s^%s", tradeNoList, list[i])
		}
	}
	data["outer_trade_no_list"] = tradeNoList
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, err
	}
	rt := make(map[string]string)
	if v, ok := rsMap["out_pay_no"]; ok {
		rt["outPayNo"] = v.(string)
	}
	if v, ok := rsMap["pay_status"]; ok {
		rt["payStatus"] = v.(string)
	}
	if v, ok := rsMap["ticket"]; ok {
		rt["ticket"] = v.(string)
	}
	if v, ok := rsMap["redirect_url"]; ok {
		rt["redirectURL"] = v.(string)
	}
	return rt, nil
}

// QueryPayResult 查询支付结果 weibopay服务名称：query_pay_result
// param: 支付订单号
// return:支付订单号，支付状态
func QueryPayResult(outPayNo string) (string, string, error) {
	data := initBaseParam()
	data["service"] = "query_pay_result"
	data["out_pay_no"] = outPayNo
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", "", err
	}
	return rsMap["out_pay_no"].(string), rsMap["pay_status"].(string), nil
}

// QueryHostingTrade 托管交易查询 weibopay服务名称：query_hosting_trade
// param: 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,用户标识类型:UID,MemberID,Email,Mobile
// return: 响应参数列表,交易记录列表
func QueryHostingTrade(tradeID, userID, startTime, endTime, pageNo, pageSize string, identityType int) (map[string]string, []map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_hosting_trade"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["start_time"], data["end_time"] = handleStartEndTime(startTime, endTime)
	if pageNo == "" {
		data["page_no"] = defaultPageNo
	}
	if pageSize == "" {
		data["page_size"] = defaultPageSize
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, nil, err
	}
	responParam := make(map[string]string)
	list := make([]map[string]string, 0)
	v, ok := rsMap["total_item"]
	if ok {
		responParam["totalItem"] = v.(string)
		totalItem, _ := strconv.Atoi(responParam["totalItem"])
		if totalItem > 0 {
			arr := strings.Split(rsMap["trade_list"].(string), "|")
			for _, v := range arr {
				vArr := strings.Split(v, "^")
				temp := make(map[string]string)
				temp["outTradeNo"] = vArr[0]
				temp["summary"] = vArr[1]
				temp["amount"] = vArr[2]
				temp["status"] = vArr[3]
				temp["createtime"] = vArr[4]
				temp["lastEditTime"] = vArr[5]
				if len(vArr) == 7 {
					temp["collectAmount"] = vArr[6]
				}
				list = append(list, temp)
			}
		}
	}
	v, ok = rsMap["page_no"]
	if ok {
		responParam["pageNo"] = v.(string)
	}
	v, ok = rsMap["page_size"]
	if ok {
		responParam["pagSize"] = v.(string)
	}
	return responParam, list, nil
}

// CreateHostingRefund 托管退款 weibopay服务名称：create_hosting_refund
// param: 交易订单号,要退款的交易订单号,退款金额(可部分退款),摘要,用户IP
// return: 交易订单号,退款状态
func CreateHostingRefund(tradeID, origtradeID, amount, summary, userIP string) (string, string, error) {
	data := initBaseParam()
	data["service"] = "create_hosting_refund"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["orig_outer_trade_no"] = strings.TrimSpace(origtradeID)
	data["refund_amount"] = strings.TrimSpace(amount)
	data["summary"] = strings.TrimSpace(summary)
	// data["split_list"] = ""
	data["user_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", "", err
	}
	return rsMap["out_trade_no"].(string), rsMap["refund_status"].(string), nil
}

// QueryHostingRefund 托管退款查询  weibopay服务名称：query_hosting_refund
// param: 退款订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,用户标识类型:UID,MemberID,Email,Mobile
// return: 响应参数列表,交易记录列表
func QueryHostingRefund(tradeID, userID, startTime, endTime, pageNo, pageSize string, identityType int) (map[string]string, []map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_hosting_refund"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["start_time"], data["end_time"] = handleStartEndTime(startTime, endTime)
	if pageNo == "" {
		data["page_no"] = defaultPageNo
	}
	if pageSize == "" {
		data["page_size"] = defaultPageSize
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, nil, err
	}
	responParam := make(map[string]string)
	list := make([]map[string]string, 0)
	v, ok := rsMap["total_item"]
	if ok {
		responParam["totalItem"] = v.(string)
		totalItem, _ := strconv.Atoi(responParam["totalItem"])
		if totalItem > 0 {
			arr := strings.Split(rsMap["trade_list"].(string), "|")
			for _, v := range arr {
				vArr := strings.Split(v, "^")
				temp := make(map[string]string)
				temp["outTradeNo"] = vArr[0]
				temp["summary"] = vArr[1]
				temp["amount"] = vArr[2]
				temp["status"] = vArr[3]
				temp["createtime"] = vArr[4]
				temp["lastEditTime"] = vArr[5]
				list = append(list, temp)
			}
		}
	}
	v, ok = rsMap["page_no"]
	if ok {
		responParam["pageNo"] = v.(string)
	}
	v, ok = rsMap["page_size"]
	if ok {
		responParam["pagSize"] = v.(string)
	}
	return responParam, list, nil
}

// CreateHostingDeposit 托管充值 weibopay服务名称：create_hosting_deposit
// param: 交易订单号,摘要,用户ID,用户IP,金额,用户手续费(可空),卡属性,卡类型,银行代码,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,是否代付冻结
// return: 交易订单号,充值状态,线下支付收款单位,线下支付收款账户,线下支付收款账号开户行,线下支付收款备注,收银台重定向地址
func CreateHostingDeposit(tradeID, summary, userID, userIP, amount, userFee, cardAttr, cardtype, bankCode string, accountType, mode, identityType int) (map[string]string, error) {
	data := initBaseParam()
	data["service"] = "create_hosting_deposit"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["amount"] = strings.TrimSpace(amount)
	data["deposit_close_time"] = rechargeTimeOut
	data["payer_ip"] = strings.TrimSpace(userIP)
	data["pay_method"] = handlePayMethod(cardAttr, cardtype, amount, bankCode)
	if strings.TrimSpace(userFee) != "" {
		data["user_fee"] = strings.TrimSpace(userFee)
	}
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, err
	}
	rt := make(map[string]string)
	if v, ok := rsMap["out_trade_no"]; ok {
		rt["outTradeNo"] = v.(string)
	}
	if v, ok := rsMap["deposit_status"]; ok {
		rt["depositStatus"] = v.(string)
	}
	if v, ok := rsMap["trans_account_name "]; ok {
		rt["transAccounName "] = v.(string)
	}
	if v, ok := rsMap["trans_account_no "]; ok {
		rt["transAccountNo "] = v.(string)
	}
	if v, ok := rsMap["trans_bank_brank "]; ok {
		rt["transBankBrank "] = v.(string)
	}
	if v, ok := rsMap["trans_trade_no "]; ok {
		rt["transTradeNo "] = v.(string)
	}
	if v, ok := rsMap["ticket"]; ok {
		rt["ticket"] = v.(string)
	}
	if v, ok := rsMap["redirect_url"]; ok {
		rt["redirectURL"] = v.(string)
	}
	return rt, nil
}

// QueryHostingDeposit 托管充值查询  weibopay服务名称：query_hosting_deposit
// param: 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,用户标识类型:UID,MemberID,Email,Mobile
// return: 响应参数列表,交易记录列表
func QueryHostingDeposit(tradeID, userID, startTime, endTime, pageNo, pageSize string, accountType, identityType int) (map[string]string, []map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_hosting_deposit"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["start_time"], data["end_time"] = handleStartEndTime(startTime, endTime)
	if pageNo == "" {
		data["page_no"] = defaultPageNo
	}
	if pageSize == "" {
		data["page_size"] = defaultPageSize
	}
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, nil, err
	}
	responParam := make(map[string]string)
	list := make([]map[string]string, 0)
	v, ok := rsMap["total_item"]
	if ok {
		responParam["totalItem"] = v.(string)
		totalItem, _ := strconv.Atoi(responParam["totalItem"])
		if totalItem > 0 {
			arr := strings.Split(rsMap["deposit_list"].(string), "|")
			for _, v := range arr {
				vArr := strings.Split(v, "^")
				temp := make(map[string]string)
				temp["outTradeNo"] = vArr[0]
				temp["amount"] = vArr[1]
				temp["status"] = vArr[2]
				temp["createtime"] = vArr[3]
				temp["lastEditTime"] = vArr[4]
				list = append(list, temp)
			}
		}
	}
	v, ok = rsMap["page_no"]
	if ok {
		responParam["pageNo"] = v.(string)
	}
	v, ok = rsMap["page_size"]
	if ok {
		responParam["pagSize"] = v.(string)
	}
	return responParam, list, nil
}

// CreateHostingWithdraw 托管提现  weibopay服务名称：create_hosting_withdraw
// param: 交易订单号,摘要,金额,用户ID,用户IP,用户手续费(可空),卡ID,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,用户标识类型:UID,MemberID,Email,Mobile 提现类型:Fast快速,General普通 返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面 提现模式:ture安全模式,转跳收银台操作
// return: 响应参数:交易订单号,提现状态,收银台重定向地址
func CreateHostingWithdraw(tradeID, summary, amount, userID, userIP, userFee, cardID string, accountType, identityType, paytoType, mode int, withdrawMode bool) (map[string]string, error) {
	data := initBaseParam()
	data["service"] = "create_hosting_withdraw"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["summary"] = strings.TrimSpace(summary)
	data["amount"] = strings.TrimSpace(amount)
	data["payto_type"] = paytoTypeList[paytoType]
	data["withdraw_close_time"] = withdrawTimeOut
	data["user_ip"] = strings.TrimSpace(userIP)
	data["card_id"] = strings.TrimSpace(cardID)
	if withdrawMode {
		data["withdraw_mode"] = "CASHDESK"
	}
	if strings.TrimSpace(userFee) != "" {
		data["user_fee"] = strings.TrimSpace(userFee)
	}
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	rt := make(map[string]string)
	if err != nil {
		return nil, err
	}
	if v, ok := rsMap["out_trade_no"]; ok {
		rt["outTrade"] = v.(string)
	}
	if v, ok := rsMap["withdraw_status"]; ok {
		rt["withdrawStatus"] = v.(string)
	}
	if v, ok := rsMap["redirect_url"]; ok {
		rt["redirectURL"] = v.(string)
	}
	return rt, nil
}

// QueryHostingWithdraw weibopay服务名称：query_hosting_withdraw
// param: 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,用户标识类型:UID,MemberID,Email,Mobile
// return: 响应参数列表,交易记录列表
func QueryHostingWithdraw(tradeID, userID, startTime, endTime, pageNo, pageSize string, accountType, identityType int) (map[string]string, []map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_hosting_withdraw"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["start_time"], data["end_time"] = handleStartEndTime(startTime, endTime)
	if pageNo == "" {
		data["page_no"] = defaultPageNo
	}
	if pageSize == "" {
		data["page_size"] = defaultPageSize
	}
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, nil, err
	}
	responParam := make(map[string]string)
	list := make([]map[string]string, 0)
	v, ok := rsMap["total_item"]
	if ok {
		responParam["totalItem"] = v.(string)
		totalItem, _ := strconv.Atoi(responParam["totalItem"])
		if totalItem > 0 {
			arr := strings.Split(rsMap["withdraw_list"].(string), "|")
			for _, v := range arr {
				vArr := strings.Split(v, "^")
				temp := make(map[string]string)
				temp["outTradeNo"] = vArr[0]
				temp["amount"] = vArr[1]
				temp["status"] = vArr[2]
				temp["createtime"] = vArr[3]
				temp["lastEditTime"] = vArr[4]
				list = append(list, temp)
			}
		}
	}
	v, ok = rsMap["page_no"]
	if ok {
		responParam["pageNo"] = v.(string)
	}
	v, ok = rsMap["page_size"]
	if ok {
		responParam["pagSize"] = v.(string)
	}
	return responParam, list, nil
}

// CreateSingleHostingPaytoCardTrade 创建单笔代付到提现卡交易 weibopay服务名称：create_single_hosting_pay_to_card_trade
// param:交易订单号,摘要,金额,用户ID,用户IP,卡ID,标的号,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 用户标识类型:UID,MemberID,Email,Mobile 提现类型:Fast快速,General普通 外部业务码
// return: 响应参数:交易订单号,提现状态
func CreateSingleHostingPaytoCardTrade(tradeID, summary, amount, userID, userIP, cardID, goodsID string, identityType, paytoType, outTradeCode int) (map[string]string, error) {
	data := initBaseParam()
	data["service"] = "create_single_hosting_pay_to_card_trade"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["summary"] = strings.TrimSpace(summary)
	data["amount"] = strings.TrimSpace(amount)
	data["payto_type"] = paytoTypeList[paytoType]
	data["collect_method"] = handleCollectMethod(userID, cardID, identityType)
	data["user_ip"] = strings.TrimSpace(userIP)
	data["goods_id"] = strings.TrimSpace(goodsID)
	data["out_trade_code"] = outTradeCodeList[outTradeCode]
	// data["creditor_info_list"] = ""
	// data["extend_param"] = ""
	rs, err := Request(&data, OrderMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	rt := make(map[string]string)
	if err != nil {
		return nil, err
	}
	rt["outTradeNo"] = rsMap["out_trade_no"].(string)
	if v, ok := rsMap["withdraw_status"]; ok {
		rt["withdrawStatus"] = v.(string)
	}
	return rt, nil
}

// FinishPreAuthTrade 代收完成 weibopay服务名称：finish_pre_auth_trade
// param:请求交易号,用户IP,交易列表
func FinishPreAuthTrade(outRequestNo, userIP string, TradeList []map[string]string) error {
	data := initBaseParam()
	data["service"] = "finish_pre_auth_trade"
	data["out_request_no"] = strings.TrimSpace(outRequestNo)
	data["user_ip"] = strings.TrimSpace(userIP)
	if TradeList != nil {
		data["trade_list"] = handleTradeList(TradeList)
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = checkResponseCode(rs)
	if err != nil {
		return err
	}
	return nil
}

// CancelPreAuthTrade 代收撤销 weibopay服务名称：cancel_pre_auth_trade
// param:请求交易号,交易列表
func CancelPreAuthTrade(outRequestNo string, TradeList []map[string]string) error {
	data := initBaseParam()
	data["service"] = "cancel_pre_auth_trade"
	data["out_request_no"] = strings.TrimSpace(outRequestNo)
	if TradeList != nil {
		data["trade_list"] = handleCancelTradeList(TradeList)
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = checkResponseCode(rs)
	if err != nil {
		return err
	}
	return nil
}
