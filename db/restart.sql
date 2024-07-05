
-- the same sql is also added to the init file to purge the DB before building again.
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS post_comments;
DROP TABLE IF EXISTS user_activity;
DROP TABLE IF EXISTS comment_replies;
DROP TABLE IF EXISTS user_reply_activity;