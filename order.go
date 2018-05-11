package sinapay

import (
	"log"
	"strings"
)

var (
	tradeTimeOut    = "1m"
	rechargeTimeOut = "15m"
	isRepayList     = []string{"Y", "N"}
)

const (
	repay = iota
	unrepay
)

//CreateHostingCollectTrade weibopay服务名称：create_hosting_collect_trade
// param:交易订单号,摘要,标的号,付款用户ID,付款用户标识,卡属性,卡类型,金额，外部业务码,是否失败重付,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面，是否代付冻结
// return: 响应参数:交易订单号,交易状态,支付状态,ticket,转跳URL
func CreateHostingCollectTrade(tradeID, summary, goodsID, userID, userIP, cardAttr, cardtype, amount string, outTradeCode, isRepay, mode int, isFreeze bool) (map[string]string, error) {
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
	data["payer_identity_type"] = "UID"
	data["payer_ip"] = strings.TrimSpace(userIP)
	data["pay_method"] = handlePayMethod(cardAttr, cardtype, amount, "")
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

// CreateHostingDeposit 托管充值 weibopay服务名称：create_hosting_deposit
// param: 交易订单号,摘要,用户ID,用户IP,金额,用户手续费(可空),卡属性,卡类型,银行代码,账户类型 BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面，是否代付冻结
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
