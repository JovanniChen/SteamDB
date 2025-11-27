package Dao

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Utils"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// GetFriendInfoByLink 通过好友链接获取好友信息
// 支持 s.team 短链接，会自动处理重定向到 steamcommunity.com
// 示例: https://s.team/p/chbn-qbdd/jbccfkfw
func (d *Dao) GetFriendInfoByLink(link string) (*Model.FriendInfo, string, error) {
	friendInfo := &Model.FriendInfo{}
	parts := strings.Split(link, "/")
	if len(parts) < 6 {
		return friendInfo, "", Errors.ErrInvalidFriendLink
	}

	inviteToken := parts[len(parts)-1] // 获取最后一个部分作为token

	// 对于 s.team 短链接，预先将 steamcommunity.com 的 cookie 存入 CookieJar
	// 这样当重定向到 steamcommunity.com 时，cookie 会自动被使用
	if strings.Contains(link, "s.team") {
		communityURL, err := url.Parse("https://steamcommunity.com")
		if err != nil {
			return friendInfo, "", err
		}

		ck := d.GetCookiesString("https://steamcommunity.com")
		if ck != nil {
			var cookies []*http.Cookie

			// 添加登录 cookie
			cookies = append(cookies, &http.Cookie{
				Name:   "steamLoginSecure",
				Value:  ck.SteamLoginSecure,
				Domain: ".steamcommunity.com", // 使用 . 前缀以支持所有子域名
				Path:   "/",
			})

			// 添加会话 cookie
			cookies = append(cookies, &http.Cookie{
				Name:   "sessionid",
				Value:  ck.SessionId,
				Domain: ".steamcommunity.com",
				Path:   "/",
			})

			// 添加语言 cookie（如果设置了）
			if d.credentials.Language != "" {
				cookies = append(cookies, &http.Cookie{
					Name:   "Steam_Language",
					Value:  d.credentials.Language,
					Domain: ".steamcommunity.com",
					Path:   "/",
				})
			}

			// 将 cookie 设置到 CookieJar 中，重定向时会自动使用
			d.httpCli.Jar.SetCookies(communityURL, cookies)
		}
	}

	// 创建请求
	req, err := d.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return friendInfo, "", err
	}

	// 发送请求，重定向会自动处理，cookie 会从 jar 中自动获取
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return friendInfo, "", err
	}
	defer resp.Body.Close()

	// 检查最终重定向后的状态码
	if resp.StatusCode != 200 {
		return friendInfo, "", Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return friendInfo, "", err
	}

	result := d.ParseFriendLinkHTML(string(body))

	if result.Status != Model.FriendLinkStatusSuccess {
		return friendInfo, "", fmt.Errorf("%s", result.Msg)
	}

	fmt.Printf("状态: %d (%s)\n", result.Status, result.Msg)
	if result.Data != nil {
		fmt.Printf("用户昵称: %s\n", result.Data.PersonName)
		fmt.Printf("好友代码: %d\n", result.Data.FriendCode)
		fmt.Printf("会话ID: %s\n", result.Data.SessionID)
		fmt.Printf("滥用ID: %s\n", result.Data.AbuseID)
	}
	return result.Data, inviteToken, nil
}

func (d *Dao) CheckFriendStatus(link string) error {
	friendInfo, _, err := d.GetFriendInfoByLink(link)
	if err != nil {
		return err
	}

	params := Param.Params{}
	params.SetString("steamids", friendInfo.AbuseID)

	req, err := d.Request(http.MethodGet, Constants.Ajaxresolveusers+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	// 发送请求，重定向会自动处理，cookie 会从 jar 中自动获取
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	// 	[
	//     {
	//         "steamid": "76561198313222178",
	//         "accountid": 352956450,
	//         "persona_name": "Call me Daddy",
	//         "avatar_url": "e6b8bd224ec0d82e391d79b6e4baedc1a7526f44",
	//         "profile_url": "jovannichen",
	//         "persona_state": 0,
	//         "city": "",
	//         "state": "",
	//         "country": "",
	//         "real_name": "",
	//         "is_friend": false,
	//         "friends_in_common": 0
	//     }
	// ]
	var result []map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		fmt.Println("没有找到好友信息")
	}

	fmt.Println(result[0]["accountid"])
	fmt.Println(result[0]["steamid"])

	return nil
}

