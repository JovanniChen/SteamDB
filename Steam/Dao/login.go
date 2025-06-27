package Dao

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"example.com/m/v2/Steam/Model"
	"example.com/m/v2/Steam/Param"
	"example.com/m/v2/Steam/Protoc"
	"example.com/m/v2/Steam/Utils"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type LoginCookie struct {
	SteamLoginSecure string `json:"steamLoginSecure"`
	SessionId        string `json:"sessionid"`
	//SteamLanguage    string `json:"Steam_Language"`
}

type Credentials struct {
	Password     string
	Username     string
	Nickname     string
	SteamID      uint64
	RSATimeStamp string
	AccessToken  string
	RefreshToken string
	Language     string
	CountryCode  string
	LoginCookies map[string]*LoginCookie
}

// AccessToken token
func (d *Dao) AccessToken() (string, error) {
	if d.credentials.AccessToken == "" {
		return "", Errors.Error("未获取到")
	}
	return d.credentials.AccessToken, nil
}

// getRSA 获取steam密钥用于加密密码
func (d *Dao) getRSA(username string) (*Model.SteamPublicKey, error) {
	publicKeySend := &Protoc.GetPasswordRSAPublicKeySend{
		AccountName: username,
	}
	data, err := proto.Marshal(publicKeySend)
	if err != nil {
		return nil, err
	}
	params := Param.Params{}
	params.SetString("origin", Steam.Origin)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))
	req, err := d.NewRequest("GET", Steam.GetPasswordRSAPublicKey+"?="+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("origin", Steam.Origin)
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
	keySendReceive := &Protoc.GetPasswordRSAPublicKeySendReceive{}
	err = proto.Unmarshal(buf.Bytes(), keySendReceive)
	if err != nil {
		return nil, err
	}
	spk := new(Model.SteamPublicKey)
	spk.Success = true
	spk.Timestamp = keySendReceive.Timestamp
	spk.PublicKeyMod = keySendReceive.PublickeyMod
	spk.PublicKeyExp = keySendReceive.PublickeyExp
	return spk, nil
}

// pollAuthSessionStatus  轮询身份验证会话状态
func (d *Dao) pollAuthSessionStatus(clientId uint64, requestId []byte) (*Protoc.PollAuthSessionStatusReceive, error) {
	loginData := &Protoc.PollAuthSessionStatusSend{
		ClientId:  clientId,
		RequestId: requestId,
	}
	data, err := proto.Marshal(loginData)
	if err != nil {
		return nil, err
	}
	params := Param.Params{}
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))
	req, err := d.NewRequest("POST", Steam.PollAuthSessionStatus, strings.NewReader(params.Encode()))
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
	eresult := resp.Header.Get("x-eresult")
	result, _ := strconv.Atoi(eresult)
	switch result {
	case 1:
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
			DeviceFriendlyName: Steam.UseAgent,
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

// Login steam登录
func (d *Dao) Login(username, password, sharedSecret string) error {
	keySendReceive, err := d.getRSA(username)
	if err != nil {
		return err
	}
	encryptedPassword, err := d.encryptPassword(password, keySendReceive)
	if err != nil {
		return err
	}
	d.credentials.Username = username
	d.credentials.Password = encryptedPassword
	d.credentials.RSATimeStamp = strconv.FormatUint(keySendReceive.Timestamp, 10)
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
