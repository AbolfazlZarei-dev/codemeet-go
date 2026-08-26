package vpndetector

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/AbolfazlZarei-dev/codemeet-go/dispatcher"
	"github.com/AbolfazlZarei-dev/codemeet-go/models"
)

type Config struct {
	CheckText                bool
	CheckDocuments           bool
	DownloadAndScanTextFiles bool
	ScanAPKContent           bool
	ScanZIPContent           bool
	ScanAllDocuments         bool

	MaxTextScanBytes    int64
	MaxAPKScanBytes     int64
	MaxZIPScanBytes     int64
	MaxZIPEntryBytes    int64
	MaxZIPEntries       int
	MaxBase64Candidates int
	DetectionThreshold  int

	DownloadAction func(
		ctx context.Context,
		fileID string,
	) (io.ReadCloser, error)

	Action func(
		ctx context.Context,
		userID string,
		chatID string,
		messageID int,
		reason string,
	)
}

func DefaultConfig() Config {
	return Config{
		CheckText:                true,
		CheckDocuments:           true,
		DownloadAndScanTextFiles: true,
		ScanAPKContent:           true,
		ScanZIPContent:           true,
		ScanAllDocuments:         false,
		MaxTextScanBytes:         4 << 20,
		MaxAPKScanBytes:          32 << 20,
		MaxZIPScanBytes:          32 << 20,
		MaxZIPEntryBytes:         2 << 20,
		MaxZIPEntries:            64,
		MaxBase64Candidates:      16,
		DetectionThreshold:       4,
	}
}

type stats struct {
	textBlocked   atomic.Int64
	apkBlocked    atomic.Int64
	configBlocked atomic.Int64
	proxyBlocked  atomic.Int64
	vpnBlocked    atomic.Int64

	scannedFiles atomic.Int64
	contentScans atomic.Int64
	zipScans     atomic.Int64
	base64Scans  atomic.Int64

	metadataHits atomic.Int64
	contentHits  atomic.Int64
	skippedFiles atomic.Int64
	errors       atomic.Int64
}

type VPNDetector struct {
	cfg   Config
	stats stats

	vpnSchemes     []string
	vpnKeywords    []string
	appKeywords    []string
	configKeywords []string
	structureKeys  []string
	strongKeywords []string
	textMIMEs      []string

	configExtensions map[string]struct{}
	skipExtensions   map[string]struct{}
}

