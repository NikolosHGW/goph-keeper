package main

import (
	"net/http"
	"os"
	"testing"
	"time"
)

// Для теста необходима развёрнутая бд с данными
// из дефолтного значения флага -d в конфиге.
func TestRunIntegration(t *testing.T) {
	done := make(chan struct{})

	go func() {
		if err := run(); err != nil {
			t.Errorf("run() вернула ошибку: %v", err)
		}
		close(done)
	}()

	time.Sleep(4 * time.Second)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatalf("Не удалось отправить HTTP-запрос: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получили: %d", resp.StatusCode)
	}

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Не удалось найти процесс: %v", err)
	}
	if err := process.Signal(os.Interrupt); err != nil {
		t.Fatalf("Не удалось отправить сигнал завершения: %v", err)
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Сервер не завершил работу в отведенное время")
	}
}
