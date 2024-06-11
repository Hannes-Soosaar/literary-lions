-- database: /home/hannes/literary-lions/db/lions.db

-- Use the ▷ button in the top right corner to run the entire file.

SELECT * FROM users;

UPDATE users SET active = 0 WHERE uuid = 'e1672398-a996-4d1c-a9ac-5d3cb5d7edee';
UPDATE users SET active = 1 WHERE uuid = 'e1672398-a996-4d1c-a9ac-5d3cb5d7edee';

SELECT id,username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE username = 'hannes';
SELECT id,username,email,password,role,created_at,modified_at,active,uuid FROM users  WHERE username = 'bob';