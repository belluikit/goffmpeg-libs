# goffmpeg-libs

Prebuilt FFmpeg static libraries for [goffmpeg](https://github.com/belluikit/goffmpeg),
published as one Go module per platform.

## Why one module per platform

goffmpeg links a purpose-built FFmpeg. Go modules distribute whatever is
committed, so vendoring every target into goffmpeg itself would force each
consumer to download all of them — roughly 245MB — to use one.

Each directory here is its own module, so the proxy serves them independently
and a build downloads only its own target (~20-30MB). goffmpeg selects the
right one with a build constraint:

```go
//go:build darwin && arm64 && !ios

import _ "github.com/belluikit/goffmpeg-libs/darwin_arm64"
```

The blank import contributes that platform's `#cgo LDFLAGS`, which propagate to
the final link. The modules expose no Go API.

The FFmpeg *headers* are identical for every target and live in goffmpeg
itself, so they are not duplicated here.

## Modules

| Module | Platform |
|---|---|
| `darwin_arm64` | macOS, Apple silicon |
| `darwin_amd64` | macOS, Intel |
| `ios_arm64` | iOS devices |
| `iossimulator_arm64` | iOS simulator (build with `-tags iossimulator`) |
| `android_arm64` | Android arm64-v8a |
| `android_arm` | Android armeabi-v7a |
| `android_amd64` | Android x86_64 |
| `linux_amd64` | Linux x86_64 |
| `linux_arm64` | Linux arm64 |
| `windows_amd64` | Windows x86_64 |
| `windows_arm64` | Windows on ARM |

## Contents

Each module carries the FFmpeg static libraries for its target, plus
`libgoffmpeg.{so,dll,dylib}` — the shared bridge the purego backend
(`CGO_ENABLED=0`) dlopens. The cgo backend links the static libraries directly
and ignores the shared object.

The Linux and Windows libraries are built with `--enable-pic`. FFmpeg's
hand-written assembly is not position-independent by default, and without it
the shared bridge cannot be linked at all.

## Provenance and licence

Built from FFmpeg `e2bc6d88cd68184288c16570e04881c3a426a05a` (`N-125841-ge2bc6d88cd`)
by the Makefile in [goffmpeg](https://github.com/belluikit/goffmpeg), which pins the configuration. Decode-only: every native decoder,
demuxer, parser and protocol, with encoders, muxers and avfilter disabled.

No `--enable-gpl` and no `--enable-nonfree`, so these artifacts are
**LGPL v2.1 or later**. Adding a GPL component would relicense everything
linking them.

The license text is in [LICENSE](LICENSE). Each platform directory has a
`NOTICE` with the exact configure line recorded in its libraries and where to
get the corresponding FFmpeg source.

Regenerate with:

```sh
make dist TARGET=linux_amd64 && make sync-libs TARGET=linux_amd64
```
