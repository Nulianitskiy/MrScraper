package db

import (
	"MrScraper/model"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

type Database struct {
	db    *sqlx.DB
	mutex sync.Mutex
}

func NewDatabase() (*Database, error) {
	connectionString := "user=dbuser password=bonobo dbname=db host=localhost port=5436 sslmode=disable"
	db, err := sqlx.Open("pgx", connectionString)
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

func (d *Database) InsertArticle(article model.Article, site string) error {
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

	_, err = d.db.Exec("INSERT INTO articles (site, title, url, author, abstract, content) VALUES ($1, $2, $3, $4, $5, $6)",
		site, article.Title, article.Link, article.Authors, article.Annotation, article.Text)
	return err
}

// GetArticles возвращает все статьи из базы данных
func (d *Database) GetArticles() ([]model.Article, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.db.Query("SELECT site, title, url, author, abstract, content FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		err := rows.Scan(&article.Site, &article.Title, &article.Link, &article.Authors, &article.Annotation, &article.Text)
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
