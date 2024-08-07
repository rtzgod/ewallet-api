FROM postgres:latest

ENV POSTGRES_USER=testuser
ENV POSTGRES_PASSWORD=123
ENV POSTGRES_DB=testdb

COPY init.sql /docker-entrypoint-initdb.d/