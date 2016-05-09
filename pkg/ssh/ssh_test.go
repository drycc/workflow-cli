package ssh

import (
	"testing"
)

const (
	pubKey      = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDIqRuNwhjntdJWuAVykr/873X8zzKo6Ms1Vx70BQx0wir8TpaLZEY6CqqKDrMbHZ2Z0ZV2ZITs9fC81GqIdmmDFXZxNfx+B1lSR3ZpmZpQpprtSCevZXihIgy+ND50Hp/wk+3VU54FxhudIlJgpPjb/o5vQFhiyM3ynR5gH3slWVaq9C0TkgXCnTzukGzSTeL7wYPNmLomkrAS0nk0yRfoUZcwmD++HMgEmYlhTbnMlkB3nxzEf/JQxhY6xCHrbtNbRkINCY21dHrsrr/MvBvD8FoKnKpxHX2+HNXZbe7Xl87L9o1OuXrtR1crvq+r+1fPjaynGir07zr9mgJPxouPL/e4ppxTL//vt1kVkWWkh/B+GyXEmP38bQGYMpEA7cAndlxPlOki35JYwDNn5CENQpDp8F4+JsKIzAF1zmkBIA7ngg0cSHHilNgwZXmX7h+7nngLgFIpP8h7A9fCpAKhHUfFUj+Zgl9Xm44+sZOwVBnVijK326TgVDFTUXjE/Xwny+3ERgYwBfOwOKmusNFnS0XHbmh+qa/+D8qge5bKilq48pKHzwngM/U6OwMxmSXTuHclLLen3Ime30TOiPzAhokrVNz/Z3VAkfBuJHby68SAKUgczUEU81wz5wFEt1n1sIJ5V49KMRGaSWb+eWvW81yA7NkDSjnsMLa/IF/ADQ== arschles@gmail.com`
	backupKeyID = "mybackup"
)

func TestParseSSHPubKey(t *testing.T) {
	info, err := ParsePubKey(backupKeyID, []byte(pubKey))
	if err != nil {
		t.Fatalf("Error parsing well formed key (%s)", err)
	}
	if info.ID != "arschles@gmail.com" {
		t.Fatalf("expected key ID arschles@gmail.com, got %s", info.ID)
	}
	if info.Public != pubKey {
		t.Fatalf("expected key contents %s, got %s", pubKey, info.Public)
	}
}
