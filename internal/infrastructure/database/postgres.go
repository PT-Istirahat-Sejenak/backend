package database

import (
	"backend/configs"
	"database/sql"
	"fmt"
)

func NewPostgresConnection(config *configs.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SQL schema for database tables
const CreateTablesSQL = `
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blood_type') THEN
    CREATE TYPE blood_type AS ENUM ('A', 'B', 'AB', 'O');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'rhesus') THEN
    CREATE TYPE rhesus AS ENUM ('positive', 'negative');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
    CREATE TYPE user_role AS ENUM ('pendonor', 'pencari', 'admin');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender') THEN
  	CREATE TYPE gender AS ENUM ('male', 'female');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
	role user_role NOT NULL,
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    profile_photo VARCHAR(255),
    phone_number VARCHAR(20) NOT NULL,
    gender gender NOT NULL,
    address TEXT NOT NULL,
    blood_type VARCHAR(2),
    rhesus VARCHAR(8),
    google_id VARCHAR(255),
	total_donation INT DEFAULT 0,
	coin INT DEFAULT 0,
	fcm_token VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS users_google_id_key ON users (google_id) WHERE google_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);

CREATE TABLE IF NOT EXISTS educations (
	id SERIAL PRIMARY KEY,
	image VARCHAR(255),
	title VARCHAR(255),
	content TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    content TEXT NOT NULL,
    is_delivered BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_delivery_status ON messages(receiver_id, is_delivered);

CREATE TABLE IF NOT EXISTS blood_requests (
	id SERIAL PRIMARY KEY,
	search_name VARCHAR(255) NOT NULL,
	location VARCHAR(255) NOT NULL,
	blood_type blood_type NOT NULL,
	rhesus rhesus NOT NULL,
	total INT NOT NULL,
	available BOOLEAN NOT NULL,
	urgency INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS histories (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	blood_request_id INT NOT NULL,
	next_donation DATE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (blood_request_id) REFERENCES blood_requests(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_histories_user_next_donation ON histories(user_id, next_donation DESC);
CREATE INDEX IF NOT EXISTS idx_histories_user_created_at_desc ON histories(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS images (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	image VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rewards (
	id SERIAL PRIMARY KEY,
	amount INT NOT NULL UNIQUE,
	price INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

INSERT INTO rewards (amount, price) 
VALUES 
  (10000, 20), 
  (20000, 30), 
  (50000, 50), 
  (100000, 75)
ON CONFLICT (amount) DO NOTHING;
`

func InitDatabase(db *sql.DB) error {
	_, err := db.Exec(CreateTablesSQL)
	return err
}
