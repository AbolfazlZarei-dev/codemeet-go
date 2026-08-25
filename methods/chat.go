package methods

import (
	"context"

	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type ChatMethods struct {
	parent *Messages
}

func newChatMethods(m *Messages) *ChatMethods { return &ChatMethods{parent: m} }

func (c *ChatMethods) GetChat(ctx context.Context, chatID string) (*models.Chat, error) {
	var ch models.Chat
	err := c.parent.doWithRetry(ctx, "getChat", map[string]string{"chat_id": chatID}, &ch)
	return &ch, err
}

func (c *ChatMethods) GetChatMember(ctx context.Context, chatID, userID string) (*models.ChatMember, error) {
	var m models.ChatMember
	err := c.parent.doWithRetry(ctx, "getChatMember", map[string]string{"chat_id": chatID, "user_id": userID}, &m)
	return &m, err
}

func (c *ChatMethods) GetChatAdministrators(ctx context.Context, chatID string) ([]models.ChatMember, error) {
	var admins []models.ChatMember
	err := c.parent.doWithRetry(ctx, "getChatAdministrators", map[string]string{"chat_id": chatID}, &admins)
	return admins, err
}

func (c *ChatMethods) GetChatMemberCount(ctx context.Context, chatID string) (int, error) {
	var count int
	err := c.parent.doWithRetry(ctx, "getChatMemberCount", map[string]string{"chat_id": chatID}, &count)
	return count, err
}

func (c *ChatMethods) PinMessage(ctx context.Context, chatID string, messageID int, disableNotif bool) error {
	return c.parent.doWithRetry(ctx, "pinChatMessage", map[string]interface{}{
		"chat_id": chatID, "message_id": messageID, "disable_notification": disableNotif,
	}, nil)
}

func (c *ChatMethods) UnpinMessage(ctx context.Context, chatID string, messageID int) error {
	return c.parent.doWithRetry(ctx, "unpinChatMessage", map[string]interface{}{
		"chat_id": chatID, "message_id": messageID,
	}, nil)
}

func (c *ChatMethods) UnpinAllMessages(ctx context.Context, chatID string) error {
	return c.parent.doWithRetry(ctx, "unpinAllChatMessages", map[string]string{"chat_id": chatID}, nil)
}

// --- متدهای جدید مدیریت اعضا ---

// BanChatMember اخراج (بن) کاربر از چت
func (c *ChatMethods) BanChatMember(ctx context.Context, chatID, userID string, untilDate int64, revokeMessages bool) error {
	return c.parent.doWithRetry(ctx, "banChatMember", map[string]interface{}{
		"chat_id":         chatID,
		"user_id":         userID,
		"until_date":      untilDate,
		"revoke_messages": revokeMessages,
	}, nil)
}

// UnbanChatMember لغو اخراج کاربر
func (c *ChatMethods) UnbanChatMember(ctx context.Context, chatID, userID string, onlyIfBanned bool) error {
	return c.parent.doWithRetry(ctx, "unbanChatMember", map[string]interface{}{
		"chat_id":        chatID,
		"user_id":        userID,
		"only_if_banned": onlyIfBanned,
	}, nil)
}

// RestrictChatMember محدود کردن کاربر
func (c *ChatMethods) RestrictChatMember(ctx context.Context, chatID, userID string, permissions *models.ChatPermissions, untilDate int64) error {
	return c.parent.doWithRetry(ctx, "restrictChatMember", map[string]interface{}{
		"chat_id":     chatID,
		"user_id":     userID,
		"permissions": permissions,
		"until_date":  untilDate,
	}, nil)
}

// PromoteChatMember ارتقای کاربر به ادمین
func (c *ChatMethods) PromoteChatMember(ctx context.Context, chatID, userID string, rights *models.ChatAdministratorRights) error {
	return c.parent.doWithRetry(ctx, "promoteChatMember", map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
		"rights":  rights,
	}, nil)
}

