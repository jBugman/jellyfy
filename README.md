## Automatically rename TV series episode files into a format [Jellyfin](https://jellyfin.org/) expects

```
> jellyfy -series="The Expanse" -season=6 ~/Downloads/The.Expanse.S06.WEBDL.1080p
Folder renamed successfully
The.Expanse.S06E01.WEBDL.1080p….mkv -> The Expanse S06E01.mkv
…
The.Expanse.S06E06.….mkv -> The Expanse S06E06.mkv
```

## Build

Replace `darwin` and `arm64` with your target platform and architecture

```
GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o build/jellyfy_mac_arm64
```

Or for local install
```
go install .
```

### Release

```
zip -9 -r "build/jellyfy_mac_arm64.zip" build/jellyfy_mac_arm64
```
