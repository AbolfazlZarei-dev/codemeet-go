package shopengine

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

// Models
type Product struct {
	ID          string
	Name        string
	Description string
	Price       int64
	Category    string
	Stock       int
}

type CartItem struct {
	Product  Product
	Quantity int
}

type Order struct {
	ID         string
	UserID     string
	Items      []CartItem
	TotalPrice int64
	Status     string // pending, paid, shipped, canceled
	CreatedAt  int64
}

// Config
type Config struct {
	AdminIDs  []string
	Currency  string // مثلاً "تومان"
	SupportID string // آیدی پشتیبانی برای ارسال سفارش

	// Actions
	SendMessageAction    func(ctx context.Context, chatID, text string, markup *models.InlineKeyboardMarkup) (int, error)
	EditMessageAction    func(ctx context.Context, chatID string, messageID int, text string, markup *models.InlineKeyboardMarkup) error
	DeleteMessageAction  func(ctx context.Context, chatID string, messageID int) error
	AnswerCallbackAction func(ctx context.Context, callbackID, text string, showAlert bool) error

	// Database Actions (برای ذخیره سازی دائمی)
	SaveProductAction   func(ctx context.Context, product Product) error
	DeleteProductAction func(ctx context.Context, productID string) error
	GetProductsAction   func(ctx context.Context) ([]Product, error)

	SaveOrderAction func(ctx context.Context, order Order) error
}

// State Management
type userState struct {
	mu    sync.Mutex
	state string
	data  map[string]string
}

type ShopEngine struct {
	cfg    Config
	admins map[string]struct{}
	stats  stats
	states sync.Map
	carts  sync.Map
}

type stats struct {
	productsAdded atomic.Int64
	ordersCreated atomic.Int64
	purchasesMade atomic.Int64
}

const (
	StateNone         = ""
	StateAddProdName  = "add_prod_name"
	StateAddProdDesc  = "add_prod_desc"
	StateAddProdPrice = "add_prod_price"
	StateAddProdCat   = "add_prod_cat"
	StateAddProdStock = "add_prod_stock"
)

// Initialization
func New(cfg Config) *ShopEngine {
	if cfg.Currency == "" {
		cfg.Currency = "تومان"
	}
	admins := make(map[string]struct{})
	for _, id := range cfg.AdminIDs {
		admins[id] = struct{}{}
	}

	return &ShopEngine{
		cfg:    cfg,
		admins: admins,
	}
}

func (se *ShopEngine) isAdmin(userID string) bool {
	_, ok := se.admins[userID]
	return ok
}

func (se *ShopEngine) getState(userID string) *userState {
	val, _ := se.states.LoadOrStore(userID, &userState{data: make(map[string]string)})
	return val.(*userState)
}

func (se *ShopEngine) clearState(userID string) {
	se.states.Delete(userID)
}

func (se *ShopEngine) getCart(userID string) map[string]*CartItem {
	val, _ := se.carts.LoadOrStore(userID, make(map[string]*CartItem))
	return val.(map[string]*CartItem)
}

// Middleware
func (se *ShopEngine) Middleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil {
				next(ctx, u)
				return
			}

			// مدیریت دکمه‌های شیشه‌ای (Callback Queries)
			if u.CallbackQuery != nil {
				se.handleCallback(ctx, u.CallbackQuery)
				return
			}

			// مدیریت پیام‌های متنی
			if u.Message != nil {
				userID := u.Message.From.ID
				text := u.Message.Text

				// 1. بررسی اگر کاربر در حالت خاصی است (مثلا ادمین دارد محصول اضافه می‌کند)
				state := se.getState(userID)
				if state.state != StateNone {
					se.handleState(ctx, u.Message, state)
					return
				}

				// 2. مدیریت دستورات
				if strings.HasPrefix(text, "/") || strings.HasPrefix(text, "#") {
					if se.handleCommand(ctx, u.Message) {
						return
					}
				}
			}

			next(ctx, u)
		}
	}
}

// Command Handlers
func (se *ShopEngine) handleCommand(ctx context.Context, msg *models.Message) bool {
	userID := msg.From.ID
	text := strings.ToLower(strings.TrimSpace(msg.Text))

	switch text {
	case "/shop", "فروشگاه":
		se.ShowMainMenu(ctx, userID)
		return true
	case "/admin", "پنل مدیریت":
		if se.isAdmin(userID) {
			se.showAdminPanel(ctx, userID)
		}
		return true
	}
	return false
}

