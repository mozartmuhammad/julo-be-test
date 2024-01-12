package repository

const (
	updateWalletBalanceQuery = `UPDATE wallets
		SET
			balance = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			id = ? AND
			balance = ?`
)
