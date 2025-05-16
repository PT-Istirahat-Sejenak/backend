# Donora

## 🔥 API

https://gsc.fahrulhehehe.my.id/api

## 📱 Overview

Donora is an app that connects blood donors and those in need of blood to help solve the blood shortage problem in Indonesia. Donora encourages people to donate blood to meet the blood supply needs in Indonesia

### ✨ Key Features

- **Authentication** : Provide a secure login and registration system to verify and protect user identity
- **Chat between Donors and Donor Seekers** : Real-time chat feature that allows donor seekers to communicate directly with blood donors
- **Blood Donor Request Form** : A form that allows users to submit specific blood donor requests based on blood type
- **Blood Donor Need Notification** : Provide quick notification to users regarding blood donation needs
- **Dora AI Chatbot** : Dora AI chatbot utilizes the Gemini API to quickly and accurately answer users questions regarding lood donation

## 🛠️ Tech Stack

- **Language** : Golang
- **AI Integration** : Gemini AI
- **Database** : PostgreSQL
- **Object Storage** : R2 Cloudflare
- **Notification** : Firebase Cloud Messaging
- **Message** : WebSocket

## 🚀 Getting Started

### Prerequisites

- Golang v1.24.2
- Docker and Docker Compose (for containerized deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/PT-Istirahat-Sejenak/backend.git
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Create a `.env` file based on `.env.example`:
   ```bash
   cp .env.example .env
   ```
4. Get Project ID and Private Key from Firebase Console, then copy to 
   ```
    .
    │   ├── infrastructure
    │   │   ├── broadcast
   ```
5. Start the development server:
   ```bash
   go run cmd/api/main.go
   ```
6. Open [http://localhost:8080](http://localhost:8080) in your browser


## 📂 Project Structure

```
┌──(fahrul㉿Fahrul)-[~/gdsc/gsc/backend]
└─$ tree
.
├── Dockerfile
├── README.md
├── cmd
│   └── api
│       └── main.go
├── configs
│   └── config.go
├── docker-compose.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.mod.bak
├── go.sum
├── go.sum.bak
├── internal
│   ├── delivery
│   │   └── http
│   │       ├── handler
│   │       │   ├── auth_handler.go
│   │       │   ├── chatbot_handler.go
│   │       │   ├── education_handler.go
│   │       │   ├── fcm_handler.go
│   │       │   ├── history_handler.go
│   │       │   ├── profile_handler.go
│   │       │   ├── reward_handle.go
│   │       │   ├── upload_evidence_handler.go
│   │       │   └── websocket_handler.go
│   │       ├── middleware
│   │       │   └── auth_middleware.go
│   │       └── routes
│   │           └── routes.go
│   ├── entity
│   │   ├── chat.go
│   │   ├── education.go
│   │   ├── fcm.go
│   │   ├── histories.go
│   │   ├── message.go
│   │   ├── reward.go
│   │   ├── token.go
│   │   ├── upload_evidence.go
│   │   └── user.go
│   ├── infrastructure
│   │   ├── broadcast
│   │   │   └── donora-f67f2-5c889d5acd0a.json
│   │   ├── database
│   │   │   └── postgres.go
│   │   ├── email
│   │   │   └── email_service.go
│   │   ├── oauth
│   │   │   └── google_oauth.go
│   │   └── storage
│   │       ├── local_storage.go
│   │       ├── s3_storage.go
│   │       └── storage.go
│   ├── repository
│   │   ├── interfaces.go
│   │   └── postgres
│   │       ├── education_repository.go
│   │       ├── history_repository.go
│   │       ├── message_repository.go
│   │       ├── reward_repository.go
│   │       ├── token_repository.go
│   │       ├── upload_evidence_repository.go
│   │       └── user_repository.go
│   └── usecase
│       ├── auth_usecase.go
│       ├── chatbot_usecase.go
│       ├── education_usecase.go
│       ├── fcm_usecase.go
│       ├── history_usecase.go
│       ├── interfaces.go
│       ├── message_usecase.go
│       ├── profile_usecase.go
│       ├── reward_usecase.go
│       └── upload_evidence.go
├── main
└── pkg
    ├── hash
    │   └── bcrypt.go
    ├── jwt
    │   └── jwt.go
    └── validator
        └── validator.go
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.