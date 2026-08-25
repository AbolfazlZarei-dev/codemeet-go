package models

type Chat struct {
	ID                  string           `json:"id"`
	Type                string           `json:"type"`
	Title               string           `json:"title,omitempty"`
	Username            string           `json:"username,omitempty"`
	Description         string           `json:"description,omitempty"`
	MembersCount        int              `json:"members_count,omitempty"`
	Bio                 string           `json:"bio,omitempty"`
	Photo               *ChatPhoto       `json:"photo,omitempty"`
	Permissions         *ChatPermissions `json:"permissions,omitempty"`
	InviteLink          string           `json:"invite_link,omitempty"`
	HasProtectedContent bool             `json:"has_protected_content,omitempty"`
}

const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeChannel    = "channel"
	ChatTypeSupergroup = "supergroup"
)

func (c *Chat) IsPrivate() bool    { return c.Type == ChatTypePrivate }
func (c *Chat) IsGroup() bool      { return c.Type == ChatTypeGroup }
func (c *Chat) IsChannel() bool    { return c.Type == ChatTypeChannel }
func (c *Chat) IsSupergroup() bool { return c.Type == ChatTypeSupergroup }

// ChatMember وضعیت عضویت کاربر
type ChatMember struct {
	User        *User  `json:"user"`
	Status      string `json:"status"`
	IsAnonymous bool   `json:"is_anonymous,omitempty"`
	CustomTitle string `json:"custom_title,omitempty"`
	UntilDate   int64  `json:"until_date,omitempty"`
	IsMember    bool   `json:"is_member,omitempty"`
}

const (
	MemberStatusCreator       = "creator"
	MemberStatusAdministrator = "administrator"
	MemberStatusMember        = "member"
	MemberStatusLeft          = "left"
	MemberStatusRestricted    = "restricted"
	MemberStatusKicked        = "kicked"
)

// ChatPermissions دسترسی‌های چت
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendAudios         bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool `json:"can_send_photos,omitempty"`
	CanSendVideos         bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
}

// ChatAdministratorRights حقوق ادمین
type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous,omitempty"`
	CanManageChat       bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool `json:"can_change_info,omitempty"`
	CanInviteUsers      bool `json:"can_invite_users,omitempty"`
	CanPostMessages     bool `json:"can_post_messages,omitempty"`
	CanEditMessages     bool `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`
	CanPostStories      bool `json:"can_post_stories,omitempty"`
	CanEditStories      bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories    bool `json:"can_delete_stories,omitempty"`
	CanManageTopics     bool `json:"can_manage_topics,omitempty"`
}

// ChatPhoto عکس چت
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatInviteLink لینک دعوت
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 *User  `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int64  `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
}

// CreateChatInviteLinkRequest درخواست ساخت لینک دعوت
type CreateChatInviteLinkRequest struct {
	ChatID             string `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int64  `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

// EditChatInviteLinkRequest درخواست ویرایش لینک دعوت
type EditChatInviteLinkRequest struct {
	ChatID             string `json:"chat_id"`
	InviteLink         string `json:"invite_link"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int64  `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}
