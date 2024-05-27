CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    site TEXT,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    author TEXT NOT NULL,
    abstract TEXT,
    content TEXT NOT NULL
);