FROM python:3.12-alpine

# Устанавливаем переменную окружения для предотвращения буферизации вывода
ENV PYTHONUNBUFFERED 1

WORKDIR /app

RUN apk update && apk add --no-cache libpq-dev gcc python3-dev g++

RUN pip install --upgrade pip setuptools wheel

COPY requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt

COPY . /app/

CMD ["python", "bayes.py"]

EXPOSE 5000