package main

import (
	"fmt"
	"grafana-screenshot/logs"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/robfig/cron/v3"
)

const iso8601 = "2006-01-02T15:04:05.000Z"

// 截图函数，返回文件路径
func captureScreenshot(cfg *Config, dash Dashboard, day time.Time) string {
	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	// 构造上海时间的前一天0点和24点
	fromShanghai := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, shanghaiLoc)
	toShanghai := fromShanghai.Add(24 * time.Hour)

	// 转换为UTC时间
	fromUTC := fromShanghai.UTC()
	toUTC := toShanghai.UTC()

	fromISO := fromUTC.Format(iso8601)
	toISO := toUTC.Format(iso8601)
	orgID := dash.OrgID
	if orgID == 0 {
		orgID = cfg.OrgID
	}

	url := fmt.Sprintf(
		"%s/render/d/%s/%s?orgId=%d&from=%s&to=%s&timezone=Asia/Shanghai&width=1920&height=1000&kiosk=1&fullPage=true",
		cfg.BaseURL, dash.DashboardUID, dash.Slug, orgID,
		fromISO, toISO,
	)
	client := &http.Client{Timeout: 180 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Accept", "image/png")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		logs.ErrorLogger.Printf("❌ %s 截图失败\n", dash.Name)
		return ""
	}
	defer resp.Body.Close()

	os.MkdirAll("screenshots", 0755)

	filePath := fmt.Sprintf("screenshots/%s-%s.png", dash.Name, day.Format("2006-01-02"))
	out, _ := os.Create(filePath)
	io.Copy(out, resp.Body)
	out.Close()
	logs.InfoLogger.Printf("%s OK", dash.Name)
	return filePath
}

// 生成 PDF
func createPDFReport(date string, dashboards []Dashboard) {
	// 按月份创建目录，如 2025-11
	monthDir := date[:7]
	if err := os.MkdirAll(monthDir, 0755); err != nil {
		logs.ErrorLogger.Println("❌ 创建目录失败:", err)
		return
	}

	// PDF 存在月份目录中
	filePath := fmt.Sprintf("%s/%s.pdf", monthDir, date)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("zh", "", "simhei.ttf")
	pdf.SetFont("zh", "", 14)
	pdf.AddPage()

	title := "流量日报 - " + date
	pdf.CellFormat(0, 12, title, "", 1, "C", false, 0, "")
	pdf.Ln(4)

	for _, dash := range dashboards {
		pdf.SetFont("zh", "", 12)
		pdf.CellFormat(0, 8, dash.Name, "", 1, "L", false, 0, "")

		img := fmt.Sprintf("screenshots/%s-%s.png", dash.Name, date)
		pdf.ImageOptions(
			img, 10, pdf.GetY(), 190, 0, false,
			fpdf.ImageOptions{ImageType: "PNG"}, 0, "",
		)
		pdf.Ln(110)

		if pdf.GetY() > 260 {
			pdf.AddPage()
			pdf.SetFont("zh", "", 14)
		}
	}

	if err := pdf.OutputFileAndClose(filePath); err != nil {
		logs.ErrorLogger.Println("❌ PDF 生成失败:", err)
	} else {
		logs.InfoLogger.Println("📄 PDF 已生成:", filePath)
	}
}

// 完整流程
func runOnce(cfg *Config) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	yesterday := now.AddDate(0, 0, -1)
	day := time.Date(
		yesterday.Year(), yesterday.Month(), yesterday.Day(),
		0, 0, 0, 0, loc,
	)
	dateStr := day.Format("2006-01-02")

	logs.InfoLogger.Println("📆", dateStr, "日报开始")

	var images []string
	for _, dash := range cfg.Dashboards {
		img := captureScreenshot(cfg, dash, day)
		if img != "" {
			images = append(images, img)
		}
	}

	createPDFReport(dateStr, cfg.Dashboards)

	// 删除截图
	for _, img := range images {
		_ = os.Remove(img)
	}
	logs.InfoLogger.Println("🧹 已清理截图文件")
	// 企业微信推送
	err := SendWeChatBotMessage(
		cfg.WeChatBotKey,
		fmt.Sprintf("📄 *Grafana 每日报告已生成*\n日期：%s", dateStr),
	)

	if err != nil {
		logs.ErrorLogger.Println("❌ 企业微信推送失败:", err)
	} else {
		logs.InfoLogger.Println("🤖 企业微信消息已推送")
	}
	logs.InfoLogger.Println("🎉 完成")
}

// 主入口
func main() {
	logs.CreateLog()
	cfg, err := LoadConfig()
	if err != nil {
		logs.ErrorLogger.Println("❌ 配置加载失败:", err)
		return
	}

	if cfg.DevMode {
		logs.InfoLogger.Println("🧪 DevMode=true => 立即执行")
		runOnce(cfg)
		return
	}

	logs.InfoLogger.Println("⏱ Cron 启动 =", cfg.CronTime)

	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc(cfg.CronTime, func() {
		runOnce(cfg)
	})
	if err != nil {
		logs.ErrorLogger.Println("❌ Cron 错误:", err)
		return
	}

	c.Start()
	select {}
}
