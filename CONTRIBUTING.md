# Contributing to GOURL

Thank you for your interest in contributing to GOURL! 🎉

## How to Contribute

### Reporting Bugs

- Use the GitHub issue tracker
- Include steps to reproduce
- Provide error messages and logs
- Mention your Go version and OS

### Suggesting Features

- Open an issue with the `enhancement` label
- Describe the use case
- Explain why it would be useful

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit with clear messages (`git commit -m 'Add amazing feature'`)
5. Push to your fork (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Code Style

- Follow Go conventions
- Use `gofmt` to format code
- Add comments for exported functions
- Keep functions focused and small

### Testing

- Add tests for new features
- Ensure all tests pass (`go test ./...`)
- Test manually before submitting

### Documentation

- Update README if needed
- Add code comments
- Update API documentation

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gourl.git
cd gourl

# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go
```

## Questions?

Feel free to open an issue for any questions!

