# README
System Authentikasi sederhana:
1. Migrasi dan Seeder

---
## Migrasi dan seeder
### 1. Migrasi
#### - Buat file migrasi
```bash
migrate create -ext sql -dir database/migrations -seq create_users_table
```
#### - Jalankan migrasi
```bash
go run ./cmd/migrate/main.go
```
### 2. Seeder
#### - Buat file seeder
```bash
go run ./cmd/seed/create.go user
```
#### - Jalankan seeder
```bash
go run ./cmd/seed/main.go
```
---

---
## Library Thank's
### 1. golang migration v4
```bash
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/mysql
go get github.com/golang-migrate/migrate/v4/source/file
```

### 2. env
```bash
go get github.com/joho/godotenv
```

### 3. (ACID)
```bash
go get github.com/gogaruda/dbtx@v1.0.1
```

### 4. ID Generator
#### - ULID
```bash
go get github.com/oklog/ulid/v2
```
---