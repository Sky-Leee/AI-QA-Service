package job

func checkIfExit() bool {
	select {
	case <-done:
		return true
	case <-semWorker:
	}
	return false
}

// 标记为睡眠
func markAsExit() {
	semWorker <- 1
}
