# Shared Queue Library (Minimal)
A pluggable, minimal Go library for queueing jobs with retry, DLQ, deduplication, and metrics hooks.

## Features
- Client-specific queue workers
- Retry and DLQ hooks
- Custom logger & metrics support
- Graceful shutdown
- Delayed jobs
- Optional uniqueness

## Usage
This is a shared library meant to be imported by a parent application.

## Git Addition

```bash
git add .
git commit -m "commit message"
git push
```

## Git Tagging
```bash
git tag v0.1.9
git push origin v0.1.9
```

Note: Increase the above counter after each push

## T0 Fetch in services
go get github.com/Dheerajvd/queue-core@v0.1.9