func New(cfg Config) *VPNDetector {
	def := DefaultConfig()

	if cfg.MaxTextScanBytes <= 0 {
		cfg.MaxTextScanBytes = def.MaxTextScanBytes
	}
	if cfg.MaxAPKScanBytes <= 0 {
		cfg.MaxAPKScanBytes = def.MaxAPKScanBytes
	}
	if cfg.MaxZIPScanBytes <= 0 {
		cfg.MaxZIPScanBytes = def.MaxZIPScanBytes
	}
	if cfg.MaxZIPEntryBytes <= 0 {
		cfg.MaxZIPEntryBytes = def.MaxZIPEntryBytes
	}
	if cfg.MaxZIPEntries <= 0 {
		cfg.MaxZIPEntries = def.MaxZIPEntries
	}
	if cfg.MaxBase64Candidates <= 0 {
		cfg.MaxBase64Candidates = def.MaxBase64Candidates
	}
	if cfg.DetectionThreshold <= 0 {
		cfg.DetectionThreshold = def.DetectionThreshold
	}

	return &VPNDetector{
		cfg: cfg,

		vpnSchemes: []string{
			"vmess://",
			"vless://",
			"trojan://",
			"ss://",
			"ssr://",
			"socks://",
			"socks4://",
			"socks5://",
			"wireguard://",
			"wg://",
			"hysteria://",
			"hysteria2://",
			"hy2://",
			"tuic://",
		},

		vpnKeywords: []string{
			"vpn",
			"v-p-n",
			"proxy",
			"proxies",
			"proxying",
			"v2ray",
			"xray",
			"sing-box",
			"singbox",
			"clash",
			"mihomo",
			"nekobox",
			"nekoray",
			"v2rayng",
			"shadowrocket",
			"outline",
			"openvpn",
			"wireguard",
			"trojan",
			"hysteria",
			"hysteria2",
			"tuic",
			"vmess",
			"vless",
			"ssr",
			"shadowsocks",
			"socks5",
			"socks4",
			"tun2socks",
			"badvpn",
			"orbot",
			"psiphon",
			"lantern",
			"cloudflare warp",
			"cloudflare-warp",
			"تونل",
			"تونلینگ",
			"وی پی ان",
			"فیلترشکن",
			"پروکسی",
			"پراکسی",
			"کانفیگ",
			"کانفیگ وی پی ان",
			"کانفیگ پروکسی",
			"سرور پروکسی",
			"سرور وی پی ان",
			"شادوکس",
			"کلش",
			"نکوباکس",
			"نکوری",
			"ویلس",
			"وی مس",
			"هستریا",
			"تروجان",
			"وایرگارد",
		},

		appKeywords: []string{
			"v2ray",
			"xray",
			"singbox",
			"sing-box",
			"clash",
			"mihomo",
			"nekobox",
			"nekoray",
			"v2rayng",
			"shadowrocket",
			"shadowsocks",
			"openvpn",
			"wireguard",
			"hysteria",
			"tuic",
			"trojan",
			"tun2socks",
			"vpnservice",
			"proxyservice",
			"badvpn",
			"orbot",
			"psiphon",
			"lantern",
			"warp",
		},

		configKeywords: []string{
			"config",
			"configuration",
			"configs",
			"proxy",
			"proxies",
			"vpn",
			"v2ray",
			"xray",
			"singbox",
			"sing-box",
			"clash",
			"mihomo",
			"nekobox",
			"nekoray",
			"shadowrocket",
			"wireguard",
			"openvpn",
			"shadowsocks",
			"ssr",
			"vmess",
			"vless",
			"trojan",
			"hysteria",
			"tuic",
			"tun",
			"tunnel",
			"کانفیگ",
			"پروکسی",
			"وی پی ان",
			"فیلترشکن",
			"سرور",
		},

		structureKeys: []string{
			"outbounds",
			"inbounds",
			"proxy-groups",
			"proxies:",
			"listeners",
			"routing",
			"rules",
			"mixed-port",
			"socks-port",
			"redir-port",
			"tun:",
			"streamsettings",
			"stream_settings",
			"reality",
			"realitysettings",
			"wssettings",
			"grpcsettings",
			"tlssettings",
		},

		strongKeywords: []string{
			"vpnservice",
			"tun2socks",
			"vmess://",
			"vless://",
			"trojan://",
			"hysteria://",
			"hysteria2://",
			"tuic://",
			"wireguard",
			"openvpn",
			"shadowsocks",
			"sing-box",
			"singbox",
			"mihomo",
			"v2ray",
			"xray",
		},

		textMIMEs: []string{
			"text/",
			"application/json",
			"application/yaml",
			"application/x-yaml",
			"application/toml",
			"application/xml",
			"application/x-config",
		},

		configExtensions: map[string]struct{}{
			".json":      {},
			".jsonc":     {},
			".yaml":      {},
			".yml":       {},
			".toml":      {},
			".conf":      {},
			".cfg":       {},
			".ini":       {},
			".ovpn":      {},
			".wg":        {},
			".wireguard": {},
			".socks":     {},
			".proxy":     {},
			".list":      {},
			".lst":       {},
			".txt":       {},
		},

		skipExtensions: map[string]struct{}{
			".jpg":  {},
			".jpeg": {},
			".png":  {},
			".gif":  {},
			".webp": {},
			".bmp":  {},
			".ico":  {},
			".mp3":  {},
			".wav":  {},
			".ogg":  {},
			".flac": {},
			".mp4":  {},
			".mkv":  {},
			".avi":  {},
			".mov":  {},
			".webm": {},
			".rar":  {},
			".7z":   {},
			".tar":  {},
			".gz":   {},
			".pdf":  {},
			".doc":  {},
			".docx": {},
			".xls":  {},
			".xlsx": {},
			".ppt":  {},
			".pptx": {},
			".exe":  {},
			".dll":  {},
			".msi":  {},
			".iso":  {},
		},
	}
}

