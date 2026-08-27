# Source Coverage

این documentation بر اساس دو منبع ارائه‌شده در همین پروژه ساخته شده است:

1. dump سورس کتابخانه
2. dump مستندات Bot API

Tree پروژه وجود packageهای زیر را تأیید می‌کند:

```text
api
cache
codemeet
contrib/antilink
contrib/antispam
contrib/forcejoin
contrib/gatekeeper
contrib/maintenancemode
contrib/pagination
contrib/profanityfilter
contrib/vpndetector
contrib/warnsystem
dispatcher
errors
logger
methods
middleware
models
polling
ratelimit
retry
webhook
```

سورس قابل بازیابی packageهای اصلی و اکثر contribها را پوشش می‌دهد. `contrib/pagination/pagination.go` در tree هست اما متن آن در dump سورس بازیابی‌شده موجود نبود؛ بنابراین API آن حدس زده نشده است.

در مستندات، قابلیت‌های موجود در Bot API documentation نیز وارد شده‌اند: پیام، رسانه، keyboard، commands، profile، groups/channels، updates، polling، webhook، models و error handling.
