package oauth

import (
	"testing"
	"time"

	"github.com/wnddd839/codebuddy-proxy/internal/accounts"
)

func TestResolvePluginBaseURL(t *testing.T) {
	if got := ResolvePluginBaseURL("domestic"); got != "https://www.codebuddy.cn" {
		t.Fatalf("domestic=%q", got)
	}
	if got := ResolvePluginBaseURL("global"); got != "https://www.codebuddy.ai" {
		t.Fatalf("global=%q", got)
	}
}

func TestShouldRefresh(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	acc := accounts.Account{RefreshToken: "r", TokenExpiresAt: now.Add(2 * time.Minute).UnixMilli()}
	if !ShouldRefresh(acc, false, 5*time.Minute, now) {
		t.Fatal("expected refresh inside window")
	}
	if ShouldRefresh(acc, false, 30*time.Second, now) {
		t.Fatal("did not expect refresh outside window")
	}
	if !ShouldRefresh(acc, true, time.Hour, now) {
		t.Fatal("force should refresh")
	}
	if ShouldRefresh(accounts.Account{TokenExpiresAt: now.Add(time.Minute).UnixMilli()}, false, time.Hour, now) {
		t.Fatal("no refresh token => false")
	}
}

func TestDecodeJWTAndAccountFromTokenData(t *testing.T) {
	token := "eyJhbGciOiAibm9uZSIsICJ0eXAiOiAiSldUIn0.eyJzdWIiOiAidTEiLCAiZW1haWwiOiAiYUBiLmMiLCAiZXhwIjogMjAwMDAwMDAwMCwgIm5hbWUiOiAiQWRhIn0.sig"
	claims := DecodeJWT(token)
	if claims["email"] != "a@b.c" || claims["sub"] != "u1" {
		t.Fatalf("claims=%v", claims)
	}
	acc := AccountFromTokenData(&TokenData{
		BearerToken:  token,
		RefreshToken: "rt",
		ExpiresIn:    3600,
		Domain:       "www.codebuddy.cn",
	}, "domestic", "L")
	if acc.Site != "domestic" || acc.RefreshToken != "rt" || acc.BearerToken != token {
		t.Fatalf("account=%+v", acc)
	}
	if acc.AuthStatus.UserID == "" {
		t.Fatal("expected user id from jwt")
	}
}

func TestIsPendingAndHostOf(t *testing.T) {
	if !isPending(map[string]any{"code": float64(11217)}) {
		t.Fatal("11217 should be pending")
	}
	if !isPending(map[string]any{"msg": "login waiting"}) {
		t.Fatal("waiting message should be pending")
	}
	if isPending(map[string]any{"code": float64(0), "msg": "ok"}) {
		t.Fatal("ok should not be pending")
	}
	if got := hostOf("https://www.codebuddy.cn/v2/x"); got != "www.codebuddy.cn" {
		t.Fatalf("hostOf=%q", got)
	}
}
