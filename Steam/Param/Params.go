package Param

import (
	"maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type Params map[string]string

func (p Params) SetString(k, s string) {
	p[k] = s
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt64(k string, i int64) {
	p[k] = strconv.FormatInt(i, 10)
}

func (p Params) GetInt(k string) int {
	i, _ := strconv.Atoi(p.GetString(k))
	return i
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

func (p Params) IsSet(key string) bool {
	_, ok := p[key]
	return ok
}

func (p Params) ToUrl() string {
	buff := ""
	for k, v := range p {
		if k != "sign" && v != "" {
			buff += k + "=" + v + "&"
		}
	}
	buff = strings.Trim(buff, "&")
	return buff
}

func (p Params) QueryEscape(s string) string {
	return url.QueryEscape(s)
}

func (p Params) Encode() string {
	if len(p) == 0 {
		return ""
	}
	var buf strings.Builder
	for _, k := range slices.Sorted(maps.Keys(p)) {
		vs := p[k]
		keyEscaped := p.QueryEscape(k)
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		buf.WriteString(p.QueryEscape(vs))
	}
	return buf.String()
}
