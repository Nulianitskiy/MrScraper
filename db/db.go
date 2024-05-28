package db

import (
	"MrScraper/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

type Database struct {
	db    *sqlx.DB
	mutex sync.Mutex
}

var (
	instance *Database
	once     sync.Once
)

func GetInstance() (*Database, error) {
	var err error
	once.Do(func() {
		instance, err = newDatabase()
		if err != nil {
			log.Fatal("Ошибка создания экземпляра Database:", err)
		}
	})
	return instance, err
}

// postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable
func newDatabase() (*Database, error) {
	connectionString := "postgres://dbuser:bonobo@db:5432/scraperdb?sslmode=disable"
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Ping базы данных для проверки подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки подключения к базе данных:", err)
	}
	log.Println("Подключение к базе данных PostgreSQL успешно")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) InsertArticle(article model.Article) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Проверка наличия статьи с таким же названием
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM articles WHERE title = $1", article.Title).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("статья с названием %s уже существует в базе данных", article.Title)
	}

	_, err = d.db.Exec("INSERT INTO articles (title, url, author, abstract, content) VALUES ($1, $2, $3, $4, $5)",
		article.Title, article.Link, article.Authors, article.Annotation, article.Text)
	return err
}

// GetArticles возвращает все статьи из базы данных
func (d *Database) GetArticles() ([]model.Article, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.db.Query("SELECT title, url, author, abstract, content FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		err := rows.Scan(&article.Title, &article.Link, &article.Authors, &article.Annotation, &article.Text)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
