# VPN Detector

برای تشخیص نشانه‌های VPN/Proxy/Config در متن و فایل.

## Defaults

```text
CheckText=true
CheckDocuments=true
DownloadAndScanTextFiles=true
ScanAPKContent=true
ScanZIPContent=true
ScanAllDocuments=false
MaxTextScanBytes=4MB
MaxAPKScanBytes=32MB
MaxZIPScanBytes=32MB
MaxZIPEntryBytes=2MB
MaxZIPEntries=64
MaxBase64Candidates=16
DetectionThreshold=4
```

## Middleware

```go
vd := vpndetector.New(vpndetector.DefaultConfig())
bot.Use(vd.VPNDetectorMiddleware())
```

## Config

شامل flagهای scan، محدودیت اندازه، `DetectionThreshold`، `DownloadAction` و `Action` است.

## Stats

```go
stats := vd.Stats()
```

شمارنده‌های داخلی شامل blocked text/APK/config/proxy/VPN، scanned files، content/ZIP/base64 scans، metadata/content hits، skipped files و errors هستند.

## DownloadAction

```go
DownloadAction func(ctx context.Context, fileID string) (io.ReadCloser, error)
```

برای دریافت محتوای فایل قبل از scan.

## نکته امنیتی

Detector نتیجه قطعی برای همه ورودی‌ها نیست. Threshold و محدودیت فایل را با workload واقعی تنظیم کنید.
