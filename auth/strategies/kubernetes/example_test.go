package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/Sosivio/libcache"
	_ "github.com/Sosivio/libcache/idle"

	"github.com/Sosivio/go-guardian/v2/auth/strategies/token"
)

func ExampleNew() {
	cache := libcache.IDLE.NewUnsafe(0)
	kube := New(cache)
	r, _ := http.NewRequest("", "/", nil)
	_, err := kube.Authenticate(r.Context(), r)
	fmt.Println(err != nil)
	// Output:
	// true
}

func ExampleGetAuthenticateFunc() {
	cache := libcache.IDLE.NewUnsafe(0)
	fn := GetAuthenticateFunc()
	kube := token.New(fn, cache)
	r, _ := http.NewRequest("", "/", nil)
	_, err := kube.Authenticate(r.Context(), r)
	fmt.Println(err != nil)
	// Output:
	// true
}

func Example() {
	st := SetServiceAccountToken("Service Account Token")
	cache := libcache.IDLE.NewUnsafe(0)
	kube := New(cache, st)
	r, _ := http.NewRequest("", "/", nil)
	_, err := kube.Authenticate(r.Context(), r)
	fmt.Println(err != nil)
	// Output:
	// true
}
