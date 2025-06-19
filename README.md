# Rattle

Rattle is a Docker log scanner that sends alerts to Telegram.  
It includes a [Telegram Mini App](https://core.telegram.org/bots/webapps) for managing alerts, filters, and container rules.

---

## 📚 Table of Contents

- [Getting Started](#️-getting-started)
- [Features](#-features)
- [Setup & Development](#️-setup--development)
- [Environment Variables](#️-environment-variables)
- [Docker (Production)](#-docker-production)
- [Log Examples](#-log-examples)
- [Tech Stack](#-tech-stack)
- [License](#-license)

---

## ▶️ Getting Started

> 💡 If you only use the light version (`scanner` + `postgres`), you don't need to configure a domain or Telegram Mini App.

### ✅ Requirements

- ✅ A domain name (e.g. `rattle.example.com`)
- ✅ Valid SSL certificate (Telegram **requires HTTPS** for WebApps)

---

### 🤖 Create a Telegram bot

1. Open [@BotFather](https://t.me/BotFather)
2. Send `/newbot`
3. Follow the prompts to name your bot (e.g. `RattleBot`) and get a **token**
4. Save the token in your `.env` file:

    ```env
    TELEGRAM_BOT_TOKEN=your_bot_token
    ```

---

### 🧩 Configure Mini App (via BotFather)

1. `/mybots` → **Select your bot**
2. → **Bot Settings** → **Configure Mini App** → **Edit Mini App URL**  
   → Enter: `https://rattle.example.com`
3. → **Bot Settings** → **Menu Button** → **Edit menu button URL**  
   → Enter: `https://rattle.example.com`

> ✅ Now your bot will show a web button that opens the Rattle interface.

---

### 🚀 Launch Rattle

Once `.env` is configured, launch the stack with:

```bash
docker compose -f docker-compose.yml up
```

> Or use the light version:

```bash
docker compose -f docker-compose.light.yml up
```

---

## 🚀 Features

- 📦 Real-time Docker container log monitoring
- 📤 Sends alerts to Telegram chats
- ⚙️ Fully configurable via `.env` or the Telegram Mini App
- 🧠 Supports regex-based pattern filtering for logs (error, info, success, etc.)
- 🔒 Per-chat access levels (admin / user)
- 🛠️ Built-in PostgreSQL backend for storing filters, access settings, and rules

---

## 🛠️ Setup & Development

Clone the repository and create a `.env` file:

```bash
cp .env.example .env
```

Then you can run the full stack using Docker Compose:

```bash
docker compose up
```

> Alternatively, use the light version if you only want scanning + Telegram alerts:

```bash
docker compose -f docker-compose.light.yml up
```

---

## ⚙️ Environment Variables

All configuration is handled via `.env` file.

Create it from the template:

```bash
cp .env.example .env
```

Common options include:

```env
# Telegram bot token (create via @BotFather)
TELEGRAM_BOT_TOKEN=your_bot_token

# Comma-separated chat IDs (optional if using Telegram Mini App)
# Tip: Your personal Telegram ID also works here.
# Use a bot like @getmyid_bot or web.telegram.org (like https://web.telegram.org/k/#-1234567890 , `-1234567890` is Group ID) to find ID.
TELEGRAM_CHAT_IDS=12345678,-98765432 # Can be your personal Telegram user ID or group ID

# Regex filters (optional if using Telegram Mini App)
INCLUDE_PATTERNS_ERROR=(?i)(^|\s|:|]|\[)ошибка(:|\s|$)
EXCLUDE_PATTERNS=(?i)timeout|heartbeat

# Database
POSTGRES_USER=rattleuser
POSTGRES_PASSWORD=your_password_here
POSTGRES_DB=rattle

# Ports
POSTGRES_PORT_EXTERNAL=52102
SERVER_PORT=52101
FRONTEND_PORT=52100
```

> 💡 You can omit most fields if you use the Telegram Mini App for configuration.

---

## 🐳 Docker (Production)

Prebuilt Docker images are available via GitHub Container Registry:

- Scanner: `ghcr.io/rattle-bot/rattle-scanner:latest`
- Server: `ghcr.io/rattle-bot/rattle-server:latest`
- Frontend: `ghcr.io/rattle-bot/frontend:latest`

---

### 🔧 docker-compose.yml (full version)

Includes: log scanner + backend server + frontend UI + PostgreSQL

```bash
docker compose -f docker-compose.yml up
```

### 🪶 docker-compose.light.yml (light version)

Includes: log scanner + PostgreSQL (no UI)

```bash
docker compose -f docker-compose.light.yml up
```

---

## 🧪 Log Examples

Here’s what Rattle sends to Telegram:

### Startup

```text
🚀 Rattle started in prod mode
```

### Container Summary (excluding rattle)

```text
📊 14 active containers:

- e133aff529d6: redis_container
- 48ab19bdb619: postgres_container
- ...
```

### Shutdown

```text
🛑 Rattle is shutting down...
```

### Error Log

```text
❌ Error in container: telegram_bot_container

2025-06-14 07:10:38,247 - aiogram.dispatcher - ERROR - Failed to fetch updates - TelegramNetworkError: HTTP Client says - ServerDisconnectedError: Server disconnected

📦 ID: c467ef7bfaf3
Name: telegram_bot_container
Image: telegram_bot

2025-06-14T07:10:38.276Z
```

> ✨ Messages use Telegram-friendly formatting for clear display and easy copying.

## 📦 Tech Stack

- [Go](https://go.dev/) (scanner)
- [Fiber](https://gofiber.io/) (API backend)
- [PostgreSQL](https://www.postgresql.org/) (stores filters, access settings, rules)
- [Vue 3](https://vuejs.org/) + [Vite](https://vite.dev/) (frontend)
- [TMA](https://core.telegram.org/bots/webapps) (Mini App UI)

---

## 📝 License

MIT — see [LICENSE](./LICENSE)
