# goffmpeg-libs

Prebuilt FFmpeg shared libraries for [goffmpeg](https://github.com/belluikit/goffmpeg),
published as one Go module per platform.

## Why one module per platform

goffmpeg links a purpose-built FFmpeg. Go modules distribute whatever is
committed, so vendoring every target into goffmpeg itself would force each
consumer to download all of them — roughly 260MB of module zips — to use one.

Each directory here is its own module, so the proxy serves them independently
and a build downloads only its own target (5-11MB on disk from v0.6.0). goffmpeg selects the
right one with a build constraint:

```go
//go:build darwin && arm64 && !ios

import _ "github.com/belluikit/goffmpeg-libs/darwin_arm64"
```

The blank import contributes that platform's `#cgo LDFLAGS`, which propagate to
the final link. The modules expose no Go API.

Those flags are in `libav.go` (the media codec set, from `shared/`) and
`libav_audio.go` (the audio set, from `shared-audio/`, under goffmpeg's
`goffmpeg_novideo` tag), which goffmpeg's `make sync-libs` generates from the
files in those directories (its `scripts/gen-libav.sh`). Don't edit them by
hand.

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

Each module carries FFmpeg for its target as shared libraries, in two codec
sets ([goffmpeg#41](https://github.com/belluikit/goffmpeg/issues/41)):

| Set | Plays | Directory | Names | goffmpeg links it |
|---|---|---|---|---|
| media, desktop (`darwin_*`, `linux_*`, `windows_*`) | every format FFmpeg decodes; webcams | `shared/` | `-goffmpeg` | by default |
| media, phones (`android_*`, `ios_*`) | MP3, WAV, MP4 and HLS with H.264 and AAC | `shared/` | `-goffmpeg` | by default |
| audio, every target | MP3 and WAV; no libswscale, no libavdevice | `shared-audio/` | `-goffmpeg-audio` | with `goffmpeg_novideo` |

On the phones the media set decodes only what apps play, so an app carries no
codec it doesn't need. On desktop it keeps every codec. On Android the media
set also has `h264_mediacodec`, MediaCodec's H.264 decoder (#21), which a
caller opens by name; a lookup by codec ID still returns the software `h264`.

A binary loads one set, and its package carries that set alone. Each module's
`NOTICE` records both configure lines, and goffmpeg's BUILDING.md ("What is in
the build") the exact decoders and demuxers.

The desktop modules also carry `lib/libgoffmpeg.{so,dll,dylib}`, the bridge the
purego backend (`CGO_ENABLED=0`) dlopens. It holds no FFmpeg code: it loads the
media set in `shared/`, looking beside itself and in `../shared`.

Up to v0.4.0, `lib/` also held FFmpeg's static archives. Nothing links them
since goffmpeg#40, so from v0.5.0 they are no longer published
([goffmpeg#45](https://github.com/belluikit/goffmpeg/issues/45)).

### Shared libraries

goffmpeg links FFmpeg as shared libraries on every platform, so an app ships
them and a user can replace them (LGPL 2.1 §6b; Player ADR-0037 D3,
[goffmpeg#40](https://github.com/belluikit/goffmpeg/issues/40)). The table
shows the media set's names; the audio set's are the same with
`-goffmpeg-audio` (`libavcodec-goffmpeg-audio.63.dylib`).

| Modules | Files in `shared/` | How a program finds them |
|---|---|---|
| `darwin_*` | `libavcodec-goffmpeg.63.dylib`, install name `@rpath/libavcodec-goffmpeg.63.dylib` | the executable's rpath |
| `ios_arm64`, `iossimulator_arm64` | `libavcodec-goffmpeg.framework`, install name `@rpath/libavcodec-goffmpeg.framework/libavcodec-goffmpeg` | `@executable_path/Frameworks`, where the app embeds them |
| `linux_*` | `libavcodec-goffmpeg-63.so`, soname the same | the executable's rpath |
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
  the search path, as in the table. On macOS and Linux the link files add an
  rpath to the set's directory, so `go run` and `go test` find them in the
  module cache; a packaged app replaces it with its own. On Windows, put
  `shared/` (or `shared-audio/` for a `goffmpeg_novideo` build) on `PATH` to
  run from source. The two sets' DLLs are named apart, so both can be there.
- **On Linux, the libraries need glibc 2.35 or newer**: Ubuntu 22.04 LTS,
  Debian 12 and later. A shared library needs at least the glibc it was built
  against, so they are built on Ubuntu 22.04, and goffmpeg's build refuses one
  that needs more
  ([goffmpeg#44](https://github.com/belluikit/goffmpeg/issues/44)).
- **Files arrive read-only, without the executable bit**: the go command
  doesn't keep file modes. macOS and Linux load them that way, so a
  development build can use them straight from the module cache. A packager
  that copies them into an app sets whatever mode it needs.

## Provenance and licence

Built from FFmpeg `e2bc6d88cd68184288c16570e04881c3a426a05a` (`N-125841-ge2bc6d88cd`)
by the Makefile in [goffmpeg](https://github.com/belluikit/goffmpeg), which pins the configuration. Decode-only, in the
two codec sets above, with encoders, muxers and avfilter disabled; every
protocol is kept.

No `--enable-gpl` and no `--enable-nonfree`, so these artifacts are
**LGPL v2.1 or later**. Adding a GPL component would relicense everything
linking them.

The license text is in [LICENSE](LICENSE). Each platform directory has a
`NOTICE` with the exact configure line recorded in its libraries, and where to
get the corresponding FFmpeg source.

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
make linux-docker TARGET=linux_amd64 DOCKER_GOALS="dist bridge-shared"
make linux-docker TARGET=linux_amd64 DOCKER_GOALS="dist SET=audio"
make sync-libs TARGET=linux_amd64
```

Publish from CI rather than a local build: goffmpeg's `check-paths` refuses
libraries that record a personal home directory, which a local Android build
does through the NDK path.
