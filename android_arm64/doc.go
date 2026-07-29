// Package android_arm64 carries the prebuilt FFmpeg static libraries for Android arm64-v8a.
//
// It exposes no Go API. goffmpeg imports it for its side effects: the
// cgo directives in libav.go contribute this platform's linker flags,
// which is how one small module per platform replaces vendoring every
// platform's artifacts into goffmpeg itself.
package android_arm64
