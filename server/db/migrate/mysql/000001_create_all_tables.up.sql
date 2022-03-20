-- This table is for registering clients.
CREATE TABLE Clients(client_id VARCHAR(255) NOT NULL PRIMARY KEY, client_secret VARCHAR(255) NOT NULL, redirect_uri VARCHAR(255) NOT NULL);

-- CREATE TABLE Users(ID INT NOT NULL UNIQUE AUTO_INCREMENT, PRIMARY KEY(ID));

-- This table is for storing temporary codes.
CREATE TABLE Codes (client_id VARCHAR(255) NOT NULL, code VARCHAR(255) NOT NULL, PRIMARY KEY (client_id, code), FOREIGN KEY (client_id) REFERENCES Clients(client_id));

-- This table is for storing valid access tokens.
CREATE TABLE Tokens (client_id VARCHAR(255) NOT NULL, token VARCHAR(255) NOT NULL, expr TIMESTAMP NOT NULL, FOREIGN KEY (client_id) REFERENCES Clients(client_id), PRIMARY KEY (client_id, token));