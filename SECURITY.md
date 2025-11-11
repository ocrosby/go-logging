# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of go-logging seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please do the following:

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report security vulnerabilities by emailing:
- **Email**: [INSERT SECURITY EMAIL]
- **Subject**: [SECURITY] Brief description of the issue

### What to include in your report:

1. **Description** of the vulnerability
2. **Steps to reproduce** the issue
3. **Potential impact** of the vulnerability
4. **Suggested fix** (if you have one)
5. **Your contact information** for follow-up

### What to expect:

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Regular Updates**: We will keep you informed about our progress
- **Fix Timeline**: We aim to release fixes for critical vulnerabilities within 30 days
- **Credit**: With your permission, we will credit you in the security advisory

## Security Best Practices

When using go-logging in your application:

### 1. Sensitive Data Handling

```go
// DO: Use built-in redaction
config := logging.NewConfig().
    AddRedactPattern(`password=\w+`).
    AddRedactPattern(`token=\w+`).
    AddRedactPattern(`apiKey=\w+`).
    Build()

// DON'T: Log sensitive data directly
logger.Info("User credentials: %s", password) // ❌ Never log passwords
```

### 2. Custom Redaction Patterns

Always configure redaction patterns for your specific sensitive data:

```go
config := logging.NewConfig().
    AddRedactPattern(`(?i)secret[_-]?key=\w+`).
    AddRedactPattern(`(?i)auth[_-]?token=\w+`).
    AddRedactPattern(`credit[_-]?card=\d+`).
    Build()
```

### 3. Environment Variables

Be careful with environment variable logging:

```go
// DO: Redact before logging
logger.Info("Config loaded from env")

// DON'T: Log entire environment
logger.Info("Env: %v", os.Environ()) // ❌ May contain secrets
```

### 4. HTTP Headers

Use selective header logging:

```go
// DO: Log specific safe headers
logging.RequestLogger(logger, "User-Agent", "Content-Type")

// DON'T: Log all headers
logging.RequestLogger(logger, "*") // ❌ May contain auth tokens
```

### 5. Error Messages

Be cautious with error logging:

```go
// DO: Log sanitized errors
logger.Fluent().Error().
    Str("operation", "database_connect").
    Msg("Connection failed")

// DON'T: Log full error with connection strings
logger.Error("DB error: %v", err) // ❌ May contain credentials
```

## Known Security Considerations

### Thread Safety

The logger is thread-safe and can be used concurrently. However:

- Dynamic level changes with `SetLevel()` are synchronized
- Field modifications create new logger instances (immutable pattern)

### Memory Safety

- No known memory leaks
- Proper cleanup of resources
- Bounded memory usage for log buffers

### Dependency Security

This library has minimal dependencies:
- Standard Go library (maintained by Go team)
- `github.com/google/uuid` (v1.6.0) - widely used and maintained

We regularly update dependencies and monitor for security advisories.

## Security Updates

Security updates will be released as:
- Patch versions for minor issues
- Minor versions for moderate issues
- Major versions for breaking security changes

Subscribe to releases to stay informed:
- Watch the repository on GitHub
- Enable notifications for security advisories

## Vulnerability Disclosure Timeline

We follow responsible disclosure practices:

1. **Day 0**: Vulnerability reported
2. **Day 2**: Report acknowledged
3. **Day 7**: Initial assessment and response plan
4. **Day 30**: Target fix release (for critical issues)
5. **Day 90**: Public disclosure (coordinated with reporter)

## Bug Bounty Program

We do not currently have a bug bounty program. However, we deeply appreciate security researchers who report vulnerabilities responsibly and will:

- Acknowledge your contribution in security advisories
- Credit you in release notes (with your permission)
- Provide reference for future security research

## Questions?

If you have questions about this security policy, please contact:
- **Email**: [INSERT CONTACT EMAIL]
- **GitHub Discussions**: For general security questions

Thank you for helping keep go-logging and its users safe!
