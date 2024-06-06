package service

import (
	"github.com/azicussdu/go-todo-app/domain"
	"github.com/azicussdu/go-todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (t *TodoListService) Create(userId int, list domain.TodoList) (int, error) {
	return t.repo.Create(userId, list)
}

func (t *TodoListService) GetAll(userId int) ([]domain.TodoList, error) {
	return t.repo.GetAll(userId)
}

func (t *TodoListService) GetById(userId, listId int) (domain.TodoList, error) {
	return t.repo.GetById(userId, listId)
}

func (t *TodoListService) Delete(userId, listId int) error {
	return t.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, input domain.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, listId, input)
}
