CREATE TABLE `User` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    firstName VARCHAR(100) NOT NULL,
    lastName  VARCHAR(100) NOT NULL,
    `password`  VARCHAR(255) NOT NULL,
    email      VARCHAR(100) NOT NULL,
    token      VARCHAR(255) NOT NULL,
    refershToken VARCHAR(255) NOT NULL,
    createdAt  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updateAt   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    userCart JSON DEFAULT NULL COMMENT 'Store user Cart information',
    userAddress JSON DEFAULT NULL COMMENT 'Store user Address information',
    userOrder JSON DEFAULT NULL COMMENT 'Store user orders information'
);




