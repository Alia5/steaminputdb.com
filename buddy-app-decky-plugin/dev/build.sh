#!/usr/bin/env bash

rm -rf plugin-dist
rm -rf *.zip

mkdir -p plugin-dist plugin-dist/bin
pnpm run build

cp -r dist plugin-dist/
cp main.py plugin-dist/main.py
cp plugin.json plugin-dist/plugin.json
cp package.json plugin-dist/package.json
cp README.md plugin-dist/README.md

cat ../LICENSE.txt THIRD-PARTY-NOTICES.txt > plugin-dist/LICENSE.txt

VERSION=$(git describe --tags --abbrev=0)
BINARY_URL="https://github.com/Alia5/steaminputdb.com/releases/download/${VERSION}/steaminputdb-buddy-linux-amd64"
echo "Downloading binary from ${BINARY_URL}..."
curl -L -o plugin-dist/bin/steaminputdb-buddy "${BINARY_URL}"
chmod +x plugin-dist/bin/steaminputdb-buddy

sed -i "s/<VERSION_SET_BY_CI>/${VERSION#v}/g" plugin-dist/package.json
sed -i "s|<URL_SET_BY_CI>|${BINARY_URL}|g" plugin-dist/package.json
CHECKSUM=$(sha256sum plugin-dist/bin/steaminputdb-buddy | awk '{ print $1 }')
sed -i "s/<CHECKSUM_SET_BY_CI>/${CHECKSUM}/g" plugin-dist/package.json

mv plugin-dist SteamInputDB-Buddy
zip -r steaminputdb-buddy-decky-plugin-${VERSION}.zip SteamInputDB-Buddy/*
mv SteamInputDB-Buddy plugin-dist