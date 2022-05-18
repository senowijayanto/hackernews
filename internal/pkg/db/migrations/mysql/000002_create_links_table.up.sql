CREATE TABLE IF NOT EXISTS links(
    id INT NOT NULL UNIQUE AUTO_INCREMENT,
    title VARCHAR (255) ,
    address VARCHAR (255) ,
    userId INT ,
    FOREIGN KEY (userId) REFERENCES users(id) ,
    PRIMARY KEY (id)
)