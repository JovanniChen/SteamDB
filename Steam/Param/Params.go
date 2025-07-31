// Param包 - HTTP请求参数处理工具
// 提供用于构建和管理HTTP请求参数的工具函数
package Param

import (
	"maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
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