func (vd *VPNDetector) VPNDetectorMiddleware() dispatcher.MiddlewareFunc {
	return func(next dispatcher.HandlerFunc) dispatcher.HandlerFunc {
		return func(ctx context.Context, u *models.Update) {
			if u == nil || u.Message == nil {
				next(ctx, u)
				return
			}

			msg := u.Message

			userID := ""
			if msg.From != nil {
				userID = msg.From.ID
			}

			chatID := msg.Chat.ID
			messageID := msg.MessageID

			if vd.cfg.CheckText {
				text := msg.Text
				if text == "" {
					text = msg.Caption
				}

				if text != "" {
					score, reason := vd.scoreText(text)

					if score >= vd.cfg.DetectionThreshold {
						vd.stats.textBlocked.Add(1)

						if hasProxyIndicator(reason) {
							vd.stats.proxyBlocked.Add(1)
						} else {
							vd.stats.vpnBlocked.Add(1)
						}

						vd.triggerAction(
							ctx,
							userID,
							chatID,
							messageID,
							reason,
						)
						return
					}
				}
			}

			if !vd.cfg.CheckDocuments || msg.Document == nil {
				next(ctx, u)
				return
			}

			doc := msg.Document
			fileName := strings.TrimSpace(doc.FileName)
			mimeType := strings.TrimSpace(doc.MimeType)

			if isAPK(fileName, mimeType) {
				score, reason := vd.scoreAPKName(fileName, mimeType)

				if score >= vd.cfg.DetectionThreshold {
					vd.stats.apkBlocked.Add(1)
					vd.triggerAction(ctx, userID, chatID, messageID, reason)
					return
				}

				if vd.cfg.ScanAPKContent && vd.cfg.DownloadAction != nil {
					vd.stats.scannedFiles.Add(1)
					vd.stats.contentScans.Add(1)

					detected, reason, err := vd.scanAPK(
						ctx,
						doc.FileID,
						fileName,
						mimeType,
					)

					if err != nil {
						vd.stats.errors.Add(1)
					}

					if detected {
						vd.stats.apkBlocked.Add(1)
						vd.stats.contentHits.Add(1)

						vd.triggerAction(
							ctx,
							userID,
							chatID,
							messageID,
							reason,
						)
						return
					}
				}

				next(ctx, u)
				return
			}

			nameScore, nameReason := vd.scoreConfigName(
				fileName,
				mimeType,
			)

			if nameScore > 0 {
				vd.stats.metadataHits.Add(1)
			}

			if isZIPName(fileName, mimeType) {
				if vd.cfg.ScanZIPContent && vd.cfg.DownloadAction != nil {
					vd.stats.scannedFiles.Add(1)
					vd.stats.zipScans.Add(1)

					detected, reason, err := vd.scanZIP(
						ctx,
						doc.FileID,
						fileName,
						mimeType,
					)

					if err != nil {
						vd.stats.errors.Add(1)
					}

					if detected {
						vd.stats.configBlocked.Add(1)
						vd.stats.contentHits.Add(1)

						vd.triggerAction(
							ctx,
							userID,
							chatID,
							messageID,
							reason,
						)
						return
					}
				}

				if nameScore >= vd.cfg.DetectionThreshold {
					vd.stats.configBlocked.Add(1)

					vd.triggerAction(
						ctx,
						userID,
						chatID,
						messageID,
						nameReason,
					)
					return
				}

				next(ctx, u)
				return
			}

			if vd.shouldScanDocument(fileName, mimeType, nameScore) &&
				vd.cfg.DownloadAction != nil {

				vd.stats.scannedFiles.Add(1)
				vd.stats.contentScans.Add(1)

				detected, contentScore, contentReason, err :=
					vd.scanTextFile(ctx, doc.FileID)

				if err != nil {
					vd.stats.errors.Add(1)
				}

				totalScore := nameScore + contentScore

				reason := nameReason
				if contentReason != "" {
					reason = contentReason
				}

				if detected || totalScore >= vd.cfg.DetectionThreshold {
					vd.stats.configBlocked.Add(1)
					vd.stats.contentHits.Add(1)

					vd.triggerAction(
						ctx,
						userID,
						chatID,
						messageID,
						reason,
					)
					return
				}
			} else if vd.cfg.ScanAllDocuments {
				vd.stats.skippedFiles.Add(1)
			}

			if nameScore >= vd.cfg.DetectionThreshold {
				vd.stats.configBlocked.Add(1)

				vd.triggerAction(
					ctx,
					userID,
					chatID,
					messageID,
					nameReason,
				)
				return
			}

			next(ctx, u)
		}
	}
}

