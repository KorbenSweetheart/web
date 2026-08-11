package service

// /connections
// `SELECT to_user_id FROM connections WHERE from_user_id = $1 AND status = 'accepted'`;

// Проблема: Сейчас ничто на уровне БД не запрещает создать несколько чатов между UserA и UserB.
// Решение: Нужен составной UNIQUE(user_one_id, user_two_id). Кроме того, чтобы избежать создания дублей вида (A, B) и (B, A),

// Отсутствие Composite Index для сортировки сообщений (messages)
// Проблема: Поле chat_id проиндексировано, но выборка сообщений чата почти всегда требует пагинации и сортировки по времени: WHERE chat_id = X ORDER BY created_at DESC. Single-column индекс на chat_id заставит БД выполнять дополнительную сортировку (Sort Node в EXPLAIN ANALYZE).
// Решение: Добавить составной индекс gorm:"index:idx_chat_messages,priority:1" на chat_id и gorm:"index:idx_chat_messages,priority:2" на id или created_at.
