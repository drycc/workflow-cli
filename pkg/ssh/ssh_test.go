package ssh

import (
	"testing"
)

type pubKey struct {
	key string
	id  string
}

var validKeys = []pubKey{
	{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDIqRuNwhjntdJWuAVykr/873X8zzKo6Ms1Vx70BQx0wir8TpaLZEY6CqqKDrMbHZ2Z0ZV2ZITs9fC81GqIdmmDFXZxNfx+B1lSR3ZpmZpQpprtSCevZXihIgy+ND50Hp/wk+3VU54FxhudIlJgpPjb/o5vQFhiyM3ynR5gH3slWVaq9C0TkgXCnTzukGzSTeL7wYPNmLomkrAS0nk0yRfoUZcwmD++HMgEmYlhTbnMlkB3nxzEf/JQxhY6xCHrbtNbRkINCY21dHrsrr/MvBvD8FoKnKpxHX2+HNXZbe7Xl87L9o1OuXrtR1crvq+r+1fPjaynGir07zr9mgJPxouPL/e4ppxTL//vt1kVkWWkh/B+GyXEmP38bQGYMpEA7cAndlxPlOki35JYwDNn5CENQpDp8F4+JsKIzAF1zmkBIA7ngg0cSHHilNgwZXmX7h+7nngLgFIpP8h7A9fCpAKhHUfFUj+Zgl9Xm44+sZOwVBnVijK326TgVDFTUXjE/Xwny+3ERgYwBfOwOKmusNFnS0XHbmh+qa/+D8qge5bKilq48pKHzwngM/U6OwMxmSXTuHclLLen3Ime30TOiPzAhokrVNz/Z3VAkfBuJHby68SAKUgczUEU81wz5wFEt1n1sIJ5V49KMRGaSWb+eWvW81yA7NkDSjnsMLa/IF/ADQ== arschles@gmail.com", "rsaId"},
	{"ssh-dss AAAAB3NzaC1kc3MAAACBAOKHxk8vLYdr25G+xha1OOjhPX8z/xAeAMbyiS6mVFrSu1mrrEaqXurJ0LVXm9Md6440noZ5j8iscfdJd5wZZ/XUfugnZ7/LNFNP0uRmLVkJsAh6RwgPZQZ8spnUtucwlWM+xOKDdVNXN7DQQp0LqNg8SsBGAJHuYw3Sd9olPGlNAAAAFQDYmFrWj23PlirKoCPjGQvWCgDfRwAAAIA1NHpuFxgo5j7R4qyb1ydKStoqOkREhQNI5kWNw1p8pksX5pMk0mVZY80VNcYw/M8LWONJ5beLJfAKxMhjfal69A7NKeD+YoY/OxT31VbDvm0cWb0RY+acCIMQ+UtfuXG27aZ6txV/AbOfA9AnhuHTyPPOyF07OHwCUS0ubn8aSgAAAIBA2Jm1k2Hxin/AB8C4N7ycpUDpGQBjIhXp69YuOTNeLcFIzCFc6sB91CorTVJdofnj+KeUAl8lIsJcEWvC4683MNewT3qeDwSClM3ojWFh6VuNuphcPKDqteX8WYnrWMJvAWEiRf0nqNNukhl9zAmAMQFc5U3Sl5TQuhc/6Ns9jA== arschles@gmail.com", "dsaId"},
	{"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBMQ/isNQFn2x7g9dIK1N4+mvEa+a01hj2LnZFBad7W+os+wc+UurVxWVoGopc/mjzqezr6vk9jgOjLdYek9T/2w= arschles@gmail.com", "ecdsaId"},
	{"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIORIdG868fEBUKoEqSQZFKfSLoHkSBmW2uXXGaZKEuus arschles@gmail.com", "ed25519Id"},
}

var invalidKeys = []pubKey{
	{"bad-key-type AAAAB3NzaC1yc2EAAAADAQABAAACAQDIqRuNwhjntdJWuAVykr/873X8zzKo6Ms1Vx70BQx0wir8TpaLZEY6CqqKDrMbHZ2Z0ZV2ZITs9fC81GqIdmmDFXZxNfx+B1lSR3ZpmZpQpprtSCevZXihIgy+ND50Hp/wk+3VU54FxhudIlJgpPjb/o5vQFhiyM3ynR5gH3slWVaq9C0TkgXCnTzukGzSTeL7wYPNmLomkrAS0nk0yRfoUZcwmD++HMgEmYlhTbnMlkB3nxzEf/JQxhY6xCHrbtNbRkINCY21dHrsrr/MvBvD8FoKnKpxHX2+HNXZbe7Xl87L9o1OuXrtR1crvq+r+1fPjaynGir07zr9mgJPxouPL/e4ppxTL//vt1kVkWWkh/B+GyXEmP38bQGYMpEA7cAndlxPlOki35JYwDNn5CENQpDp8F4+JsKIzAF1zmkBIA7ngg0cSHHilNgwZXmX7h+7nngLgFIpP8h7A9fCpAKhHUfFUj+Zgl9Xm44+sZOwVBnVijK326TgVDFTUXjE/Xwny+3ERgYwBfOwOKmusNFnS0XHbmh+qa/+D8qge5bKilq48pKHzwngM/U6OwMxmSXTuHclLLen3Ime30TOiPzAhokrVNz/Z3VAkfBuJHby68SAKUgczUEU81wz5wFEt1n1sIJ5V49KMRGaSWb+eWvW81yA7NkDSjnsMLa/IF/ADQ== arschles@gmail.com", "rsaId"},
}

func TestParseValidSSHPubKey(t *testing.T) {
	for _, keyAndID := range validKeys {
		key := keyAndID.key
		id := keyAndID.id
		info, err := ParsePubKey(id, []byte(key))
		if err != nil {
			t.Fatalf("Error parsing well formed key (%s)", err)
		}
		if info.ID != "arschles@gmail.com" {
			t.Fatalf("expected key ID arschles@gmail.com, got %s", info.ID)
		}
		if info.Public != key {
			t.Fatalf("expected key contents %s, got %s", key, info.Public)
		}

	}
}

func TestParseInvalidSSHPubKey(t *testing.T) {
	for _, keyAndID := range invalidKeys {
		key := keyAndID.key
		id := keyAndID.id
		_, err := ParsePubKey(id, []byte(key))
		if err == nil {
			t.Fatalf("Key should be invalid but was not: (%s)", key)
		}
	}
}
