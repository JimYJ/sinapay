[![Build Status](https://travis-ci.org/JimYJ/sinapay.svg?branch=master)](https://travis-ci.org/JimYJ/sinapay)
[![Go Report Card](https://goreportcard.com/badge/github.com/JimYJ/sinapay)](https://goreportcard.com/report/github.com/JimYJ/sinapay)
 [![GoDoc](https://godoc.org/github.com/JimYJ/sinapay?status.svg)](https://www.godoc.org/github.com/JimYJ/sinapay)


# sinapay(weibopay)新浪支付 SDK 1.2
SDK for sinapay(weibopay) API v1.2<br>
根据sinapay 1.2版本接口开发

# How to get

```go
go get -u github.com/JimYJ/sinapay
```

# import
```go
import "github.com/JimYJ/sinapay"

```


# Usage:

## init 
```go
sinapay.InitSinaPay(pid, pubPEM, privPEM, sinapayPubPEM)
sinapay.TestMode() //use sinapay test URl
sinapay.DebugMode() //if you want
```
**创建激活会员**  weibopay服务名称：create_activate_member<br>
**param:** 用户请求IP,用户ID,用户账户标识类型:UID<br>
```go
	sinapay.CreateActiveMember(userIP, userID, identityType)
```

**设置用户实名信息**  weibopay服务名称：set_real_name<br>
**param:** 用户ID,真实姓名,身份证号,用户请求IP,用户账户标识类型:UID<br>
```go
	sinapay.SetRealName(userID, realname, IDNumber, userIP string, identityType int)
```

**设置支付密码**  weibopay服务名称：set_pay_password<br>
**param:** 用户ID,委托扣款展示方式(可空),同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC:返回PC页面,用户账户标识类型:UID<br>
**return:** 转跳页面<br>
```go
	sinapay.SetPayPassword(userID, withholdParam, returnURL, notifyURL string, mode, identityType int)
```

**修改支付密码**  weibopay服务名称：modify_pay_password
**param:** 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
**return:** 转跳页面
```go
	sinapay.ModifyPayPassword(userID, returnURL, notifyURL string, mode, identityType int)
```

**找回支付密码**  weibopay服务名称：find_pay_password
**param:** 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
**return:** 转跳页面
```go
	sinapay.FindPayPass(userID, returnURL, notifyURL string, mode, identityType int)
```

**检测是否已设置支付密码**  weibopay服务名称： query_is_set_pay_password
**param:** 用户ID,同步回跳页面(可空),异步通知接口(可空),用户账户标识类型:UID
**return:** 转跳页面
```go
	sinapay.QueryisSetPayPassword(userID, returnURL, notifyURL string, identityType int)
```

**绑定银行卡**  weibopay服务名称： binding_bank_card 默认SIGN模式
**param:** 用户ID,用户IP,银行编号,银行卡号,账户姓名,账户绑定手机,省份,城市,用户账户标识类型:UID
**return:** 卡ID(SIGN模式不返回),是否验证银行信息(SIGN模式不返回),绑卡推进ticket
```go
	sinapay.BindingBankCard(userID, userIP, bankCode, bankAccountNo, accounName, phone, province, city string, identityType int)
```

**绑定银行卡推进**  weibopay服务名称：binding_bank_card_advance
**param:** 绑定银行卡返回的ticket(SIGN模式),短信验证码,用户IP
**return:** 卡ID,银行卡信息验证状态
```go
	sinapay.BindingBankCardAdvance(ticket, validCode, userIP string)
```

**解绑银行卡**  weibopay服务名称：unbinding_bank_card
**param:** 用户ID,用户IP,卡ID,用户账户标识类型:UID
**return:** 解绑推进ticket
```go
	sinapay.UnbindingBankCard(userID, userIP, cardID string, identityType int)
```
**解绑银行卡推进**  weibopay服务名称：unbinding_bank_card_advance
**param:** 用户ID,用户IP,解绑银行卡返回的ticket,短信验证码,用户账户标识类型:UID
```go
	sinapay.UnbindingBankCardAdvance(userID, userIP, ticket, validCode string, identityType int)
```

**查询银行卡**  weibopay服务名称：query_bank_card
**param:** 用户ID,卡ID,用户账户标识类型:UID
**return:** 卡列表
```go
	sinapay.QueryBankCard(userID, cardID string, identityType int)
```

**查询余额/基金份额**  weibopay服务名称：query_balance
**param:** 用户ID,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 ,用户账户标识类型:UID,MemberID,Email,Mobile
**return:** 余额/基金份额,可用余额/基金份额,存钱罐收益(非查询存钱罐时为nil)
```go
	sinapay.QueryBalance(userID string, accountType, identityType int)
```

**查询收支明细**  weibopay服务名称：query_account_details
**param:** 用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 ,用户账户标识类型:UID,MemberID,Email,Mobile
**return:** 参数列表,收支明细列表
```go
	sinapay.QueryAccountDetails(userID, startTime, endTime, pageNo, pageSize string, accountType, identityType int) 
```

**冻结余额**  weibopay服务名称：balance_freeze
**param:** 用户ID,用户IP,摘要,金额,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 ,用户账户标识类型:UID
**return:** 冻结单号(查询状态和解冻用)
```go
	sinapay.BalanceFreeze(userID, userIP, summary, amount string, accountType, identityType int)
```

**解冻余额**  weibopay服务名称：balance_unfreeze
**param:用户ID,用户IP,原冻结单号,摘要,金额(为空表示全额解冻),用户账户标识类型:UID
**return:** 解冻单号(查询状态用)
```go
	sinapay.BalanceUnfreeze(userID, userIP, outFreezeNo, summary, amount string, identityType int)
```

**查询冻结解冻结果**  weibopay服务名称：query_ctrl_result
**param:** 冻结解冻订单号
**return:** 是否成功，失败原因(或请求接口报错)
```go
	sinapay.QueryCtrlResult(outCtrlNo string)
```

**查询企业会员信息**  weibopay服务名称：query_member_infos
**param:** 用户ID,用户账户标识类型:UID
**return:** 查询结果
```go
	sinapay.QueryMemberInfos(userID string, identityType int)
```

**查询企业会员审核结果**  weibopay服务名称：query_audit_result
**param:** 用户ID,用户账户标识类型:UID
**return:** 是否成功，失败原因(或请求接口报错)
```go
	sinapay.QueryAuditResult(userID string, identityType int)
```

**查询中间账户**  weibopay服务名称：query_middle_account
**param:** 外部业务码
**return:** 查询结果列表
```go
	sinapay.QueryMiddleAccount(outTradeCode int)
```

**修改认证手机**  weibopay服务名称：modify_verify_mobile
**param:** 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
**return:** 转跳页面
```go
	sinapay.ModifyVerifyMobile(userID, notifyURL, returnURL string, mode, identityType int)
```

**修改认证手机**  weibopay服务名称：find_verify_mobile
**param:** 用户ID,同步回跳页面(可空),异步通知接口(可空),返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID
**return:** 转跳页面
```go
	sinapay.FindVerifyMobile(userID, notifyURL, returnURL string, mode, identityType int)
```

**修改银行预留手机**  weibopay服务名称：change_bank_mobile
**param:** 用户ID,卡ID,手机号,用户账户标识类型:UID
**return:** ticket(推进接口用)
```go
	sinapay.ChangeBankMobile(userID, cardID, phone string, identityType int)
```

**修改银行预留手机推进**  weibopay服务名称：change_bank_mobile_advance
**param:** ticket(上一个接口获得),手机验证码
**return:** 卡Id,银行卡是否验证
```go
	sinapay.ChangeBankMobileAdvance(ticket, validCode string)
```

**创建托管代收交易**  weibopay服务名称：create_hosting_collect_trade
**param:交易订单号,摘要,标的号,付款用户ID,付款用户IP,卡属性,卡类型,金额,外部业务码,是否失败重付,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,用户账户标识类型:UID,是否代付冻结
**return:** 响应参数:交易订单号,交易状态,支付状态,ticket,转跳URL
```go
	sinapay.CreateHostingCollectTrade(tradeID, summary, goodsID, userID, userIP, cardAttr, cardType, amount string, outTradeCode, isRepay, mode, identityType int, isFreeze bool)
```

**创建托管代付交易**  weibopay服务名称：create_single_hosting_pay_trade
**param:** 交易订单号,摘要,标的号,付款用户ID,收款用户ID,付款用户IP,金额,备注,付款用户标识类型，收款用户标识类型:UID,MemberID,Email,Mobile,付款账户类型,收款账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户，代收分账列表:分账不可超过10笔,分账信息中的付款人必须为收款信息中的收款人，或分账信息中的所有收款人
**return:** 交易订单号,交易状态
```go
	sinapay.CreateSingleHostingPayTrade(tradeID, summary, goodsID, payerID, payeeID, userIP, amount, remarks string, payerIdentityType, payeeIdentityType, payerAccountType, payeeAccountType, outTradeCode int, splitList []map[string]string)
```
**托管交易支付**  weibopay服务名称：pay_hosting_trade
**param:** 支付订单号,付款用户IP,卡属性,卡类型,金额,交易订单号列表
**return:** 响应参数:支付订单号,支付状态,ticket,转跳URL
```go
	sinapay.PayHostingTrade(outPayNo, userIP, cardAttr, cardType, amount string, list []string)
```

**查询支付结果**  weibopay服务名称：query_pay_result
**param:** 支付订单号
**return:支付订单号，支付状态
```go
	sinapay.QueryPayResult(outPayNo string)
```

**托管交易查询**  weibopay服务名称：query_hosting_trade
**param:** 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,用户标识类型:UID,MemberID,Email,Mobile
**return:** 响应参数列表,交易记录列表
```go
	sinapay.QueryHostingTrade(tradeID, userID, startTime, endTime, pageNo, pageSize string, identityType int)
```

**托管退款**  weibopay服务名称：create_hosting_refund
**param:** 交易订单号,要退款的交易订单号,退款金额(可部分退款),摘要,用户IP
**return:** 交易订单号,退款状态
```go
	sinapay.CreateHostingRefund(tradeID, origtradeID, amount, summary, userIP string)
```

**托管退款查询**  weibopay服务名称：query_hosting_refund
**param:** 退款订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,用户标识类型:UID,MemberID,Email,Mobile
**return:** 响应参数列表,交易记录列表
```go
	sinapay.QueryHostingRefund(tradeID, userID, startTime, endTime, pageNo, pageSize string, identityType int)
```

**托管充值**  weibopay服务名称：create_hosting_deposit
**param:** 交易订单号,摘要,用户ID,用户IP,金额,用户手续费(可空),卡属性,卡类型,银行代码,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面,是否代付冻结
**return:** 交易订单号,充值状态,线下支付收款单位,线下支付收款账户,线下支付收款账号开户行,线下支付收款备注,收银台重定向地址
```go
	sinapay.CreateHostingDeposit(tradeID, summary, userID, userIP, amount, userFee, cardAttr, cardtype, bankCode string, accountType, mode, identityType int)
```

**托管充值查询** weibopay服务名称：query_hosting_deposit
**param:** 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,用户标识类型:UID,MemberID,Email,Mobile
**return:** 响应参数列表,交易记录列表
```go
	sinapay.QueryHostingDeposit(tradeID, userID, startTime, endTime, pageNo, pageSize string, accountType, identityType int)
```

**托管提现** weibopay服务名称：create_hosting_withdraw
**param:** 交易订单号,摘要,金额,用户ID,用户IP,用户手续费(可空),卡ID,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,用户标识类型:UID,MemberID,Email,Mobile 提现类型:Fast快速,General普通 返回页面类型 RedirectURLMobile:返回移动页面,RedirectURLPC返回PC页面 提现模式:ture安全模式,转跳收银台操作
**return:** 响应参数:交易订单号,提现状态,收银台重定向地址
```go
	sinapay.CreateHostingWithdraw(tradeID, summary, amount, userID, userIP, userFee, cardID string, accountType, identityType, paytoType, mode int, withdrawMode bool)
```

**托管提现查询** weibopay服务名称：query_hosting_withdraw
**param:** 交易订单号,用户ID,开始时间,结束时间(格式2006-01-02 15:04:05),页数,每页记录数,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户,用户标识类型:UID,MemberID,Email,Mobile
**return:** 响应参数列表,交易记录列表
```go
	sinapay.QueryHostingWithdraw(tradeID, userID, startTime, endTime, pageNo, pageSize string, accountType, identityType int)
```

**创建单笔代付到提现卡交易**  weibopay服务名称：create_single_hosting_pay_to_card_trade
**param:** 交易订单号,摘要,金额,用户ID,用户IP,卡ID,标的号,账户类型: BASIC基本户 ENSURE保证金户 RESERVE准备金 SAVING_POT存钱罐 BANK银行账户 用户标识类型:UID,MemberID,Email,Mobile 提现类型:Fast快速,General普通 外部业务码
**return:** 响应参数:交易订单号,提现状态
```go
	sinapay.CreateSingleHostingPaytoCardTrade(tradeID, summary, amount, userID, userIP, cardID, goodsID string, identityType, paytoType, outTradeCode int)
```

**代收完成**  weibopay服务名称：finish_pre_auth_trade
**param:** 请求交易号,用户IP,交易列表
```go
	sinapay.FinishPreAuthTrade(outRequestNo, userIP string, TradeList []map[string]string)
```

**代收撤销**  weibopay服务名称：cancel_pre_auth_trade
**param:** 请求交易号,交易列表
```go
	sinapay.CancelPreAuthTrade(outRequestNo string, TradeList []map[string]string)
```