func (d *Dao) AddFriendByLink(link string) error {
	friendInfo, inviteToken, err := d.GetFriendInfoByLink(link)
	if err != nil {
		return err
	}

	sessionid := d.GetLoginCookies()["steamcommunity.com"].SessionId

	params := Param.Params{}
	params.SetString("invite_token", inviteToken)
	params.SetString("sessionid", sessionid)
	params.SetString("steamid_user", friendInfo.AbuseID)

	req, err := d.Request(http.MethodGet, Constants.AddFriendByLink+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	// 发送请求，重定向会自动处理，cookie 会从 jar 中自动获取
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return Errors.ErrAddFriendFailed
	}
	result := &Model.AddFriendByLinkResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Success != 1 {
		return Errors.ErrAddFriendFailed
	}

	return nil
}

func (d *Dao) AddFriendByFriendCode(friendCode uint32) error {
	sessionid := d.GetLoginCookies()["steamcommunity.com"].SessionId

	// 创建 multipart/form-data 请求体
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加表单字段
	writer.WriteField("steamid", strconv.FormatUint(Utils.FriendCodeToSteamID64(friendCode), 10))
	writer.WriteField("sessionID", sessionid)
	writer.WriteField("accept_invite", "0")

	// 关闭 writer 以完成 multipart 消息
	writer.Close()

	// 创建请求
	req, err := d.Request(http.MethodPost, Constants.AddFriendAjax, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Origin", Constants.CommunityOrigin)
	req.Header.Set("Referer", Constants.CommunityOrigin+"/profiles/"+strconv.FormatUint(d.GetSteamID(), 10)+"/friends/add")

	// 设置正确的 Content-Type（包含 boundary）
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应是否被 gzip 压缩
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	// 读取响应内容
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应内容:", string(body))

	if resp.StatusCode != 200 {
		return fmt.Errorf("添加好友失败: %s", string(body))
	}

	// 解析响应 JSON
	result := &Model.AddFriendByCodeResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查添加好友是否成功
	if !result.Success {
		return fmt.Errorf("添加好友失败: %s", result.ErrorText)
	}

	return nil
}

func (d *Dao) RemoveFriend(steamID uint64) error {
	sessionid := d.GetLoginCookies()["steamcommunity.com"].SessionId

	params := Param.Params{}
	params.SetString("sessionID", sessionid)
	params.SetString("steamid", strconv.FormatUint(steamID, 10))

	req, err := d.Request(http.MethodPost, Constants.RemoveFriendAjax, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应是否被 gzip 压缩
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	fmt.Println("状态码:", resp.StatusCode)
	fmt.Println("响应内容:", string(body))

	return nil
}

// ParseFriendLinkHTML 解析好友链接页面的 HTML 内容
// 从 Steam 好友邀请页面中提取用户信息和链接状态
// 参数：htmlContent - HTML 页面内容
// 返回值：解析结果，包含状态、消息和用户信息
func (d *Dao) ParseFriendLinkHTML(htmlContent string) *Model.FriendLinkParseResult {
	// 解析 HTML 文档
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusFailed,
			Msg:    fmt.Sprintf("HTML 解析失败: %v", err),
		}
	}

	userInfo := &Model.FriendInfo{}

	// 1. 提取用户昵称
	// XPath: //div[@class="persona_name"]/span[@class="actual_persona_name"]/text()
	personNameNode := htmlquery.FindOne(doc, `//div[@class="persona_name"]/span[@class="actual_persona_name"]`)
	if personNameNode == nil {
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusFailed,
			Msg:    "无法找到指定的个人资料（请检查链接是否可用）",
		}
	}
	userInfo.PersonName = strings.TrimSpace(htmlquery.InnerText(personNameNode))

	// 2. 提取好友代码（miniprofile ID）
	// 尝试多种可能的选择器
	var avatarNode *html.Node

	// 尝试1: 标准的 playerAvatar
	avatarNode = htmlquery.FindOne(doc, `//div[contains(@class,'playerAvatar')]`)
	if avatarNode == nil {
		// 尝试2: profile_avatar 类
		avatarNode = htmlquery.FindOne(doc, `//div[contains(@class,'profile_avatar')]`)
	}
	if avatarNode == nil {
		// 尝试3: 直接查找任何有 data-miniprofile 属性的元素
		avatarNode = htmlquery.FindOne(doc, `//*[@data-miniprofile]`)
	}

	if avatarNode == nil {
		// 输出所有包含 avatar 的 div
		avatarNodes := htmlquery.Find(doc, `//div[contains(@class,'avatar')]`)
		for i, node := range avatarNodes {
			if i < 3 { // 只打印前3个
				for _, attr := range node.Attr {
					if attr.Key == "class" {
						fmt.Printf("%s\n", attr.Val)
					}
				}
			}
		}

		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusFailed,
			Msg:    "好友链接页面结构发生变化(friendCode获取失败)，请联系开发人员。",
		}
	}

	miniProfileAttr := ""
	for _, attr := range avatarNode.Attr {
		if attr.Key == "data-miniprofile" {
			miniProfileAttr = attr.Val
			break
		}
	}

	// 如果父节点上没有 data-miniprofile，尝试在子元素中查找
	if miniProfileAttr == "" {
		// 查找该节点下所有带 data-miniprofile 的子元素
		childNode := htmlquery.FindOne(avatarNode, `.//*[@data-miniprofile]`)
		if childNode != nil {
			for _, attr := range childNode.Attr {
				if attr.Key == "data-miniprofile" {
					miniProfileAttr = attr.Val
					break
				}
			}
		}
	}

	// 如果还是没找到，尝试直接在整个文档中搜索
	if miniProfileAttr == "" {
		allMiniProfileNodes := htmlquery.Find(doc, `//*[@data-miniprofile]`)

		if len(allMiniProfileNodes) > 0 {
			// 使用第一个找到的
			for _, attr := range allMiniProfileNodes[0].Attr {
				if attr.Key == "data-miniprofile" {
					miniProfileAttr = attr.Val
					break
				}
			}
		}
	}

	if miniProfileAttr == "" {
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusFailed,
			Msg:    "好友链接页面结构发生变化(friendCode获取失败)，请联系开发人员。",
		}
	}

	friendCode, err := strconv.ParseUint(miniProfileAttr, 10, 64)
	if err != nil {
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusFailed,
			Msg:    fmt.Sprintf("friendCode 解析失败: %v", err),
		}
	}
	userInfo.FriendCode = friendCode

	// 3. 从表单中提取 sessionid 和 abuseID
	// XPath: //form[@id='abuseForm']/input
	inputNodes := htmlquery.Find(doc, `//form[@id='abuseForm']/input`)
	sessionIDFound := false
	abuseIDFound := false

	for _, inputNode := range inputNodes {
		var name, value string
		for _, attr := range inputNode.Attr {
			switch attr.Key {
			case "name":
				name = attr.Val
			case "value":
				value = attr.Val
			}
		}

		switch name {
		case "sessionid":
			userInfo.SessionID = value
			sessionIDFound = true
		case "abuseID":
			userInfo.AbuseID = value
			abuseIDFound = true
		}
	}

	// 如果关键字段未找到，返回仅包含信息的结果
	if !sessionIDFound || !abuseIDFound {
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusOnlyInfo,
			Msg:    "好友链接页面结构发生变化(链接有效判断失败)，请联系开发人员。",
			Data:   userInfo,
		}
	}

	// 4. 检查按钮状态判断链接类型
	// XPath: //div[@class='profile_header_actions']/a[1]/@href
	buttonNode := htmlquery.FindOne(doc, `//div[@class='profile_header_actions']/a[1]`)
	if buttonNode == nil {
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusOnlyInfo,
			Msg:    "好友链接页面结构发生变化(链接有效判断失败)，请联系开发人员。",
			Data:   userInfo,
		}
	}

	var hrefValue string
	for _, attr := range buttonNode.Attr {
		if attr.Key == "href" {
			hrefValue = attr.Val
			break
		}
	}

	// 判断链接状态
	if strings.Contains(hrefValue, "OpenFriendChat") {
		// 已经是好友
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusIsFriend,
			Msg:    "您关注了一份好友邀请 - 但你们已经是好友了。",
			Data:   userInfo,
		}
	}

	if strings.Contains(hrefValue, "AddFriend") {
		// 链接已过期
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusOnlyInfo,
			Msg:    "您关注的好友邀请已过期。",
			Data:   userInfo,
		}
	}

	if hrefValue != "#" {
		// 未知的按钮类型
		return &Model.FriendLinkParseResult{
			Status: Model.FriendLinkStatusOnlyInfo,
			Msg:    "好友链接页面结构发生变化(链接有效判断失败)，请联系开发人员。",
			Data:   userInfo,
		}
	}

	// 5. 链接有效，可以添加好友
	return &Model.FriendLinkParseResult{
		Status: Model.FriendLinkStatusSuccess,
		Msg:    "好友链接有效，可以添加好友。",
		Data:   userInfo,
	}
}
