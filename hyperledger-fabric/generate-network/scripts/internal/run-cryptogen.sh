#!/usr/bin/env bash

set -e

CFG_DIR="organizations/cryptogen"

for yaml in "$CFG_DIR"/crypto-config-*.yaml; do
  echo "🔧  Генерируем материалы для $(basename "$yaml")"
  set -x
  cryptogen generate --config="$yaml" --output="organizations"
  { set +x; } 2>/dev/null
done

echo "✅  Certificates & MSP сгенерированы"
