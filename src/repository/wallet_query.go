package repository

const (
	updateWalletBalanceQuery = `UPDATE wallets
		SET
			balance = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			id = ? AND
			balance = ?`

	updateTransactionStatusQuery = `UPDATE transactions
		SET
			status = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			id = ?`

	insertTransactionQuery = `INSERT INTO transactions
		(id, wallet_id, customer_xid, transaction_type, amount, reference_id, status, created_at, updated_at)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`

	getTransactionsQuery = `SELECT 
		id, wallet_id, customer_xid, transaction_type, amount, reference_id, status, created_at, updated_at 
		FROM transactions WHERE wallet_id = ? order by created_at`

	getWalletQuery = `SELECT 	
		id, customer_xid, status, enabled_at, balance, created_at, updated_at FROM wallets 
		WHERE customer_xid = ?`

	updateWalletStatusQuery = `UPDATE wallets
		SET
			status = ?,
			enabled_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			customer_xid = ?`

	insertWalletQuery = `INSERT INTO wallets
		(id, customer_xid)
		VALUES(?, ?)`
)
