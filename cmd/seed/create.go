package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

const tmpl = `package seeders
import (
	"database/sql"
	"fmt"
)

func {{.FuncName}}(db *sql.DB) error {
	fmt.Println("ayo buat seeder...")
	return nil
}
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("gunakan: go run ./cmd/seed/create.go NamaSeeder")
		return
	}

	name := os.Args[1]
	funcName := strings.Title(name)
	fileName := strings.ToLower(name) + "Seeder.go"

	f, err := os.Create("database/seeders/" + fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t := template.Must(template.New("seeder").Parse(tmpl))
	t.Execute(f, map[string]string{"FuncName": funcName})

	fmt.Println("seeder", fileName, "berhasil dibuat")
}
