package Model

// UpdateEventInfo 提取的更新事件信息（简化版）
type UpdateEventInfo struct {
	UniqueID  string `json:"unique_id"`
	AppID     int    `json:"appid"`
	StartTime int64  `json:"start_time"`
	EventName string `json:"event_name"`
}

// GameUpdateEvents 游戏更新事件的主结构
type GameUpdateEvents struct {
	ForwardComplete  bool             `json:"forwardComplete"`
	BackwardComplete bool             `json:"backwardComplete"`
	Documents        []EventDocument  `json:"documents"`
	Apps             []AppInfo        `json:"apps"`
	Events           []EventDetail    `json:"events"`
}

// EventDocument 事件文档概要信息
type EventDocument struct {
	ClanID    int64  `json:"clanid"`
	UniqueID  string `json:"unique_id"`
	EventType int    `json:"event_type"`
	AppID     int    `json:"appid"`
	StartTime int64  `json:"start_time"`
	Score     int    `json:"score"`
}

// AppInfo 应用信息
type AppInfo struct {
	Source int `json:"source"`
	AppID  int `json:"appid"`
}

// EventDetail 事件详细信息
type EventDetail struct {
	GID                   string            `json:"gid"`
	ClanSteamID           string            `json:"clan_steamid"`
	EventName             string            `json:"event_name"`
	EventType             int               `json:"event_type"`
	AppID                 int               `json:"appid"`
	ServerAddress         *string           `json:"server_address"`
	ServerPassword        *string           `json:"server_password"`
	Rtime32StartTime      int64             `json:"rtime32_start_time"`
	Rtime32EndTime        int64             `json:"rtime32_end_time"`
	CommentCount          int               `json:"comment_count"`
	CreatorSteamID        string            `json:"creator_steamid"`
	LastUpdateSteamID     string            `json:"last_update_steamid"`
	EventNotes            string            `json:"event_notes"`
	JSONData              string            `json:"jsondata"`
	AnnouncementBody      *AnnouncementBody `json:"announcement_body"`
	Published             int               `json:"published"`
	Hidden                int               `json:"hidden"`
	Rtime32VisibilityStart int64            `json:"rtime32_visibility_start"`
	Rtime32VisibilityEnd   int64            `json:"rtime32_visibility_end"`
	BroadcasterAccountID  int               `json:"broadcaster_accountid"`
	FollowerCount         int               `json:"follower_count"`
	IgnoreCount           int               `json:"ignore_count"`
	ForumTopicID          string            `json:"forum_topic_id"`
	Rtime32LastModified   int64             `json:"rtime32_last_modified"`
	NewsPostGID           string            `json:"news_post_gid"`
	RtimeModReviewed      int64             `json:"rtime_mod_reviewed"`
	FeaturedAppTagID      int               `json:"featured_app_tagid"`
	ReferencedAppIDs      []int             `json:"referenced_appids"`
	BuildID               int               `json:"build_id,omitempty"`
	BuildBranch           string            `json:"build_branch,omitempty"`
	Unlisted              int               `json:"unlisted"`
	VotesUp               int               `json:"votes_up"`
	VotesDown             int               `json:"votes_down"`
	CommentType           string            `json:"comment_type"`
	GIDFeature            string            `json:"gidfeature"`
	GIDFeature2           string            `json:"gidfeature2"`
	ClanSteamIDOriginal   string            `json:"clan_steamid_original"`
}

// AnnouncementBody 公告正文信息
type AnnouncementBody struct {
	GID             string   `json:"gid"`
	ClanID          string   `json:"clanid"`
	PosterID        string   `json:"posterid"`
	Headline        string   `json:"headline"`
	PostTime        int64    `json:"posttime"`
	UpdateTime      int64    `json:"updatetime"`
	Body            string   `json:"body"`
	CommentCount    int      `json:"commentcount"`
	Tags            []string `json:"tags"`
	Language        int      `json:"language"`
	Hidden          int      `json:"hidden"`
	ForumTopicID    string   `json:"forum_topic_id"`
	EventGID        string   `json:"event_gid"`
	VoteUpCount     int      `json:"voteupcount"`
	VoteDownCount   int      `json:"votedowncount"`
	BanCheckResult  int      `json:"ban_check_result"`
	Banned          int      `json:"banned"`
}

// ExtractUpdateEvents 提取前三条 event_type 为 12 的更新事件信息
func (g *GameUpdateEvents) ExtractUpdateEvents() []UpdateEventInfo {
	// 创建 events 的 gid -> event_name 映射
	eventMap := make(map[string]string)
	for _, event := range g.Events {
		eventMap[event.GID] = event.EventName
	}

	var results []UpdateEventInfo
	count := 0

	// 遍历 documents，筛选 event_type 为 12 的数据
	for _, doc := range g.Documents {
		if doc.EventType == 12 {
			// 通过 unique_id 从 eventMap 中获取 event_name
			eventName := eventMap[doc.UniqueID]

			results = append(results, UpdateEventInfo{
				UniqueID:  doc.UniqueID,
				AppID:     doc.AppID,
				StartTime: doc.StartTime,
				EventName: eventName,
			})

			count++
			if count >= 3 {
				break
			}
		}
	}

	return results
}

// ExtractUpdateEventsWithLimit 提取指定数量的 event_type 为 12 的更新事件信息
// 返回提取的事件列表和总共找到的 event_type=12 的数量
func (g *GameUpdateEvents) ExtractUpdateEventsWithLimit(limit int) ([]UpdateEventInfo, int) {
	// 创建 events 的 gid -> event_name 映射
	eventMap := make(map[string]string)
	for _, event := range g.Events {
		eventMap[event.GID] = event.EventName
	}

	var results []UpdateEventInfo
	totalCount := 0 // 总共找到的 event_type=12 的数量

	// 遍历 documents，筛选 event_type 为 12 的数据
	for _, doc := range g.Documents {
		if doc.EventType == 12 {
			totalCount++

			// 只添加到结果中，如果还没有达到限制
			if len(results) < limit {
				// 通过 unique_id 从 eventMap 中获取 event_name
				eventName := eventMap[doc.UniqueID]

				results = append(results, UpdateEventInfo{
					UniqueID:  doc.UniqueID,
					AppID:     doc.AppID,
					StartTime: doc.StartTime,
					EventName: eventName,
				})
			}
		}
	}

	return results, totalCount
}
