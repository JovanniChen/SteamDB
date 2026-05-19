// Steam数据库操作主程序
// 本程序用于连接Steam平台，进行用户登录、获取令牌代码和添加反应等操作
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/JovanniChen/SteamDB/Steam"
)

// var config *Steam.Config = Steam.NewConfig("jiochenhao:JSJ3ZAN9tc@198.64.246.6:50100")

// var config *Steam.Config = Steam.NewConfig("127.0.0.1:7893")

// var config *Steam.Config = Steam.NewConfig("jiochenhao:JSJ3ZAN9tc@198.64.247.190:50100")

var config *Steam.Config = Steam.DefaultConfig()

func main() {
	accountIndex := flag.Int("account", 3, "account index used by case functions")
	caseName := flag.String("case", "", "case name, e.g. TestLogin")
	listCases := flag.Bool("list-cases", false, "print all available case names")
	gameID := flag.Int("game-id", 1879330, "game id for TestGetGameUpdateInofs")
	flag.Parse()

	if *listCases {
		printCaseNames()
		return
	}

	if *caseName == "" {
		fmt.Println("未指定 --case，使用 --list-cases 查看可用 case")
		return
	}

	if *caseName == "TestGetGameUpdateInofs" {
		TestGetGameUpdateInofs(*gameID)
		return
	}

	runner, ok := caseRegistry[*caseName]
	if !ok {
		fmt.Printf("未知 case: %s\n\n", *caseName)
		printCaseNames()
		os.Exit(1)
	}

	runner(*accountIndex)
}

func printCaseNames() {
	names := make([]string, 0, len(caseRegistry)+1)
	for name := range caseRegistry {
		names = append(names, name)
	}
	names = append(names, "TestGetGameUpdateInofs")
	sort.Strings(names)

	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	fmt.Println("可用 case 及调用方式:")
	fmt.Println("--------------------------------------------------------------------------")
	for i, name := range names {
		cmd := fmt.Sprintf("go run . --case %s --account 3", name)
		if name == "TestGetGameUpdateInofs" {
			cmd = fmt.Sprintf("go run . --case %s --game-id 1879330", name)
		}
		fmt.Printf("%2d. %-*s  %s\n", i+1, maxLen, name, cmd)
	}
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("提示: 将命令中的 --account 3 替换为你的账号索引")
}
