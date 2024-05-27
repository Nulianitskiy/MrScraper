FROM postgres:alpine3.18

ENV POSTGRES_PASSWORD=bonobo \
    POSTGRES_USER=dbuser \
    POSTGRES_DB=db

COPY ./db.sql /docker-entrypoint-initdb.d/db.sql

EXPOSE 5432