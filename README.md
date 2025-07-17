# Auth Service
System Authentikasi sederhana:
1. Migrasi dan Seeder
2. Manajemen User
3. Auth System (login, logout, register tamu) 
4. middleware
    - Cors
    - Auth
    - Role
    - Verify
5. Verify Email
6. Refresh Token

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
```bash
# 1. golang migration v4
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/mysql
go get github.com/golang-migrate/migrate/v4/source/file

# 2. env
go get github.com/joho/godotenv

# 3. (ACID)
go get github.com/gogaruda/dbtx@v1.0.1

# 4. ID Generator
go get github.com/oklog/ulid/v2
go get github.com/google/uuid

# 5. Error Handling System
go get github.com/gogaruda/apperror@v1.3.0

# 6. GIN
go get -u github.com/gin-gonic/gin

# 7. Validasi
go get github.com/gogaruda/valigo@v1.0.2

# 8. JWT
go get github.com/golang-jwt/jwt/v5

# 9. Email
go get gopkg.in/gomail.v2

# 10. cors
go get github.com/gin-contrib/cors

# 11. Swagger
go install github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
```