func (se *ShopEngine) ShowMainMenu(ctx context.Context, chatID string) {
	keyboard := models.NewInlineKeyboard(
		models.InlineRow(models.Btn("🛍 محصولات", "shop_list")),
		models.InlineRow(models.Btn("🛒 سبد خرید", "shop_cart")),
		models.InlineRow(models.Btn("📦 سفارشات من", "shop_orders")),
		models.InlineRow(models.Btn("📞 پشتیبانی", "shop_support")),
	)
	se.cfg.SendMessageAction(ctx, chatID, "👋 به ربات فروشگاه ما خوش آمدید!\nلطفاً یک گزینه را انتخاب کنید:", keyboard)
}

func (se *ShopEngine) showAdminPanel(ctx context.Context, chatID string) {
	keyboard := models.NewInlineKeyboard(
		models.InlineRow(models.Btn("➕ افزودن محصول", "admin_add_prod")),
		models.InlineRow(models.Btn("📋 لیست محصولات", "admin_list_prod")),
		models.InlineRow(models.Btn("📦 سفارشات جدید", "admin_list_orders")),
	)
	se.cfg.SendMessageAction(ctx, chatID, "🛠 <b>پنل مدیریت</b>\nیک گزینه را انتخاب کنید:", keyboard)
}

// State Handlers (For Admin Product Creation)
func (se *ShopEngine) handleState(ctx context.Context, msg *models.Message, state *userState) {
	userID := msg.From.ID
	text := msg.Text

	state.mu.Lock()
	defer state.mu.Unlock()

	switch state.state {
	case StateAddProdName:
		state.data["name"] = text
		state.state = StateAddProdDesc
		se.cfg.SendMessageAction(ctx, userID, "📝 لطفاً توضیحات محصول را ارسال کنید:", nil)

	case StateAddProdDesc:
		state.data["desc"] = text
		state.state = StateAddProdPrice
		se.cfg.SendMessageAction(ctx, userID, "💰 لطفاً قیمت محصول را ارسال کنید (فقط عدد):", nil)

	case StateAddProdPrice:
		state.data["price"] = text
		state.state = StateAddProdStock
		se.cfg.SendMessageAction(ctx, userID, "📦 لطفاً موجودی انبار را ارسال کنید (فقط عدد):", nil)

	case StateAddProdStock:
		state.data["stock"] = text
		state.state = StateAddProdCat
		se.cfg.SendMessageAction(ctx, userID, "🏷 لطفاً دسته‌بندی محصول را ارسال کنید (مثلاً دیجیتال):", nil)

	case StateAddProdCat:
		state.data["cat"] = text

		// ساخت محصول نهایی
		prod := Product{
			ID:          fmt.Sprintf("p_%d", time.Now().Unix()),
			Name:        state.data["name"],
			Description: state.data["desc"],
			Price:       parseInt64(state.data["price"]),
			Stock:       parseInt(state.data["stock"]),
			Category:    state.data["cat"],
		}

		if se.cfg.SaveProductAction != nil {
			if err := se.cfg.SaveProductAction(ctx, prod); err != nil {
				se.cfg.SendMessageAction(ctx, userID, "❌ خطا در ذخیره‌سازی محصول.", nil)
			} else {
				se.stats.productsAdded.Add(1)
				se.cfg.SendMessageAction(ctx, userID, fmt.Sprintf("✅ محصول <b>%s</b> با موفقیت ذخیره شد!", prod.Name), nil)
			}
		} else {
			se.cfg.SendMessageAction(ctx, userID, "✅ محصول ساخته شد (حالت دمو، دیتابیس متصل نیست).", nil)
		}

		se.clearState(userID)
		se.showAdminPanel(ctx, userID)
	}
}

// Callback Handlers
func (se *ShopEngine) handleCallback(ctx context.Context, cq *models.CallbackQuery) {
	userID := cq.From.ID
	data := cq.Data
	chatID := cq.Message.Chat.ID

	// پاسخ به کالبک برای رفع لودینگ دکمه
	if se.cfg.AnswerCallbackAction != nil {
		se.cfg.AnswerCallbackAction(ctx, cq.ID, "", false)
	}

	if strings.HasPrefix(data, "shop_") {
		se.handleShopCallbacks(ctx, userID, chatID, data, cq.Message.MessageID)
	} else if strings.HasPrefix(data, "admin_") && se.isAdmin(userID) {
		se.handleAdminCallbacks(ctx, userID, chatID, data, cq.Message.MessageID)
	} else {
		// اگر کاربر عادی روی دکمه ادمین کلیک کرد
		se.cfg.SendMessageAction(ctx, chatID, "🚫 شما دسترسی ادمین ندارید.", nil)
	}
}

