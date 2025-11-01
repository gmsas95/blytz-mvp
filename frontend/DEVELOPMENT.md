# Development Setup Guide

This guide provides comprehensive instructions for setting up the development
environment for the Blytz Frontend application.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Development Tools](#development-tools)
- [Code Quality](#code-quality)
- [Testing](#testing)
- [Git Workflow](#git-workflow)
- [IDE Setup](#ide-setup)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- **Node.js**: Version 18.0.0 or higher
- **npm**: Version 9.0.0 or higher
- **Git**: Latest version

### Recommended Tools

- **VS Code**: Latest version with recommended extensions
- **Docker**: For containerized development
- **Postman**: For API testing

## Quick Start

### 1. Clone and Install

```bash
git clone <repository-url>
cd frontend
npm install
```

### 2. Environment Setup

```bash
cp .env.local.example .env.local
# Edit .env.local with your configuration
```

### 3. Start Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to view the application.

## Development Tools

### Package Manager Scripts

```bash
# Development
npm run dev              # Start development server
npm run build           # Build for production
npm run start           # Start production server

# Code Quality
npm run lint            # Run ESLint
npm run lint:fix        # Fix ESLint issues
npm run format          # Format code with Prettier
npm run format:check    # Check code formatting
npm run type-check      # TypeScript type checking

# Testing
npm run test            # Run tests
npm run test:watch      # Run tests in watch mode
npm run test:coverage   # Run tests with coverage

# Pre-commit hooks
npm run prepare         # Setup git hooks
npm run pre-commit      # Run pre-commit checks
```

### Development Dependencies

The project includes the following development tools:

- **ESLint**: JavaScript/TypeScript linting with Next.js, React, and
  accessibility rules
- **Prettier**: Code formatting with consistent style
- **Husky**: Git hooks for pre-commit checks
- **lint-staged**: Run linters on staged files
- **Jest**: Testing framework with React Testing Library
- **jest-axe**: Accessibility testing
- **TypeScript**: Static type checking

## Code Quality

### ESLint Configuration

Our ESLint setup includes:

- **Next.js recommended rules**: Best practices for Next.js applications
- **React rules**: React-specific linting and hooks validation
- **TypeScript rules**: Strict type checking
- **Import/Export rules**: Proper import organization and sorting
- **Accessibility rules**: jsx-a11y for better accessibility
- **Custom rules**: Project-specific conventions

#### Key Rules

- No unused variables or imports
- Strict TypeScript checking
- Import sorting and organization
- React hooks dependencies validation
- Accessibility compliance

### Prettier Configuration

Consistent code formatting with:

- Single quotes for strings
- Trailing commas in ES5 compatible format
- 100-character line width
- 2-space indentation
- Consistent bracket and spacing rules

### Pre-commit Hooks

We use Husky and lint-staged to ensure code quality:

1. **Pre-commit hook**: Runs linters and formatters on staged files
2. **Commit-msg hook**: Validates commit message format

#### Staged Files Processing

```json
{
  "*.{js,jsx,ts,tsx}": ["eslint --fix", "prettier --write"],
  "*.{json,md,yml,yaml}": ["prettier --write"],
  "*.{css,scss,sass}": ["prettier --write"]
}
```

## Testing

### Test Structure

```
src/
├── __tests__/          # Test files
├── components/
│   └── __tests__/      # Component tests
├── lib/
│   └── __tests__/      # Utility tests
├── hooks/
│   └── __tests__/      # Hook tests
└── services/
    └── __tests__/      # Service tests
```

### Running Tests

```bash
# Run all tests
npm run test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run specific test file
npm run test -- BidButton.test.tsx
```

### Testing Guidelines

#### Component Testing

- Use React Testing Library for component tests
- Test user interactions and behavior
- Test accessibility with jest-axe
- Mock external dependencies

#### Utility Testing

- Test pure functions with various inputs
- Test edge cases and error conditions
- Ensure type safety

#### Test Coverage

- Target 80% coverage minimum
- Focus on critical business logic
- Test error handling and edge cases

### Accessibility Testing

We use jest-axe for automated accessibility testing:

```typescript
import { axe, toHaveNoViolations } from 'jest-axe';

expect.extend(toHaveNoViolations);

test('component has no accessibility violations', async () => {
  const { container } = render(<MyComponent />);
  const results = await axe(container);
  expect(results).toHaveNoViolations();
});
```

## Git Workflow

### Commit Message Format

We follow conventional commits:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

#### Examples

```bash
feat(auction): add real-time bidding functionality
fix(auth): resolve login validation issue
docs(readme): update installation instructions
test(components): add BidButton unit tests
```

### Branch Strategy

- `main`: Production-ready code
- `develop`: Integration branch for features
- `feature/*`: Feature development branches
- `hotfix/*`: Emergency fixes

### Pull Request Process

1. Create feature branch from `develop`
2. Implement changes with tests
3. Ensure all checks pass
4. Create pull request to `develop`
5. Request code review
6. Address feedback and merge

## IDE Setup

### VS Code Configuration

The project includes comprehensive VS Code settings:

#### Extensions

Recommended extensions are listed in `.vscode/extensions.json`:

- **ESLint**: Real-time linting
- **Prettier**: Code formatting
- **Tailwind CSS**: IntelliSense for Tailwind
- **Jest**: Test runner and debugging
- **GitLens**: Enhanced Git capabilities
- **Auto Rename Tag**: HTML/XML tag renaming

#### Settings

Key VS Code settings (`.vscode/settings.json`):

- Format on save enabled
- Auto-fix ESLint issues on save
- Organize imports on save
- TypeScript strict mode
- IntelliSense optimizations

#### Debugging

Debug configurations (`.vscode/launch.json`):

- Debug Next.js client-side code
- Debug server-side code
- Debug tests
- Attach to running process

#### Tasks

Build tasks (`.vscode/tasks.json`):

- Start development server
- Build for production
- Run linting and tests
- Full quality check

### Other IDEs

For other editors, ensure:

- ESLint integration
- Prettier integration
- TypeScript language support
- Git integration

## Environment Configuration

### Development Variables

Create `.env.local` with:

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# Authentication
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key

# Feature Flags
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_DEBUG=true

# External Services
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...
NEXT_PUBLIC_GOOGLE_ANALYTICS_ID=GA-...
```

### Environment Types

- `.env.local`: Local development (ignored by Git)
- `.env.development`: Development environment
- `.env.production`: Production environment

## Performance

### Build Optimization

The build process includes:

- Tree shaking for unused code elimination
- Code splitting for optimal loading
- Image optimization
- Bundle analysis

### Monitoring

- Performance metrics in Next.js
- Bundle size monitoring
- Lighthouse CI integration
- Error tracking

## Troubleshooting

### Common Issues

#### TypeScript Errors

```bash
# Clear TypeScript cache
rm -rf .next/types
rm tsconfig.tsbuildinfo

# Rebuild
npm run build
```

#### ESLint Issues

```bash
# Reset ESLint cache
npx eslint --clear-cache

# Reinstall dependencies
npm ci
```

#### Test Failures

```bash
# Clear Jest cache
npx jest --clearCache

# Update snapshots
npx jest --updateSnapshot
```

#### Development Server Issues

```bash
# Clear Next.js cache
rm -rf .next

# Clear node_modules
rm -rf node_modules package-lock.json
npm install
```

### Getting Help

- Check the [Next.js documentation](https://nextjs.org/docs)
- Review
  [React Testing Library docs](https://testing-library.com/docs/react-testing-library/intro/)
- Join our development team discussions
- Create an issue in the project repository

## Best Practices

### Code Organization

- Use absolute imports with `@/` prefix
- Group related files in directories
- Keep components focused and small
- Separate business logic from UI

### Performance

- Use React.memo for expensive components
- Implement proper loading states
- Optimize images and assets
- Use Next.js Image component

### Security

- Validate all inputs
- Use environment variables for secrets
- Implement proper authentication
- Keep dependencies updated

### Testing

- Write tests for new features
- Test edge cases and error handling
- Use descriptive test names
- Keep tests simple and focused

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

For detailed contribution guidelines, see [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE)
file for details.
