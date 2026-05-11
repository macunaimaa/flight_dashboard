#!/bin/bash
# Seed script: creates a default tenant and admin user for development.
# Usage: bash scripts/seed_tenants.sh
#
# Requires: mongosh connected to the dashboard database.
# Password for admin user: admin123

MONGO_URI="${MONGO_URI:-mongodb://localhost:27017}"
MONGO_DB="${MONGO_DB:-dashboard}"

# Generate bcrypt hash for "admin123" using the Go helper
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PASSWORD_HASH=$(cd "$SCRIPT_DIR/../backend" && go run ./cmd/genhash 2>/dev/null)

if [ -z "$PASSWORD_HASH" ]; then
  echo "ERROR: Could not generate password hash. Make sure Go is installed and you're in the project root."
  exit 1
fi

mongosh "$MONGO_URI/$MONGO_DB" --quiet --eval "
  db.tenants.updateOne(
    { tenant_id: 'demo' },
    {
      \$setOnInsert: {
        tenant_id: 'demo',
        name: 'Demo Tenant',
        config: {
          poll_regions: [
            { lamin: -34.0, lomin: -74.0, lamax: 5.5, lomax: -34.0 },
            { lamin: 35.0, lomin: -10.0, lamax: 60.0, lomax: 30.0 }
          ],
          max_tracked_aircraft: 5000,
          track_retention_days: 30
        },
        active: true,
        created_at: new Date(),
        updated_at: new Date()
      }
    },
    { upsert: true }
  );

  db.users.updateOne(
    { tenant_id: 'demo', email: 'admin@demo.com' },
    {
      \$set: {
        tenant_id: 'demo',
        email: 'admin@demo.com',
        password_hash: '$PASSWORD_HASH',
        role: 'admin',
        active: true,
        created_at: new Date(),
        last_login: new Date()
      }
    },
    { upsert: true }
  );

  print('Seed complete.');
  print('Tenant: demo');
  print('User: admin@demo.com / admin123');
"