func (se *ShopEngine) handleShopCallbacks(ctx context.Context, userID, chatID, data string, msgID int) {
	switch data {
	case "shop_list":
		// فرض میکنیم محصولاتی از دیتابیس می‌گیریم (در اینجا دمو می‌زنیم)
		// در عمل: products, _ := se.cfg.GetProductsAction(ctx)
		text := "🛍 <b>لیست محصولات:</b>\n\n"
		keyboard := models.NewInlineKeyboard(
			models.InlineRow(models.Btn("خرید محصول تست ۱", "shop_buy_p1")),
			models.InlineRow(models.Btn("خرید محصول تست ۲", "shop_buy_p2")),
			models.InlineRow(models.Btn("🔙 بازگشت", "shop_main")),
		)
		se.cfg.EditMessageAction(ctx, chatID, msgID, text, keyboard)

	case "shop_cart":
		cart := se.getCart(userID)
		text := "🛒 <b>سبد خرید شما:</b>\n\n"
		var total int64
		if len(cart) == 0 {
			text += "سبد خرید شما خالی است."
		} else {
			for _, item := range cart {
				text += fmt.Sprintf("▫️ %s (x%d) - %d %s\n", item.Product.Name, item.Quantity, item.Product.Price*int64(item.Quantity), se.cfg.Currency)
				total += item.Product.Price * int64(item.Quantity)
			}
			text += fmt.Sprintf("\n💰 <b>مبلغ کل:</b> %d %s", total, se.cfg.Currency)
		}

		keyboard := models.NewInlineKeyboard()
		if len(cart) > 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, models.InlineRow(models.Btn("✅ نهایی کردن خرید", "shop_checkout")))
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, models.InlineRow(models.Btn("🗑 خالی کردن سبد", "shop_clear_cart")))
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, models.InlineRow(models.Btn("🔙 بازگشت", "shop_main")))

		se.cfg.EditMessageAction(ctx, chatID, msgID, text, keyboard)

	case "shop_main":
		keyboard := models.NewInlineKeyboard(
			models.InlineRow(models.Btn("🛍 محصولات", "shop_list")),
			models.InlineRow(models.Btn("🛒 سبد خرید", "shop_cart")),
			models.InlineRow(models.Btn("📦 سفارشات من", "shop_orders")),
			models.InlineRow(models.Btn("📞 پشتیبانی", "shop_support")),
		)
		se.cfg.EditMessageAction(ctx, chatID, msgID, "👋 به ربات فروشگاه ما خوش آمدید!\nلطفاً یک گزینه را انتخاب کنید:", keyboard)

	case "shop_support":
		supportLink := "https://codemeet.chat/your_support"
		if se.cfg.SupportID != "" {
			supportLink = "https://codemeet.chat/" + se.cfg.SupportID
		}
		keyboard := models.NewInlineKeyboard(
			models.InlineRow(models.URLBtn("💬 ارتباط با پشتیبانی", supportLink)),
			models.InlineRow(models.Btn("🔙 بازگشت", "shop_main")),
		)
		se.cfg.EditMessageAction(ctx, chatID, msgID, "📞 برای ارتباط با پشتیبانی روی دکمه زیر کلیک کنید:", keyboard)
	}
}

func (se *ShopEngine) handleAdminCallbacks(ctx context.Context, userID, chatID, data string, msgID int) {
	switch data {
	case "admin_add_prod":
		state := se.getState(userID)
		state.mu.Lock()
		state.state = StateAddProdName
		state.data = make(map[string]string)
		state.mu.Unlock()

		se.cfg.EditMessageAction(ctx, chatID, msgID, "➕ <b>افزودن محصول جدید</b>\n\nلطفاً نام محصول را ارسال کنید:", nil)
	}
}

// Helpers
func parseInt64(s string) int64 {
	var n int64
	fmt.Sscanf(s, "%d", &n)
	return n
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
