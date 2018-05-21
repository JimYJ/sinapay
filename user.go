package sinapay

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//返回转跳地址类型 移动端网页或PC网页
const (
	RedirectURLMobile = iota
	RedirectURLPC
)

//查询账户类型 个人账户推荐使用存钱罐SavingPot操作资金,基本户存在限额限次的问题
const (
	Basic = iota
	Ensure
	Reserve
	SavingPot
	Bank
	Default
)

//外部业务码(交易码)
const (
	CollectionMore       = iota //代收-其他
	CollectionInvestment        //代收投资
	CollectionRepayment         //代收还款
	CollectionMoreII            //代付-其他
	CollectionLoan              //代付-借款
	CollectionPrincipal         //代付-本金收益
)

// 用户账户标识类型 identity_type
const (
	UID = iota
	MemberID
	Email
	Mobile
)

var (
	acountTypeList   = []string{"BASIC", "ENSURE", "RESERVE", "SAVING_POT", "BANK"}
	defaultPageNo    = "1"
	defaultPageSize  = "20"
	freezeNoPrefix   = "FZ"
	unfreezeNoPrefix = "NF"
	outTradeCodeList = []string{"1000", "1001", "1002", "2000", "2001", "2001"}
	identityTypeList = []string{"UID", "MEMBER_ID", "MOBILE", "EMAIL"}
)