func (vd *VPNDetector) scoreText(text string) (int, string) {
	text = normalizeText(text)
	if text == "" {
		return 0, ""
	}

	lower := strings.ToLower(text)
	score := 0
	reasons := make([]string, 0, 5)

	if containsAny(lower, vd.vpnSchemes) {
		score += 9
		reasons = append(reasons, "لینک VPN یا پروکسی")
	}

	if hits := countAny(lower, vd.strongKeywords, 4); hits > 0 {
		score += hits * 5
		reasons = append(reasons, "نشانه قوی VPN یا پروکسی")
	}

	if hits := countAny(lower, vd.vpnKeywords, 5); hits > 0 {
		score += hits * 3
		reasons = append(reasons, "عبارت مرتبط با VPN یا پروکسی")
	}

	if hits := countAny(lower, vd.configKeywords, 3); hits > 0 {
		score += hits * 3
		reasons = append(reasons, "نشانه‌های کانفیگ یا پروکسی")
	}

	if containsConfigMarker(lower) {
		score += 4
		reasons = append(reasons, "نشانه‌های تنظیمات پروکسی")
	}

	if hasHostPortPattern(lower) {
		score += 2
		reasons = append(reasons, "ساختار سرور و پورت")
	}

	if score == 0 {
		return 0, ""
	}

	return score, "محتوای مشکوک: " + strings.Join(reasons, "، ")
}

func (vd *VPNDetector) scoreAPKName(fileName, mimeType string) (int, string) {
	lower := strings.ToLower(normalizeText(fileName))

	score := 0

	if hits := countAny(lower, vd.appKeywords, 3); hits > 0 {
		score += hits * 4
	}

	if hits := countAny(lower, vd.configKeywords, 2); hits > 0 {
		score += hits * 2
	}

	if strings.HasSuffix(lower, ".apk") {
		score++
	}

	if strings.Contains(strings.ToLower(mimeType), "android") {
		score++
	}

	if score >= vd.cfg.DetectionThreshold {
		return score, "APK مشکوک به VPN یا پروکسی: " + fileName
	}

	return score, ""
}

func (vd *VPNDetector) scoreConfigName(fileName, mimeType string) (int, string) {
	name := normalizeText(fileName)
	lower := strings.ToLower(name)

	score := 0
	reasons := make([]string, 0, 3)

	ext := strings.ToLower(filepath.Ext(lower))

	_, isConfigExt := vd.configExtensions[ext]
	if isConfigExt {
		score += 2
		reasons = append(reasons, "پسوند فایل کانفیگ")
	}

	configHits := countAny(lower, vd.configKeywords, 3)
	if configHits > 0 {
		score += configHits * 3
		reasons = append(reasons, "نام فایل مرتبط با VPN یا پروکسی")
	}

	mimeLower := strings.ToLower(strings.TrimSpace(mimeType))
	if isTextMIME(mimeLower, vd.textMIMEs) {
		score++
		reasons = append(reasons, "نوع فایل متنی یا کانفیگ")
	}

	if isConfigExt && configHits > 0 {
		score += 5
	}

	if score == 0 {
		return 0, ""
	}

	return score,
		"فایل مشکوک به کانفیگ VPN یا پروکسی: " +
			fileName +
			" (" +
			strings.Join(reasons, "، ") +
			")"
}

func (vd *VPNDetector) shouldScanDocument(
	fileName string,
	mimeType string,
	nameScore int,
) bool {
	if !vd.cfg.DownloadAndScanTextFiles {
		return false
	}

	name := strings.ToLower(strings.TrimSpace(fileName))
	mime := strings.ToLower(strings.TrimSpace(mimeType))

	if vd.isConfigFile(name, mime) {
		return true
	}

	if nameScore >= 2 {
		return true
	}

	if isTextMIME(mime, vd.textMIMEs) {
		return true
	}

	if !vd.cfg.ScanAllDocuments {
		return false
	}

	ext := strings.ToLower(filepath.Ext(name))

	if _, skip := vd.skipExtensions[ext]; skip {
		return false
	}

	return true
}

