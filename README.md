# 🐘 Go + PostgreSQL Template

> Production-ready template for Go backend developers using PostgreSQL, `sqlc`, `migrate`, Docker, and full integration tests.

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org)

## ✨ Features

- ✅ **Type-safe SQL** with [`sqlc`](https://sqlc.dev/)
- ✅ **Database migrations** via [`golang-migrate`](https://github.com/golang-migrate/migrate)
- ✅ **Dockerized PostgreSQL** for local development
- ✅ **Integration tests** with real DB using [`testcontainers-go`](https://github.com/testcontainers/testcontainers-go)
- ✅ **Input validation** with `validator`
- ✅ **Structured logging** with `zerolog`
- ✅ **Graceful shutdown** and connection pooling
- ✅ **Makefile** for common tasks

## 🚀 Quick Start

1. **Start PostgreSQL**
   ```bash
   make db-up
