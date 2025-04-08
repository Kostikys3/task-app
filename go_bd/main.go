package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Импортируем PostgreSQL-драйвер
	"log"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=postgres1 dbname=db_tasks sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка при подключении к БД:", err)
	}
	defer db.Close()
	err = db.Ping() // проверяем, можем ли установить соединение
	if err != nil {
		fmt.Println("Ошибка при установки соединения:", err)
	} else {
		fmt.Println("Успешное подключение к БД")
	}

	addTasks(db, 1, 1, "test1", "Test2")

}

// Удалять задачу по id.
func deleteTasks(db *sql.DB, tasks_id int) error {
	res := "DELETE FROM tasks WHERE id = $1;"

	_, err := db.Exec(res, tasks_id)

	if err != nil {
		fmt.Println("Ошибка при удалении строки: ", err)
	} else {
		fmt.Println("Текст успешно удалён!")
	}
	return nil
}

// Обновлять задачу по id
func updateTasks(db *sql.DB, tasks_id, author_id, assigned_id int, title, content string) error {
	res := "UPDATE tasks SET author_id = $2, assigned_id = $3, title = $4, content = $5 WHERE id = $1;"

	_, err := db.Exec(res, tasks_id, author_id, assigned_id, content, title)

	if err != nil {
		fmt.Println("Ошибка при редактировании текста:", err)
	} else {
		fmt.Println("Текст успешно добавлен!")
	}
	return nil
}

// Создавать новые задачи
func addTasks(db *sql.DB, author_id, assigned_id int, title string, content string) error {
	query := "INSERT INTO tasks( author_id, assigned_id, title, content) VALUES ($1, $2, $3, $4);"
	_, err := db.Exec(query, author_id, assigned_id, title, content)
	if err != nil {
		return fmt.Errorf("err3")
	}

	fmt.Println("Текст успешно добавлен!")
	return nil
}

// Получать список всех задач
func selectTasks(db *sql.DB, author_id, assigned_id int, title string, content string) error {
	rows, err := db.Query("SELECT tasks.id, tasks.author_id, tasks.assigned_id, tasks.title , tasks.content FROM tasks")
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, author_id, assigned_id int
		var title, content string

		err = rows.Scan(&id, &author_id, &assigned_id, &title, &content)
		if err != nil {
			log.Println("Ошибка при сканировании строки", err)
			continue
		}

		fmt.Printf("id: %d author_id: %d assigned_id: %d title: %s content: %s\n",
			id, author_id, assigned_id, title, content)

		if err = rows.Err(); err != nil {
			log.Fatal("Ошибка при обходе строк:", err)

		}

	}

	return nil
}

// Получать список задач по автору
func selectAuthorTasks(db *sql.DB, author_id int) error {
	query := `SELECT users.id, users.name, tasks.id, tasks.assigned_id, tasks.title, tasks.content FROM users JOIN tasks ON users.id = tasks.assigned_id WHERE users.id = $1`

	rows, err := db.Query(query, author_id)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var users_id, tasks_id, assigned_id int
		var name, title, content string

		err = rows.Scan(&users_id, &name, &tasks_id, &assigned_id, &title, &content)
		if err != nil {
			log.Println("Ошибка при сканировании строки", err)
			continue
		}

		fmt.Printf("users id: %d name: %s tasks id: %d assigned_id %d title: %s content, %s\n",
			users_id, name, tasks_id, assigned_id, title, content)

		if err = rows.Err(); err != nil {
			log.Fatal("Ошибка при обходе строк:", err)

		}
	}
	return nil
}

// Получать список задач по метке
func selectLabelsTasks(db *sql.DB, label string) error {
	query := `SELECT t.id, t.author_id, t.assigned_id, t.title , t.content 
FROM tasks t
JOIN tasks_labels tl ON t.id = tl.task_id
JOIN labels l ON tl.label_id = l.id
WHERE l.name = $1`

	rows, err := db.Query(query, label)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tasks_id, assigned_id, author_id int
		var title, content string

		err = rows.Scan(&tasks_id, &author_id, &assigned_id, &title, &content)
		if err != nil {
			log.Println("Ошибка при сканировании строки", err)
			continue
		}

		fmt.Printf("tasks id: %d author id: %d assigned id: %d title: %s content, %s\n",
			tasks_id, author_id, assigned_id, title, content)

		if err = rows.Err(); err != nil {
			log.Fatal("Ошибка при обходе строк:", err)

		}
	}
	return nil
}
