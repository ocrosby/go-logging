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
- `task ci`: Run full CI pipeline locally
- `task lint`: Run linting checks
- `task fmt`: Format code