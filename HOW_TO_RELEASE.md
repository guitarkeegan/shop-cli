# 1. Push feature branch & merge to main (via PR or direct merge)
```bash
git checkout main
git pull origin main
```

# 2. Prepare for release
```bash
go mod tidy                    # Clean up dependencies
git add go.mod go.sum
git commit -m "Prepare for release"  # Only if there were changes
git push origin main           # Push the commit to remote
```

# 3. Tag the release (git tag to see all of them)
> use the following format
```bash
git tag v0.0.0 
git push origin v0.0.0
```

# 4. Create the release
```bash
goreleaser release --clean
```
