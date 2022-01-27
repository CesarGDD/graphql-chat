package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"cesargdd/graph-subcriptions/graph/generated"
	"cesargdd/graph-subcriptions/jwt"
	"cesargdd/graph-subcriptions/pg"
	"cesargdd/graph-subcriptions/userAuth"
	"context"
	"errors"
	"fmt"
	"strconv"
)

func (r *mutationResolver) NewMessage(ctx context.Context, input pg.NewMessageInput) (*pg.Message, error) {
	user := userAuth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	id, _ := strconv.Atoi(input.UserID)
	res, err := db.CreateMessage(ctx, pg.CreateMessageParams{
		UserID:  int32(id),
		Content: input.Content,
	})
	if err != nil {
		fmt.Println(err)
	}
	msg := &pg.Message{
		ID:      res.ID,
		UserID:  res.UserID,
		Content: res.Content,
	}
	chatMsg = append(chatMsg, msg)

	for _, observer := range chat {
		observer <- chatMsg
	}
	return msg, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input pg.RegisterInput) (*pg.AuthResponse, error) {
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		fmt.Println(err)
	}
	res, err := db.CreateUser(ctx, pg.CreateUserParams{
		Username: input.Username,
		Password: hashedPassword,
	})
	token, err := jwt.GenerateToken(res.Username)
	if err != nil {
		fmt.Println(err)
	}
	return &pg.AuthResponse{
		AuthToken: &pg.AuthToken{
			AccessToken: token,
		},
		User: &pg.User{
			ID:       res.ID,
			Username: res.Username,
			Password: res.Password,
		},
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input pg.LoginInput) (*pg.AuthResponse, error) {
	res, err := db.GetIdUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, errors.New("invalid username")
	}

	if !CheckPasswordHash(input.Password, res.Password) {
		return nil, errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(res.Username)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &pg.AuthResponse{
		AuthToken: &pg.AuthToken{
			AccessToken: token,
		},
		User: &res,
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*pg.User, error) {
	inputId, err := strconv.Atoi(id)
	res, err := db.GetUserById(ctx, int32(inputId))
	if err != nil {
		fmt.Println(err)
	}
	return &pg.User{
		ID:       res.ID,
		Username: res.Username,
	}, nil
}

func (r *queryResolver) Messages(ctx context.Context) ([]pg.Message, error) {
	user := userAuth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	res, err := db.ListMessages(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return res, err
}

func (r *subscriptionResolver) Messages(ctx context.Context) (<-chan []*pg.Message, error) {
	user := userAuth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	id := randString(8)
	events := make(chan []*pg.Message, 1)

	go func() {
		<-ctx.Done()
		delete(chat, id)
	}()

	chat[id] = events

	return events, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
