package accounts_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/wnddd839/codebuddy-proxy/internal/accounts"
)

func TestEnabledDefaultsTrueWhenMissing(t *testing.T) {
	var account accounts.Account
	if err := json.Unmarshal([]byte(`{"bearerToken":"eyJhbGciOiJIUzI1NiJ9.e30.xx","label":"demo"}`), &account); err != nil {
		t.Fatal(err)
	}
	normalized := accounts.CreateAccount(account)
	if !normalized.Enabled {
		t.Fatalf("expected enabled=true by default, got false")
	}
	if !accounts.HasCredentials(normalized) {
		t.Fatal("expected credentials")
	}
}

func TestEnabledFalseRespected(t *testing.T) {
	var account accounts.Account
	if err := json.Unmarshal([]byte(`{"enabled":false,"bearerToken":"eyJhbGciOiJIUzI1NiJ9.e30.xx","id":"abc","createdAt":1}`), &account); err != nil {
		t.Fatal(err)
	}
	normalized := accounts.NormalizeAccount(account, 1)
	if normalized.Enabled {
		t.Fatalf("expected enabled=false")
	}
}

func TestPoolRoundTripSelectAndMark(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "accounts.json")
	pool := accounts.NewPool(path)
	account := accounts.CreateAccount(accounts.Account{
		Label:       "one",
		BearerToken: "eyJhbGciOiJIUzI1NiJ9.e30.one",
		Site:        "domestic",
	})
	if _, _, err := pool.Upsert(account); err != nil {
		t.Fatal(err)
	}
	selection, err := pool.Select(accounts.SelectOptions{Site: "domestic"})
	if err != nil {
		t.Fatal(err)
	}
	if selection.Account.Label != "one" {
		t.Fatalf("unexpected account: %+v", selection.Account)
	}
	if err := pool.MarkResult(selection, true, ""); err != nil {
		t.Fatal(err)
	}
	store, err := pool.Read()
	if err != nil {
		t.Fatal(err)
	}
	if store.Accounts[0].SuccessRequests != 1 {
		t.Fatalf("expected successRequests=1, got %d", store.Accounts[0].SuccessRequests)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestSummarizeStoreForSitePrefersActiveRegion(t *testing.T) {
	store := accounts.EmptyStore()
	store.Accounts = []accounts.Account{
		accounts.CreateAccount(accounts.Account{Label: "cn", Site: "domestic", BearerToken: "cn-token", Enabled: true}),
		accounts.CreateAccount(accounts.Account{Label: "us", Site: "global", BearerToken: "us-token", Enabled: true}),
	}
	summary := accounts.SummarizeStoreForSite(store, "mem", "global")
	if summary.ActiveSite != "global" {
		t.Fatalf("activeSite=%q", summary.ActiveSite)
	}
	if summary.DomesticCount != 1 || summary.GlobalCount != 1 {
		t.Fatalf("counts domestic=%d global=%d", summary.DomesticCount, summary.GlobalCount)
	}
	if summary.Primary == nil || summary.Primary.Site != "global" {
		t.Fatalf("primary=%+v", summary.Primary)
	}
	if summary.ActiveEnabledCount != 1 {
		t.Fatalf("activeEnabled=%d", summary.ActiveEnabledCount)
	}
}
