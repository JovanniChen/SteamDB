// login.go - Steam登录相关功能实现
// 包含用户认证、RSA加密、令牌管理等核心登录逻辑
package Dao

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"example.com/m/v2/Steam/Model"
	"example.com/m/v2/Steam/Param"
	"example.com/m/v2/Steam/Protoc"
	"example.com/m/v2/Steam/Utils"
	"google.golang.org/protobuf/proto"
)

// LoginCookie Steam登录Cookie结构体
// 存储登录成功后的会话信息
type LoginCookie struct {
	SteamLoginSecure string `json:"steamLoginSecure"` // Steam安全登录令牌
	SessionId        string `json:"sessionid"`        // 会话ID
	//SteamLanguage    string `json:"Steam_Language"`    // Steam语言设置（已注释）
}

// Credentials 用户凭据结构体
// 存储用户的登录信息和认证状态
type Credentials struct {
	Password     string                    // 用户密码
	Username     string                    // 用户名
	Nickname     string                    // 用户昵称
	SteamID      uint64                    // Steam用户ID
	RSATimeStamp string                    // RSA时间戳，用于密码加密
	AccessToken  string                    // 访问令牌
	RefreshToken string                    // 刷新令牌
	Language     string                    // 语言偏好设置
	CountryCode  string                    // 国家代码
	LoginCookies map[string]*LoginCookie   // 各域名的登录Cookie映射
}

// AccessToken 获取访问令牌
// 返回当前有效的访问令牌，用于API调用认证
// 返回值：访问令牌字符串和可能的错误
func (d *Dao) AccessToken() (string, error) {
	if d.credentials.AccessToken == "" {
		return "", Errors.Error("未获取到")
	}
	return d.credentials.AccessToken, nil
}

// getRSA 获取Steam RSA公钥用于密码加密
// Steam使用RSA加密来保护用户密码在传输过程中的安全性
// 参数：username - 要登录的用户名
// 返回值：Steam公钥信息和可能的错误
func (d *Dao) getRSA(username string) (*Model.SteamPublicKey, error) {
	// 构建获取公钥的请求参数
	publicKeySend := &Protoc.GetPasswordRSAPublicKeySend{
		AccountName: username,
	}
	
	// 将请求参数序列化为protobuf格式
	data, err := proto.Marshal(publicKeySend)
	if err != nil {
		return nil, err
	}
	
	// 构建URL参数
	params := Param.Params{}
	params.SetString("origin", Steam.Origin)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))
	
	// 创建HTTP请求
	req, err := d.NewRequest("GET", Steam.GetPasswordRSAPublicKey+"?="+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}
	
	// 设置必要的请求头
	req.Header.Add("origin", Steam.Origin)
	req.Header.Set("referer", Steam.Origin+"/")
	
	// 发送请求并重试
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 检查响应状态码
	if resp.StatusCode != 200 {
		return nil, Errors.ResponseError(resp.StatusCode)
	}
	
	// 读取响应体
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	
	// 解析响应数据为protobuf格式
	keySendReceive := &Protoc.GetPasswordRSAPublicKeySendReceive{}
	err = proto.Unmarshal(buf.Bytes(), keySendReceive)
	if err != nil {
		return nil, err
	}
	
	// 构建Steam公钥对象并返回
	spk := new(Model.SteamPublicKey)
	spk.Success = true                               // 标记获取成功
	spk.Timestamp = keySendReceive.Timestamp         // RSA时间戳
	spk.PublicKeyMod = keySendReceive.PublickeyMod   // 公钥模数
	spk.PublicKeyExp = keySendReceive.PublickeyExp   // 公钥指数
	return spk, nil
}

