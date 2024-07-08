-- NO foreign key constraints have been implemented for the init of the db.
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS user_activity;
DROP TABLE IF EXISTS comment_replies;
DROP TABLE IF EXISTS user_reply_activity;

--create a users table
CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    email   TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TEXT NOT NULL,
    modified_at TEXT NOT NULL,
    active INTEGER NOT NULL, 
    uuid TEXT NOT NULL
);
--populate example data to the users table
INSERT INTO users (username,email,password, role,created_at,modified_at,active,uuid) VALUES
('hannes','a@a.a','a','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1,'e1672398-a996-4d1c-a9ac-5d3cb5d7edee'),
('bob','bob','b','U',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,0,'5ff0d1da-2ba9-4c3f-881f-4e27ddb418ff'),
('paul','paul','1','SU',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1,'dbaea0b7-5ae9-4ca4-8889-054a1dba33da'),  -- could use integer and have it hardcoded in the config for the rolls
('fred','fred','2','A',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1,'c421eadd-c44e-4512-9ae7-b9865a036a61');
--create the posts table
CREATE TABLE IF NOT EXISTS posts(
    id INTEGER PRIMARY  KEY,
    title TEXT NOT NULL,
    body   TEXT NOT NULL,
    likes INTEGER NOT NULL, 
    dislikes INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    modified_at TEXT NOT NULL,
    active INTEGER NOT NULL -- this was originally status
);
--populate example data to the posts table 
INSERT INTO posts(title,body,likes,dislikes,user_id,category_id,created_at,modified_at,active)VALUES
('POST 1','Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',1,1,1,1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('POST 2','Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',10,0,1,3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('POST 3','Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',74,0,3,2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('POST 4','Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',0,0,2,4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,0),
('POST 5','Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',99,100,1,3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('POST 6','Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',100,-0,1,2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1);
-- creates the comments table
CREATE TABLE IF NOT EXISTS  comments(
    id INTEGER PRIMARY  KEY,
    body   TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    likes INTEGER NOT NULL, -- CHECK how the encryptions works
    dislikes INTEGER NOT NULL,
    post_id INTEGER NOT NULL, -- dunno if we will need this.
    created_at TEXT NOT NULL,
    modified_at TEXT NOT NULL,
    active INTEGER NOT NULL
);
-- populate example data to comments
INSERT INTO comments (body,user_id,likes,dislikes,post_id,created_at,modified_at,active) VALUES
('Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',1,0,0,1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',2,0,0,1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',1,0,0,2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1),
('Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',1,0,0,2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,0),
('Lorem ipsum dolor sit amet, consectetur adipiscing elit.. ',3,2,-10,3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1);
-- create the categories table 
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY  KEY,
    category   TEXT NOT NULL,
    active INTEGER NOT NULL,
    created_at TEXT NOT NULL
);
--populate example data to categories table 
INSERT INTO categories (category,active,created_at) VALUES
('Comedy',1,CURRENT_TIMESTAMP),
('Horror',1,CURRENT_TIMESTAMP),
('Sci-Fi',1,CURRENT_TIMESTAMP),
('Action',1,CURRENT_TIMESTAMP),
('Adventure',1,CURRENT_TIMESTAMP),
('Romance',1,CURRENT_TIMESTAMP),
('Miscellaneous',1,CURRENT_TIMESTAMP),
('Kids',1,CURRENT_TIMESTAMP);

-- user activity table table to manage likes
CREATE TABLE IF NOT EXISTS user_activity (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    like_activity BOOLEAN DEFAULT FALSE,
    dislike_activity BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, post_id)
);
-- create table to manage likes
CREATE TABLE IF NOT EXISTS user_reply_activity (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    comment_id INT NOT NULL,
    like_activity BOOLEAN DEFAULT FALSE,
    dislike_activity BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, comment_id)
);
-- create table to manage replies to comments
CREATE TABLE IF NOT EXISTS comment_replies (
    id INTEGER PRIMARY  KEY,
    body   TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL, 
    post_id INTEGER NOT NULL, 
    created_at TEXT NOT NULL,
    active INTEGER NOT NULL
);
