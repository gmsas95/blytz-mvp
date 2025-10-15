#!/bin/bash
set -e

echo "🔍 Validating Microservices Architecture..."
echo "=========================================="

# Define expected services
SERVICES=(
  "auction-service"
  "auth-service"
  "chat-service"
  "gateway"
  "logistics-service"
  "order-service"
  "payment-service"
  "product-service"
)

echo "📋 Checking service structure..."
echo

all_good=true

for service in "${SERVICES[@]}"; do
  echo "🔍 Checking $service:"
  service_dir="services/$service"

  if [ ! -d "$service_dir" ]; then
    echo "  ❌ Directory not found: $service_dir"
    all_good=false
    continue
  fi

  # Check essential files
  if [ -f "$service_dir/cmd/main.go" ]; then
    echo "  ✅ cmd/main.go (executable entry point)"
  else
    echo "  ❌ Missing cmd/main.go"
    all_good=false
  fi

  if [ -f "$service_dir/Dockerfile" ]; then
    echo "  ✅ Dockerfile (containerization)"
  else
    echo "  ❌ Missing Dockerfile"
    all_good=false
  fi

  if [ -f "$service_dir/go.mod" ]; then
    echo "  ✅ go.mod (module definition)"
  else
    echo "  ❌ Missing go.mod"
    all_good=false
  fi

  if [ -f "$service_dir/go.sum" ]; then
    echo "  ✅ go.sum (dependency checksums)"
  else
    echo "  ⚠️  Missing go.sum (may need go mod tidy)"
  fi

  echo
done

# Check shared module
echo "🔍 Checking shared module:"
if [ -d "services/shared" ]; then
  echo "  ✅ services/shared directory exists"
  if [ -f "services/shared/go.mod" ]; then
    echo "  ✅ shared/go.mod exists"
  else
    echo "  ❌ shared/go.mod missing"
    all_good=false
  fi
else
  echo "  ❌ services/shared directory missing"
  all_good=false
fi

echo
echo "🔍 Checking GitHub Actions workflow:"
if grep -q "services/" /home/sas/blytzmvp-clean/.github/workflows/deploy.yml; then
  echo "  ✅ Workflow uses services/ paths"
else
  echo "  ❌ Workflow still references backend/ paths"
  all_good=false
fi

echo
echo "🔍 Checking OpenAPI specifications:"
opapi_count=$(ls -1 /home/sas/blytzmvp-clean/openapi/*.yaml 2>/dev/null | wc -l)
if [ "$opapi_count" -eq 8 ]; then
  echo "  ✅ All 8 services have OpenAPI specs"
else
  echo "  ⚠️  Found $opapi_count OpenAPI specs (expected 8)"
fi

echo
echo "=========================================="
if [ "$all_good" = true ]; then
  echo "🎉 Architecture validation PASSED!"
  echo "✅ All microservices are properly structured"
  echo "✅ Directory structure follows CLAUDE.md specification"
  echo "✅ Ready for local development and CI/CD"
else
  echo "❌ Architecture validation FAILED!"
  echo "🔧 Please fix the issues above before proceeding"
  exit 1
fi

echo
echo "🚀 Next steps:"
echo "1. Run: ./fix-go-mods.sh (requires Go 1.23+)"
echo "2. Test locally: cd services/auth-service && go run main.go"
echo "3. Run full CI: make ci-pipeline"
echo "4. Deploy: git push origin main""" file_path="/home/sas/blytzmvp-clean/validate-architecture.sh"}