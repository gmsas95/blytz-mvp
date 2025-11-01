#!/bin/bash

# Comprehensive Test Runner for Payment Service
# This script runs all tests with different configurations and coverage

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVICE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COVERAGE_DIR="$SERVICE_DIR/coverage"
REPORT_DIR="$SERVICE_DIR/test_reports"

# Create directories
mkdir -p "$COVERAGE_DIR"
mkdir -p "$REPORT_DIR"

echo -e "${BLUE}🧪 Payment Service Test Suite${NC}"
echo -e "${BLUE}===================================${NC}"

# Function to print section header
print_section() {
    echo -e "\n${YELLOW}📋 $1${NC}"
    echo -e "${YELLOW}-----------------------------------${NC}"
}

# Function to run tests with coverage
run_test_with_coverage() {
    local test_name="$1"
    local test_path="$2"
    local coverage_file="$COVERAGE_DIR/${test_name}.out"
    local html_file="$REPORT_DIR/${test_name}.html"

    echo -e "Running ${test_name} tests..."

    if go test -v -coverprofile="$coverage_file" "$test_path" 2>&1 | tee "$REPORT_DIR/${test_name}.log"; then
        echo -e "${GREEN}✅ ${test_name} tests passed${NC}"

        # Generate HTML coverage report
        go tool cover -html="$coverage_file" -o "$html_file"
        echo -e "📊 Coverage report generated: ${html_file}"

        # Show coverage percentage
        coverage_pct=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
        echo -e "📈 Coverage: ${coverage_pct}"

        return 0
    else
        echo -e "${RED}❌ ${test_name} tests failed${NC}"
        return 1
    fi
}

# Function to run benchmark tests
run_benchmarks() {
    local test_name="$1"
    local test_path="$2"
    local bench_file="$REPORT_DIR/${test_name}_bench.txt"

    echo -e "Running ${test_name} benchmarks..."

    if go test -bench=. -benchmem "$test_path" | tee "$bench_file"; then
        echo -e "${GREEN}✅ ${test_name} benchmarks completed${NC}"
        return 0
    else
        echo -e "${RED}❌ ${test_name} benchmarks failed${NC}"
        return 1
    fi
}

# Function to run race condition tests
run_race_tests() {
    local test_name="$1"
    local test_path="$2"

    echo -e "Running ${test_name} race condition tests..."

    if go test -race -v "$test_path"; then
        echo -e "${GREEN}✅ ${test_name} race tests passed${NC}"
        return 0
    else
        echo -e "${RED}❌ ${test_name} race tests failed${NC}"
        return 1
    fi
}

# Change to service directory
cd "$SERVICE_DIR"

# Check if Go modules are up to date
print_section "Environment Setup"
echo "Checking Go modules..."
if ! go mod tidy; then
    echo -e "${RED}❌ Failed to tidy Go modules${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go modules are up to date${NC}"

# Check Go version
go_version=$(go version)
echo "Go version: $go_version"

# Initialize test counters
total_tests=0
passed_tests=0
failed_tests=0

# 1. Unit Tests
print_section "Unit Tests"

if run_test_with_coverage "fiuu_client" "./pkg/fiuu/"; then
    ((passed_tests++))
else
    ((failed_tests++))
fi
((total_tests++))

if run_test_with_coverage "payment_service" "./internal/services/"; then
    ((passed_tests++))
else
    ((failed_tests++))
fi
((total_tests++))

# 2. Integration Tests
print_section "Integration Tests"

if run_test_with_coverage "api_handlers" "./internal/api/handlers/"; then
    ((passed_tests++))
else
    ((failed_tests++))
fi
((total_tests++))

# 3. Performance Tests
print_section "Performance Tests"

if run_benchmarks "fiuu_client" "./pkg/fiuu/"; then
    ((passed_tests++))
else
    ((failed_tests++))
fi
((total_tests++))

# 4. Race Condition Tests
print_section "Race Condition Tests"