// CreateActiveMember 创建激活会员 weibopay服务名称：create_activate_member
// param :用户请求IP,用户ID,用户账户标识类型:UID
func CreateActiveMember(userIP, userID string, identityType int) error {
	data := initBaseParam()
	data["service"] = "create_activate_member"
	data["client_ip"] = strings.TrimSpace(userIP)
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["member_type"] = "1"
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

// SetRealName 设置用户实名信息  weibopay服务名称：set_real_name
// param :用户ID,真实姓名,身份证号,用户请求IP,用户账户标识类型:UID
func SetRealName(userID, realname, IDNumber, userIP string, identityType int) error {
	data := initBaseParam()
	data["service"] = "set_real_name"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["real_name"] = strings.TrimSpace(realname)
	data["cert_type"] = "IC"
	data["cert_no"] = strings.TrimSpace(IDNumber)
	data["client_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	// data["need_confirm"] = ""
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

// SetPayPassword 设置支付密码 weibopay服务名称：set_pay_password
// param: 用户ID,委托扣款展示方式(可空),同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC:返回PC页面,用户账户标识类型:UID
// return: 转跳页面
func SetPayPassword(userID, withholdParam, returnURL, notifyURL string, mode, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "set_pay_password"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	if notifyURL != "" {
		data["notify_url"] = strings.TrimSpace(notifyURL)
	}
	if returnURL != "" {
		data["return_url"] = strings.TrimSpace(returnURL)
	}
	data["return_url"] = returnURL
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	if withholdParam != "" {
		data["withhold_param"] = withholdParam
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["redirect_url"].(string), nil
}

// ModifyPayPassword 修改支付密码 weibopay服务名称：modify_pay_password
// param: 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
// return: 转跳页面
func ModifyPayPassword(userID, returnURL, notifyURL string, mode, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "modify_pay_password"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	if notifyURL != "" {
		data["notify_url"] = strings.TrimSpace(notifyURL)
	}
	if returnURL != "" {
		data["return_url"] = strings.TrimSpace(returnURL)
	}
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["redirect_url"].(string), nil
}

// FindPayPass 找回支付密码 weibopay服务名称：find_pay_password
// param: 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
// return: 转跳页面
func FindPayPass(userID, returnURL, notifyURL string, mode, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "find_pay_password"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	if notifyURL != "" {
		data["notify_url"] = strings.TrimSpace(notifyURL)
	}
	if returnURL != "" {
		data["return_url"] = strings.TrimSpace(returnURL)
	}
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["redirect_url"].(string), nil
}

// QueryisSetPayPassword 检测是否已设置支付密码 weibopay服务名称： query_is_set_pay_password
// param: 用户ID,同步回跳页面(可空),异步通知接口(可空),用户账户标识类型:UID
// return: 转跳页面
func QueryisSetPayPassword(userID, returnURL, notifyURL string, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "query_is_set_pay_password"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["is_set_paypass"].(string), nil
}

// BindingBankCard 绑定银行卡 weibopay服务名称： binding_bank_card 默认SIGN模式
// param: 用户ID,用户IP,银行编号,银行卡号,账户姓名,账户绑定手机,省份,城市,用户账户标识类型:UID
// return: 卡ID(SIGN模式不返回),是否验证银行信息(SIGN模式不返回),绑卡推进ticket
func BindingBankCard(userID, userIP, bankCode, bankAccountNo, accounName, phone, province, city string, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "binding_bank_card"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["request_no"] = strconv.FormatInt(time.Now().Local().UnixNano(), 10)
	data["bank_code"] = strings.TrimSpace(bankCode)
	data["bank_account_no"] = strings.TrimSpace(bankAccountNo)
	data["account_name"] = strings.TrimSpace(accounName)
	data["card_type"] = "DEBIT"  //只允许借记卡
	data["card_attribute"] = "C" //只允许对私账户
	// data["cert_type"] = ""
	// data["cert_no"] = ""
	data["phone_no"] = strings.TrimSpace(phone)
	// data["validity_period"] = ""    //信用卡-有效期
	// data["verification_value"] = "" //信用卡-CVV2
	data["province"] = strings.TrimSpace(province)
	data["city"] = strings.TrimSpace(city)
	// data["bank_branch"] = "" //支行名称
	data["verify_mode"] = "SIGN"
	data["client_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["ticket"].(string), nil //rsMap["card_id"].(string), rsMap["is_verified"].(string)
}

// BindingBankCardAdvance 绑定银行卡推进 weibopay服务名称：binding_bank_card_advance
// param: 绑定银行卡返回的ticket(SIGN模式),短信验证码,用户IP
// return: 卡ID,银行卡信息验证状态
func BindingBankCardAdvance(ticket, validCode, userIP string) (string, string, error) {
	data := initBaseParam()
	data["service"] = "binding_bank_card_advance"
	data["ticket"] = strings.TrimSpace(ticket)
	data["valid_code"] = strings.TrimSpace(validCode)
	data["client_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", "", err
	}
	return rsMap["card_id"].(string), rsMap["is_verified"].(string), nil
}

// UnbindingBankCard 解绑银行卡 weibopay服务名称：unbinding_bank_card
// param: 用户ID,用户IP,卡ID,用户账户标识类型:UID
// return: 解绑推进ticket
func UnbindingBankCard(userID, userIP, cardID string, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "unbinding_bank_card"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["card_id"] = cardID
	data["advance_flag"] = "Y"
	data["client_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["ticket"].(string), nil
}

// UnbindingBankCardAdvance 解绑银行卡推进 weibopay服务名称：unbinding_bank_card_advance
// param: 用户ID,用户IP,解绑银行卡返回的ticket,短信验证码,用户账户标识类型:UID
func UnbindingBankCardAdvance(userID, userIP, ticket, validCode string, identityType int) error {
	data := initBaseParam()
	data["service"] = "unbinding_bank_card_advance"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["ticket"] = strings.TrimSpace(ticket)
	data["valid_code"] = strings.TrimSpace(validCode)
	data["client_ip"] = strings.TrimSpace(userIP)
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

// QueryBankCard 查询银行卡 weibopay服务名称：query_bank_card
// param: 用户ID,卡ID,用户账户标识类型:UID
// return: 卡列表
func QueryBankCard(userID, cardID string, identityType int) ([]map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_bank_card"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	if cardID != "" {
		data["card_id"] = cardID
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]string, 0)
	arr := strings.Split(rsMap["card_list"].(string), "|")
	for _, v := range arr {
		vArr := strings.Split(v, "^")
		temp := make(map[string]string)
		temp["cardID"] = vArr[0]
		temp["bankCode"] = vArr[1]
		temp["bankAccountNo"] = vArr[2]
		temp["accountName"] = vArr[3]
		temp["cardType"] = vArr[4]
		temp["cardAttribute"] = vArr[5]
		temp["verifyMode"] = vArr[6]
		temp["createDate"] = vArr[7]
		temp["isSafeCard"] = vArr[8]
		list = append(list, temp)
	}
	// log.Println(rsMap["card_list"], list)
	return list, nil
}

// QueryBalance 查询余额/基金份额 weibopay服务名称：query_balance
// param:用户ID,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 ,用户账户标识类型:UID,MemberID,Email,Mobile
// return: 余额/基金份额,可用余额/基金份额,存钱罐收益(非查询存钱罐时为nil)
func QueryBalance(userID string, accountType, identityType int) (string, string, map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_balance"
	data["identity_id"] = strings.TrimSpace(userID)
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	data["identity_type"] = identityTypeList[identityType]
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", "", nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", "", nil, err
	}
	bonus := make(map[string]string)
	if accountType == SavingPot {
		temp := strings.Split(rsMap["bonus"].(string), "^")
		bonus["lastday"] = temp[0]
		bonus["lastmonth"] = temp[1]
		bonus["all"] = temp[2]
	} else {
		bonus = nil
	}
	return rsMap["balance"].(string), rsMap["available_balance"].(string), bonus, nil
}

// QueryAccountDetails 查询收支明细 weibopay服务名称：query_account_details
// param: 用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 ,用户账户标识类型:UID,MemberID,Email,Mobile
// return: 参数列表,收支明细列表
func QueryAccountDetails(userID, startTime, endTime, pageNo, pageSize string, accountType, identityType int) (map[string]string, []map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_account_details"
	data["identity_id"] = strings.TrimSpace(userID)
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	data["identity_type"] = identityTypeList[identityType]
	data["start_time"], data["end_time"] = handleStartEndTime(startTime, endTime)
	if pageNo == "" {
		data["page_no"] = defaultPageNo
	}
	if pageSize == "" {
		data["page_size"] = defaultPageSize
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
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
			arr := strings.Split(rsMap["detail_list"].(string), "|")
			for _, v := range arr {
				vArr := strings.Split(v, "^")
				temp := make(map[string]string)
				temp["summary"] = vArr[0]
				temp["createtime"] = vArr[1]
				temp["incordec"] = vArr[2] //增减方向
				temp["change"] = vArr[3]
				temp["balance"] = vArr[4]
				if accountType == SavingPot {
					temp["type"] = vArr[5]
				}
				list = append(list, temp)
			}
		}
	}
	v, ok = rsMap["total_income"]
	if ok {
		responParam["totalIncome"] = v.(string)
	}
	v, ok = rsMap["total_payout"]
	if ok {
		responParam["totalPayout"] = v.(string)
	}
	v, ok = rsMap["total_bonus"]
	if ok {
		responParam["totalBonus"] = v.(string)
	}
	// log.Println(responParam, list)
	return responParam, list, nil
}

// BalanceFreeze 冻结余额 weibopay服务名称：balance_freeze
// param: 用户ID,用户IP,摘要,金额,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 ,用户账户标识类型:UID
// return: 冻结单号(查询状态和解冻用)
func BalanceFreeze(userID, userIP, summary, amount string, accountType, identityType int) (string, error) {
	data := initBaseParam()
	outFreezeNo := fmt.Sprintf("%s%s", freezeNoPrefix, strconv.FormatInt(time.Now().Local().UnixNano(), 10))
	data["service"] = "balance_freeze"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["out_freeze_no"] = outFreezeNo
	data["amount"] = strings.TrimSpace(amount)
	data["summary"] = strings.TrimSpace(summary)
	data["client_ip"] = strings.TrimSpace(userIP)
	if accountType != Default {
		data["account_type"] = acountTypeList[accountType]
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return outFreezeNo, nil
}

// BalanceUnfreeze 解冻余额 weibopay服务名称：balance_unfreeze
// param:用户ID,用户IP,原冻结单号,摘要,金额(为空表示全额解冻),用户账户标识类型:UID
// return: 解冻单号(查询状态用)
func BalanceUnfreeze(userID, userIP, outFreezeNo, summary, amount string, identityType int) (string, error) {
	data := initBaseParam()
	outUnfreezeNo := fmt.Sprintf("%s%s", freezeNoPrefix, strconv.FormatInt(time.Now().Local().UnixNano(), 10))
	data["service"] = "balance_unfreeze"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["out_unfreeze_no"] = outUnfreezeNo
	data["out_freeze_no"] = strings.TrimSpace(outFreezeNo)
	if amount != "" {
		data["amount"] = strings.TrimSpace(amount)
	}
	data["summary"] = strings.TrimSpace(summary)
	data["client_ip"] = strings.TrimSpace(userIP)
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return outUnfreezeNo, nil
}

// QueryCtrlResult 查询冻结解冻结果 weibopay服务名称：query_ctrl_result
// param: 冻结解冻订单号
// return: 是否成功，失败原因(或请求接口报错)
func QueryCtrlResult(outCtrlNo string) (bool, error) {
	data := initBaseParam()
	data["service"] = "query_ctrl_result"
	data["out_ctrl_no"] = outCtrlNo
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return false, err
	}
	rsMsg, err := checkResponseCode(rs)
	if err != nil {
		return false, err
	}
	var ctrlErr error
	var status bool
	if rsMsg["ctrl_status"].(string) != "SUCCESS" {
		ctrlErr = errors.New(rsMsg["error_msg"].(string))
		status = true
	} else {
		status = false
	}
	return status, ctrlErr
}

// QueryMemberInfos 查询企业会员信息 weibopay服务名称：query_member_infos
// param: 用户ID,用户账户标识类型:UID
// return: 查询结果
func QueryMemberInfos(userID string, identityType int) (map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_member_infos"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	// data["member_type"] = "2"
	data["is_mask"] = "Y"
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, err
	}
	decryptlist := map[string]bool{"email": true, "cert_no": true, "legal_person": true, "legal_person_phone": true}
	decryptResponData(&rsMap, decryptlist)
	responParam := make(map[string]string)
	responParam["companyName"] = rsMap["company_name"].(string)
	responParam["address"] = rsMap["address"].(string)
	responParam["licenseNo"] = rsMap["license_no"].(string)
	responParam["licenseExpireDate"] = rsMap["license_expire_date"].(string)
	responParam["businessScope"] = rsMap["business_scope"].(string)
	responParam["telephone"] = rsMap["telephone"].(string)
	responParam["organizationNo"] = rsMap["organization_no"].(string)
	responParam["summary"] = rsMap["summary"].(string)
	responParam["legalPerson"] = rsMap["legal_person"].(string)
	responParam["certNo"] = rsMap["cert_no"].(string)
	responParam["certType"] = rsMap["cert_type"].(string)
	responParam["legalPersonPhone"] = rsMap["legal_person_phone"].(string)
	responParam["extendParam"] = rsMap["extend_param"].(string)
	v, ok := rsMap["email"]
	if ok {
		responParam["email"] = v.(string)
	}
	v, ok = rsMap["website"]
	if ok {
		responParam["website"] = v.(string)
	}
	v, ok = rsMap["extend_param"]
	if ok {
		responParam["extendParam"] = v.(string)
	}
	return responParam, nil
}

// QueryAuditResult 查询企业会员审核结果 weibopay服务名称：query_audit_result
// param: 用户ID,用户账户标识类型:UID
// return: 是否成功，失败原因(或请求接口报错)
func QueryAuditResult(userID string, identityType int) (bool, error) {
	data := initBaseParam()
	data["service"] = "query_audit_result"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return false, err
	}
	rsMsg, err := checkResponseCode(rs)
	if err != nil {
		return false, err
	}
	var ctrlErr error
	var status bool
	if rsMsg["audit_result"].(string) != "SUCCESS" {
		ctrlErr = errors.New(rsMsg["audit_mgs"].(string))
		status = true
	} else {
		status = false
	}
	return status, ctrlErr
}

// QueryMiddleAccount 查询中间账户 weibopay服务名称：query_middle_account
// param: 外部业务码
// return: 查询结果列表
func QueryMiddleAccount(outTradeCode int) ([]map[string]string, error) {
	data := initBaseParam()
	data["service"] = "query_middle_account"
	data["out_trade_code"] = outTradeCodeList[outTradeCode]
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]string, 0)
	arr := strings.Split(rsMap["account_list"].(string), "|")
	for _, v := range arr {
		vArr := strings.Split(v, "^")
		temp := make(map[string]string)
		temp["outTradeCode"] = vArr[0]
		temp["middleAccountNo"] = vArr[1]
		temp["balance"] = vArr[2]
		list = append(list, temp)
	}
	return list, nil
}