func (vd *VPNDetector) isConfigFile(fileName, mimeType string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))

	if _, ok := vd.configExtensions[ext]; ok {
		return true
	}

	return isTextMIME(mimeType, vd.textMIMEs)
}

func (vd *VPNDetector) scanTextFile(
	ctx context.Context,
	fileID string,
) (bool, int, string, error) {
	reader, err := vd.cfg.DownloadAction(ctx, fileID)
	if err != nil {
		return false, 0, "", err
	}

	defer reader.Close()

	data, err := readLimitedContext(
		ctx,
		reader,
		vd.cfg.MaxTextScanBytes,
	)
	if err != nil {
		return false, 0, "", err
	}

	if len(data) == 0 {
		return false, 0, "", nil
	}

	return vd.scanContent(data)
}

func (vd *VPNDetector) scanContent(
	data []byte,
) (bool, int, string, error) {
	if len(data) == 0 {
		return false, 0, "", nil
	}

	if looksBinary(data) {
		if !containsStrongBinarySignature(data, vd.strongKeywords) {
			return false, 0, "", nil
		}
	}

	text := normalizeBytes(data)
	if text == "" {
		return false, 0, "", nil
	}

	lower := strings.ToLower(text)

	score := 0
	reasons := make([]string, 0, 7)

	if hits := countAny(lower, vd.vpnSchemes, 4); hits > 0 {
		score += hits * 8
		reasons = append(reasons, "لینک VPN یا پروکسی داخل فایل")

		if score >= vd.cfg.DetectionThreshold {
			return true, score, strings.Join(reasons, "، "), nil
		}
	}

	if hits := countAny(lower, vd.strongKeywords, 5); hits > 0 {
		score += hits * 4
		reasons = append(reasons, "نشانه‌های قوی کانفیگ VPN")

		if score >= vd.cfg.DetectionThreshold {
			return true, score, "کانفیگ VPN یا پروکسی داخل فایل پیدا شد", nil
		}
	}

	if hits := countAny(lower, vd.structureKeys, 8); hits >= 2 {
		score += 5
		reasons = append(reasons, "ساختار کانفیگ VPN یا Proxy")
	}

	if countJSONKeys(lower) >= 2 {
		score += 5
		reasons = append(reasons, "ساختار JSON کانفیگ")
	}

	if countYAMLMarkers(lower) >= 2 {
		score += 5
		reasons = append(reasons, "ساختار YAML کانفیگ")
	}

	if containsConfigMarker(lower) {
		score += 3
		reasons = append(reasons, "کلیدهای شبکه و پروکسی")
	}

	if hasHostPortPattern(lower) {
		score += 2
		reasons = append(reasons, "ساختار Host و Port")
	}

	if score < vd.cfg.DetectionThreshold && vd.containsEncodedVPN(lower) {
		score += 8
		reasons = append(reasons, "کانفیگ Base64 شده")
	}

	if score >= vd.cfg.DetectionThreshold {
		return true,
			score,
			"محتوای مشکوک به VPN یا Proxy: " +
				strings.Join(reasons, "، "),
			nil
	}

	return false, score, strings.Join(reasons, "، "), nil
}

func (vd *VPNDetector) containsEncodedVPN(text string) bool {
	if text == "" {
		return false
	}

	checked := 0
	start := -1

	for i := 0; i < len(text); i++ {
		if isBase64Char(text[i]) {
			if start == -1 {
				start = i
			}
			continue
		}

		if start != -1 {
			if vd.checkBase64Candidate(text[start:i]) {
				return true
			}

			checked++
			if checked >= vd.cfg.MaxBase64Candidates {
				return false
			}

			start = -1
		}
	}

	if start != -1 && checked < vd.cfg.MaxBase64Candidates {
		return vd.checkBase64Candidate(text[start:])
	}

	return false
}

func (vd *VPNDetector) checkBase64Candidate(value string) bool {
	value = strings.TrimSpace(value)

	if len(value) < 32 || len(value) > 8192 {
		return false
	}

	if len(value)%4 == 1 {
		return false
	}

	vd.stats.base64Scans.Add(1)

	decoded, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(value)
	}
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(value)
	}
	if err != nil {
		decoded, err = base64.RawURLEncoding.DecodeString(value)
	}
	if err != nil || len(decoded) < 12 {
		return false
	}

	if looksBinary(decoded) {
		return false
	}

	lower := strings.ToLower(normalizeBytes(decoded))

	if containsAny(lower, vd.strongKeywords) {
		return true
	}

	if containsAny(lower, vd.vpnSchemes) {
		return true
	}

	return countJSONKeys(lower) >= 2 ||
		countYAMLMarkers(lower) >= 2
}