// SetChatAdministratorCustomTitle تنظیم عنوان سفارشی برای ادمین
func (c *ChatMethods) SetChatAdministratorCustomTitle(ctx context.Context, chatID, userID, customTitle string) error {
	return c.parent.doWithRetry(ctx, "setChatAdministratorCustomTitle", map[string]string{
		"chat_id":      chatID,
		"user_id":      userID,
		"custom_title": customTitle,
	}, nil)
}

// --- تنظیمات چت ---

// SetChatTitle تغییر عنوان چت
func (c *ChatMethods) SetChatTitle(ctx context.Context, chatID, title string) error {
	return c.parent.doWithRetry(ctx, "setChatTitle", map[string]string{
		"chat_id": chatID, "title": title,
	}, nil)
}

// SetChatDescription تغییر توضیحات چت
func (c *ChatMethods) SetChatDescription(ctx context.Context, chatID, description string) error {
	return c.parent.doWithRetry(ctx, "setChatDescription", map[string]string{
		"chat_id": chatID, "description": description,
	}, nil)
}

// SetChatPhoto تغییر عکس چت (از فایل)
func (c *ChatMethods) SetChatPhoto(ctx context.Context, chatID, photoPath string) error {
	return c.parent.doWithRetryMultipart(ctx, "setChatPhoto",
		map[string]string{"chat_id": chatID},
		map[string]string{"photo": photoPath}, nil)
}

// DeleteChatPhoto حذف عکس چت
func (c *ChatMethods) DeleteChatPhoto(ctx context.Context, chatID string) error {
	return c.parent.doWithRetry(ctx, "deleteChatPhoto", map[string]string{"chat_id": chatID}, nil)
}

// SetChatPermissions تنظیم دسترسی‌های چت
func (c *ChatMethods) SetChatPermissions(ctx context.Context, chatID string, permissions *models.ChatPermissions) error {
	return c.parent.doWithRetry(ctx, "setChatPermissions", map[string]interface{}{
		"chat_id": chatID, "permissions": permissions,
	}, nil)
}

// LeaveChat خروج ربات از چت
func (c *ChatMethods) LeaveChat(ctx context.Context, chatID string) error {
	return c.parent.doWithRetry(ctx, "leaveChat", map[string]string{"chat_id": chatID}, nil)
}

// --- لینک دعوت ---

// ExportChatInviteLink ساخت لینک دعوت
func (c *ChatMethods) ExportChatInviteLink(ctx context.Context, chatID string) (string, error) {
	var link string
	err := c.parent.doWithRetry(ctx, "exportChatInviteLink", map[string]string{"chat_id": chatID}, &link)
	return link, err
}

// CreateChatInviteLink ساخت لینک دعوت با پارامترها
func (c *ChatMethods) CreateChatInviteLink(ctx context.Context, req *models.CreateChatInviteLinkRequest) (*models.ChatInviteLink, error) {
	var link models.ChatInviteLink
	err := c.parent.doWithRetry(ctx, "createChatInviteLink", req, &link)
	return &link, err
}

// EditChatInviteLink ویرایش لینک دعوت
func (c *ChatMethods) EditChatInviteLink(ctx context.Context, req *models.EditChatInviteLinkRequest) (*models.ChatInviteLink, error) {
	var link models.ChatInviteLink
	err := c.parent.doWithRetry(ctx, "editChatInviteLink", req, &link)
	return &link, err
}

// RevokeChatInviteLink لغو لینک دعوت
func (c *ChatMethods) RevokeChatInviteLink(ctx context.Context, chatID, inviteLink string) (*models.ChatInviteLink, error) {
	var link models.ChatInviteLink
	err := c.parent.doWithRetry(ctx, "revokeChatInviteLink", map[string]string{
		"chat_id": chatID, "invite_link": inviteLink,
	}, &link)
	return &link, err
}