// ModifyVerifyMobile 修改认证手机 weibopay服务名称：modify_verify_mobile
// param: 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
// return: 转跳页面
func ModifyVerifyMobile(userID, notifyURL, returnURL string, mode, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "modify_verify_mobile"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	if notifyURL != "" {
		data["notify_url"] = strings.TrimSpace(notifyURL)
	}
	if returnURL != "" {
		data["return_url"] = strings.TrimSpace(returnURL)
	}
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["redirect_url"].(string), nil
}

// FindVerifyMobile 修改认证手机 weibopay服务名称：find_verify_mobile
// param: 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
// return: 转跳页面
func FindVerifyMobile(userID, notifyURL, returnURL string, mode, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "find_verify_mobile"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	if notifyURL != "" {
		data["notify_url"] = strings.TrimSpace(notifyURL)
	}
	if returnURL != "" {
		data["return_url"] = strings.TrimSpace(returnURL)
	}
	if mode == RedirectURLMobile {
		data["cashdesk_addr_category"] = "MOBILE"
	}
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["redirect_url"].(string), nil
}

// ChangeBankMobile 修改银行预留手机 weibopay服务名称：change_bank_mobile
// param: 用户ID,卡ID,手机号,用户账户标识类型:UID
// return: ticket(推进接口用)
func ChangeBankMobile(userID, cardID, phone string, identityType int) (string, error) {
	data := initBaseParam()
	data["service"] = "change_bank_mobile"
	data["identity_id"] = strings.TrimSpace(userID)
	data["identity_type"] = identityTypeList[identityType]
	data["card_id"] = strings.TrimSpace(cardID)
	data["phone_no"] = strings.TrimSpace(phone)
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", err
	}
	return rsMap["ticket"].(string), nil
}

// ChangeBankMobileAdvance 修改银行预留手机推进 weibopay服务名称：change_bank_mobile_advance
// param: ticket(上一个接口获得),手机验证码
// return: 卡Id,银行卡是否验证
func ChangeBankMobileAdvance(ticket, validCode string) (string, string, error) {
	data := initBaseParam()
	data["service"] = "change_bank_mobile_advance"
	data["ticket"] = strings.TrimSpace(ticket)
	data["valid_code"] = strings.TrimSpace(validCode)
	// data["extend_param"] = ""
	rs, err := Request(&data, UserMode)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	rsMap, err := checkResponseCode(rs)
	if err != nil {
		return "", "", err
	}
	return rsMap["card_id"].(string), rsMap["is_verified"].(string), nil
}
