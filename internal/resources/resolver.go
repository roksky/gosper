package resources

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"strings"

	"github.com/roksky/gosper/internal/diag"
)

// Resolver opens report resources by reference.
//
// Open resolves ref through the resolver's configured root and returns a
// reader the caller owns and must Close. Open honors ctx and returns ctx.Err()
// when the context is already canceled. Open must never return (nil, nil).
//
// By default only relative paths inside the root are accepted; URI schemes
// (file://, http://, ...) and paths that escape the root are rejected. See
// WithAllowedSchemes and New.
type Resolver interface {
	Open(ctx context.Context, ref string) (io.ReadCloser, error)
}

// Option configures a Resolver built by New or FromDir.
type Option func(*config)

type config struct {
	allowedSchemes map[string]struct{}
}

// WithAllowedSchemes permits references carrying one of the named URI schemes.
// Scheme names are case-insensitive; a trailing colon is ignored, so "file",
// "File" and "file:" are equivalent. No scheme is permitted by default
// (OPS-002).
//
// Only the "file" scheme currently resolves, and its path is interpreted
// relative to the resolver root. Other allowed schemes are still rejected with
// code diag.CodeUnsupported until a transport is configured (DAT-011,
// OPS-002).
func WithAllowedSchemes(schemes ...string) Option {
	return func(c *config) {
		for _, s := range schemes {
			s = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(s, ":")))
			if s == "" {
				continue
			}
			if c.allowedSchemes == nil {
				c.allowedSchemes = make(map[string]struct{})
			}
			c.allowedSchemes[s] = struct{}{}
		}
	}
}

type fsResolver struct {
	fsys    fs.FS
	allowed map[string]struct{}
}

// New returns a Resolver rooted at fsys. References are resolved as paths
// inside fsys, which may be an in-memory filesystem (such as
// testing/fstest.MapFS) or an on-disk tree (see FromDir).
func New(fsys fs.FS, opts ...Option) Resolver {
	var c config
	for _, o := range opts {
		o(&c)
	}
	return &fsResolver{fsys: fsys, allowed: c.allowedSchemes}
}

// FromDir returns a Resolver rooted at the directory dir on the local
// filesystem. References may not escape dir, and URI schemes remain denied
// unless enabled with WithAllowedSchemes.
func FromDir(dir string, opts ...Option) Resolver {
	return New(os.DirFS(dir), opts...)
}

// Open implements Resolver.
func (r *fsResolver) Open(ctx context.Context, ref string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, validationErr("empty resource reference")
	}

	if scheme, ok := schemeOf(ref); ok {
		if _, allowed := r.allowed[scheme]; !allowed {
			return nil, unsupportedErr(fmt.Sprintf("resource scheme %q is not allowed", scheme))
		}
		if scheme != "file" {
			return nil, unsupportedErr(fmt.Sprintf("resource scheme %q has no configured transport", scheme))
		}
		ref = trimFileScheme(ref)
	}

	if !fs.ValidPath(ref) {
		return nil, validationErr(fmt.Sprintf("resource reference %q is not a valid rooted path", ref))
	}

	f, err := r.fsys.Open(ref)
	if err != nil {
		return nil, fmt.Errorf("resources: open %q: %w", ref, err)
	}
	return f, nil
}

// schemeOf reports the lowercase URI scheme in ref, if it has one. Single
// letters are treated as Windows drive letters rather than schemes.
func schemeOf(ref string) (string, bool) {
	u, err := url.Parse(ref)
	if err != nil || u.Scheme == "" {
		return "", false
	}
	scheme := strings.ToLower(u.Scheme)
	if len(scheme) == 1 {
		return "", false
	}
	return scheme, true
}

// trimFileScheme removes a leading "file:" scheme (and any authority/leading
// slashes) so the remainder resolves as a path inside the root.
func trimFileScheme(ref string) string {
	s := ref[len("file:"):]
	s = strings.TrimPrefix(s, "//")
	return strings.TrimLeft(s, "/")
}

func unsupportedErr(msg string) error {
	return diag.New(diag.CodeUnsupported, diag.SeverityError, msg, nil)
}

func validationErr(msg string) error {
	return diag.New(diag.CodeValidation, diag.SeverityError, msg, nil)
}
