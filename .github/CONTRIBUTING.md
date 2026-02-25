# Contributing to pxBackupManager

Thank you for considering contributing to pxBackupManager! We appreciate your help in making this project better.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. Please be courteous and constructive in all interactions.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps which reproduce the problem**
- **Provide specific examples to demonstrate the steps**
- **Describe the behavior you observed after following the steps**
- **Explain which behavior you expected to see instead and why**
- **Include screenshots and animated GIFs if possible**
- **Include your environment details** (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- **Use a clear and descriptive title**
- **Provide a step-by-step description of the suggested enhancement**
- **Provide specific examples to demonstrate the steps**
- **Describe the current behavior and the expected behavior**
- **Explain why this enhancement would be useful**

### Pull Requests

- Follow the Go code style guidelines
- Include appropriate test coverage
- Update documentation and README if applicable
- Use clear commit messages that explain the changes
- Link related issues in the PR description
- Ensure the code builds successfully on both Linux and Windows

## Development Setup

### Prerequisites

- Go 1.23 or later
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/CodeMeAPixel/pxBackupManager.git
cd pxBackupManager

# Build the project
make build

# Run in development mode
make dev

# Run tests (if available)
go test ./...
```

## Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Keep functions small and focused
- Add comments for exported functions
- Use meaningful variable names

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### Example commit message:

```
Add S3 backup upload functionality

- Implement S3 upload with custom endpoint support
- Add configuration flags for S3 credentials
- Handle S3 errors gracefully

Fixes #123
```

## Testing

Before submitting a pull request:

1. Ensure your code compiles without warnings or errors
2. Test both successful and failure scenarios
3. Test on multiple platforms if possible (Windows, Linux, macOS)
4. Verify backward compatibility

## Documentation

- Update the README if you add new features or change existing ones
- Add comments to your code explaining non-obvious logic
- Update command-line flag documentation in the README
- Include examples for new features

## License

By contributing to pxBackupManager, you agree that your contributions will be licensed under its GNU Affero General Public License v3.0 (AGPL 3.0).

## Questions?

Feel free to reach out to the maintainers:

- **Email**: hey@codemeapixel.dev
- **GitHub Issues**: https://github.com/CodeMeAPixel/pxBackupManager/issues

## Acknowledgments

We appreciate all contributors who help make pxBackupManager better. Your efforts are invaluable!