// pollAuthSessionStatus 轮询身份验证会话状态
// 在登录过程中定期检查认证状态，等待用户完成双因素认证等步骤
// 参数：
//   clientId - 客户端ID
//   requestId - 请求ID字节数组
// 返回值：认证状态响应和可能的错误
func (d *Dao) pollAuthSessionStatus(clientId uint64, requestId []byte) (*Protoc.PollAuthSessionStatusReceive, error) {
	// 构建轮询请求数据
	loginData := &Protoc.PollAuthSessionStatusSend{
		ClientId:  clientId,
		RequestId: requestId,
	}
	
	// 序列化请求数据
	data, err := proto.Marshal(loginData)
	if err != nil {
		return nil, err
	}
	
	// 构建POST请求参数
	params := Param.Params{}
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))
	
	// 创建HTTP POST请求
	req, err := d.NewRequest("POST", Steam.PollAuthSessionStatus, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	// 发送请求并获取响应
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 检查响应状态码
	if resp.StatusCode != 200 {
		return nil, Errors.ResponseError(resp.StatusCode)
	}
	
	// 获取Steam的响应结果代码
	eresult := resp.Header.Get("x-eresult")
	result, _ := strconv.Atoi(eresult)
	
	// 根据结果代码处理不同情况
	switch result {
	case 1: // 成功状态
		credentialsReceive := &Protoc.PollAuthSessionStatusReceive{}
		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}
		err = proto.Unmarshal(buf.Bytes(), credentialsReceive)
		return credentialsReceive, nil
	}
	return nil, Errors.Unavailable()
}

