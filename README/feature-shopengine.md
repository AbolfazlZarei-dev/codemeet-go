# Shop Engine

موتور فروش ساده با state داخلی و callbackهای قابل اتصال.

## Models

```go
type Product struct {
    ID string
    Name string
    Description string
    Price int64
    Category string
    Stock int
}

type CartItem struct {
    Product Product
    Quantity int
}

type Order struct {
    ID string
    UserID string
    Items []CartItem
    TotalPrice int64
    Status string
    CreatedAt int64
}
```

Status:

```text
pending
paid
shipped
canceled
```

## Config

شامل:

```text
AdminIDs
Currency
SupportID
SendMessageAction
EditMessageAction
DeleteMessageAction
AnswerCallbackAction
SaveProductAction
DeleteProductAction
GetProductsAction
SaveOrderAction
```

## استفاده

```go
shop := shopengine.New(shopengine.Config{
    AdminIDs: []string{"admin"},
    Currency: "تومان",
    SupportID: "support",
    // callbacks...
})

bot.Use(shop.Middleware())
shop.ShowMainMenu(ctx, chatID)
```

## Stateها

```text
StateNone
StateAddProdName
StateAddProdDesc
StateAddProdPrice
StateAddProdCat
StateAddProdStock
```

Persistence با Actionهای شما انجام می‌شود؛ بنابراین می‌توانید Database یا storage دیگر را وصل کنید.
