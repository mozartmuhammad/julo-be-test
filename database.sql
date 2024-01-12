CREATE DATABASE miniwallet;

USE miniwallet;

CREATE TABLE IF NOT EXISTS `wallets` (
    id VARCHAR(36) NOT NULL,
    customer_xid VARCHAR(36),
    status VARCHAR(20) DEFAULT 'disabled',
    enabled_at TIMESTAMP,
    balance INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE(`customer_xid`)
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `transactions` (
    id VARCHAR(36) NOT NULL,
    wallet_id VARCHAR(36),
    customer_xid VARCHAR(36),
    transaction_type ENUM('deposit', 'withdrawal'),
    amount INT NOT NULL,
    reference_id VARCHAR(75) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE(`transaction_type`, `reference_id`)
) ENGINE=INNODB;