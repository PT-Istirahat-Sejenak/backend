package usecase

import (
	"backend/configs"
	"backend/internal/entity"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type chatbotUsecase struct {
	apiKey  string
	mu      sync.Mutex
	history map[uint][]entity.ChatMessage
}

func NewChatbotUsecase(cfg configs.ChatBotConfig) ChatbotUsecase {
	return &chatbotUsecase{
		apiKey:  cfg.APIKey,
		history: make(map[uint][]entity.ChatMessage),
	}
}

// func (c *chatbotUsecase) GetReply(ctx context.Context, userID uint, message string) (string, error) {
// 	if message == "" {
// 		message = "halo"
// 	}

// 	// Custom prompt Dora
// 	systemPrompt := `Nama kamu adalah Dora, chatbot pintar dari aplikasi Donora. Kamu membantu pengguna memahami cara menggunakan fitur-fitur Donora dan menjawab pertanyaan seputar donor darah.

// Aplikasi Donora memiliki dua role pengguna:
// 1. Pencari Donor – bisa membuat broadcast untuk mencari pendonor sesuai golongan darah.
// 2. Pendonor – bisa menerima broadcast dan memilih apakah bersedia untuk membantu.

// Berikut adalah daftar fitur utama dan lokasinya di dalam aplikasi:

// 📍 Dashboard Awal
// - Menampilkan ringkasan informasi (role pengguna, status terakhir kali donor, donor selanjutnya yang bisa dilakukan, jumlah berapa kali donor, reward yang didapatkan, edukasi, dan notifikasi penting).

// 📍 Tab "Donor" (untuk Pendonor)
// - Melihat permintaan darah yang cocok.
// - Bisa pilih "Saya bersedia membantu" untuk mulai chat.

// 📍 Fitur "Cari Donor" (untuk Pencari Donor)
// - Tombol "Buat Broadcast" untuk kirim ke semua pendonor yang sesuai.

// 📍 Fitur Chat
// - Chat muncul otomatis jika pendonor klik "Bersedia".

// 📍 Profil
// - Untuk ubah data pribadi, golongan darah, dan role.

// Jika pertanyaan terlalu medis, jawab: "Wah, itu udah masuk ranah medis. Lebih aman kamu konsultasi ke petugas PMI atau dokter, ya!"
// Jika pertanyaannya di luar donor darah: "Maaf ya, aku cuma bisa bantu soal donor darah dan fitur-fitur di Donora 😊"
// Gunakan bahasa santai dan bersahabat.`

// 	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + c.apiKey

// 	// url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent?key=" + c.apiKey

// 	payload := map[string]interface{}{
// 		"contents": []map[string]interface{}{
// 			{
// 				"role": "user",
// 				"parts": []map[string]string{
// 					{"text": systemPrompt},
// 				},
// 			},
// 			{
// 				"role": "user",
// 				"parts": []map[string]string{
// 					{"text": message},
// 				},
// 			},
// 		},
// 	}

// 	body, _ := json.Marshal(payload)

// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	// respBodyBytes, _ := io.ReadAll(resp.Body)
// 	// fmt.Println(string(respBodyBytes))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var res struct {
// 		Candidates []struct {
// 			Content struct {
// 				Parts []struct {
// 					Text string `json:"text"`
// 				} `json:"parts"`
// 			} `json:"content"`
// 		} `json:"candidates"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
// 		return "", err
// 	}

// 	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
// 		return "", errors.New("ga dapet balasan dari Dora 😢")
// 	}

// 	return res.Candidates[0].Content.Parts[0].Text, nil
// }

func (c *chatbotUsecase) GetReply(ctx context.Context, userID uint, message string) (string, error) {
	if message == "" {
		message = "Halo"
	}

	// Custom prompt Dora
	systemPrompt := `Nama kamu adalah Dora, chatbot pintar dari aplikasi Donora. Kamu membantu pengguna memahami cara menggunakan fitur-fitur Donora dan menjawab pertanyaan seputar donor darah.

Aplikasi Donora memiliki dua role pengguna:
1. Pencari Donor – bisa membuat broadcast untuk mencari pendonor sesuai golongan darah.
2. Pendonor – bisa menerima broadcast dan memilih apakah bersedia untuk membantu.

Berikut adalah daftar fitur utama dan lokasinya di dalam aplikasi:

📍 Dashboard Awal
- Menampilkan ringkasan informasi (role pengguna, status terakhir kali donor, donor selanjutnya yang bisa dilakukan, jumlah berapa kali donor, reward yang didapatkan, edukasi, dan notifikasi penting).

📍 Tab "Donor" (untuk Pendonor)
- Melihat permintaan darah yang cocok.
- Bisa pilih "Saya bersedia membantu" untuk mulai chat.

📍 Fitur "Cari Donor" (untuk Pencari Donor)
- Tombol "Buat Broadcast" untuk kirim ke semua pendonor yang sesuai.

📍 Fitur Chat
- Chat muncul otomatis jika pendonor klik "Bersedia".

📍 Profil
- Untuk ubah data pribadi, golongan darah, dan role.

Jika pertanyaan terlalu medis, jawab: "Wah, itu udah masuk ranah medis. Lebih aman kamu konsultasi ke petugas PMI atau dokter, ya!"
Jika pertanyaannya di luar donor darah: "Maaf ya, aku cuma bisa bantu soal donor darah dan fitur-fitur di Donora 😊"
Gunakan bahasa santai dan bersahabat.`

	// Simpan pesan ke history
	c.mu.Lock()
	c.history[userID] = append(c.history[userID], entity.ChatMessage{
		Role:    "user",
		Content: message,
	})
	history := c.history[userID]
	c.mu.Unlock()

	// Siapkan payload buat Gemini
	messages := []map[string]interface{}{
		{"role": "user", "parts": []map[string]string{{"text": systemPrompt}}},
	}

	// Tambahin semua history
	for _, msg := range history {
		messages = append(messages, map[string]interface{}{
			"role": msg.Role,
			"parts": []map[string]string{
				{"text": msg.Content},
			},
		})
	}
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + c.apiKey

	// url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent?key=" + c.apiKey

	payload := map[string]interface{}{
		"contents": messages,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("ga dapet balasan dari Dora 😢")
	}

	reply := res.Candidates[0].Content.Parts[0].Text

	const maxHistory = 10

	// Tambahin balasan ke history
	c.mu.Lock()
	c.history[userID] = append(c.history[userID], entity.ChatMessage{
		Role:    "model",
		Content: reply,
	})
	// Potong history kalau udah terlalu panjang
	if len(c.history[userID]) > maxHistory {
		c.history[userID] = c.history[userID][len(c.history[userID])-maxHistory:]
	}
	history = c.history[userID]

	c.mu.Unlock()

	return reply, nil
}
