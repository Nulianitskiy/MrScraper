from flask import Flask, request, jsonify
from flask_cors import CORS
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.preprocessing import LabelEncoder
from sklearn.model_selection import train_test_split
from sklearn.naive_bayes import MultinomialNB
from sklearn.metrics import classification_report
import re
import nltk
import psycopg2
import pandas as pd
from sqlalchemy import create_engine
from nltk.corpus import stopwords
from nltk import word_tokenize
from nltk.stem import WordNetLemmatizer
from urllib.parse import unquote

# загрузка необходимых ресурсов
nltk.download('punkt')
nltk.download('stopwords')
nltk.download('wordnet')

app = Flask(__name__)
CORS(app)

vectorizer = None
label_encoder = None
model = None


def get_data_from_db():
    """
    Извлекает данные из базы данных PostgreSQL.
    """
    engine = create_engine("postgresql+psycopg2://dbuser:bonobo@db:5432/scraperdb")

    query = "SELECT content, theme FROM articles"
    df = pd.read_sql(query, engine)

    # Применяем функцию clean к столбцу content
    df['content'] = df['content'].apply(clean)

    return df


def clean(text):
    # Очистка по регулярке: оставляем только английские буквы и пробелы
    text = re.sub(r"[^\nA-Za-z -]", "", text)

    # приведение к lowercase
    tokens = word_tokenize(text.lower())

    # Удаляем знаки препинания
    punctuation_marks = ['!', ',', '(', ')', ':', '-', '?', '.', '..', '...']
    only_words = [token for token in tokens if token not in punctuation_marks]

    # Приводим слова к нормальной форме
    lemmatizer = WordNetLemmatizer()
    lemmas = [lemmatizer.lemmatize(token) for token in only_words]

    # Удаление стоп-слов
    stop_words = stopwords.words("english")
    filtered_words = [token for token in lemmas if token not in stop_words]

    return ' '.join(filtered_words)


@app.route('/train', methods=['POST'])
def train():
    """
    Обучает модель на данных из базы данных.
    """
    global vectorizer, label_encoder, model

    # Извлекаем данные из базы данных
    df = get_data_from_db()

    print("Данные из базы данных загружены")
    print(df.head())

    vectorizer = TfidfVectorizer()

    x = vectorizer.fit_transform(df['content'])
    print("Данные векторизованы")

    # Преобразуем строковые метки в числовые
    label_encoder = LabelEncoder()
    y = label_encoder.fit_transform(df['theme'])

    # Проверим баланс классов
    class_counts = pd.Series(y).value_counts()
    print("Распределение классов:")
    print(class_counts)

    # Разделяем данные на обучающий и тестовый наборы
    x_train, x_test, y_train, y_test = train_test_split(
        x, y, test_size=0.2, random_state=43
    )

    print("Данные разделены на обучающую и тестовую выборки")

    model = MultinomialNB()

    model.fit(x_train, y_train)
    print("Модель обучена")

    # Прогнозируем метки классов для тестового набора
    y_pred = model.predict(x_test)

    report = classification_report(y_test, y_pred, zero_division=0)
    print(report)

    return jsonify({'accuracy': report})

@app.route('/predict', methods=['POST'])
def predict():
    """
    Возвращает предсказание для нового текста статьи.
    """
    global vectorizer, label_encoder, model

    if model is None or vectorizer is None or label_encoder is None:
        return jsonify({'error': 'Model is not trained yet.'}), 400

    # Получаем текст статьи из запроса
    data = request.get_json()
    raw_text = data.get('text')
    if not raw_text:
        return jsonify({'error': 'No text provided.'}), 400

    # Декодируем URL-закодированный текст
    new_article = clean(unquote(raw_text))

    if not new_article:
        return jsonify({'error': 'No text provided after cleaning.'}), 400

    # Векторизуем новую статью
    new_article_vector = vectorizer.transform([new_article])

    # Предсказываем метку класса для новой статьи
    prediction_proba = model.predict_proba(new_article_vector)
    prediction_confidence = prediction_proba.max()
    predicted_label_idx = prediction_proba.argmax()

    if prediction_confidence < 0.5:
        return jsonify({'prediction': 'Uncertain', 'confidence': prediction_confidence})

    predicted_label = label_encoder.inverse_transform([predicted_label_idx])[0]

    return jsonify({'prediction': predicted_label, 'confidence': prediction_confidence})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
