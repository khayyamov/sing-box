CREATE DATABASE IF NOT EXISTS users_db
    CHARACTER SET utf8
    COLLATE utf8_general_ci;

USE users_db;


CREATE TABLE vless
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);


CREATE TABLE vmess
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);


CREATE TABLE tuic
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);


CREATE TABLE trojan
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE PRIMARY KEY
);


CREATE TABLE shadowtls
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);


CREATE TABLE shadowsocks_multi
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);

CREATE TABLE shadowsocks_relay
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);

CREATE TABLE naive
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);


CREATE TABLE hysteria2
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);


CREATE TABLE hysteria
(
    id           INT AUTO_INCREMENT     NOT NULL PRIMARY KEY,
    user_json     VARCHAR(255)            NOT NULL UNIQUE
);
