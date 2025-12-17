// Param包 - HTTP请求参数处理工具
// 提供用于构建和管理HTTP请求参数的工具函数
package Param

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	SIGN_TYPE_MD5         = "MD5"
	SIGN_TYPE_HMAC_SHA256 = "HMAC-SHA256"
)

// Params HTTP请求参数映射类型
// 基于map[string]string，用于存储和操作HTTP请求参数
type Params map[string]string

// SetString 设置字符串参数
// 参数：k - 参数名, s - 参数值
func (p Params) SetString(k, s string) {
	p[k] = s
}

// GetString 获取字符串参数
// 参数：k - 参数名
// 返回值：参数值，不存在时返回空字符串
func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

// SetInt64 设置int64类型参数
// 将整数转换为字符串存储
// 参数：k - 参数名, i - 整数值
func (p Params) SetInt64(k string, i int64) {
	p[k] = strconv.FormatInt(i, 10)
}

// GetInt 获取int类型参数
// 将字符串参数转换为整数
// 参数：k - 参数名
// 返回值：整数值，转换失败时返回0
func (p Params) GetInt(k string) int {
	i, _ := strconv.Atoi(p.GetString(k))
	return i
}

// GetInt64 获取int64类型参数
// 将字符串参数转换为64位整数
// 参数：k - 参数名
// 返回值：64位整数值，转换失败时返回0
func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

// IsSet 检查参数是否存在
// 参数：key - 要检查的参数名
// 返回值：true表示参数存在，false表示不存在
func (p Params) IsSet(key string) bool {
	_, ok := p[key]
	return ok
}

// ToUrl 将参数转换为URL查询字符串
// 排除sign参数和空值参数，生成key=value&key=value格式
// 返回值：URL查询字符串
func (p Params) ToUrl() string {
	buff := ""
	for k, v := range p {
		// 跳过sign参数和空值
		if k != "sign" && v != "" {
			buff += k + "=" + v + "&"
		}
	}
	// 移除末尾的&符号
	buff = strings.Trim(buff, "&")
	return buff
}

// QueryEscape URL编码函数
// 对字符串进行URL编码，处理特殊字符
// 参数：s - 要编码的字符串
// 返回值：编码后的字符串
func (p Params) QueryEscape(s string) string {
	return url.QueryEscape(s)
}

// Encode 将参数编码为HTTP表单格式
// 按字母顺序排列参数，并对键值进行URL编码
// 返回值：编码后的表单数据字符串
func (p Params) Encode() string {
	if len(p) == 0 {
		return ""
	}

	var buf strings.Builder
	// 按字母顺序遍历所有键
	for _, k := range slices.Sorted(maps.Keys(p)) {
		vs := p[k]
		keyEscaped := p.QueryEscape(k) // 对键进行URL编码

		// 在非第一个参数前添加&分隔符
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		// 构建 key=value 格式
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		buf.WriteString(p.QueryEscape(vs)) // 对值进行URL编码
	}
	return buf.String()
}

func (p Params) EncodeBy(param []string) string {
	if len(p) == 0 {
		return ""
	}
	var buf strings.Builder
	for _, k := range param {
		vs := p[k]
		keyEscaped := p.QueryEscape(k) // 对键进行URL编码
		// 在非第一个参数前添加&分隔符
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		// 构建 key=value 格式
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		buf.WriteString(p.QueryEscape(vs)) // 对值进行URL编码
	}
	return buf.String()
}

func (p Params) CreateTimeStamp() {
	p["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().Unix())
}

func (p Params) CreateSign(secret string) {
	p["sign"], _ = p.Sign(SIGN_TYPE_MD5, secret)
}

func (p Params) Sign(signType, secret string) (string, error) {
	type kv struct {
		key   string
		value string
	}
	kvlist := make([]kv, 0)
	var totalLen int
	totalLen += len(secret) * 2 // 前后各拼接一次 secret
	for k, v := range p {
		if k == "sign" {
			continue
		}
		kvlist = append(kvlist, kv{key: k, value: v})
		totalLen += len(k) + len(v)
	}
	sort.Slice(kvlist, func(i, j int) bool {
		return kvlist[i].key < kvlist[j].key
	})
	var sb strings.Builder
	sb.Grow(totalLen) // 精确预分配，性能最优
	sb.WriteString(secret)
	for _, kv := range kvlist {
		sb.WriteString(kv.key)
		sb.WriteString(kv.value)
	}
	sb.WriteString(secret)

	switch signType {
	case SIGN_TYPE_MD5:
		return p.md5(sb.String())
	case SIGN_TYPE_HMAC_SHA256:
		return p.sha256(sb.String(), secret)
	}
	return "", errors.New("")
}

func (p Params) MakeSign(secret string) string {
	sign, _ := p.Sign(SIGN_TYPE_MD5, secret)
	return sign
}

func (p Params) CheckSign(secret string) (bool, error) {
	if !p.IsSet("sign") {
		return false, errors.New("签名不存在")
	} else if p.GetString("sign") == "" {
		return false, errors.New("签名存在但不合法")
	}
	returnSign := p.GetString("sign")
	calSign := p.MakeSign(secret)
	if calSign == returnSign {
		return true, nil
	}
	return false, fmt.Errorf("签名验证失败: 预期签名[%s], 实际签名[%s]", calSign, returnSign)
}

func (p Params) md5(str string) (string, error) {
	md5Hash := md5.New()
	if _, err := md5Hash.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("计算 MD5 失败: %w", err)
	}
	return hex.EncodeToString(md5Hash.Sum(nil)), nil
}

func (p Params) sha256(str, secret string) (string, error) {
	hmacHash := hmac.New(sha256.New, []byte(secret))
	if _, err := hmacHash.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("计算 HMAC-SHA256 失败: %w", err)
	}
	sign := hex.EncodeToString(hmacHash.Sum(nil))
	return sign, nil
}

// Get 请求
func (p Params) Get(URL string) ([]byte, error) {
	// 构建查询参数
	values := url.Values{}
	for k, v := range p {
		values.Set(k, v)
	}
	// 拼接URL和查询参数
	fullURL := URL
	if query := values.Encode(); query != "" {
		if strings.Contains(URL, "?") {
			fullURL += "&" + query
		} else {
			fullURL += "?" + query
		}
	}
	// 发送GET请求
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("GET请求失败: %w", err)
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	return body, nil
}

// Post 发送POST请求 (Content-Type: application/x-www-form-urlencoded)
func (p Params) Post(URL string) ([]byte, error) {
	// http.DefaultTransport
	// 构建表单数据
	values := url.Values{}
	for k, v := range p {
		values.Set(k, v)
	}
	resp, err := http.PostForm(URL, values)
	if err != nil {
		return nil, fmt.Errorf("POST请求失败: %w", err)
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}

func (p Params) SetParams(values url.Values) {
	for k, v := range values {
		if len(v) == 0 {
			continue
		}
		p[k] = v[0]
	}
}
