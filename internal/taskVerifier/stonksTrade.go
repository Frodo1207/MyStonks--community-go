package taskVerifier

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	solanaRPCURL = "https://api.mainnet-beta.solana.com"
)

// 定义所需的结构体
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SignatureInfo struct {
	Signature string `json:"signature"`
	Slot      int    `json:"slot"`
	Err       any    `json:"err"`
	Memo      string `json:"memo"`
	BlockTime *int64 `json:"blockTime"`
}

type Transaction struct {
	Meta        TransactionMeta `json:"meta"`
	Transaction struct {
		Message struct {
			AccountKeys []struct {
				Pubkey string `json:"pubkey"`
			} `json:"accountKeys"`
			Instructions []struct {
				Parsed struct {
					Info struct {
						Destination string `json:"destination"`
						Lamports    int    `json:"lamports"`
						Source      string `json:"source"`
						Authority   string `json:"authority"`
						Mint        string `json:"mint"`
						Amount      string `json:"amount"`
						TokenAmount struct {
							Amount   string  `json:"amount"`
							Decimals int     `json:"decimals"`
							UIAmount float64 `json:"uiAmount"`
						} `json:"tokenAmount"`
					} `json:"info"`
					Type string `json:"type"`
				} `json:"parsed"`
				Program string `json:"program"`
			} `json:"instructions"`
		} `json:"message"`
	} `json:"transaction"`
	BlockTime int64 `json:"blockTime"`
}

type TransactionMeta struct {
	PreTokenBalances  []TokenBalance `json:"preTokenBalances"`
	PostTokenBalances []TokenBalance `json:"postTokenBalances"`
	LogMessages       []string       `json:"logMessages"`
}

type TokenBalance struct {
	AccountIndex  int    `json:"accountIndex"`
	Mint          string `json:"mint"`
	Owner         string `json:"owner"`
	UITokenAmount struct {
		Amount   string  `json:"amount"`
		Decimals int     `json:"decimals"`
		UIAmount float64 `json:"uiAmount"`
	} `json:"uiTokenAmount"`
}

// TransactionDetail 用于返回的交易详情
type TransactionDetail struct {
	Signature   string    `json:"signature"`
	Time        time.Time `json:"time"`
	Type        string    `json:"type"` // "in" or "out"
	TokenMint   string    `json:"token_mint"`
	Amount      float64   `json:"amount"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	IsToken     bool      `json:"is_token"`
}

// GetTodayTransactions 获取从今天0点到现在的交易信息
func GetTodayTransactions(walletAddress string) ([]TransactionDetail, error) {
	// 获取今天0点时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 获取最近签名
	signatures, err := getRecentSignatures(walletAddress, 100) // 适当增加limit值以确保覆盖
	if err != nil {
		return nil, fmt.Errorf("获取交易签名失败: %v", err)
	}

	var transactions []TransactionDetail

	// 分析每笔交易
	for _, sig := range signatures {
		if sig.BlockTime == nil {
			continue
		}

		txTime := time.Unix(*sig.BlockTime, 0)
		if txTime.Before(startOfDay) {
			continue
		}

		tx, err := getTransaction(sig.Signature)
		if err != nil {
			log.Printf("获取交易详情失败: %v (签名: %s)", err, sig.Signature)
			continue
		}

		// 分析交易
		details := analyzeTransactionDetails(walletAddress, tx)
		for i := range details {
			details[i].Signature = sig.Signature
			details[i].Time = txTime
		}

		transactions = append(transactions, details...)
	}

	return transactions, nil
}

// analyzeTransactionDetails 分析交易详情
func analyzeTransactionDetails(walletAddress string, tx *Transaction) []TransactionDetail {
	var details []TransactionDetail

	// 检查SOL转账
	for _, instr := range tx.Transaction.Message.Instructions {
		if instr.Parsed.Type == "transfer" {
			// SOL转账
			amount := float64(instr.Parsed.Info.Lamports) / 1e9
			if instr.Parsed.Info.Source == walletAddress {
				// 转出SOL
				details = append(details, TransactionDetail{
					Type:        "out",
					TokenMint:   "SOL",
					Amount:      amount,
					FromAddress: walletAddress,
					ToAddress:   instr.Parsed.Info.Destination,
					IsToken:     false,
				})
			} else if instr.Parsed.Info.Destination == walletAddress {
				// 转入SOL
				details = append(details, TransactionDetail{
					Type:        "in",
					TokenMint:   "SOL",
					Amount:      amount,
					FromAddress: instr.Parsed.Info.Source,
					ToAddress:   walletAddress,
					IsToken:     false,
				})
			}
		} else if instr.Parsed.Type == "transferChecked" {
			// 代币转账
			amount := instr.Parsed.Info.TokenAmount.UIAmount
			if instr.Parsed.Info.Authority == walletAddress {
				// 转出代币
				details = append(details, TransactionDetail{
					Type:        "out",
					TokenMint:   instr.Parsed.Info.Mint,
					Amount:      amount,
					FromAddress: walletAddress,
					ToAddress:   "", // 代币转账的接收方需要从日志或其他地方获取
					IsToken:     true,
				})
			} else {
				// 转入代币
				details = append(details, TransactionDetail{
					Type:        "in",
					TokenMint:   instr.Parsed.Info.Mint,
					Amount:      amount,
					FromAddress: "", // 代币转账的发送方需要从日志或其他地方获取
					ToAddress:   walletAddress,
					IsToken:     true,
				})
			}
		}
	}

	// 检查代币余额变化
	for _, balance := range tx.Meta.PostTokenBalances {
		if balance.Owner == walletAddress {
			// 检查是否已经有对应的转账记录
			found := false
			for _, d := range details {
				if d.TokenMint == balance.Mint {
					found = true
					break
				}
			}

			if !found {
				// 如果没有找到对应的转账记录，可能是其他类型的代币交易
				details = append(details, TransactionDetail{
					Type:      "unknown",
					TokenMint: balance.Mint,
					Amount:    balance.UITokenAmount.UIAmount,
					IsToken:   true,
				})
			}
		}
	}

	return details
}

// 以下是辅助函数，保持不变
func getRecentSignatures(address string, limit int) ([]SignatureInfo, error) {
	params := []interface{}{
		address,
		map[string]interface{}{
			"limit": limit,
		},
	}

	request := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getSignaturesForAddress",
		Params:  params,
	}

	var response RPCResponse
	if err := sendRPCRequest(request, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: %d - %s", response.Error.Code, response.Error.Message)
	}

	var signatures []SignatureInfo
	if err := json.Unmarshal(response.Result, &signatures); err != nil {
		return nil, fmt.Errorf("解析签名失败: %v", err)
	}

	return signatures, nil
}

func getTransaction(signature string) (*Transaction, error) {
	params := []interface{}{
		signature,
		map[string]interface{}{
			"encoding": "jsonParsed",
		},
	}

	request := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getTransaction",
		Params:  params,
	}

	var response RPCResponse
	if err := sendRPCRequest(request, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: %d - %s", response.Error.Code, response.Error.Message)
	}

	var tx Transaction
	if err := json.Unmarshal(response.Result, &tx); err != nil {
		return nil, fmt.Errorf("解析交易失败: %v (原始数据: %s)", err, string(response.Result))
	}

	return &tx, nil
}

func sendRPCRequest(request RPCRequest, response interface{}) error {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(solanaRPCURL, "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status not OK: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
