package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func LoadEnv() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла, используем значения по умолчанию")
	}
}

var (
	webDir = "./web" // Путь к директории с веб-ресурсами
)

func openDB(scheduler bool) (*sql.DB, error) {
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
	if scheduler {
		var scheduler = `
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
		_, err = db.Exec(scheduler)
		if err != nil {
			log.Printf("Ошибка при создании таблицы: %v", err)
			return nil, err
		}
		log.Println("Таблица scheduler создана или уже существует.")
	}

	return db, nil
}

func main() {
	// Загружаем переменные окружения
	LoadEnv()

	// Получаем порт из переменной окружения или используем значение по умолчанию
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // Значение по умолчанию
	}
	log.Println("Используем порт:", portStr) // Логируем порт для отладки

	// Преобразуем порт в целое число
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Ошибка преобразования порта: %v", err)
	}

	// Открываем базу данных и создаем таблицы
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
