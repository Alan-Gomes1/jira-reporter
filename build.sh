#!/bin/bash

# Para a execução se ocorrer um erro
set -e

# Cria um diretório para os arquivos de release, se não existir
mkdir -p releases

# Define o nome base do executável
BINARY_NAME="jira-reporter"

# Arquivos essenciais que precisam ir junto com o binário
FILES_TO_PACKAGE="template.html env-example"

# ----- Build para Windows (amd64) -----
echo "Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -o "./releases/${BINARY_NAME}.exe" .
cd releases
zip -q "${BINARY_NAME}-windows-amd64.zip" "${BINARY_NAME}.exe" ${FILES_TO_PACKAGE}
rm "${BINARY_NAME}.exe"
cd ..
echo "Windows build packaged."

# ----- Build para Linux (amd64) -----
echo "Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o "./releases/${BINARY_NAME}-linux-amd64" .
cd releases
zip -q "${BINARY_NAME}-linux-amd64.zip" "${BINARY_NAME}-linux-amd64" ${FILES_TO_PACKAGE}
rm "${BINARY_NAME}-linux-amd64"
cd ..
echo "Linux build packaged."

# ----- Build para macOS (Intel amd64) -----
echo "Building for macOS amd64 (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o "./releases/${BINARY_NAME}-macos-amd64" .
cd releases
zip -q "${BINARY_NAME}-macos-amd64.zip" "${BINARY_NAME}-macos-amd64" ${FILES_TO_PACKAGE}
rm "${BINARY_NAME}-macos-amd64"
cd ..
echo "macOS (Intel) build packaged."

# ----- Build para macOS (Apple Silicon arm64) -----
echo "Building for macOS arm64 (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o "./releases/${BINARY_NAME}-macos-arm64" .
cd releases
zip -q "${BINARY_NAME}-macos-arm64.zip" "${BINARY_NAME}-macos-arm64" ${FILES_TO_PACKAGE}
rm "${BINARY_NAME}-macos-arm64"
cd ..
echo "macOS (Apple Silicon) build packaged."

echo "All builds are complete and packaged in the 'releases' directory."