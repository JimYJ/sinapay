# sinapay(weibopay)新浪支付 SDK 1.2
SDK for sinapay(weibopay) API v1.2(WIP)<br>
根据sinapay 1.2版本接口开发(正在开发中...)

# How to get

```
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
sinapay.TestMode() //if you want to test
```

```go
    sinapay.CreateActiveMember("127.0.0.1", "0e5626d142cb4c46995d35f9a7cc1b09")
	sinapay.SetRealName("0e5626d142cb4c46995d35f9a7cc1b09", "姓名", "身份证号", "127.0.0.1")
	sinapay.SetPayPassword("0e5626d142cb4c46995d35f9a7cc1b09", "", "", "", sinapay.RedirectURLMobile)
	sinapay.ModifyPayPassword("0e5626d142cb4c46995d35f9a7cc1b09", sinapay.RedirectURLMobile)
	sinapay.FindPayPass("0e5626d142cb4c46995d35f9a7cc1b09", sinapay.RedirectURLMobile)
	sinapay.QueryisSetPayPassword("0e5626d142cb4c46995d35f9a7cc1b09")
	ticket, err := sinapay.BindingBankCard("0e5626d142cb4c46995d35f9a7cc1b09",
		"127.0.0.1", "CCB", "身份证号",
		"姓名", "13000000000", "浙江省", "杭州市")
	ticket, err)
	sinapay.BindingBankCardAdvance("6f7eb93c96084f419bc110c2bcdd2003", "884316", "127.0.0.1")
	sinapay.UnbindingBankCard("0e5626d142cb4c46995d35f9a7cc1b09", "127.0.0.1", "新浪卡ID")
	sinapay.UnbindingBankCardAdvance("0e5626d142cb4c46995d35f9a7cc1b09", "127.0.0.1", "8f3dd96a40fc48f1a4cb17492a282516", "641423"))
	sinapay.QueryBankCard("0e5626d142cb4c46995d35f9a7cc1b09", "")
	sinapay.QueryBalance("200004595271", sinapay.Basic, sinapay.MemberAccount)
	sinapay.QueryAccountDetails("200004595271", "", "", "", "", sinapay.Basic, sinapay.MemberAccount)
	sinapay.BalanceFreeze("0e5626d142cb4c46995d35f9a7cc1b09", "127.0.0.1", "summary", "800", sinapay.Basic)
	sinapay.BalanceUnfreeze("0e5626d142cb4c46995d35f9a7cc1b09", "127.0.0.1", "drfg", "summary", "800")
	sinapay.QueryCtrlResult("outCtrlNo")
	sinapay.QueryMemberInfos("200004595271")
	sinapay.QueryAuditResult("0e5626d142cb4c46995d35f9a7cc1b09")
	sinapay.QueryMiddleAccount(sinapay.CollectionInvestment)
	sinapay.ModifyVerifyMobile("0e5626d142cb4c46995d35f9a7cc1b09", "", "", sinapay.RedirectURLMobile)
	sinapay.FindVerifyMobile("0e5626d142cb4c46995d35f9a7cc1b09", "", "", sinapay.RedirectURLMobile)
```