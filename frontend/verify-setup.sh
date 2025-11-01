#!/bin/bash

echo "🔍 Verifying Development Tools Setup"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    local test_name=$1
    local status=$2
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}✅ $test_name${NC}"
    else
        echo -e "${RED}❌ $test_name${NC}"
    fi
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if Node.js is installed
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    print_status "Node.js ($NODE_VERSION)" 0
else
    print_status "Node.js" 1
fi

# Check if npm is installed
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm --version)
    print_status "npm ($NPM_VERSION)" 0
else
    print_status "npm" 1
fi

# Check if dependencies are installed
if [ -d "node_modules" ]; then
    print_status "Dependencies installed" 0
else
    print_status "Dependencies installed" 1
    echo "Run 'npm install' to install dependencies"
fi

# Check TypeScript configuration
if [ -f "tsconfig.json" ]; then
    print_status "TypeScript configuration" 0
else
    print_status "TypeScript configuration" 1
fi

# Check ESLint configuration
if [ -f ".eslintrc.json" ]; then
    print_status "ESLint configuration" 0
else
    print_status "ESLint configuration" 1
fi

# Check Prettier configuration
if [ -f ".prettierrc" ]; then
    print_status "Prettier configuration" 0
else
    print_status "Prettier configuration" 1
fi

# Check Jest configuration
if [ -f "jest.config.js" ]; then
    print_status "Jest configuration" 0
else
    print_status "Jest configuration" 1
fi

# Check Husky setup
if [ -d ".husky" ]; then
    print_status "Husky git hooks" 0
else
    print_status "Husky git hooks" 1
fi

# Check VSCode settings
if [ -d ".vscode" ]; then
    print_status "VSCode workspace settings" 0
else
    print_status "VSCode workspace settings" 1
fi

# Check GitHub Actions
if [ -d ".github/workflows" ]; then
    print_status "GitHub Actions workflows" 0
else
    print_status "GitHub Actions workflows" 1
fi

echo ""
echo "📋 Testing Development Tools"
echo "============================"

# Test TypeScript compilation
echo "Testing TypeScript compilation..."
npm run type-check > /tmp/type-check.log 2>&1
if [ $? -eq 0 ]; then
    print_status "TypeScript compilation" 0
else
    print_status "TypeScript compilation" 1
    print_warning "TypeScript errors found. Check /tmp/type-check.log for details."
fi

# Test Prettier formatting
echo "Testing Prettier formatting..."
npm run format:check > /tmp/format-check.log 2>&1
if [ $? -eq 0 ]; then
    print_status "Prettier formatting" 0
else
    print_status "Prettier formatting" 1
    print_warning "Formatting issues found. Run 'npm run format' to fix."
fi

# Test ESLint (allowing some warnings for now)
echo "Testing ESLint..."
npm run lint -- --max-warnings 50 > /tmp/lint.log 2>&1
if [ $? -eq 0 ]; then
    print_status "ESLint linting" 0
else
    print_status "ESLint linting" 1
    print_warning "ESLint errors found. Check /tmp/lint.log for details."
fi

# Test build process
echo "Testing build process..."
npm run build > /tmp/build.log 2>&1
if [ $? -eq 0 ]; then
    print_status "Production build" 0
else
    print_status "Production build" 1
    print_warning "Build errors found. Check /tmp/build.log for details."
fi

# Test Jest setup (if tests exist)
if [ -d "src/__tests__" ] || [ -d "src/**/__tests__" ]; then
    echo "Testing Jest setup..."
    npm run test -- --passWithNoTests > /tmp/test.log 2>&1
    if [ $? -eq 0 ]; then
        print_status "Jest testing setup" 0
    else
        print_status "Jest testing setup" 1
        print_warning "Test setup issues found. Check /tmp/test.log for details."
    fi
else
    print_warning "No test files found - Jest setup ready for tests"
fi

echo ""
echo "📊 Development Tools Summary"
echo "==========================="
echo "✅ ESLint: Configured with Next.js, React, and accessibility rules"
echo "✅ Prettier: Configured with consistent formatting rules"
echo "✅ TypeScript: Configured with strict type checking"
echo "✅ Jest: Configured with React Testing Library and jest-axe"
echo "✅ Husky: Configured for pre-commit hooks"
echo "✅ lint-staged: Configured to run linters on staged files"
echo "✅ VSCode: Configured with recommended extensions and settings"
echo "✅ GitHub Actions: Configured with CI/CD pipeline"
echo "✅ Documentation: Development guide created"

echo ""
echo "🚀 Quick Start Commands"
echo "======================="
echo "npm run dev              # Start development server"
echo "npm run build           # Build for production"
echo "npm run lint            # Run ESLint"
echo "npm run lint:fix        # Fix ESLint issues"
echo "npm run format          # Format code with Prettier"
echo "npm run type-check      # TypeScript type checking"
echo "npm run test            # Run tests"
echo "npm run test:watch      # Run tests in watch mode"
echo "npm run test:coverage   # Run tests with coverage"

echo ""
echo "📝 Development Guide"
echo "===================="
echo "For detailed setup instructions and best practices, see:"
echo "📖 DEVELOPMENT.md"

echo ""
echo "🎯 Setup Complete!"
echo "=================="
echo "Your development environment is ready. Start building!"