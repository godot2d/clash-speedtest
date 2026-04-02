package ip

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// DefaultNameTemplate is the built-in format when -rename-template is not set.
const DefaultNameTemplate = `{{.OriginalName}} | ⬇️ {{.DownloadSpeedMBps}}MB/s | ⚡{{.LatencyMs}}ms | 📦{{.PacketLoss}}%`

// NodeNameData is the data passed to the rename template.
type NodeNameData struct {
	OriginalName      string // original proxy name
	Flag              string // country flag emoji (empty when location is not queried)
	CountryCode       string // e.g. US, HK (empty when location is not queried)
	Index             string // padded global sequence number, e.g. 001
	Direction         string // ⬇️, ⬆️, or ⚡
	Speed             string // primary metric value
	SpeedUnit         string // MB/s or ms
	LatencyMs         string // latency in milliseconds
	DownloadSpeedMBps string // download MB/s
	UploadSpeedMBps   string // upload MB/s
	PacketLoss        string // packet loss percentage, e.g. 0.0
}

// GenerateNodeNameFromTemplate renders name from a text/template. Placeholders:
// {{.OriginalName}}, {{.Index}}, {{.Direction}}, {{.Speed}}, {{.SpeedUnit}},
// {{.LatencyMs}}, {{.DownloadSpeedMBps}}, {{.UploadSpeedMBps}}, {{.PacketLoss}}.
// {{.Flag}} and {{.CountryCode}} are empty unless set externally.
// If template is empty, DefaultNameTemplate is used. On execute error, falls back to default format.
func GenerateNodeNameFromTemplate(tmpl string, originalName string, latency time.Duration, downloadSpeed, uploadSpeed float64, packetLoss float64, nameCount map[string]int) (string, error) {
	if tmpl == "" {
		tmpl = DefaultNameTemplate
	}
	t, err := template.New("name").Parse(tmpl)
	if err != nil {
		return "", err
	}
	data := buildNodeNameData(originalName, latency, downloadSpeed, uploadSpeed, packetLoss, nameCount)
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Sprintf("%s | %s %s%s | ⚡%sms | 📦%s%%", data.OriginalName, data.Direction, data.Speed, data.SpeedUnit, data.LatencyMs, data.PacketLoss), nil
	}
	return buf.String(), nil
}

func buildNodeNameData(originalName string, latency time.Duration, downloadSpeed, uploadSpeed float64, packetLoss float64, nameCount map[string]int) NodeNameData {
	speed := downloadSpeed
	direction := "⬇️"
	speedUnit := "MB/s"
	if downloadSpeed <= 0 {
		speed = uploadSpeed
		direction = "⬆️"
	}
	if downloadSpeed <= 0 && uploadSpeed <= 0 && latency > 0 {
		speed = float64(latency.Milliseconds())
		direction = "⚡"
		speedUnit = "ms"
	}
	speedMBps := speed / (1024 * 1024)
	if speedUnit == "ms" {
		speedMBps = speed
	}
	count := nameCount[""] + 1
	nameCount[""] = count
	dlMBps := downloadSpeed / (1024 * 1024)
	ulMBps := uploadSpeed / (1024 * 1024)
	latencyMs := "N/A"
	if latency > 0 {
		latencyMs = fmt.Sprintf("%d", latency.Milliseconds())
	}
	return NodeNameData{
		OriginalName:      originalName,
		Flag:              "",
		CountryCode:       "",
		Index:             fmt.Sprintf("%03d", count),
		Direction:         direction,
		Speed:             fmt.Sprintf("%.2f", speedMBps),
		SpeedUnit:         speedUnit,
		LatencyMs:         latencyMs,
		DownloadSpeedMBps: fmt.Sprintf("%.2f", dlMBps),
		UploadSpeedMBps:   fmt.Sprintf("%.2f", ulMBps),
		PacketLoss:        fmt.Sprintf("%.1f", packetLoss),
	}
}

func GenerateNodeName(originalName string, latency time.Duration, downloadSpeed float64, uploadSpeed float64, packetLoss float64, nameCount map[string]int) string {
	name, _ := GenerateNodeNameFromTemplate("", originalName, latency, downloadSpeed, uploadSpeed, packetLoss, nameCount)
	return name
}
