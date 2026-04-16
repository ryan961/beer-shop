package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type Card struct {
	Id      int64
	CardNo  string
	CCV     string
	Expires string
	Name    string
}

type CardRepo interface {
	CreateCard(ctx context.Context, c *Card) (*Card, error)
	GetCard(ctx context.Context, id int64) (*Card, error)
	ListCard(ctx context.Context, id int64) ([]*Card, error)
	DeleteCard(ctx context.Context, id int64) error
}

type CardUseCase struct {
	repo CardRepo
	log  *log.Helper
}

func NewCardUseCase(i do.Injector) (*CardUseCase, error) {
	repo := do.MustInvoke[CardRepo](i)
	logger := do.MustInvoke[log.Logger](i)
	return &CardUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/card"))}, nil
}

func (uc *CardUseCase) Create(ctx context.Context, c *Card) (*Card, error) {
	return uc.repo.CreateCard(ctx, c)
}

func (uc *CardUseCase) Get(ctx context.Context, id int64) (*Card, error) {
	return uc.repo.GetCard(ctx, id)
}

func (uc *CardUseCase) List(ctx context.Context, uid int64) ([]*Card, error) {
	return uc.repo.ListCard(ctx, uid)
}

func (uc *CardUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteCard(ctx, id)
}
