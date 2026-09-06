# User Models

Package: `models`

## `User`

```go
type User struct {
    ID string `json:"id"`
    IsBot bool `json:"is_bot"`
    FirstName string `json:"first_name"`
    Username string `json:"username,omitempty"`
    LastName string `json:"last_name,omitempty"`
    LanguageCode string `json:"language_code,omitempty"`
    IsPremium bool `json:"is_premium,omitempty"`
    AddedToAttachmentMenu bool `json:"added_to_attachment_menu,omitempty"`
    CanJoinGroups bool `json:"can_join_groups,omitempty"`
    CanReadAllGroupMessages bool `json:"can_read_all_group_messages,omitempty"`
    SupportsInlineQueries bool `json:"supports_inline_queries,omitempty"`
}
```

## Helperها

```go
user.IsUser()
user.FullName()
user.Mention()
user.HTMLMention()
```

اگر username موجود باشد `Mention()` آن را با `@` برمی‌گرداند؛ در غیر این صورت FullName استفاده می‌شود.
