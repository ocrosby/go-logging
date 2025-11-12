# Test Coverage

This project enforces a **85%** test coverage threshold for the main logging package (`pkg/logging`) to ensure code quality and reliability.

## Coverage Requirements

- **Minimum threshold**: 85%
- **Scope**: `pkg/logging` package only (excludes examples and mocks)
- **Enforcement**: CI pipeline will fail if coverage drops below the threshold

## Running Coverage Checks

### Local Development

```bash
# Run tests with coverage validation
task test-coverage-check

# Run traditional coverage report with HTML output
task test-coverage

# Run all quality checks (includes coverage validation)
task check
```

### CI Pipeline

The coverage check runs automatically in the CI pipeline using the local build system:
- Installs the Task runner
- Executes `task ci-test` which includes coverage validation
- Uses the same threshold and logic as local development
- The pipeline will fail if tests don't pass or coverage drops below 85%

## Coverage Configuration

The coverage threshold is controlled by the **local build system** in `Taskfile.yml`:

```yaml
vars:
  COVERAGE_THRESHOLD: 85.0
```

The CI workflow uses the same local build commands, ensuring consistency between local development and CI environments. To change the threshold, update only the `Taskfile.yml` file - the CI will automatically use the new value.

## Coverage Exclusions

The following are excluded from coverage calculations:
- Example packages (`examples/*`)
- Mock implementations (`pkg/logging/mocks/*`)
- Test files (`*_test.go`)

## Understanding Coverage Reports

When coverage is below the threshold, the system will show:
- Current coverage percentage
- Required threshold
- List of functions/files with the lowest coverage
- Specific areas that need more tests

Example output:
```
❌ Coverage 83.2% is below the required 85% threshold

Functions with low coverage:
pkg/logging/handler_composition.go:96:   NewConditionalHandler    0.0%
pkg/logging/outputs.go:322:             NewRotatingFileOutput    0.0%
...
```

## Best Practices

1. **Write tests for new features** before the coverage drops
2. **Focus on critical paths** when improving coverage
3. **Use meaningful test scenarios** rather than just achieving numbers
4. **Test edge cases** and error conditions
5. **Keep tests maintainable** and readable

## Architecture

The coverage validation system is designed with these principles:
- **Single source of truth**: Coverage threshold controlled by `Taskfile.yml`
- **Local-first**: CI uses the same commands developers run locally
- **Consistent behavior**: Same validation logic in all environments
- **Easy maintenance**: Change threshold in one place only

## Current Status

The project currently maintains **85.4%** test coverage, meeting the required threshold.

## Available Tasks

```bash
# Individual tasks
task test-coverage-check    # Coverage validation only
task test-coverage         # Traditional coverage with HTML report

# Composite tasks  
task check                 # All quality checks (format, lint, coverage, build)
task ci-test              # CI test pipeline (mocks, wire, coverage)
task ci                   # Full CI pipeline (all checks)
```