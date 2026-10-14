package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// ValidPriorities son los únicos valores de prioridad aceptados.
var ValidPriorities = []string{"urgente", "alta", "media", "baja"}

// priorityLabels mapea cada valor de prioridad a su etiqueta visible.
var priorityLabels = map[string]string{
	"urgente": "Urgente",
	"alta":    "Alta",
	"media":   "Media",
	"baja":    "Baja",
}

// Task representa una tarea de la lista.
type Task struct {
	ID        uuid.UUID `json:"id" db:"id" form:"-"`
	Title     string    `json:"title" db:"title" form:"title"`
	Priority  string    `json:"priority" db:"priority" form:"priority"`
	Completed bool      `json:"completed" db:"completed" form:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" form:"-"`
}

// Tasks es una colección de Task.
type Tasks []Task

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	if t.Priority == "" {
		t.Priority = "media"
	}
	if !isValidPriority(t.Priority) {
		verrs.Add("priority", "La prioridad no es válida.")
	}

	return verrs, nil
}

// isValidPriority indica si value pertenece al conjunto de prioridades permitidas.
func isValidPriority(value string) bool {
	for _, p := range ValidPriorities {
		if p == value {
			return true
		}
	}
	return false
}

// PriorityLabel devuelve la etiqueta visible de la prioridad de la tarea.
func (t Task) PriorityLabel() string {
	if label, ok := priorityLabels[t.Priority]; ok {
		return label
	}
	return t.Priority
}
