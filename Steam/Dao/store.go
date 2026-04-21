package Dao

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Param"
)

func (d *Dao) GetPackageDetails(subID int) error {
	params := Param.Params{}
	params.SetString("packageids", strconv.Itoa(subID))
	params.SetString("cc", "cn")
	params.SetString("l", "schinese")

	// 创建请求
	req, err := d.NewRequest(http.MethodGet, Constants.PackageDetails+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	// 发送请求，重定向会自动处理，cookie 会从 jar 中自动获取
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查最终重定向后的状态码
	if resp.StatusCode != 200 {
		return errors.New("status != 200")
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
