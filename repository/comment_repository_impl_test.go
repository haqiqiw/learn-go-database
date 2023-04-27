package repository

import (
	"context"
	"fmt"
	learn_go_database "learn-go-database"
	"learn-go-database/entity"

	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCommentInsert(t *testing.T) {
	commentRepository := NewCommentRepositoryImpl(learn_go_database.GetConnection())

	ctx := context.Background()
	comment := entity.Comment{
		Email:   "repository@test.comm",
		Comment: "Test repository",
	}

	result, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func TestFindById(t *testing.T) {
	commentRepository := NewCommentRepositoryImpl(learn_go_database.GetConnection())

	ctx := context.Background()
	comment, err := commentRepository.FindById(ctx, 35)
	// comment, err := commentRepository.FindById(ctx, 95) // not found
	if err != nil {
		panic(err)
	}

	fmt.Println(comment)
}

func TestFindAll(t *testing.T) {
	commentRepository := NewCommentRepositoryImpl(learn_go_database.GetConnection())

	ctx := context.Background()
	comments, err := commentRepository.FindAll(ctx)
	if err != nil {
		panic(err)
	}

	for _, comment := range comments {
		fmt.Println(comment)
	}
}
