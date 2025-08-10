#!/bin/bash

# TigerGraph CLI Release Script
# Usage: ./scripts/release.sh v1.0.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if version is provided
if [ $# -eq 0 ]; then
    log_error "Please provide a version tag (e.g., v1.0.0)"
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    log_error "Version must be in format vX.Y.Z or vX.Y.Z-suffix (e.g., v1.0.0, v1.0.0-beta)"
    exit 1
fi

log_info "Starting release process for version: $VERSION"

# Check if we're on the right branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "master" && "$CURRENT_BRANCH" != "main" ]]; then
    log_warning "You're not on master/main branch. Current branch: $CURRENT_BRANCH"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Release cancelled"
        exit 1
    fi
fi

# Check if working directory is clean
if ! git diff-index --quiet HEAD --; then
    log_error "Working directory is not clean. Please commit or stash your changes."
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    log_error "Tag $VERSION already exists"
    exit 1
fi

# Run tests
log_info "Running tests..."
if command -v make >/dev/null 2>&1; then
    make test-short
else
    go test -v -short ./...
fi
log_success "Tests passed!"

# Update version in constants if needed
log_info "Updating version in constants..."
VERSION_NUMBER=${VERSION#v}
sed -i.bak "s/VERSION_CLI.*=.*\".*\"/VERSION_CLI = \"$VERSION_NUMBER\"/" pkg/constants/constants.go
rm -f pkg/constants/constants.go.bak

# Build and test
log_info "Building and testing..."
if command -v make >/dev/null 2>&1; then
    make build
else
    go build -o tg ./cmd/main.go
fi

# Test the built binary
./tg version
log_success "Build successful!"

# Commit version update if there are changes
if ! git diff-index --quiet HEAD --; then
    log_info "Committing version update..."
    git add pkg/constants/constants.go
    git commit -m "chore: bump version to $VERSION"
fi

# Create and push tag
log_info "Creating and pushing tag..."
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$CURRENT_BRANCH"
git push origin "$VERSION"

log_success "Tag $VERSION created and pushed!"
log_info "GitHub Actions will now build and create the release automatically."
log_info "Check the progress at: https://github.com/zrougamed/tgCli/actions"

# Clean up
rm -f tg

log_success "Release process completed!"
echo
echo "Next steps:"
echo "1. Monitor the GitHub Actions workflow"
echo "2. Review the generated release notes"
echo "3. Test the release artifacts"
echo "4. Announce the release!"