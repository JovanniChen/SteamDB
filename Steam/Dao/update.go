package Dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func (d *Dao) GetGameUpdateInofs(gameID int) error {
	fmt.Println("gameID = ", gameID)
	// 构建请求参数
	params := Param.Params{}
	params.SetString("updates", "true")

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Constants.GetGameUpdateInofs+"/"+strconv.Itoa(gameID)+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := d.RetryRequest(1, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 输出原始响应(调试用) - 只输出前500个字符
	htmlContent := buf.String()
	if len(htmlContent) > 500 {
		Logger.Debugf("HTML内容长度: %d 字符（只显示前500字符）: %s...", len(htmlContent), htmlContent[:500])
	} else {
		Logger.Debugf("HTML内容: %s", htmlContent)
	}

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		Logger.Errorf("解析HTML失败: %v", err)
		return err
	}
	Logger.Info("HTML解析成功")

	// 使用用户提供的xpath路径
	xpathQuery := "/html/body/div[1]/div[7]/div[6]/div[1]"
	Logger.Infof("使用xpath提取内容: %s", xpathQuery)

	node := htmlquery.FindOne(doc, xpathQuery)
	if node != nil {
		Logger.Info("✓ 成功找到节点")

		// 获取节点的HTML内容（包含所有子元素）
		htmlOutput := htmlquery.OutputHTML(node, true)

		// 获取节点的文本内容
		textContent := htmlquery.InnerText(node)

		// 获取节点的所有属性
		Logger.Info("节点属性:")
		for _, attr := range node.Attr {
			if strings.EqualFold(attr.Key, "data-initialevents") {
				Logger.Infof("  %s = %s", attr.Key, attr.Val)
			} else {
				Logger.Infof("  %s = %s", attr.Key, attr.Val[:min(len(attr.Val), 100)])
			}
		}

		// 如果节点有id属性，特别处理
		nodeID := getAttrValue(node, "id")
		if nodeID != "" {
			Logger.Infof("节点ID: %s", nodeID)
		}

		// 尝试提取data-config属性（Steam页面常用）
		dataConfig := getAttrValue(node, "data-config")
		if dataConfig != "" {
			Logger.Info("找到data-config属性，尝试解析JSON...")
			var config map[string]any
			if err := json.Unmarshal([]byte(dataConfig), &config); err == nil {
				// 格式化输出JSON
				prettyJSON, _ := json.MarshalIndent(config, "", "  ")
				Logger.Infof("解析的配置数据:\n%s", string(prettyJSON))
			} else {
				Logger.Errorf("JSON解析失败: %v", err)
				Logger.Infof("原始data-config内容（前2000字符）: %s", dataConfig[:min(len(dataConfig), 2000)])
			}
		}

		// 尝试提取data-initialEvents属性
		dataEvents := getAttrValue(node, "data-initialEvents")
		if dataEvents != "" {
			Logger.Info("找到data-initialEvents属性")
			Logger.Infof("事件数据长度: %d 字符", len(dataEvents))
			// 可选：解析events数据
		}

		// 输出基本信息
		Logger.Infof("文本内容长度: %d 字符", len(textContent))
		Logger.Infof("HTML内容长度: %d 字符", len(htmlOutput))

		// 如果需要查看HTML结构，显示前1000字符
		if len(htmlOutput) > 0 && len(htmlOutput) <= 2000 {
			Logger.Infof("完整HTML内容:\n%s", htmlOutput)
		} else if len(htmlOutput) > 2000 {
			Logger.Infof("HTML内容（前2000字符）:\n%s...", htmlOutput[:2000])
		}

	} else {
		Logger.Warn("✗ 未找到匹配的节点")
		Logger.Info("尝试查找页面中可能有用的div...")

		// 如果找不到，尝试提取一些常见的有用节点
		usefulXpaths := map[string]string{
			"application_config": "//div[@id='application_config']",
			"page_content":       "//div[@class='responsive_page_content']",
			"main_content":       "//div[@role='main']",
		}

		for name, xpath := range usefulXpaths {
			testNode := htmlquery.FindOne(doc, xpath)
			if testNode != nil {
				Logger.Infof("✓ 找到 %s 节点: %s", name, xpath)
			}
		}
	}

	return nil
}

// countChildren 计算节点的子节点数量
func countChildren(node *html.Node) int {
	count := 0
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		count++
	}
	return count
}

// getAttrValue 获取节点指定属性的值
func getAttrValue(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}
