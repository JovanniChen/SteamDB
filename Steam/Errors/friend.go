package Errors

import "errors"

var (
	ErrInvalidFriendLink       = errors.New("不可用的好友链接")
	ErrProfileNotFound         = errors.New("无法找到指定的个人资料（请检查链接是否可用）")
	ErrProfileStructureChanged = errors.New("好友链接页面结构发生变化，请联系开发人员")
	ErrFriendLinkExpired       = errors.New("您关注的好友邀请已过期")
	ErrAlreadyFriend           = errors.New("您关注了一份好友邀请 - 但你们已经是好友了")
	ErrAddFriendFailed         = errors.New("添加好友失败")
)

func IsInvalidFriendLink(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrInvalidFriendLink)
}
