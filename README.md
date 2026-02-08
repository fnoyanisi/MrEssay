# MrEssay

MrEssay is a Go based system that helps create, and review student writing using LLMs and Telegram.

Students submit handwritten work via Telegram. The system extracts the text using a vision-capable LLM (OCR), reviews the content, and sends structured feedback back to the student.

## How it works

1. An **Agent** sends an essay topic to Telegram account (the user gets a notification)
2. The user finishes the essay, sends the photo to **Telegram bot**.
3. Another **Agent** performs OCR
   * The photo is sent to an LLM
   * Exact text is transcribed (no correction or inference)
4. Another **Reviewer agent** reviews student's response against the topic.
5. The **final feedback** is returned to the student via Telegram

Each step is a sequential LLM call managed explicitly in Go (no automated orchestration, fancy DAGs or Temporal workflows).

## Project structure

```
cmd/mressay/
  main.go            # Application entry point

lib/openai/
  agent.go           # Agent definition (model, memory, config)
  client.go          # OpenAI HTTP client
  requests.go        # Request/response types and helpers

lib/telegram/
  bot.go              # Telegram bot integration
```

Use `% go build -o bin/chat cmd/mressay/main.go` to build the project and run `% ./bin/chat` to execute. You can bundle this into a Docker container, use cron to schedule the executions.

## Some design decisions

* **Agents have memory** (conversation history is stored client-side)
* **Models are stateless**; all context is sent per request
* **No tool calling or orchestrator agent** jsut KISS
* **Sequential, explicit workflow** for clarity and control (no DAGs, no Temporal)
* OCR is treated as a **vision LLM call**, not an external tool (again, me being lazy)


## Requirements

* Go 1.21+
* OpenAI API key (store in `OPENAI_API_KEY` environment variabel)
* Telegram Bot token (store in `TELEGRAM_API_TOKEN` environment variabel)
* Telegram User ID (store in `TELEGRAM_ID` environment variabel)


## Status

Fun project - probably needs a lot of additions.

* Not thread safe
* The Telegram loop blocks (goroutines can be used)
* No tooling support for Agent
* Sequential workflow
* Possibly can use the raw Telegram api (I have used a wrapper around a wrapper :/)
