// Package resources loads report resources — templates now, images and fonts
// later — from injected, deny-by-default sources.
//
// A Resolver maps a caller-supplied reference to a reader. By default only
// relative paths inside the resolver's root are allowed; URI schemes such as
// file:// and http:// are rejected. Enabling a scheme is an explicit, audited
// decision (OPS-002).
//
// Resolution is rooted: references are interpreted as paths inside the
// injected fs.FS, and traversal outside that root (../ or absolute paths) is
// rejected. Symlinks that escape the root are not detected here; see
// docs/architecture/data-and-resources.md.
package resources
