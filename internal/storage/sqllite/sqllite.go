package sqllite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sachin141097/students-api/internal/config"
	"github.com/sachin141097/students-api/internal/types"
)

type Sqllite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqllite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqllite{
		Db: db,
	}, nil

}
func (s *Sqllite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age) VALUES(?,?,?)")
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil

}
func (s *Sqllite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age from students WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil

}
func (s *Sqllite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
func (s *Sqllite) UpdateStudent(id int64, name string, email string, age int) error {
	query := `UPDATE students SET name=?,email=?,age=? WHERE id=?`
	_, err := s.Db.Exec(query, name, email, age, id)
	if err != nil {
		return fmt.Errorf("error updating student: %w", err)
	}
	return nil
}
func (s *Sqllite) DeleteStudent(id int64) error {
	query := `DELETE FROM students WHERE id=?`
	_, err := s.Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting student: %w", err)
	}
	return nil
}
