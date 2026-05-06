package main

import (
	"encoding/json" // import: 引入外部套件。此套件用於解析 JSON 格式資料
	"fmt"           // 用於格式化輸出與字串處理 (如 Printf, Sprintf)
	"net/http"      // 用於發送 HTTP 請求 (如 Get 請求)
	"sync"          // 用於併發控制 (如 WaitGroup 同步機制)
)

// RateResponse 模擬 API 回傳結構
type RateResponse struct {
	Source string  `json:"source"`
	Base   string  `json:"base"`
	Target string  `json:"target"`
	Rate   float64 `json:"rate"`
	Error  error   `json:"-"`
}

// FetchRate 模擬從不同 API 抓取匯率的函式。
// 參數:
// - source: 資料來源名稱
// - base: 基準貨幣 (如 USD)
// - target: 目標貨幣 (如 TWD)
// - ch: 用於傳回結果的唯寫通道 (chan<-)
// - wg: 用於同步任務狀態的 WaitGroup 指標
// func: 定義函式 (Function) 的關鍵字，後接函式名稱與參數
func FetchRate(source string, base string, target string, ch chan<- RateResponse, wg *sync.WaitGroup) {
	// defer: 延遲執行。當函式執行結束前（不論成功或失敗），會自動執行 defer 後的指令。
	// 這裡用於確保 WaitGroup 計數減一，標記此併發任務已完成。
	defer wg.Done()

	// 使用公開的 ExchangeRate-API (不需要 API Key)
	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", base)

	resp, err := http.Get(url)
	if err != nil {
		// <-: 通道 (Channel) 發送操作。將結果「送入」 ch 通道中。
		ch <- RateResponse{Source: source, Error: fmt.Errorf("網路請求失敗: %w", err)}
		return
	}
	// 再次使用 defer 確保 http 回應的 Body 會在函式結束時關閉，釋放網路資源。
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- RateResponse{Source: source, Error: fmt.Errorf("API 回傳錯誤狀態碼: %d", resp.StatusCode)}
		return
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		ch <- RateResponse{Source: source, Error: fmt.Errorf("解析 JSON 失敗: %w", err)}
		return
	}

	rate, ok := result.Rates[target]
	if !ok {
		ch <- RateResponse{Source: source, Error: fmt.Errorf("找不到目標貨幣 %s", target)}
		return
	}

	ch <- RateResponse{
		Source: source,
		Base:   base,
		Target: target,
		Rate:   rate,
	}
}

func main() {
	sources := []string{"Bank_A", "Bank_B", "Bank_C"}
	baseCurrency := "USD"
	targetCurrency := "TWD"

	ch := make(chan RateResponse, len(sources))
	var wg sync.WaitGroup

	fmt.Printf("正在從 %d 個來源抓取 %s/%s 匯率...\n", len(sources), baseCurrency, targetCurrency)

	for _, source := range sources {
		wg.Add(1)
		go FetchRate(source, baseCurrency, targetCurrency, ch, &wg)
	}

	// 啟動一個 Goroutine 來關閉 channel，避免主執行緒阻塞
	go func() {
		// 等待所有 FetchRate 的任務都執行完畢
		wg.Wait()
		// close: 關閉通道，告訴接收端 (range ch) 資料已經傳送完畢
		close(ch)
	}()

	// 處理結果
	// for ... range ch: 從通道中「接收」資料的迴圈。
	// 它會不斷接收資料直到 ch 被 close(ch) 關閉為止。
	for res := range ch {
		if res.Error != nil {
			fmt.Printf("[%s] 錯誤: %v\n", res.Source, res.Error)
			continue
		}
		fmt.Printf("[%s] 匯率: %.4f %s/%s\n", res.Source, res.Rate, res.Base, res.Target)
	}

	fmt.Println("所有資料抓取完成。")
}
