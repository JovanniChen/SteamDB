package main

import (
	"fmt"

	"example.com/m/v2/Steam/Dao"
)

func main() {
	d := Dao.New("")
	code, _ := d.GetTokenCode("ELYLgQOKGniuFg3tVxajtstv6kM=")
	fmt.Println(code)
	// err := d.Login("rgckq82191", "vxlu26493E", "")
	err := d.Login("xuszv439", "kS6llWROUvxh", "ELYLgQOKGniuFg3tVxajtstv6kM=")

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
