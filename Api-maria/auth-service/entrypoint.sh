#!/bin/sh
set -e

echo "⏳ Waiting for PostgreSQL to be ready..."

python - <<'EOF'
import sys, time, os
try:
    import psycopg2
except ImportError:
    print("psycopg2 not available, skipping DB wait")
    sys.exit(0)

url = os.environ.get('DATABASE_URL', '')
if not url:
    print("DATABASE_URL not set, skipping DB wait")
    sys.exit(0)

retries = 30
while retries > 0:
    try:
        conn = psycopg2.connect(url)
        conn.close()
        print("✅ PostgreSQL is ready!")
        sys.exit(0)
    except Exception as e:
        retries -= 1
        print(f"DB not ready ({e}), retrying in 2s... ({retries} attempts left)")
        time.sleep(2)

print("❌ Could not connect to DB after 30 attempts")
sys.exit(1)
EOF

echo "🔄 Running Alembic migrations..."
alembic upgrade head

echo "🚀 Starting uvicorn on port 8001..."
exec uvicorn app.main:app --host 0.0.0.0 --port 8001 --workers 2
