# 34 — `contrib/vpndetector`

VPNDetector برای شناسایی نشانه‌های VPN/Proxy در متن و فایل‌ها طراحی شده است.

## Config

```go
type Config struct {
    CheckText bool
    CheckDocuments bool
    DownloadAndScanTextFiles bool
    ScanAPKContent bool
    ScanZIPContent bool
    ScanAllDocuments bool

    MaxTextScanBytes int64
    MaxAPKScanBytes int64
    MaxZIPScanBytes int64
    MaxZIPEntryBytes int64
    MaxZIPEntries int
    MaxBase64Candidates int
    DetectionThreshold int

    DownloadAction func(ctx context.Context, fileID string) (io.ReadCloser, error)
    Action func(ctx context.Context, userID, chatID string, messageID int, reason string)
}
```

## Defaults

```text
CheckText = true
CheckDocuments = true
DownloadAndScanTextFiles = true
ScanAPKContent = true
ScanZIPContent = true
ScanAllDocuments = false

MaxTextScanBytes = 4 MiB
MaxAPKScanBytes = 32 MiB
MaxZIPScanBytes = 32 MiB
MaxZIPEntryBytes = 2 MiB
MaxZIPEntries = 64
MaxBase64Candidates = 16
DetectionThreshold = 4
```

## تشخیص

سیستم می‌تواند موارد زیر را بررسی کند:

- VPN/proxy schemes
- keywordها
- configuration markers
- host:port patterns
- encoded content
- Base64 candidates
- APK content
- ZIP content
- configuration files

## APK و ZIP

ZIP با محدودیت تعداد entry و اندازه‌ی entry بررسی می‌شود تا resource exhaustion کاهش یابد.

APK نیز به عنوان archive بررسی می‌شود و entryهای مهم با ruleهای داخلی امتیاز می‌گیرند.

## Scoring

تشخیص بر اساس جمع نشانه‌ها و `DetectionThreshold` انجام می‌شود.

## Stats

- textBlocked
- apkBlocked
- configBlocked
- proxyBlocked
- vpnBlocked
- scannedFiles
- contentScans
- zipScans
- base64Scans
- metadataHits
- contentHits
- skippedFiles
- errors

## مثال

```go
detector := vpndetector.New(vpndetector.Config{
    CheckText: true,
    CheckDocuments: true,
    ScanAPKContent: true,
    ScanZIPContent: true,
    DetectionThreshold: 4,
})

bot.Use(detector.VPNDetectorMiddleware())
```

## نکته

این ابزار heuristic است؛ Detection نباید به‌عنوان اثبات قطعی VPN بودن یک کاربر تفسیر شود.
