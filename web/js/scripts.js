async function fetchArticles() {
    const response = await fetch('/articles');
    const articles = await response.json();
    const articlesList = document.getElementById('articles');
    articles.forEach(article => {
        const listItem = document.createElement('li');
        listItem.innerHTML = `
            <h2>${article.title}</h2>
            <p><strong>Authors:</strong> ${article.authors}</p>
            <p><strong>Topic:</strong> ${article.theme}</p>
            <p>${article.annotation}</p>
        `;
        listItem.addEventListener('click', () => {
            window.location.href = `article?id=${article.id}`;
        });
        articlesList.appendChild(listItem);
    });
}

async function fetchArticle() {
    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get('id');
    const response = await fetch(`/articles/${id}`);
    const article = await response.json();
    document.getElementById('title').textContent = article.title;
    document.getElementById('authors').textContent = article.authors;
    document.getElementById('topic').textContent = article.theme;
    document.getElementById('annotation').textContent = article.annotation;
    document.getElementById('text').textContent = article.text;
    document.getElementById('link').href = article.link;
}

async function checkBelonging() {
    const inputText = document.getElementById('inputText').value;
    const encodedText = encodeURIComponent(inputText); // Кодируем текст
    const response = await fetch('http://localhost:5000/predict', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ text: encodedText }),
    });
    const result = await response.json();
    document.getElementById('result').textContent = `Предсказание: ${result.prediction} Уверенность: ${result.confidence}`;
}

if (document.getElementById('articles')) {
    fetchArticles();
} else if (document.getElementById('title')) {
    fetchArticle();
} else if (document.getElementById('checkButton')) {
    document.getElementById('checkButton').addEventListener('click', checkBelonging);
}
