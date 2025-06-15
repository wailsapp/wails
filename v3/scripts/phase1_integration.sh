#!/bin/bash

# Phase 1 Integration Script
# This script runs all Phase 1 integration tests and generates a comprehensive report

set -e

echo "=== Wails v3 Phase 1 Integration ==="
echo "Date: $(date)"
echo ""

# Change to project root
cd "$(dirname "$0")/.."

# Create results directory
mkdir -p test-results

echo "1. Running integration tests..."
go test -v ./tests/integration -run TestPhase1 2>&1 | tee test-results/phase1_tests.log

echo ""
echo "2. Running performance benchmarks..."
go test -bench=BenchmarkPhase1Integration -benchmem -run=^$ -count=3 ./tests/integration 2>&1 | tee test-results/phase1_benchmarks.log

echo ""
echo "3. Running HTTP asset serving benchmarks..."
go test -bench=BenchmarkHTTPAssetServing -benchmem -run=^$ -count=3 ./tests/integration 2>&1 | tee test-results/phase1_http_benchmarks.log

echo ""
echo "4. Running memory pressure benchmarks..."
go test -bench=BenchmarkMemoryPressure -benchmem -run=^$ -count=3 ./tests/integration 2>&1 | tee test-results/phase1_memory_benchmarks.log

echo ""
echo "5. Running concurrent load benchmarks..."
go test -bench=BenchmarkConcurrentLoad -benchmem -run=^$ -count=3 ./tests/integration 2>&1 | tee test-results/phase1_concurrent_benchmarks.log

echo ""
echo "6. Generating integration report..."
go test -v ./tests/integration -run TestGenerateReport

echo ""
echo "7. Running full test suite to ensure stability..."
go test ./... -short 2>&1 | tee test-results/phase1_full_tests.log

echo ""
echo "=== Phase 1 Integration Complete ==="
echo ""
echo "Results saved to:"
echo "  - test-results/phase1_tests.log"
echo "  - test-results/phase1_benchmarks.log"
echo "  - test-results/phase1_http_benchmarks.log"
echo "  - test-results/phase1_memory_benchmarks.log"
echo "  - test-results/phase1_concurrent_benchmarks.log"
echo "  - test-results/phase1_full_tests.log"
echo "  - PHASE1_INTEGRATION_REPORT.md"
echo ""

# Check if all tests passed
if grep -q "FAIL" test-results/*.log; then
    echo "⚠️  Some tests failed. Please review the logs."
    exit 1
else
    echo "✅ All tests passed!"
fi

# Display summary from report if it exists
if [ -f "PHASE1_INTEGRATION_REPORT.md" ]; then
    echo ""
    echo "=== Executive Summary ==="
    grep -A 4 "## Executive Summary" PHASE1_INTEGRATION_REPORT.md | tail -n +2
fi