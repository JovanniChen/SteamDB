package main

import (
	"example.com/m/v2/Steam/Dao"
	"fmt"
)

func main() {
	d := Dao.New("")
	//code, _ := d.GetTokenCode("Q9h1SqNq+4jk5M+P6UCykXdoFmU=")
	//fmt.Println(code)
	err := d.Login("xwabg995", "5uBKjPMRxImK", "Q9h1SqNq+4jk5M+P6UCykXdoFmU=")
	//err := d.Login("yl32940", "32940abcd", "")

	//err := d.SetLoginInfo("yl32940", "32940abcd", "eyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwNV8yNjdCOEY2NF8xMDA0RiIsICJzdWIiOiAiNzY1NjExOTgxMzI4MjQ0NTgiLCAiYXVkIjogWyAid2ViIiBdLCAiZXhwIjogMTc1MDQ3NjE1MywgIm5iZiI6IDE3NDE3NDk0MDYsICJpYXQiOiAxNzUwMzg5NDA2LCAianRpIjogIjAwMDZfMjY3RTI0ODRfN0EwQTkiLCAib2F0IjogMTc1MDM4OTQwNiwgInJ0X2V4cCI6IDE3Njg0OTI5NDIsICJwZXIiOiAwLCAiaXBfc3ViamVjdCI6ICIxMTIuMzIuMjIuMjAwIiwgImlwX2NvbmZpcm1lciI6ICIxMTIuMzIuMjIuMjAwIiB9.T3YnJMoPHZ9wultuCChJAqr46Swe4dKV8gSJhx4duyBfUA-82pAbnLqXacnlrcUS30wumJDmohx2iXWhKCTTBg",
	//	"CN", `{
	// "checkout.steampowered.com" : {
	//   "steamLoginSecure" : "76561198132824458%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAxMl8yNjdCOEY1N18yMDAxMiIsICJzdWIiOiAiNzY1NjExOTgxMzI4MjQ0NTgiLCAiYXVkIjogWyAid2ViOmNoZWNrb3V0IiBdLCAiZXhwIjogMTc1MDQ3NjM1MywgIm5iZiI6IDE3NDE3NDk2MzMsICJpYXQiOiAxNzUwMzg5NjMzLCAianRpIjogIjAwMDZfMjY3RTI0ODRfODU3NjAiLCAib2F0IjogMTc1MDM4OTYzMiwgInJ0X2V4cCI6IDE3Njg2Njc1NTAsICJwZXIiOiAwLCAiaXBfc3ViamVjdCI6ICIxMTIuMzIuMjIuMjAwIiwgImlwX2NvbmZpcm1lciI6ICIxMTIuMzIuMjIuMjAwIiB9.5bsT90ysuKZF1WuGTIIQ7gZ5QSRr6DvydmgTWStnyg_zJ04XsdtVPchTmcde_LiC20WdfjUnpI46ESwKGUEDBA",
	//   "sessionid" : "df862cc9e2c447c4b9ff4d00"
	// },
	// "help.steampowered.com" : {
	//   "steamLoginSecure" : "76561198132824458%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAxMl8yNjdCOEY1N18yMDAxMiIsICJzdWIiOiAiNzY1NjExOTgxMzI4MjQ0NTgiLCAiYXVkIjogWyAid2ViOmhlbHAiIF0sICJleHAiOiAxNzUwNDc2NjAwLCAibmJmIjogMTc0MTc0OTYzMywgImlhdCI6IDE3NTAzODk2MzMsICJqdGkiOiAiMDAwNl8yNjdFMjQ4NF84NTc1RiIsICJvYXQiOiAxNzUwMzg5NjMyLCAicnRfZXhwIjogMTc2ODY2NzU1MCwgInBlciI6IDAsICJpcF9zdWJqZWN0IjogIjExMi4zMi4yMi4yMDAiLCAiaXBfY29uZmlybWVyIjogIjExMi4zMi4yMi4yMDAiIH0.lhqfYLgYWn-CuG_p3G_fEQ9rXO3zv58wElymrA7Xan7q7c7t68zAala9gxOPsBIOiXGz-8NJC5IOvcKqVdkRCg",
	//   "sessionid" : "8e640ba6e6cfa234fd0e2398"
	// },
	// "steam.tv" : {
	//   "steamLoginSecure" : "76561198132824458%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAxMl8yNjdCOEY1N18yMDAxMiIsICJzdWIiOiAiNzY1NjExOTgxMzI4MjQ0NTgiLCAiYXVkIjogWyAid2ViOnN0ZWFtdHYiIF0sICJleHAiOiAxNzUwNDc3NzI0LCAibmJmIjogMTc0MTc0OTYzMywgImlhdCI6IDE3NTAzODk2MzMsICJqdGkiOiAiMDAwNl8yNjdFMjQ4NF84NTc2MSIsICJvYXQiOiAxNzUwMzg5NjMyLCAicnRfZXhwIjogMTc2ODY2NzU1MCwgInBlciI6IDAsICJpcF9zdWJqZWN0IjogIjExMi4zMi4yMi4yMDAiLCAiaXBfY29uZmlybWVyIjogIjExMi4zMi4yMi4yMDAiIH0.slhCe11HMxDrg-28RnSYIar0DGtVQCTEIzvajPIL7CrNZmWtKlk5CGDrttish9fL20ZmNrIFOdRQuh_ZKfdADg",
	//   "sessionid" : "40fd860ec7171b9128a7ed77"
	// },
	// "steamcommunity.com" : {
	//   "steamLoginSecure" : "76561198132824458%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAxMl8yNjdCOEY1N18yMDAxMiIsICJzdWIiOiAiNzY1NjExOTgxMzI4MjQ0NTgiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NTA0NzYxMzEsICJuYmYiOiAxNzQxNzQ5NjMzLCAiaWF0IjogMTc1MDM4OTYzMywgImp0aSI6ICIwMDA2XzI2N0UyNDg0Xzg1NzVFIiwgIm9hdCI6IDE3NTAzODk2MzIsICJydF9leHAiOiAxNzY4NjY3NTUwLCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiMTEyLjMyLjIyLjIwMCIsICJpcF9jb25maXJtZXIiOiAiMTEyLjMyLjIyLjIwMCIgfQ.oIJgmyoZhjzRXjzsUibvTsycfKi2jxjPOoGgUTHnt7eJC6KGh_yKuz2G9g2nB391sTu_kTXCH2sOBZZNeMmQBQ",
	//   "sessionid" : "8b64597fbeed99fc69175eaf"
	// },
	// "store.steampowered.com" : {
	//   "steamLoginSecure" : "76561198132824458%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAxMl8yNjdCOEY1N18yMDAxMiIsICJzdWIiOiAiNzY1NjExOTgxMzI4MjQ0NTgiLCAiYXVkIjogWyAid2ViOnN0b3JlIiBdLCAiZXhwIjogMTc1MDQ3NzUwMywgIm5iZiI6IDE3NDE3NDk2MzMsICJpYXQiOiAxNzUwMzg5NjMzLCAianRpIjogIjAwMDZfMjY3RTI0ODRfODU3NUQiLCAib2F0IjogMTc1MDM4OTYzMiwgInJ0X2V4cCI6IDE3Njg2Njc1NTAsICJwZXIiOiAwLCAiaXBfc3ViamVjdCI6ICIxMTIuMzIuMjIuMjAwIiwgImlwX2NvbmZpcm1lciI6ICIxMTIuMzIuMjIuMjAwIiB9.YbToccFB01YPHRFcjxc6N7PbN_Og4a4JTxX7Gj2I3uAmTrRDx5nXyjrJjAdn9kfx_IdBnLx2FApkbguGgSf2CA",
	//   "sessionid" : "5e48da3db0da69abf1b8e6cd"
	// }
	//}`)
	if err != nil {
		fmt.Println(err)
		return
	}
	str, err := d.GetUserCookies()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = d.UserInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = d.SetLanguage("schinese")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	fmt.Println(string(str))

}