func (vd *VPNDetector) scanAPK(
	ctx context.Context,
	fileID string,
	fileName string,
	mimeType string,
) (bool, string, error) {
	reader, err := vd.cfg.DownloadAction(ctx, fileID)
	if err != nil {
		return false, "", err
	}

	defer reader.Close()

	data, err := readLimitedContext(
		ctx,
		reader,
		vd.cfg.MaxAPKScanBytes,
	)
	if err != nil {
		return false, "", err
	}

	if len(data) < 4 || !bytes.HasPrefix(data, []byte("PK")) {
		return false, "", nil
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return false, "", nil
	}

	score, _ := vd.scoreAPKName(fileName, mimeType)
	checked := 0
	totalRead := int64(0)

	for _, file := range zr.File {
		if checked >= vd.cfg.MaxZIPEntries {
			break
		}

		if err := ctx.Err(); err != nil {
			return false, "", err
		}

		name := strings.ToLower(file.Name)

		if !isImportantAPKEntry(name) {
			continue
		}

		if file.UncompressedSize64 >
			uint64(vd.cfg.MaxZIPEntryBytes) {
			continue
		}

		if int64(file.UncompressedSize64)+totalRead >
			vd.cfg.MaxAPKScanBytes {
			break
		}

		rc, err := file.Open()
		if err != nil {
			continue
		}

		limit := minInt64(
			vd.cfg.MaxZIPEntryBytes,
			2<<20,
		)

		entryData, err := readLimitedContext(ctx, rc, limit)
		_ = rc.Close()

		if err != nil || len(entryData) == 0 {
			continue
		}

		totalRead += int64(len(entryData))
		checked++

		lower := bytes.ToLower(entryData)

		if containsBytesAny(lower, vd.strongKeywords) {
			score += 8
		}
		if bytes.Contains(lower, []byte("vpnservice")) {
			score += 7
		}
		if bytes.Contains(lower, []byte("tun2socks")) {
			score += 7
		}
		if bytes.Contains(lower, []byte("socks5")) {
			score += 3
		}
		if bytes.Contains(lower, []byte("wireguard")) {
			score += 5
		}
		if bytes.Contains(lower, []byte("openvpn")) {
			score += 5
		}
		if bytes.Contains(lower, []byte("proxyservice")) {
			score += 5
		}

		if score >= vd.cfg.DetectionThreshold {
			return true,
				"محتوای APK مشکوک به VPN، پروکسی یا تونل شبکه است: " +
					fileName,
				nil
		}
	}

	return false, "", nil
}

func (vd *VPNDetector) scanZIP(
	ctx context.Context,
	fileID string,
	fileName string,
	mimeType string,
) (bool, string, error) {
	reader, err := vd.cfg.DownloadAction(ctx, fileID)
	if err != nil {
		return false, "", err
	}

	defer reader.Close()

	data, err := readLimitedContext(
		ctx,
		reader,
		vd.cfg.MaxZIPScanBytes,
	)
	if err != nil {
		return false, "", err
	}

	if len(data) < 4 || !bytes.HasPrefix(data, []byte("PK")) {
		return false, "", nil
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return false, "", nil
	}

	checked := 0
	var totalUncompressed int64

	maxTotal := vd.cfg.MaxTextScanBytes * 4

	for _, file := range zr.File {
		if checked >= vd.cfg.MaxZIPEntries {
			break
		}

		if err := ctx.Err(); err != nil {
			return false, "", err
		}

		name := strings.ToLower(file.Name)

		if !isPotentialConfigEntry(name) {
			continue
		}

		if file.UncompressedSize64 >
			uint64(vd.cfg.MaxZIPEntryBytes) {
			continue
		}

		size := int64(file.UncompressedSize64)

		if size > maxTotal-totalUncompressed {
			break
		}

		totalUncompressed += size

		rc, err := file.Open()
		if err != nil {
			continue
		}

		entryData, err := readLimitedContext(
			ctx,
			rc,
			vd.cfg.MaxZIPEntryBytes,
		)

		_ = rc.Close()

		if err != nil {
			continue
		}

		if len(entryData) == 0 {
			continue
		}

		checked++

		detected, score, reason, _ := vd.scanContent(entryData)

		if detected || score >= vd.cfg.DetectionThreshold {
			return true,
				"فایل فشرده حاوی کانفیگ VPN یا پروکسی است: " +
					fileName +
					" / " +
					file.Name +
					" - " +
					reason,
				nil
		}
	}

	return false, "", nil
}

