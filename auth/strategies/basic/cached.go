package basic

import (
	"context"
	"net/http"

	"github.com/Sosivio/go-guardian/v2/auth"
	"github.com/Sosivio/go-guardian/v2/auth/internal"
)

// ExtensionKey represents a key for the password in info extensions.
// Typically used when basic strategy cache the authentication decisions.
const ExtensionKey = "x-go-guardian-basic-password"

// NewCached return new auth.Strategy.
// The returned strategy, caches the invocation result of authenticate function.
func NewCached(f AuthenticateFunc, cache auth.Cache, opts ...auth.Option) auth.Strategy {
	cb := new(cachedBasic)
	cb.fn = f
	cb.cache = cache
	cb.comparator = plainText{}
	cb.hasher = internal.PlainTextHasher{}
	for _, opt := range opts {
		opt.Apply(cb)
	}
	return New(cb.authenticate, opts...)
}

type cachedBasic struct {
	fn         AuthenticateFunc
	comparator Comparator
	cache      auth.Cache
	hasher     internal.Hasher
}

func (c *cachedBasic) authenticate(ctx context.Context, r *http.Request, userName, pass string) (auth.Info, error) { // nolint:lll
	hash := c.hasher.Hash(userName)
	v, ok := c.cache.Load(hash)

	// if info not found invoke user authenticate function
	if !ok {
		return c.authenticatAndHash(ctx, r, hash, userName, pass)
	}

	if _, ok := v.(auth.Info); !ok {
		return nil, auth.NewTypeError("strategies/basic:", (*auth.Info)(nil), v)
	}

	info := v.(auth.Info)
	ext := info.GetExtensions()

	if !ext.Has(ExtensionKey) {
		return c.authenticatAndHash(ctx, r, hash, userName, pass)
	}

	return info, c.comparator.Compare(ext.Get(ExtensionKey), pass)
}

func (c *cachedBasic) authenticatAndHash(ctx context.Context, r *http.Request, hash string, userName, pass string) (auth.Info, error) { //nolint:lll
	info, err := c.fn(ctx, r, userName, pass)
	if err != nil {
		return nil, err
	}

	hashedPass, _ := c.comparator.Hash(pass)
	info.GetExtensions().Set(ExtensionKey, hashedPass)
	c.cache.Store(hash, info)

	return info, nil
}
