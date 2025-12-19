#!/bin/bash

# Создаем папку bin на уровень выше (как ожидает CI/CD)
mkdir -p ../../bin

# Определяем тег сборки
if [[ "$1" == "dev" ]]; then
    echo "Building dev version..."
    go build -tags dev -o ../../bin/service ./cmd/app
else
    echo "Building prod version..."
    go build -o ../../bin/service ./cmd/app
fi

# Проверяем
echo "Built executable:"
ls -la ../../bin/service
../../bin/service