func isPotentialConfigEntry(name string) bool {
	name = strings.ToLower(name)

	ext := filepath.Ext(name)

	switch ext {
	case ".json", ".jsonc", ".yaml", ".yml",
		".toml", ".conf", ".cfg", ".ini",
		".ovpn", ".wg", ".txt":
		return true
	}

	return strings.Contains(name, "config") ||
		strings.Contains(name, "proxy") ||
		strings.Contains(name, "vpn") ||
		strings.Contains(name, "v2ray") ||
		strings.Contains(name, "xray") ||
		strings.Contains(name, "clash") ||
		strings.Contains(name, "mihomo") ||
		strings.Contains(name, "singbox")
}

func isImportantAPKEntry(name string) bool {
	if name == "androidmanifest.xml" {
		return true
	}

	if strings.HasSuffix(name, ".dex") ||
		strings.HasSuffix(name, ".so") {
		return true
	}

	return strings.HasPrefix(name, "assets/") ||
		strings.HasPrefix(name, "res/raw/")
}

func isTextMIME(mime string, list []string) bool {
	if mime == "" {
		return false
	}

	for _, prefix := range list {
		if strings.HasSuffix(prefix, "/") {
			if strings.HasPrefix(mime, prefix) {
				return true
			}
		} else if strings.HasPrefix(mime, prefix) {
			return true
		}
	}

	return false
}

func isAPK(fileName, mimeType string) bool {
	if strings.EqualFold(filepath.Ext(fileName), ".apk") {
		return true
	}

	return strings.Contains(
		strings.ToLower(mimeType),
		"android.package-archive",
	)
}

func isZIPName(fileName, mimeType string) bool {
	if strings.EqualFold(filepath.Ext(fileName), ".zip") {
		return true
	}

	return strings.Contains(
		strings.ToLower(mimeType),
		"zip",
	)
}

func containsAny(text string, items []string) bool {
	for _, item := range items {
		if strings.Contains(text, item) {
			return true
		}
	}

	return false
}

func countAny(text string, items []string, max int) int {
	count := 0

	for _, item := range items {
		if strings.Contains(text, item) {
			count++

			if count >= max {
				return count
			}
		}
	}

	return count
}

func containsBytesAny(data []byte, items []string) bool {
	for _, item := range items {
		if bytes.Contains(data, []byte(item)) {
			return true
		}
	}

	return false
}

func countJSONKeys(text string) int {
	keys := [...]string{
		`"server"`,
		`"address"`,
		`"host"`,
		`"port"`,
		`"protocol"`,
		`"network"`,
		`"outbounds"`,
		`"inbounds"`,
		`"streamsettings"`,
		`"proxy-groups"`,
		`"proxies"`,
		`"dns"`,
		`"routing"`,
		`"tls"`,
		`"reality"`,
		`"sni"`,
		`"uuid"`,
		`"password"`,
		`"cipher"`,
	}

	count := 0

	for _, key := range keys {
		if strings.Contains(text, key) {
			count++

			if count >= 6 {
				return count
			}
		}
	}

	return count
}

func countYAMLMarkers(text string) int {
	keys := [...]string{
		"proxies:",
		"proxy-groups:",
		"listeners:",
		"outbounds:",
		"inbounds:",
		"mixed-port:",
		"socks-port:",
		"redir-port:",
		"tun:",
		"rules:",
		"dns:",
		"server:",
		"port:",
		"uuid:",
		"cipher:",
		"password:",
	}

	count := 0

	for _, key := range keys {
		if strings.Contains(text, key) {
			count++

			if count >= 6 {
				return count
			}
		}
	}

	return count
}

