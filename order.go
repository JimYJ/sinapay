package sinapay

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	tradeTimeOut    = "1m"
	rechargeTimeOut = "15m"
	isRepayList     = []string{"Y", "N"}
)

// isRepay 支付失败是否重付
const (
	Repay = iota
	Unrepay
)

//CreateHostingCollectTrade 创建托管代收交易 weibopay服务名称：create_hosting_collect_trade
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
	if outTradeCode != CollectionAll {
		data["out_trade_code"] = outTradeCodeList[outTradeCode]
	}
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
	if outTradeCode != CollectionAll {
		data["out_trade_code"] = outTradeCodeList[outTradeCode]
	}
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
// param: 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数
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

// CreateHostingDeposit 托管充值 weibopay服务名称：create_hosting_deposit
// param: 交易订单号,摘要,用户ID,用户IP,金额,用户手续费(可空),卡属性,卡类型,银行代码,账户类型() BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面，是否代付冻结
// return: 交易订单号,充值状态,线下支付收款单位,线下支付收款账户,线下支付收款账号开户行,线下支付收款备注,收银台重定向地址
func CreateHostingDeposit(tradeID, summary, userID, userIP, amount, userFee, cardAttr, cardtype, bankCode string, accountType, mode int) (map[string]string, error) {
	data := initBaseParam()
	data["service"] = "create_hosting_deposit"
	data["out_trade_no"] = strings.TrimSpace(tradeID)
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = "UID"
	data["amount"] = amount
	data["payer_ip"] = userIP
	data["deposit_close_time"] = rechargeTimeOut
	data["payer_ip"] = userIP
	data["pay_method"] = handlePayMethod(cardAttr, cardtype, amount, bankCode)
	if userFee != "" {
		data["user_fee"] = userFee
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