// ajaxRefresh 刷新 仅用来获取steam ak_bmsc
func (d *Dao) ajaxRefresh() (*Model.RefreshResponse, error) {
	params := Param.Params{}
	params.SetString("redir", Steam.Origin+"/")
	req, err := d.NewRequest("POST", Steam.AjaxRefresh, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	m := map[string]string{}
	cookieArr := resp.Header["Set-Cookie"]
	for _, cookie := range cookieArr {
		cookies := strings.Split(cookie, ";")
		for _, item := range cookies {
			kv := strings.Split(item, "=")
			if len(kv) != 2 {
				continue
			}
			m[kv[0]] = kv[1]
		}
	}

	response := &Model.RefreshResponse{
		Cookie: m,
	}
	response.Success = true
	return response, nil
}

// finalizeLogin 登录最终阶段，返回steam登录各个域名所需要的信息
func (d *Dao) finalizeLogin(ak_bmsc, refreshToken, sessionid string) (*Model.FinalizeResponse, error) {
	params := Param.Params{}
	params.SetString("redir", "https://store.steampowered.com/login/?redir=&redir_ssl=1&snr=1_4_600__global-header")
	params.SetString("nonce", refreshToken)
	params.SetString("sessionid", sessionid)
	req, err := d.NewRequest("POST", Steam.FinalizeLogin, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("cookie", ak_bmsc)
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, Errors.ResponseError(resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	body := &Model.FinalizeResponse{}
	body.Success = true
	err = json.Unmarshal(buf.Bytes(), body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// AutoLogin 登录具体域名
func (d *Dao) AutoLogin(url, nonce, auth string, steamID uint64, reDir string) (*Model.CheckLoginResponse, error) {
	params := Param.Params{}
	params.SetString("nonce", nonce)
	params.SetString("auth", auth)
	params.SetInt64("steamID", int64(steamID))
	params.SetString("redir", reDir)
	req, err := d.NewRequest("POST", url, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", "")
	req.Header.Set("origin", Steam.Origin)
	req.Header.Set("referer", Steam.Origin+"/")
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, Errors.ResponseError(resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	body := &struct {
		Result   int `json:"result"`
		RtExpiry int `json:"rtExpiry"`
	}{}
	err = json.Unmarshal(buf.Bytes(), body)
	if err != nil {
		return nil, err
	}
	switch body.Result {
	case 1:
		cookieArr := resp.Header["Set-Cookie"]
		steamLoginSecure := ""
		sessionid := ""
		for _, cookie := range cookieArr {
			cookies := strings.Split(cookie, ";")
			for _, item := range cookies {
				kv := strings.Split(item, "=")
				if len(kv) != 2 {
					continue
				}
				key := kv[0]
				value := kv[1]
				if key == "steamLoginSecure" {
					steamLoginSecure = value
				}
				if key == "sessionid" {
					sessionid = value
				}
			}
		}
		if steamLoginSecure == "" {
			fmt.Println("steamLoginSecure is empty")
		}
		if sessionid == "" {
			sessionid = Utils.SafeHexString(12)
		}

		response := &Model.CheckLoginResponse{
			Url: url,
			Data: struct {
				SteamLoginSecure string `json:"steamLoginSecure"`
				SessionId        string `json:"sessionid"`
			}{SteamLoginSecure: steamLoginSecure, SessionId: sessionid},
		}
		response.Success = true
		response.Message = "成功加载"
		return response, nil
	}
	return nil, nil
}

// afterVerificationLogin 最后一步校验登录
func (d *Dao) afterVerificationLogin(clientId, steamId uint64, requestId []byte) error {
	sessionStatusReceive, err := d.pollAuthSessionStatus(clientId, requestId)
	if err != nil {
		return err
	}
	accessToken := sessionStatusReceive.AccessToken
	refreshToken := sessionStatusReceive.RefreshToken
	refresh, err := d.ajaxRefresh()
	if err != nil {
		return err
	}
	sessionid := Utils.SafeHexString(12)
	akBmsc := refresh.Cookie["ak_bmsc"]
	finalze, err := d.finalizeLogin(akBmsc, refreshToken, sessionid)
	if err != nil {
		return err
	}
	if !finalze.Success {
		return Errors.Unavailable()
	}
	wait := sync.WaitGroup{}
	wait.Add(len(finalze.TransferInfo))
	cookieList := make([]*Model.CheckLoginResponse, 0)
	for _, info := range finalze.TransferInfo {
		go func() {
			response, err := d.AutoLogin(info.Url, info.Params.Nonce, info.Params.Auth, steamId, "")
			if err == nil {
				if response.Success {
					cookieList = append(cookieList, response)
				}
			}
			wait.Done()
		}()
	}
	wait.Wait()
	d.credentials.SteamID = steamId
	d.credentials.AccessToken = accessToken
	d.credentials.RefreshToken = refreshToken
	d.credentials.LoginCookies = map[string]*LoginCookie{}
	for _, cookie := range cookieList {
		ur, _ := url.Parse(cookie.Url)

		d.credentials.LoginCookies[ur.Host] = &LoginCookie{
			SteamLoginSecure: cookie.Data.SteamLoginSecure,
			SessionId:        cookie.Data.SessionId,
		}
	}
	return nil
}

// submitVerificationCode 验证码提交
func (d *Dao) submitVerificationCode(clientId uint64, steamId uint64, confirmationType int32, code string) error {
	emailCode := &Protoc.EmailCode{
		ClientId: clientId,
		Steamid:  steamId,
		Code:     code,
		CodeType: confirmationType,
	}
	data, err := proto.Marshal(emailCode)
	if err != nil {
		return err
	}

	params := Param.Params{}
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))
	req, err := d.NewRequest("POST", Steam.UpdateCode, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Errors.ResponseError(resp.StatusCode)
	}
	eresult := resp.Header.Get("x-eresult")
	result, _ := strconv.Atoi(eresult)
	switch result {
	case 1:
		return nil // 登录成功
	case 29:
		return Errors.Error("验证码无效")
	case 84:
		return Errors.Error("提交验证码 请求过于频繁，请稍后再试")
	case 88:
		return Errors.Error("验证码输入错误")
	case 65:
		return Errors.Error("验证码输入错误")
	default:
		return Errors.Error("提交验证码时发生错误 x_eresult:" + eresult)
	}
}

// generateTokenCode 手机令牌校验
func (d *Dao) generateTokenCode(clientId uint64, requestId []byte, steamId uint64, confirmationType int32, sharedSecret string) error {
	stime, err := d.SteamTime()
	if err != nil {
		return err
	}
	code := Utils.GenerateAuthCode(sharedSecret, stime)
	err = d.submitVerificationCode(clientId, steamId, confirmationType, code)
	if err != nil {
		return err
	}
	return d.afterVerificationLogin(clientId, steamId, requestId)
}

// GetTokenCode 获取token Code
func (d *Dao) GetTokenCode(sharedSecret string) (string, error) {
	stime, err := d.SteamTime()
	if err != nil {
		return "", err
	}
	code := Utils.GenerateAuthCode(sharedSecret, stime)
	return code, nil
}

// beginAuthSessionViaCredentials 开始处理Auth登录问题
func (d *Dao) beginAuthSessionViaCredentials(sharedSecret string) error {
	timestamp, _ := strconv.ParseInt(d.credentials.RSATimeStamp, 10, 64)
	loginData := &Protoc.BeginAuthSessionViaCredentialsSend{
		AccountName:         d.credentials.Username,
		EncryptedPassword:   d.credentials.Password,
		EncryptionTimestamp: timestamp,
		RememberLogin:       false,
		Persistence:         1,
		WebsiteId:           "Store",
		DeviceDetails: &Protoc.DeviceDetails{
			DeviceFriendlyName: Steam.UserAgent,
			PlatformType:       2,
		},
		Language: 6,
	}
	data, err := proto.Marshal(loginData)
	if err != nil {
		return err
	}
	params := Param.Params{}
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))
	req, err := d.NewRequest("POST", Steam.BeginAuthSessionViaCredentials, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Errors.ResponseError(resp.StatusCode)
	}
	eresult := resp.Header.Get("x-eresult")
	result, _ := strconv.Atoi(eresult)
	switch result {
	case 1:
		credentialsReceive := &Protoc.BeginAuthSessionViaCredentialsReceive{}
		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return err
		}
		err = proto.Unmarshal(buf.Bytes(), credentialsReceive)
		// 判断是否需要二次验证
		allowedConfirmations := credentialsReceive.AllowedConfirmations
		steamId := credentialsReceive.SteamId
		requestId := credentialsReceive.RequestId
		clientId := credentialsReceive.ClientId
		for _, item := range allowedConfirmations {
			switch item.ConfirmationType {
			case 1: // 免校验
				return d.afterVerificationLogin(clientId, steamId, requestId)
			case 6: // 邮箱验证码
				return Errors.Error("需要邮箱验证码,功能暂未实现") // 暂未实现
			case 3: // 需要手机验证码
				if len(sharedSecret) < 10 {
					return Errors.Error("需要手机验证码,或者手机令牌")
				}
				return d.generateTokenCode(clientId, requestId, steamId, item.ConfirmationType, sharedSecret)
			default:
				return Errors.Unavailable()
			}
		}
		break
	case 5:
		return Errors.Error("密码错误")
	case 84:
		return Errors.Error("请求过于频繁,请稍后再试")
	default:
		return Errors.Error("未捕获的异常,result:" + strconv.Itoa(result))
	}
	return Errors.Unavailable()
}

// encryptPassword 密码加密
func (d *Dao) encryptPassword(password string, spk *Model.SteamPublicKey) (string, error) {
	pk := new(rsa.PublicKey)
	exp, err := spk.Exponent()
	if err != nil {
		return "", Errors.Error(err.Error())
	}
	pk.E = int(exp)
	if pk.N, err = spk.Modulus(); err != nil {
		return "", Errors.Error(err.Error())
	}
	out, err := rsa.EncryptPKCS1v15(rand.Reader, pk, []byte(password))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

// Login Steam用户登录
// 执行完整的Steam登录流程，包括RSA加密密码和身份验证
// 参数：
//   username - Steam用户名
//   password - Steam密码（明文）
//   sharedSecret - Steam Guard共享密钥（用于生成验证码）
// 返回值：登录成功返回nil，失败返回错误信息
func (d *Dao) Login(username, password, sharedSecret string) error {
	// 1. 获取RSA公钥用于密码加密
	keySendReceive, err := d.getRSA(username)
	if err != nil {
		return err
	}
	
	// 2. 使用RSA公钥加密密码
	encryptedPassword, err := d.encryptPassword(password, keySendReceive)
	if err != nil {
		return err
	}
	
	// 3. 保存用户凭据信息
	d.credentials.Username = username
	d.credentials.Password = encryptedPassword
	d.credentials.RSATimeStamp = strconv.FormatUint(keySendReceive.Timestamp, 10)
	
	// 4. 开始通过凭据进行身份验证
	return d.beginAuthSessionViaCredentials(sharedSecret)
}

// SetLoginInfo 设置登录信息
func (d *Dao) SetLoginInfo(username, password, accessToken, countryCode string, cookies string) error {
	//keySendReceive, err := d.getRSA(username)
	//if err != nil {
	//	return err
	//}
	//encryptedPassword, err := d.encryptPassword(password, keySendReceive)
	//if err != nil {
	//	return err
	//}
	d.credentials.Username = username
	d.credentials.Password = password
	d.credentials.AccessToken = accessToken
	d.credentials.CountryCode = countryCode
	d.credentials.Language = "english"
	err := json.Unmarshal([]byte(cookies), &d.credentials.LoginCookies)
	if err != nil {
		return Errors.Error(err.Error())
	}
	return nil
}
