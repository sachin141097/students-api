package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sachin141097/students-api/internal/storage"
	"github.com/sachin141097/students-api/internal/types"
	"github.com/sachin141097/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}
func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("updating a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return

		}
		var student types.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err := validator.New().Struct(student); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}
		err = storage.UpdateStudent(intId, student.Name, student.Email, student.Age)
		if err != nil {
			slog.Error("error updating student", slog.String("id", id), slog.Any("error", err))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		}
		slog.Info("student updated successfully", slog.String("id", id))
		response.WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"})

	}
}
func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("deleting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		err = storage.DeleteStudent(intId)
		if err != nil {
			slog.Error("error deleting a student", slog.String("id", id), slog.Any("error", err))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("student deleted successfully", slog.String("id", id))
		response.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}
