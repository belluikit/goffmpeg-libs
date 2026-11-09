# goffmpeg-libs

Prebuilt FFmpeg libraries, static and shared, for [goffmpeg](https://github.com/belluikit/goffmpeg),
published as one Go module per platform.

## Why one module per platform

goffmpeg links a purpose-built FFmpeg. Go modules distribute whatever is
committed, so vendoring every target into goffmpeg itself would force each
consumer to download all of them — roughly 260MB of module zips — to use one.

Each directory here is its own module, so the proxy serves them independently
and a build downloads only its own target (16-31MB, `shared/` included). goffmpeg selects the
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

### Shared libraries

`shared/` holds the same FFmpeg built as shared libraries. Player ships these
in apps, so a user can replace them (LGPL 2.1 §6b; Player ADR-0037 D3). They
come from a configure of their own, recorded in each `NOTICE`. goffmpeg links
the static libraries until it switches
([goffmpeg#40](https://github.com/belluikit/goffmpeg/issues/40)).

| Modules | Files in `shared/` | How a program finds them |
|---|---|---|
| `darwin_*` | `libavcodec-goffmpeg.63.dylib`, install name `@rpath/libavcodec-goffmpeg.63.dylib` | the executable's rpath |
| `ios_arm64`, `iossimulator_arm64` | `libavcodec-goffmpeg.framework`, install name `@rpath/libavcodec-goffmpeg.framework/libavcodec-goffmpeg` | `@executable_path/Frameworks`, where the app embeds them |
| `linux_*` | `libavcodec-goffmpeg.so.63`, soname the same | the executable's rpath |
| `android_*` | `libavcodec-goffmpeg.so`, soname the same | the APK's `lib/<abi>/` |
| `windows_*` | `avcodec-goffmpeg-63.dll`, import library `libavcodec-goffmpeg.dll.a` | the executable's directory |

The other libraries follow the same pattern: avutil, swresample, swscale,
avformat, and avdevice on desktop.

- **The `-goffmpeg` suffix** keeps them from colliding with another FFmpeg in
  the same process or package: a distro `libavcodec.so.62` on Linux, an
  `avcodec-62.dll` another component already loaded on Windows, another SDK's
  `libavcodec.so` in an APK. FFmpeg's checklist allows it: "Do not rename
  FFmpeg dlls to some obfuscated name, but adding a suffix or prefix is fine".
- **Each file has one name**, the one the loader looks for. Go module zips
  leave out symlinks, so the usual `libavcodec.dylib` → `libavcodec.63.dylib`
  links aren't there. Link against the file itself.
- **Android sonames carry no version**, because an APK installs only files
  named `lib*.so`. The libraries are aligned for 16KB pages, as Google Play
  requires.
- **Each iOS framework** has an `Info.plist` (bundle id
  `com.trendvidia.goffmpeg.<library>`, minimum iOS 13.0) and an ad-hoc
  signature. The app's packager re-signs it when it embeds it.
- **The libraries carry no rpath of their own.** Whatever loads them supplies
  the search path, as in the table.
- **Files arrive read-only, without the executable bit**: the go command
  doesn't keep file modes. macOS and Linux load them that way, so a
  development build can use them straight from the module cache. A packager
  that copies them into an app sets whatever mode it needs.

## Provenance and licence

Built from FFmpeg `e2bc6d88cd68184288c16570e04881c3a426a05a` (`N-125841-ge2bc6d88cd`)
by the Makefile in [goffmpeg](https://github.com/belluikit/goffmpeg), which pins the configuration. Decode-only: every native decoder,
demuxer, parser and protocol, with encoders, muxers and avfilter disabled.

No `--enable-gpl` and no `--enable-nonfree`, so these artifacts are
**LGPL v2.1 or later**. Adding a GPL component would relicense everything
linking them.

The license text is in [LICENSE](LICENSE). Each platform directory has a
`NOTICE` with the exact configure lines recorded in its libraries, one for
`lib/` and one for `shared/`, and where to get the corresponding FFmpeg source.

That source is [belluikit/ffmpeg](https://github.com/belluikit/ffmpeg), a
mirror of FFmpeg, at the tag `goffmpeg-e2bc6d88cd`. LGPL-2.1 asks for the
source to be offered from the same place as the binaries, and an upstream link
alone depends on a repository we do not control keeping that commit. Every
FFmpeg revision a published module was built from keeps its
`goffmpeg-<short sha>` tag there.

## https and TLS

FFmpeg needs a TLS library to open `https`, and these builds use
`--disable-autodetect`, so a target has only the one its configure line names.
Each module's `NOTICE` records that line.

| Modules | TLS library |
|---|---|
| `darwin_*`, `ios_arm64`, `iossimulator_arm64` | SecureTransport |
| `linux_*` | OpenSSL |
| `windows_*` | Schannel |
| `android_*` | **none** |

So on Android an `https` URL cannot open, including an `https` HLS playlist. It
fails at `avformat_open_input`, which goffmpeg reports as
`failed to open decoder: -1`. Plain `http` works, and so does `https` read
through Go and handed to goffmpeg's `OpenReader`. HLS cannot go that way,
because FFmpeg has to fetch the segments itself.

Regenerate with:

```sh
make dist TARGET=linux_amd64 && make dist-shared TARGET=linux_amd64 && make sync-libs TARGET=linux_amd64
```

Publish from CI rather than a local build: goffmpeg's `check-paths` refuses
libraries that record a personal home directory, which a local Android build
does through the NDK path.
