# Opencode Configuration

## Commit Message Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) with the Angular convention.

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation

### Scopes (optional)
- **logging**: Core logging functionality
- **config**: Configuration related changes
- **async**: Async worker and related functionality
- **handlers**: Handler interfaces and implementations
- **formatters**: Log formatters (JSON, text, console)
- **middleware**: HTTP middleware
- **providers**: Dependency injection providers
- **examples**: Example code
- **ci**: Continuous integration
- **docs**: Documentation

### Examples
```
feat(logging): add structured logging support
fix(async): resolve worker shutdown race condition  
docs(readme): update installation instructions
test(handlers): add comprehensive handler tests
chore(ci): update GitHub Actions workflow
```

### Breaking Changes
For breaking changes, add `!` after the type/scope:
```
feat(config)!: change default log level to INFO
```

Or include `BREAKING CHANGE:` in the footer:
```
feat(config): add new configuration options

BREAKING CHANGE: Default log level changed from DEBUG to INFO
```

## Development Commands

- `task test-coverage-check`: Run tests with 85% coverage validation
- `task ci`: Run full CI pipeline locally (recommended for CI systems)
- `task lint`: Run linting checks
- `task fmt`: Format code

## CI Integration

### GitHub Actions
The repository includes a properly configured GitHub Actions workflow (`.github/workflows/ci.yml`) that:
- Uses `task ci-test` for proper coverage validation
- Includes fallback validation for CI systems without task runner
- Publishes coverage reports to Codecov
- Uses appropriate timeouts and error handling

### Jenkins
A complete Jenkinsfile is provided with:
- Proper Go setup and dependency management
- Task runner installation and usage  
- Coverage validation with 85% threshold
- HTML coverage report publishing
- Integration test execution with proper timeouts
- Comprehensive error handling and troubleshooting messages

### For Jenkins/CI Systems
**Important**: Use `task ci` or `task test-coverage-check` instead of `go test ./...`

The project has a 85% test coverage requirement. Running `go test ./...` will include the mocks package (which has low coverage by design) and cause CI failures. 

**Correct CI Commands:**
```bash
# Full CI pipeline (recommended)
task ci

# Or just test with coverage validation
task test-coverage-check
```

**Fallback for systems without Task runner:**
```bash
# Test only the main logging package with proper timeout
go test -v -timeout=60s -coverprofile=logging-coverage.out ./pkg/logging

# Validate coverage meets 85% threshold
COVERAGE=$(go tool cover -func=logging-coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 85" | bc -l) )); then
  echo "❌ Coverage ${COVERAGE}% is below required 85% threshold"
  exit 1
fi
echo "✅ Coverage ${COVERAGE}% meets the 85% threshold requirement"
```

**Incorrect Commands (will fail in CI):**
```bash
go test ./...  # Includes mocks, reports combined low coverage
go test -coverprofile=coverage.out ./...  # Same issue
```

The coverage validation specifically tests `./pkg/logging` which maintains 86%+ coverage.