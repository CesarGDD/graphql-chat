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
	"math/rand"
)

var chatMsg []*pg.Message
var chat map[string]chan []*pg.Message

func init() {
	chat = map[string]chan []*pg.Message{}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (r *messageResolver) User(ctx context.Context, obj *pg.Message) (*pg.User, error) {
	res, err := db.GetUserById(context.Background(), int32(obj.UserID))
	if err != nil {
		fmt.Println("Error getting user", err)
	}
	return &pg.User{
		ID:       res.ID,
		Username: res.Username,
	}, nil
}

func (r *mutationResolver) NewMessage(ctx context.Context, input pg.NewMessageInput) (*pg.Message, error) {
	user := userAuth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	res, err := db.CreateMessage(context.Background(), pg.CreateMessageParams{
		UserID:  int32(input.UserID),
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
	res, err := db.CreateUser(context.Background(), pg.CreateUserParams{
		Username: input.Username,
		Password: hashedPassword,
	})
	if err != nil {
		fmt.Println(err)
	}
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
	res, err := db.GetIdUserByUsername(context.Background(), input.Username)
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

func (r *queryResolver) User(ctx context.Context, id int) (*pg.User, error) {
	// inputId, err := strconv.Atoi(id)
	res, err := db.GetUserById(context.Background(), int32(id))
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
	res, err := db.ListMessages(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return res, err
}

func (r *subscriptionResolver) Messages(ctx context.Context) (<-chan []*pg.Message, error) {
	id := randString(8)
	events := make(chan []*pg.Message, 1)
	user := userAuth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	fmt.Println("==========================================================================================================================", user)

	go func() {
		<-context.Background().Done()
		delete(chat, id)
	}()

	chat[id] = events

	return events, nil
}

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type messageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