if run_race_tests "payment_service" "./internal/services/"; then
    ((passed_tests++))
else
    ((failed_tests++))
fi
((total_tests++))

# 5. Combined Coverage Report
print_section "Combined Coverage Report"

echo "Generating combined coverage report..."
combined_coverage="$COVERAGE_DIR/combined.out"

# Combine all coverage files
find "$COVERAGE_DIR" -name "*.out" -type f -exec cat {} \; > "$combined_coverage"

# Generate combined HTML report
go tool cover -html="$combined_coverage" -o "$REPORT_DIR/combined_coverage.html"

# Show total coverage
total_coverage=$(go tool cover -func="$combined_coverage" | grep total | awk '{print $3}')
echo -e "📊 Total Coverage: ${total_coverage}"

# 6. Test Summary
print_section "Test Summary"

echo -e "Total Test Suites: ${total_tests}"
echo -e "${GREEN}Passed: ${passed_tests}${NC}"
echo -e "${RED}Failed: ${failed_tests}${NC}"

if [ $failed_tests -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed successfully!${NC}"

    # Coverage thresholds
    coverage_num=$(echo $total_coverage | sed 's/%//')
    if (( $(echo "$coverage_num >= 80" | bc -l) )); then
        echo -e "${GREEN}✅ Excellent coverage (${total_coverage})${NC}"
    elif (( $(echo "$coverage_num >= 60" | bc -l) )); then
        echo -e "${YELLOW}⚠️  Good coverage (${total_coverage}) - Aim for 80%+${NC}"
    else
        echo -e "${RED}❌ Low coverage (${total_coverage}) - Add more tests!${NC}"
    fi
else
    echo -e "\n${RED}💥 Some tests failed. Check the logs above.${NC}"
    exit 1
fi

# 7. Generate Test Report
print_section "Test Report Generation"

report_file="$REPORT_DIR/test_summary_$(date +%Y%m%d_%H%M%S).md"

cat > "$report_file" << EOF
# Payment Service Test Report

**Generated:** $(date)
**Go Version:** $go_version

## Test Results

| Test Suite | Status | Coverage |
|------------|--------|----------|
EOF

# Add test results to report
for test_suite in fiuu_client payment_service api_handlers; do
    if [ -f "$COVERAGE_DIR/${test_suite}.out" ]; then
        coverage=$(go tool cover -func="$COVERAGE_DIR/${test_suite}.out" | grep total | awk '{print $3}')
        echo "| ${test_suite} | ✅ Passed | ${coverage} |" >> "$report_file"
    else
        echo "| ${test_suite} | ❌ Failed | N/A |" >> "$report_file"
    fi
done

cat >> "$report_file" << EOF

## Overall Coverage: ${total_coverage}

## Files Generated

- Combined Coverage Report: [combined_coverage.html](file://$REPORT_DIR/combined_coverage.html)
- Test Logs: Available in $REPORT_DIR/*.log
- Benchmark Results: Available in $REPORT_DIR/*_bench.txt

## Recommendations

EOF

if (( $(echo "$coverage_num >= 80" | bc -l) )); then
    echo "✅ Excellent test coverage. Ready for production deployment." >> "$report_file"
elif (( $(echo "$coverage_num >= 60" | bc -l) )); then
    echo "⚠️ Good test coverage, but consider adding more tests to reach 80%+." >> "$report_file"
else
    echo "❌ Low test coverage. Add comprehensive unit and integration tests before production deployment." >> "$report_file"
fi

echo -e "📄 Test report generated: ${report_file}"

# 8. Cleanup (optional)
print_section "Cleanup"

echo "Cleaning up temporary files..."
find "$SERVICE_DIR" -name "*.test" -delete 2>/dev/null || true
echo -e "${GREEN}✅ Cleanup completed${NC}"

echo -e "\n${BLUE}🏁 Test suite completed!${NC}"
echo -e "View reports: ${REPORT_DIR}/"

exit 0