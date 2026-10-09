#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"
HOOK_PATH="$HOOKS_DIR/pre-commit"

mkdir -p "$HOOKS_DIR"

cat > "$HOOK_PATH" <<'HOOK_EOF'
#!/usr/bin/env bash
set -euo pipefail
REPO_ROOT="$(git rev-parse --show-toplevel)"
exec "$REPO_ROOT/scripts/pre-commit.sh"
HOOK_EOF

chmod +x "$HOOK_PATH"
chmod +x "$REPO_ROOT/scripts/pre-commit.sh" || true

echo "Installed Git pre-commit hook at: $HOOK_PATH"

