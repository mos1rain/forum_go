#!/bin/bash

# Удаляем старую базу данных если она существует
rm -f forum.db

# Создаем новую базу данных
sqlite3 forum.db << EOF
.read migrations/000001_init_auth.up.sql
.read migrations/000002_forum_schema.up.sql
EOF

echo "База данных успешно инициализирована!" 