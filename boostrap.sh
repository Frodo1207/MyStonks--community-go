#!/bin/bash

set -e

echo "🚀 Starting bootstrap..."

go install github.com/air-verse/air@latest

# 1. 创建 tmp 目录
if [ ! -d tmp ]; then
  echo "📁 Creating tmp directory..."
  mkdir tmp
else
  echo "📁 tmp directory already exists."
fi

# 2. 生成 .air.toml 配置文件
echo "📝 Generating .air.toml configuration..."

cat > .air.toml <<'EOF'
[build]
cmd = "go build -o tmp/server . && chmod +x tmp/server"
bin = "tmp/server"
full_bin = "tmp/server start --config=config/config-dev.yaml"
delay = 1000
exclude_dir = ["tmp", "vendor"]
exclude_file = ["*_test.go"]

log = "air.log"
send_interrupt = true
color = true
debug = false

[watch]
include_ext = ["go", "yaml", "toml"]
exclude_dir = ["tmp", "vendor"]
EOF

echo "✅ .air.toml generated."

# 3. 给 tmp/server 授权（保险起见）
if [ -f tmp/server ]; then
  echo "🔑 Setting executable permission on tmp/server..."
  chmod +x tmp/server
fi

echo "🎉 Bootstrap complete. You can now run:"
echo "sh ./dev-start.sh"
