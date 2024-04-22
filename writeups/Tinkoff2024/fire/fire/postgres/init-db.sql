CREATE TABLE IF NOT EXISTS Sacred (
    artefact TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Users (
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    age INT,
    PRIMARY KEY (username),
    UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS Messages (
    title VARCHAR(255) NOT NULL,
    body VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    FOREIGN KEY (username) REFERENCES Users(username) ON DELETE SET NULL
);

INSERT INTO Sacred (artefact) VALUES ('⚡️'), ('${ARTIFACT_2}'), ('${ARTIFACT_3}');
INSERT INTO Users (username, password, age) VALUES ('Prometheus', '${PROMETHEUS_PASSWORD}', 22);
INSERT INTO Messages (title, body, username) VALUES ('NEW DROP 💨💨💨', 'Have you seen this new Theogony by Hesiod? So fresh and hot 💨💨💨', 'Prometheus');
INSERT INTO Messages (title, body, username) VALUES ('The New Hope', 'Oh, ship, here we go again. They hid fire. Once again. I will bring it back to us!', 'Prometheus');