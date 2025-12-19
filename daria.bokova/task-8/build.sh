#!/bin/bash

# Создаем папку bin
mkdir -p bin

# Собираем с тегом dev
go build -tags dev -o bin/service ./cmd/app

# Проверяем
./bin/service