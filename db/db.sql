CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    author TEXT NOT NULL,
    abstract TEXT,
    content TEXT NOT NULL,
    theme TEXT
);