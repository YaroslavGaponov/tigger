package oracle

import "github.com/YaroslavGaponov/tigger/internal/ollama"

var (
	AssistantPrompt = ollama.Message{
		Role: "system",
		Content: `You are an assistant that answers questions exclusively based on the provided context.
Answer using only one sentence.
If the answer cannot be found in the context, respond: "Information not found in the text."
Do not use external knowledge.`,
	}

	EstimatePrompt = ollama.Message{
		Role: "system",
		Content: `Given a question and a text, evaluate how relevant the text is for answering the question.
Return only a single integer from 0 to 100:
- 0 — completely irrelevant.
- 50 — partially relevant but insufficient.
- 100 — directly and fully answers the question.
Do not include explanations, comments, or any other text.`,
	}
)
