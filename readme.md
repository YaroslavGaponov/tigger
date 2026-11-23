tigger
=========
A simple tool for asking questions based on a given document.

# Demo

## Run as cli
```sh
make demo
```

```sh
❓question: Who is the author of this story?
🙈 answer: The author of this story is Anton Pavlovich Chekhov (Чехов Антон Павлович).

❓question: Name the characters in this story.
🙈 answer: The characters mentioned in this excerpt from Chekhov's "The Foolish Frenchman" are Henry Pourquoа (a clown from the Ginz brothers' circus), a full and respectable gentleman at a neighboring table, and a neighbor who was served selyanka.

❓question:  What do the characters eat?
🙈 answer: The characters eat blini (crepes), solyanka (a Russian soup), salmon, and caviar.
```


## Run as service
```sh
make up
```

```sh
curl --url 'http://localhost:8989/ask' --data-raw $'{\n  "kind": "fast",\n  "question":"почему француз пошел в ресторан"\n}'
{"answer":"Француз Генри Пуркуа, клоун из цирка братьев Гинц, зашел в московский трактир Тестова позавтракать."}

curl --url 'http://localhost:8989/ask' --data-raw $'{\n  "kind": "fast",\n  "question":"какая фамилия Саши"\n}'
{"answer":"Саша Енякина — её фамилия Енякина."}
```


## Requirements

- **Ollama** installed and running.
- **Chat model:** `qwen3.8:27b`
- **Embedding model:** `nomic-embed-text`
- **Golang:** `1.27.3`