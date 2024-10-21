INSERT INTO apps (id, name, secret)
VALUES (1, 'myApp', 'mySecret')
ON CONFLICT DO NOTHING