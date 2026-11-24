package server

const ChatHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tigger Chat</title>
    <style>
        * { box-sizing: border-box; }

        body {
            margin: 0;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
            background: #f5f5f5;
            height: 100vh;
            display: flex;
            flex-direction: column;
        }

        header {
            padding: 16px 24px;
            background: white;
            border-bottom: 1px solid #ddd;
            font-size: 18px;
            font-weight: 600;
        }

        #messages {
            flex: 1;
            overflow-y: auto;
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 16px;
        }

        .message {
            max-width: 75%;
            padding: 12px 16px;
            border-radius: 16px;
            line-height: 1.5;
            white-space: pre-wrap;
            overflow-wrap: anywhere;
        }

        .user {
            align-self: flex-end;
            background: #2563eb;
            color: white;
        }

        .assistant {
            align-self: flex-start;
            background: white;
            border: 1px solid #e5e5e5;
        }

        .error { color: #dc2626; }

        form {
            display: flex;
            gap: 12px;
            padding: 16px 24px;
            background: white;
            border-top: 1px solid #ddd;
        }

        input {
            flex: 1;
            padding: 12px 16px;
            border: 1px solid #ccc;
            border-radius: 12px;
            font-size: 16px;
            outline: none;
        }

        input:focus { border-color: #2563eb; }

        button {
            padding: 12px 24px;
            border: none;
            border-radius: 12px;
            background: #2563eb;
            color: white;
            font-size: 15px;
            cursor: pointer;
        }

        button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
    </style>
</head>
<body>

<header>Tigger Chat</header>

<div id="messages"></div>

<form id="chat-form">
    <input
        id="question"
        type="text"
        placeholder="Ask a question..."
        autocomplete="off"
        required
    >
    <button id="send" type="submit">Send</button>
</form>

<script>
    const API_URL = "/ask";

    const form = document.getElementById("chat-form");
    const input = document.getElementById("question");
    const messages = document.getElementById("messages");
    const sendButton = document.getElementById("send");

    function addMessage(text, type) {
        const element = document.createElement("div");
        element.className = "message " + type;
        element.textContent = text;
        messages.appendChild(element);
        messages.scrollTop = messages.scrollHeight;
        return element;
    }

    form.addEventListener("submit", async (event) => {
        event.preventDefault();

        const question = input.value.trim();
        if (!question) return;

        addMessage(question, "user");

        input.value = "";
        input.disabled = true;
        sendButton.disabled = true;

        const loading = addMessage("Thinking...", "assistant");

        try {
            const response = await fetch(API_URL, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    kind: "fast",
                    question: question
                })
            });

            if (!response.ok) {
                throw new Error("HTTP " + response.status);
            }

            const data = await response.json();
            loading.textContent = data.answer || "Empty response";

        } catch (error) {
            loading.textContent = "Error: " + error.message;
            loading.classList.add("error");
        } finally {
            input.disabled = false;
            sendButton.disabled = false;
            input.focus();
        }
    });
</script>

</body>
</html>`
