tigger
=========
A simple tool for asking questions based on a given document.

# Demo

## Run 
```sh
go run cmd/cli/main.go
```

## Chat example
```sh
🚩loading assets/test.txt file...done
🚩indexing assets/test.txt file...done
❓question: Who is the author of this story?
🙈 answer: Anton Pavlovich Chekhov (Антон Павлович Чехов).
❓question: Name the characters in this story.
🙈 answer: Based on the context, the characters are:

1. **Генри Пуркуа** (Henri Purqua) – a clown from the circus of the Ginz brothers, the main character.
2. **Благообразный господин** (the dignified gentleman) – Purqua's neighbor at the restaurant table, who is eating blini.

The story is "Глупый француз" by Anton Chekhov. No other characters are named in the provided text.
❓question: What do the characters eat?
🙈 answer: Based on the text, the characters eat **блины (blini/crepes)**, **селянку (selyanka, a type of soup)**, **семгу (salmon)**, and **икру (caviar)**. The neighbor also peppers his selyanka with **cayenne pepper (кайенский перец)**.
```

## Requirements

- **Ollama** installed and running.
- **Chat model:** `qwen3.8:27b`
- **Embedding model:** `nomic-embed-text`
- **Golang:** `1.27.3`