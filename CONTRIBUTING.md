# Contributing to go-logging

Thank you for considering contributing to go-logging! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Go version** and **operating system**
- **Code snippets** or **minimal reproducible example**
- **Logs or error messages** (if applicable)

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

- **Clear title and description** of the enhancement
- **Use cases** explaining why this would be useful
- **Proposed API** or interface changes
- **Examples** of how it would be used

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes**:
   - Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
   - Write or update tests for your changes
   - Update documentation as needed
   - Ensure your code follows the existing style
3. **Test your changes**:
   ```bash
   go test ./...
   go fmt ./...
   go vet ./...
   ```
4. **Commit your changes**:
   - Use clear, descriptive commit messages
   - Follow [Conventional Commits](https://www.conventionalcommits.org/) format
   - Example: `feat: add async logging support` or `fix: resolve race condition in logger`
5. **Push to your fork** and submit a pull request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/go-logging.git
cd go-logging

# Install dependencies
go mod download

# Run tests
go test ./pkg/logging/...

# Run tests with coverage
go test ./pkg/logging/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Coding Standards

### Go Code Style

- Follow standard Go formatting: `go fmt ./...`
- Follow Go best practices and idioms
- Use meaningful variable and function names
- Keep functions small and focused
- Avoid premature optimization

### Code Quality

- **Test Coverage**: Aim for 80%+ test coverage for new code
- **Documentation**: Add godoc comments for public APIs
- **Error Handling**: Handle all errors appropriately
- **Thread Safety**: Ensure concurrent safety where needed
- **Performance**: Consider performance implications

### Testing Requirements

All code changes must include tests:

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // Act
    result := logger.DoSomething()
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Documentation

- Add godoc comments for all exported functions, types, and constants
- Update README.md if adding new features
- Add examples to `examples/` directory for significant features
- Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/) format

## Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(fluent): add Float64 method to fluent interface

Add support for logging float64 values in the fluent interface.

Closes #42
```

```
fix(logger): resolve race condition in level change

Use mutex to protect concurrent level changes.

Fixes #38
```

## Pull Request Process

1. **Update documentation** for any API changes
2. **Add tests** for new functionality
3. **Ensure CI passes** all checks
4. **Update CHANGELOG.md** with your changes
5. **Request review** from maintainers
6. **Address feedback** and make necessary changes
7. Once approved, a maintainer will merge your PR

## Code Review Guidelines

### For Contributors

- Be open to feedback and suggestions
- Respond to comments promptly
- Make requested changes or explain why you disagree
- Keep discussions focused and professional

### For Reviewers

- Be respectful and constructive
- Explain the reasoning behind suggestions
- Distinguish between required changes and suggestions
- Approve promptly when requirements are met

## Project Structure

```
go-logging/
├── pkg/
│   └── logging/          # Main package
│       ├── logger.go     # Core interfaces
│       ├── *_test.go     # Tests
│       └── ...
├── examples/             # Example applications
├── docs/                 # Additional documentation
├── .github/              # GitHub workflows and templates
├── README.md             # Project documentation
├── CONTRIBUTING.md       # This file
├── CHANGELOG.md          # Version history
├── LICENSE               # MIT License
└── go.mod                # Go module definition
```

## Release Process

Releases are managed by maintainers:

1. Update CHANGELOG.md with release notes
2. Update version in relevant files
3. Create and push a git tag: `git tag v1.x.x`
4. GitHub Actions will automatically create a release

## Questions?

If you have questions about contributing, please:

1. Check existing documentation
2. Search closed issues for similar questions
3. Open a new discussion on GitHub Discussions
4. Reach out to maintainers if needed

## Recognition

All contributors will be recognized in:
- The project README
- Release notes
- GitHub contributors list

Thank you for contributing to go-logging! 🎉
