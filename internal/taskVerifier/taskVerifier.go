package taskVerifier

import (
	"fmt"
	"sync"
)

// TaskVerifier 任务验证器
type TaskVerifier struct {
}

var (
	instance *TaskVerifier
	once     sync.Once
)

func GetInstance() *TaskVerifier {
	once.Do(func() {
		instance = &TaskVerifier{}
	})
	return instance
}

func (tv *TaskVerifier) Verify(task interface{}) (bool, error) {

	if task == nil {
		return false, fmt.Errorf("task cannot be nil")
	}

	return true, nil
}

func (tv *TaskVerifier) VerifyStonksTradeTask(sol_addr string) (bool, error) {
	transactions, err := GetTodayTransactions(sol_addr)
	fmt.Printf("transactions: %v\n", transactions)
	for _, transaction := range transactions {
		fmt.Printf("transaction_mint: %v\n", transaction.TokenMint)
		fmt.Printf("in or out of: %v\n", transaction.Type)
		fmt.Printf("amount: %v\n", transaction.Amount)
		if transaction.Type == "in" && transaction.TokenMint == "6NcdiK8B5KK2DzKvzvCfqi8EHaEqu48fyEzC8Mm9pump" {
			return true, nil
		}
	}
	if err != nil {
		return false, err
	}
	return false, nil
}
