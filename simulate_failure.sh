#!/bin/bash
# simulate_failure.sh

echo "==========================================================="
echo "   Initiating Simulated Catastrophic Failure Event"
echo "==========================================================="
echo ""
echo "Creating mock incident payloads..."

# Payload 1: RDBMS Outage (Severity 0 - Critical)
cat << 'EOF' > rdbms_payload.json
{
  "component_id": "DB_CLUSTER_PRIMARY",
  "error_code": "CONNECTION_REFUSED_503",
  "metadata": {
    "severity": 0,
    "region": "us-east-1",
    "details": "RDBMS connection pool exhausted. Primary node unresponsive."
  }
}
EOF

# Payload 2: Downstream MCP Failure (Severity 1 - High)
cat << 'EOF' > mcp_payload.json
{
  "component_id": "MCP_CONTROLLER_99",
  "error_code": "MCP_FAILURE_RETRY",
  "metadata": {
    "severity": 1,
    "region": "eu-central-1",
    "details": "Downstream Multi-Cloud Provisioner uncommunicative. Heartbeat timeout."
  }
}
EOF

echo "Blasting http://localhost:8080/api/v1/signals with 10,000 requests..."
echo "Watch the backend terminal for:"
echo " 1. Rate limiter drops (if exceeding 12,000/sec burst)"
echo " 2. Real-time Ingestion Throughput (Signals/sec logs)"
echo " 3. Redis/Mutex Debouncing (Only 1 Postgres Incident per component_id)"
echo "-----------------------------------------------------------"

# Send 5000 requests for each payload as fast as possible using background jobs
for i in {1..5000}; do
  curl -s -X POST http://localhost:8080/api/v1/signals -H "Content-Type: application/json" -d @rdbms_payload.json > /dev/null &
  curl -s -X POST http://localhost:8080/api/v1/signals -H "Content-Type: application/json" -d @mcp_payload.json > /dev/null &
done

# Wait for all background cURL jobs to finish
wait

echo "-----------------------------------------------------------"
echo "Simulation Complete!"
echo "Check your frontend dashboard (http://localhost:5173). You should see EXACTLY two new incidents (one for RDBMS, one for MCP) despite sending 10,000 signals."
echo "Check your MongoDB raw_signals collection to verify all 10,000 payloads were stored asynchronously."

rm rdbms_payload.json mcp_payload.json
