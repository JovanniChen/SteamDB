package main

import (
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

func TestRemoveFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err := client.RemoveFriend(76561198313222178); err != nil {
		Logger.Error("删除好友失败: ", err)
		return
	}
	Logger.Info("删除好友成功")
}

func TestCreateFriendLink(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	friendLink, err := client.CreateFriendLink()
	if err != nil {
		Logger.Error("创建好友链接失败: ", err)
		return
	}
	Logger.Info("创建好友链接成功: ", friendLink)
}

func TestAddFriendByFriendCode(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err := client.AddFriendByFriendCode(352956450); err != nil {
		Logger.Error("添加好友失败: ", err)
		return
	}
	Logger.Info("添加好友成功")
}

func TestAcceptFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	if err := client.AcceptFriend("76561198313222178"); err != nil {
		Logger.Error("接受好友请求失败: ", err)
		return
	}
	Logger.Info("接受好友请求成功")
}

func TestCheckIsFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	isFriend, err := client.CheckIsFriend("76561198313222178")
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info("是否是好友: ", isFriend)
}

func TestCheckFriendStatus(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err = client.CheckFriendStatus("https://s.team/p/jtgt-jrbr/CDPKKTCF"); err != nil {
		Logger.Error("检查好友状态失败: ", err)
		return
	}
	Logger.Info("检查好友状态成功")
}

func TestAddFriendByLink(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	var links = make([]string, 0)
	links = append(links, "https://s.team/p/hvfw-wmtw/BGDVKBND")

	for i := 0; i < len(links); i++ {
		steamID, err := client.AddFriendByLink(links[i])
		if err != nil {
			Logger.Error(i, err)
		}

		Logger.Info("添加好友成功: ", steamID)
		time.Sleep(2 * time.Second)
	}

}

func TestGetFriendInfoByLink(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	friendInfo, inviteToken, err := client.GetFriendInfoByLink("https://s.team/p/chbn-qbdd/BHMGGQBR")
	if err != nil {
		Logger.Info(friendInfo)
		Logger.Error(err)
		return
	}
	Logger.Info(friendInfo)
	Logger.Info(inviteToken)
}

func TestGetFriendInfoByLinkAndAddFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	friendInfo, inviteToken, err := client.GetFriendInfoByLink("https://s.team/p/chbn-qbdd/QVRQMRJK")
	if err != nil {
		Logger.Error(err)
		return
	}

	_, err = client.AddFriendByInviteTokenAndSteamID(inviteToken, friendInfo.AbuseID)
	if err != nil {
		Logger.Error("通过邀请token和steamID添加好友失败,错误:", err)
		return
	}

	Logger.Info("通过邀请token和steamID添加好友成功")
}
