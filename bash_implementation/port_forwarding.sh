#!/usr/bin/env bash

# Ensure we are using bash
if [ -z "$BASH_VERSION" ]; then
  echo "❌ Please run this script with bash, not sh."
  exit 1
fi

set -euo pipefail

# Load environment variables
# Some variables are set in .env(eg. NAMESPACE, PORT_FORWARD_RULES)
set -a

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/.env"
set +a



# Loop through rules
echo "$PORT_FORWARD_RULES" | while IFS=: read -r PREFIX PORT; do
  [ -z "$PREFIX" ] && continue
  
  POD=$(kubectl get pods -n "$NAMESPACE" --no-headers -o custom-columns=":metadata.name" | grep "^${PREFIX}" | head -n 1 || true)

  if [[ -n "$POD" ]]; then
    echo "Found pod: $POD (prefix: $PREFIX)"
    echo "Forwarding local port $PORT to $PORT..."

    # Run in background
    kubectl port-forward -n "$NAMESPACE" pod/"$POD" "$PORT":"$PORT" >/dev/null 2>&1 &
    PF_PID=$!
    echo "Port-forwarding started for $POD on port $PORT"
    echo "  └─ Process ID: $PF_PID"
    echo "  └─ Command: kubectl port-forward -n $NAMESPACE pod/$POD $PORT:$PORT"
  else
    echo "⚠️  No pod found with prefix: $PREFIX"
  fi
done

echo "✅ All matching port-forwards started (running in background)."
