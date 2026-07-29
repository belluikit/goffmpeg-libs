// Package iossimulator_arm64 carries the prebuilt FFmpeg static libraries for the iOS simulator on Apple silicon.
//
// It exposes no Go API. goffmpeg imports it for its side effects: the
// cgo directives in libav.go contribute this platform's linker flags,
// which is how one small module per platform replaces vendoring every
// platform's artifacts into goffmpeg itself.
package iossimulator_arm64
