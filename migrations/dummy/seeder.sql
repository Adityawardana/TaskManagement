INSERT INTO public.users (id,email,password_hash,deleted_at,created_at,updated_at) VALUES
	 ('019f8a27-9cb2-7af0-9313-1598b35f4349','user1@gmail.com','$2a$10$fa8aelMQd39bJtXXDZznU.rqRU6cvzy1VIYXI9ZCPVZUzGIItY9Eq',NULL,'2026-07-22 14:07:59.689409','2026-07-22 14:07:59.689409'),
	 ('019f8a27-d345-73c8-b55e-2cfc6aeee19a','user2@gmail.com','$2a$10$8zkXUxMRHlLKVgVq1fCFn.xonn5tWxHQh2XCvqzwfHj3ZsE1wSUgq',NULL,'2026-07-22 14:08:13.641004','2026-07-22 14:08:13.641004'),
	 ('019f8a27-e73d-7e01-b8ba-9a512743d471','user3@gmail.com','$2a$10$tEdrPpnjQ6G0GcOAW19qfesjb6Ob4JewqkTGvYx3xqbG.ORZi2xn.',NULL,'2026-07-22 14:08:18.762491','2026-07-22 14:08:18.762491'),
	 ('019f8b01-3524-72d2-99b2-63dc164aeb5d','user4@gmail.com','$2a$10$NRnnPN2nF4MDk.BsGnphw.a36FTCit4DaSb2qh9aVvOx0JDuxOabe',NULL,'2026-07-22 18:05:40.021828','2026-07-22 18:05:40.021828'),
	 ('019f8b0a-e668-743a-a43d-5fdc3804ec69','user5@gmail.com','$2a$10$Wm0vyYK2OfB0OOlDjpib/e1s4uZz6Xc8Augpz.e4c1cMD39BeG6wW',NULL,'2026-07-22 18:16:15.22003','2026-07-22 18:16:15.22003'),
	 ('019f8b0b-7591-7080-b61a-26c4ae6d1600','user6@gmail.com','$2a$10$4h0d3J3XGR0xXOf1756HguQfOb02dk79RteqjRDwiK7XaM7sOKJma',NULL,'2026-07-22 18:16:51.871302','2026-07-22 18:16:51.871302');

INSERT INTO public.teams (id, "name",created_at) VALUES
	 (1, 'default','2026-07-22 14:07:59.689409'),
	 (2, 'Team2','2026-07-22 18:16:51.871302');

INSERT INTO public.team_members (team_id,user_id,created_at) VALUES
	 (1,'019f8a27-9cb2-7af0-9313-1598b35f4349','2026-07-22 14:07:59.689409'),
	 (1,'019f8a27-d345-73c8-b55e-2cfc6aeee19a','2026-07-22 14:08:13.641004'),
	 (1,'019f8a27-e73d-7e01-b8ba-9a512743d471','2026-07-22 14:08:18.762491'),
	 (1,'019f8b01-3524-72d2-99b2-63dc164aeb5d','2026-07-22 18:05:40.021828'),
	 (1,'019f8b0a-e668-743a-a43d-5fdc3804ec69','2026-07-22 18:16:15.22003'),
	 (2,'019f8b0b-7591-7080-b61a-26c4ae6d1600','2026-07-22 18:16:51.871302');

INSERT INTO public.tasks (id, owner_id,title,description,status,assigned_to,created_at,updated_at) VALUES
	 (3, '019f8a27-9cb2-7af0-9313-1598b35f4349','Task 3 - Created by User1','task 3 description','todo',NULL,'2026-07-22 15:07:24.044288','2026-07-22 15:07:24.044288'),
	 (5, '019f8a27-d345-73c8-b55e-2cfc6aeee19a','Task 5 - Created by User2','task 5 description','todo',NULL,'2026-07-22 15:08:28.504686','2026-07-22 15:08:28.504686'),
	 (6, '019f8a27-d345-73c8-b55e-2cfc6aeee19a','Task 6 - Created by User2','task 6 description','todo',NULL,'2026-07-22 15:08:36.805172','2026-07-22 15:08:36.805172'),
	 (7, '019f8a27-e73d-7e01-b8ba-9a512743d471','Task 7 - Created by User3','task 7 description','todo',NULL,'2026-07-22 15:10:15.76132','2026-07-22 15:10:15.76132'),
	 (2, '019f8a27-d345-73c8-b55e-2cfc6aeee19a','Task 2 - Created by User2','task 2 description','doing',NULL,'2026-07-22 14:10:51.233089','2026-07-22 15:16:29.678232'),
	 (4, '019f8a27-9cb2-7af0-9313-1598b35f4349','Task 4 - Created by User1 - Assigned To User2','task 4 description','doing','019f8a27-d345-73c8-b55e-2cfc6aeee19a','2026-07-22 15:07:32.278808','2026-07-22 17:47:49.257279'),
	 (8, '019f8a27-e73d-7e01-b8ba-9a512743d471','Task 8 - Created by User3','task 8 description updated','doing','019f8a27-9cb2-7af0-9313-1598b35f4349','2026-07-22 15:10:23.857308','2026-07-22 18:04:20.595756'),
	 (9, '019f8b0b-7591-7080-b61a-26c4ae6d1600','Task Team 2 - Created by User6','task team 2 description updated','todo','019f8b0b-7591-7080-b61a-26c4ae6d1600','2026-07-22 18:17:59.80475','2026-07-22 18:21:06.829746'),
	 (1, '019f8a27-9cb2-7af0-9313-1598b35f4349','Task 1 - Created by User1','task 1 description','doing','019f8a27-e73d-7e01-b8ba-9a512743d471','2026-07-22 14:09:07.391481','2026-07-22 18:31:05.866184');

INSERT INTO public.task_logs (task_id,"action",old_value,new_value,performed_by,created_at) VALUES
	 (1,'assign',null,'019f8a27-9cb2-7af0-9313-1598b35f4349','019f8a27-9cb2-7af0-9313-1598b35f4349','2026-07-22 14:55:04.181554'),
	 (4,'assign',null,'019f8a27-d345-73c8-b55e-2cfc6aeee19a','019f8a27-9cb2-7af0-9313-1598b35f4349','2026-07-22 17:05:16.140866'),
	 (8,'assign',null,'019f8a27-9cb2-7af0-9313-1598b35f4349','019f8a27-e73d-7e01-b8ba-9a512743d471','2026-07-22 18:04:20.595756'),
	 (9,'assign',null,'019f8b0b-7591-7080-b61a-26c4ae6d1600','019f8b0b-7591-7080-b61a-26c4ae6d1600','2026-07-22 18:21:06.829746'),
	 (1,'assign','019f8a27-9cb2-7af0-9313-1598b35f4349','019f8a27-d345-73c8-b55e-2cfc6aeee19a','019f8a27-9cb2-7af0-9313-1598b35f4349','2026-07-22 18:29:25.139069'),
	 (1,'assign','019f8a27-d345-73c8-b55e-2cfc6aeee19a','019f8a27-e73d-7e01-b8ba-9a512743d471','019f8a27-9cb2-7af0-9313-1598b35f4349','2026-07-22 18:31:05.866184');

SELECT setval('teams_id_seq', (SELECT MAX(id) FROM teams));
SELECT setval('tasks_id_seq', (SELECT MAX(id) FROM tasks));
SELECT setval('task_logs_id_seq', (SELECT MAX(id) FROM task_logs));