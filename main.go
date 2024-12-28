package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"go_final_project/tests"
)

var webDir = "./web" // Путь к директории с веб-ресурсами

func openDB(createTable bool) (*sql.DB, error) {
	// Определяем путь к файлу базы данных через переменную окружения или по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	// Открываем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Если необходимо, создаем таблицу и индекс
	if createTable {
		createTableDB := `
            CREATE TABLE IF NOT EXISTS scheduler (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date TEXT NOT NULL,
                title TEXT NOT NULL,
                comment TEXT,
                repeat TEXT CHECK(length(repeat) <= 128)
            );
            CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
        `
		// Выполняем запрос
		_, err = db.Exec(createTableDB)
		if err != nil {
			return nil, err
		}
		log.Println("Таблица scheduler создана.")
	}

	return db, nil
}

func main() {
	port := tests.Port // Используем порт из тестового пакета

	db, err := openDB(true) // Передаем true, чтобы создать таблицы
	if err != nil {
		log.Fatalf("Ошибка при открытии базы данных: %v", err)
	}
	defer db.Close()

	// Создание файлового сервера
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запуск HTTP-сервера
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