func containsConfigMarker(text string) bool {
	return strings.Contains(text, "server=") ||
		strings.Contains(text, "host=") ||
		strings.Contains(text, "address=") ||
		strings.Contains(text, "port=") ||
		strings.Contains(text, "server:") ||
		strings.Contains(text, "host:") ||
		strings.Contains(text, "address:") ||
		strings.Contains(text, "port:") ||
		strings.Contains(text, "proxy=") ||
		strings.Contains(text, "proxy:") ||
		strings.Contains(text, "socks5://") ||
		strings.Contains(text, "socks://") ||
		strings.Contains(text, "tun:") ||
		strings.Contains(text, "mixed-port:") ||
		strings.Contains(text, "socks-port:") ||
		strings.Contains(text, "redir-port:")
}

func hasHostPortPattern(text string) bool {
	hits := 0

	if strings.Contains(text, "server") {
		hits++
	}
	if strings.Contains(text, "host") {
		hits++
	}
	if strings.Contains(text, "port") {
		hits++
	}
	if strings.Contains(text, "address") {
		hits++
	}

	return hits >= 2
}

func looksBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	n := len(data)
	if n > 8192 {
		n = 8192
	}

	control := 0

	for i := 0; i < n; i++ {
		c := data[i]

		if c == 0 {
			return true
		}

		if c < 9 || (c >= 14 && c < 32) {
			control++
		}
	}

	return control > n/100
}

func containsStrongBinarySignature(
	data []byte,
	keywords []string,
) bool {
	if len(data) == 0 {
		return false
	}

	for _, keyword := range keywords {
		if bytes.Contains(data, []byte(keyword)) {
			return true
		}
	}

	return false
}

func isBase64Char(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '+' ||
		c == '/' ||
		c == '-' ||
		c == '_' ||
		c == '='
}

func normalizeText(text string) string {
	text = strings.TrimSpace(text)

	if text == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(text))

	changed := false

	for _, r := range text {
		switch r {
		case '\u200c':
			b.WriteByte(' ')
			changed = true
		case '\u200d', '\ufeff':
			changed = true
		case 'ي', 'ى':
			b.WriteRune('ی')
			changed = true
		case 'ك':
			b.WriteRune('ک')
			changed = true
		case 'ۀ':
			b.WriteRune('ه')
			changed = true
		default:
			b.WriteRune(r)
		}
	}

	if !changed {
		return text
	}

	return strings.TrimSpace(b.String())
}

func normalizeBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	return normalizeText(string(data))
}

func readLimitedContext(
	ctx context.Context,
	reader io.Reader,
	maxBytes int64,
) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(
		io.LimitReader(reader, maxBytes),
	)

	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

func hasProxyIndicator(reason string) bool {
	lower := strings.ToLower(reason)

	return strings.Contains(lower, "proxy") ||
		strings.Contains(lower, "پروکسی") ||
		strings.Contains(lower, "لینک")
}

func (vd *VPNDetector) triggerAction(
	ctx context.Context,
	userID string,
	chatID string,
	messageID int,
	reason string,
) {
	if vd.cfg.Action == nil {
		return
	}

	vd.cfg.Action(
		ctx,
		userID,
		chatID,
		messageID,
		reason,
	)
}

func (vd *VPNDetector) Stats() map[string]int64 {
	return map[string]int64{
		"text_blocked":   vd.stats.textBlocked.Load(),
		"apk_blocked":    vd.stats.apkBlocked.Load(),
		"config_blocked": vd.stats.configBlocked.Load(),
		"proxy_blocked":  vd.stats.proxyBlocked.Load(),
		"vpn_blocked":    vd.stats.vpnBlocked.Load(),

		"scanned_files": vd.stats.scannedFiles.Load(),
		"content_scans": vd.stats.contentScans.Load(),
		"zip_scans":     vd.stats.zipScans.Load(),
		"base64_scans":  vd.stats.base64Scans.Load(),

		"metadata_hits": vd.stats.metadataHits.Load(),
		"content_hits":  vd.stats.contentHits.Load(),
		"skipped_files": vd.stats.skippedFiles.Load(),
		"errors":        vd.stats.errors.Load(),
	}
}

var errInvalidContent = errors.New("invalid content")
