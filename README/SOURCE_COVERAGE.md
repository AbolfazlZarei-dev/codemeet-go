
# Source Coverage

این Documentation بر اساس دو منبع اصلی ارائه‌شده در همین پروژه تهیه شده است:

1. **Dump سورس کتابخانه**
2. **Dump مستندات CodeMeet Bot API**

محتوای مستندات تا حد امکان بر اساس ساختار، APIها، packageها، typeها، methodها و قابلیت‌هایی نوشته شده است که در منابع ارائه‌شده قابل بررسی و تأیید بوده‌اند.

## Package Coverage

Tree پروژه وجود packageهای اصلی زیر را تأیید می‌کند:

```text
api
cache
codemeet
contrib/antilink
contrib/antispam
contrib/forcejoin
contrib/gatekeeper
contrib/maintenancemode
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
````

سورس قابل بازیابی، packageهای اصلی و مجموعه‌ی packageهای `contrib` موجود در ساختار پروژه را پوشش می‌دهد. Documentation هر package بر اساس API و implementation قابل مشاهده در منابع پروژه تهیه شده است.

در مواردی که بخشی از implementation یا جزئیات یک API در منابع قابل بازیابی نباشد، از ساختن signature، parameter، behavior یا implementation فرضی خودداری شده است.

## Bot API Coverage

قابلیت‌های مستندشده در **CodeMeet Bot API** نیز در Documentation کتابخانه پوشش داده شده‌اند و در بخش‌های مربوطه به API سطح کتابخانه ارتباط داده شده‌اند.

این موارد شامل حوزه‌های زیر هستند:

* Messages
* Media
* Keyboards
* Callback Queries
* Bot Commands
* Bot Profile
* Chats
* Groups
* Channels
* Updates
* Long Polling
* Webhook
* Models
* Files
* Error Handling

## API Compatibility

وجود یک قابلیت در Documentation یا تعریف یک Method در کتابخانه، لزوماً به معنی فعال بودن همان قابلیت در نسخه‌ی فعلی CodeMeet Bot API نیست.

برخی APIها ممکن است در کتابخانه با هدف **Future Compatibility** و آماده‌سازی ساختار SDK برای قابلیت‌های آینده تعریف شده باشند، در حالی که endpoint یا قابلیت متناظر آن‌ها هنوز در Bot API فعلی ارائه نشده باشد.

در چنین مواردی Documentation وضعیت قابلیت را به‌صورت شفاف مشخص می‌کند و از معرفی APIهای آینده به‌عنوان قابلیت قطعی نسخه‌ی فعلی خودداری می‌شود.

به‌عنوان اصل کلی:

> **Bot API فعال، منبع نهایی تعیین‌کننده‌ی پشتیبانی واقعی یک قابلیت در زمان اجرا است.**

اگر یک Method در کتابخانه وجود داشته باشد اما Bot API فعلی endpoint یا schema متناظر آن را ارائه نکند، ممکن است اجرای آن درخواست با خطای API مواجه شود.

این موارد به‌عنوان **Compatibility Limitation** یا **Future / Prepared API** در نظر گرفته می‌شوند و در نسخه‌های بعدی کتابخانه، در صورت انتشار قابلیت مربوطه در Bot API، قابل هماهنگ‌سازی هستند.

## Documentation Scope

این Documentation با هدف پوشش کامل موارد قابل استناد از:

* سورس CodeMeet Go
* ساختار packageها
* API Methods
* Models
* Runtime Components
* قابلیت‌های Bot API
* ابزارهای Production
* packageهای `contrib`
* Examples
* Error Handling
* Performance و Reliability

تهیه شده است.

هرجا اطلاعات کافی در منابع وجود نداشته باشد، Documentation از حدس زدن جزئیات فنی خودداری می‌کند تا مستندات با implementation واقعی پروژه و specification موجود Bot API همخوان باقی بمانند.


