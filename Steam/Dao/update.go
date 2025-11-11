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
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func (d *Dao) GetGameUpdateInofs(gameID int) (*Model.GameUpdateEvents, error) {
	// 构建请求参数
	params := Param.Params{}
	params.SetString("updates", "true")

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Constants.GetGameUpdateInofs+"/"+strconv.Itoa(gameID)+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, err := d.RetryRequest(1, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(buf.String()))
	if err != nil {
		Logger.Errorf("解析HTML失败: %v", err)
		return nil, err
	}

	// 使用xpath查找目标节点
	xpathQuery := "/html/body/div[1]/div[7]/div[6]/div[1]"
	node := htmlquery.FindOne(doc, xpathQuery)
	if node == nil {
		Logger.Warn("未找到目标节点")
		return nil, fmt.Errorf("未找到xpath节点: %s", xpathQuery)
	}

	// 提取data-initialevents属性
	dataEvents := getAttrValue(node, "data-initialevents")
	if dataEvents == "" {
		// 尝试大小写变体
		dataEvents = getAttrValue(node, "data-initialEvents")
	}

	fmt.Println(dataEvents)

	if dataEvents == "" {
		Logger.Warn("未找到data-initialevents属性")
		return nil, fmt.Errorf("节点中不存在data-initialevents属性")
	}

	// 解析 JSON 数据为结构体
	var events Model.GameUpdateEvents
	if err := json.Unmarshal([]byte(dataEvents), &events); err != nil {
		Logger.Errorf("解析JSON失败: %v", err)
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &events, nil
}

// GetGameUpdateEvents 获取游戏更新事件（简化版，只返回关键字段）
// 参数：gameID - 游戏ID，limit - 提取数量限制
// 返回：更新事件列表、总共找到的event_type=12的数量、是否需要更新
func (d *Dao) GetGameUpdateEvents(gameID int, limit int) ([]Model.UpdateEventInfo, int, bool, error) {
	// 调用完整的获取方法
	events, err := d.GetGameUpdateInofs(gameID)
	if err != nil {
		return nil, 0, false, err
	}

	// 提取指定数量的更新事件
	updateEvents, totalFound := events.ExtractUpdateEventsWithLimit(limit)

	// 如果没有找到任何更新事件，返回
	if len(updateEvents) == 0 {
		return updateEvents, totalFound, false, nil
	}

	// 检查并保存最新的更新事件（取第一条，因为是按时间排序的）
	needUpdate, err := CheckAndSaveUpdateEvent(gameID, &updateEvents[0])
	if err != nil {
		return nil, 0, false, fmt.Errorf("检查更新事件失败: %w", err)
	}

	return updateEvents, totalFound, needUpdate, nil
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
