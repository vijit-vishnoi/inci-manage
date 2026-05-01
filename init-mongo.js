// Connect to the specific database
db = db.getSiblingDB('inci_mongo_db');

// Ensure the application user exists
db.createUser({
  user: "inci_app_user",
  pwd: "inci_app_password",
  roles: [
    {
      role: "readWrite",
      db: "inci_mongo_db"
    }
  ]
});

// Drop old collections if they exist to clear out the standard ticketing schema
db.audit_logs.drop();
db.attachments.drop();

// Create raw_signals collection specifically designed for high-volume JSON error payloads.
// Using a timeseries collection provides optimized storage and querying for sequential high-throughput data.
db.createCollection(
    "raw_signals",
    {
       timeseries: {
          timeField: "timestamp",
          metaField: "metadata",
          granularity: "seconds"
       }
    }
);

// Create indexes for efficient timeseries aggregation
// Note: Timeseries collections automatically create an index on the timeField and metaField,
// but adding compound indexes can speed up specific aggregation patterns.
db.raw_signals.createIndex({ "timestamp": -1 });
db.raw_signals.createIndex({ "metadata.source": 1, "timestamp": -1 });
db.raw_signals.createIndex({ "metadata.severity": 1, "timestamp": -1 });
