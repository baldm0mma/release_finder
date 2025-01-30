# Release Finder

A command-line tool to identify which released versions of Grafana contain a specific commit. This tool is particularly useful for developers and users who want to track which Grafana releases include specific features, bug fixes, or changes.

## Features

- Finds all released versions containing a specific commit
- Supports both full and shortened commit hashes
- Handles Grafana's versioning scheme, including:
  - Standard semantic versions (v1.2.3)
  - Pre-release versions (v1.2.3-beta)
  - Security patches (v1.2.3+security)
- Shows results with clear visual indicators (✓)

## Prerequisites

- Go 1.x or higher
- Git
- Access to a local copy of the Grafana repository

## Installation

1. Clone this repository:

```bash
git clone https://github.com/yourusername/release-finder
cd release-finder
```

2. Sync dependencies:

```bash
go mod tidy
```

3. Build the program:

```bash
go build
```

## Usage

Run the program by providing two arguments:

1. Path to your local Grafana repository
2. The commit hash you want to check

```bash
./release-finder <grafana-repo-path> <commit-hash>
```

For example:

```bash
./release-finder ~/grafana a1b2c3d4e5f6
```

You can also run it directly without building:

```bash
go run main.go ~/grafana a1b2c3d4e5f6
```

### Example Output

```
Results for commit a1b2c3d4e5f6:

✓ v9.0.0
✓ v9.0.1
✓ v9.0.2
✓ v9.1.0
```

This indicates that the specified commit is present in all the listed releases.

If a commit isn't in any release yet, you'll see:

```
This commit is not in any publicly released version yet.
```

## Version Handling

The tool understands Grafana's versioning scheme:

- Regular versions (e.g., v9.0.0)
- Pre-release versions (e.g., v9.0.0-beta1) are treated as older than their release version
- Security patches (e.g., v9.0.0+security) are treated as newer than their base version

Version comparison follows these rules:

- v1.0.0-beta < v1.0.0 (pre-release is older)
- v1.0.0 < v1.0.0+security (security patch is newer)
- v1.0.0+security < v1.0.0+security2 (ordered numerically)
- v1.0.0+security < v1.1.0 (next minor version is newer)

## Development

### Running Tests

Run all tests with verbose output:

```bash
go test -v
```

The test suite includes:

- Version comparison tests
- Release finding tests with mock repositories
- Various versioning scheme scenarios

### Code Structure

- `main.go`: Core program logic
- `main_test.go`: Test suite

Key functions:

- `findReleases`: Identifies which releases contain a commit
- `compareVersions`: Handles version comparison logic
- `displayReleases`: Formats and displays results

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

[Your chosen license]

## Acknowledgments

- Grafana team for their versioning scheme
- [Any other acknowledgments]
