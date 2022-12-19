#!/usr/bin/env bash

PLATFORMS=windows/amd64

wails build -clean -platform=$PLATFORMS -webview2=embed -trimpath -obfuscated
