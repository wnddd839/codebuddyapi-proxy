package billing

import (
	"testing"
)

func TestOfficialUsageAndBillingBaseURL(t *testing.T) {
	if got := OfficialUsageURL("domestic"); got != "https://www.codebuddy.cn/profile/plan" {
		t.Fatalf("domestic usage=%q", got)
	}
	if got := OfficialUsageURL("global"); got != "https://www.codebuddy.ai/profile/plan" {
		t.Fatalf("global usage=%q", got)
	}
	t.Setenv("CODEBUDDY_BILLING_BASE_URL", "")
	if got := BillingBaseURL("domestic"); got != "https://www.codebuddy.cn" {
		t.Fatalf("domestic billing=%q", got)
	}
	t.Setenv("CODEBUDDY_BILLING_BASE_URL", "https://billing.example/")
	if got := BillingBaseURL("domestic"); got != "https://billing.example" {
		t.Fatalf("override billing=%q", got)
	}
}

func TestMapDosageNotify(t *testing.T) {
	n := mapDosageNotify(map[string]any{"dosageNotifyCode": float64(0)})
	if n.Level != "ok" || n.DosageNotifyCode != 0 {
		t.Fatalf("notify=%+v", n)
	}
	n = mapDosageNotify(map[string]any{"dosageNotifyCode": float64(3)})
	if n.Level != "bad" {
		t.Fatalf("exhausted=%+v", n)
	}
	n = mapDosageNotify(map[string]any{"dosageNotifyCode": float64(99), "dosageNotifyZh": "自定义"})
	if n.Level != "warn" || n.DosageNotifyZh != "自定义" {
		t.Fatalf("custom=%+v", n)
	}
}

func TestSummarizeResourceAccounts(t *testing.T) {
	credits := summarizeResourceAccounts([]any{
		map[string]any{
			"CapacityRemain": float64(30),
			"CapacitySize":   float64(100),
			"CapacityUsed":   float64(70),
			"CapacityType":   float64(1),
			"PackageName":    "Basic",
			"OriginUnit":     "credits",
		},
	})
	if credits.Remaining == nil || *credits.Remaining != 30 {
		t.Fatalf("remaining=%v", credits.Remaining)
	}
	if credits.Total == nil || *credits.Total != 100 {
		t.Fatalf("total=%v", credits.Total)
	}
	if len(credits.Packages) != 1 {
		t.Fatalf("packages=%d", len(credits.Packages))
	}
}

func TestFormatNumber(t *testing.T) {
	if got := formatNumber(12); got != "12" {
		t.Fatalf("int=%q", got)
	}
	if got := formatNumber(12.5); got != "12.5" {
		t.Fatalf("float=%q", got)
	}
}
