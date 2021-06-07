#!/bin/bash
set -xe
apt-get update -yqq
apt-get install make build-essential libwebkit2gtk-4.0-dev libgtk-3-dev gcc-multilib g++-multilib sshpass p7zip-full gcc-mingw-w64-x86-64 gcc-mingw-w64-i686 -yqq
