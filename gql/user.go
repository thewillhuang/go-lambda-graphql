package gql

import (
	"context"
	"errors"
	"go-lambda-graphql/config"
	"go-lambda-graphql/models"
	"go-lambda-graphql/services/auth"
	"strconv"
	"time"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

// User struct
type User struct {
	Entity
	Name  string
	Email string
}

// UserResolver struct
type UserResolver struct {
	V *User
	U *User
}

// ID struct
type ID struct {
	ID string
}

// Signup mutation
func (r *Resolver) Signup(ctx context.Context, args struct {
	Email    string
	Name     string
	Password string
}) (*UserResolver, error) {
	tx, error := boil.Begin()
	if error != nil {
		return nil, error
	}
	hasEmail, _ := models.Usrs(tx, Where("email = ?", args.Email)).Exists()
	if hasEmail {
		return nil, errors.New("email taken")
	}
	err := validation.ValidateStruct(&args,
		validation.Field(&args.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&args.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&args.Password, validation.Required, validation.Length(5, 0)),
	)
	if err != nil {
		return nil, err
	}
	hash, _ := auth.HashPassword(args.Password)
	var newUser models.Usr
	newUser.Name = args.Name
	newUser.Email = args.Email
	newUser.PasswordHash = hash
	insertErr := newUser.Insert(tx)
	if insertErr != nil {
		tx.Rollback()
		return nil, insertErr
	}
	tx.Commit()
	usr := &User{
		Entity: Entity{
			ID:      relay.MarshalID("usr", ID{strconv.FormatInt(newUser.ID, 10)}),
			Created: graphql.Time{Time: newUser.CreatedAt},
			Updated: graphql.Time{Time: newUser.UpdatedAt},
		},
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	return &UserResolver{
		U: usr,
		V: usr,
	}, nil
}

// Jwt query
func (r *Resolver) Jwt(ctx context.Context, args struct {
	Email    string
	Password string
}) (*string, error) {
	usr, err := models.UsrsG(Where("email = ?", args.Email)).One()
	if err != nil {
		return nil, errors.New("wrong email or password combination")
	}
	validPassword := auth.CheckPasswordHash(args.Password, usr.PasswordHash)
	if !validPassword {
		return nil, errors.New("wrong email or password combination")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      usr.ID,
		"email":   usr.Email,
		"created": usr.CreatedAt,
		"updated": usr.UpdatedAt,
		"name":    usr.Name,
		"nbf":     time.Date(2017, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	return &tokenString, nil
}

// Viewer field
func (r *Resolver) Viewer(ctx context.Context, args struct {
	Jwt string
}) (*UserResolver, error) {
	token, err := auth.GetToken(args.Jwt)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		UserID := claims["id"].(float64)
		created, _ := time.Parse(time.RFC3339, claims["created"].(string))
		updated, _ := time.Parse(time.RFC3339, claims["updated"].(string))
		email := claims["email"].(string)
		name := claims["name"].(string)
		usr := &User{
			Entity: Entity{
				ID:      relay.MarshalID("usr", ID{strconv.FormatInt(int64(UserID), 10)}),
				Created: graphql.Time{Time: created},
				Updated: graphql.Time{Time: updated},
			},
			Name:  name,
			Email: email,
		}

		return &UserResolver{
			U: usr,
			V: usr,
		}, nil
	}

	return nil, err
}

// UpdateUser mutation
func (r *Resolver) UpdateUser(ctx context.Context, args struct {
	Email    *string
	Name     *string
	Password *string
	Jwt      string
}) (*UserResolver, error) {
	token, err := auth.GetToken(args.Jwt)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tx, error := boil.Begin()
		if error != nil {
			return nil, error
		}
		var dbOverrides []string
		id := int64(claims["id"].(float64))
		updatedUser, _ := models.FindUsr(tx, id)
		updatedUser.ID = id
		err := validation.ValidateStruct(&args,
			validation.Field(&args.Email, validation.Length(5, 50), is.Email),
			validation.Field(&args.Name, validation.Length(5, 50)),
			validation.Field(&args.Password, validation.Length(5, 0)),
		)
		if err != nil {
			return nil, err
		}
		if args.Name != nil {
			dbOverrides = append(dbOverrides, "name")
			updatedUser.Name = *args.Name
		}
		if args.Email != nil {
			dbOverrides = append(dbOverrides, "email")
			updatedUser.Email = *args.Email
		}
		if args.Password != nil {
			dbOverrides = append(dbOverrides, "password_hash")
			hash, _ := auth.HashPassword(*args.Password)
			updatedUser.PasswordHash = hash
		}

		usr := &User{
			Entity: Entity{
				ID:      relay.MarshalID("usr", ID{strconv.FormatInt(updatedUser.ID, 10)}),
				Created: graphql.Time{Time: updatedUser.CreatedAt},
				Updated: graphql.Time{Time: updatedUser.UpdatedAt},
			},
			Name:  updatedUser.Name,
			Email: updatedUser.Email,
		}
		dbError := updatedUser.Upsert(tx, true, []string{"id"}, dbOverrides)
		if dbError != nil {
			tx.Rollback()
		}
		tx.Commit()
		return &UserResolver{
			V: usr,
			U: usr,
		}, nil
	}
	return nil, err
}

// ID returns the id from User resolver
func (r *UserResolver) ID(ctx context.Context) (graphql.ID, error) {
	return r.U.ID, nil
}

// Created returns the Created from User resolver
func (r *UserResolver) Created(ctx context.Context) (graphql.Time, error) {
	return r.U.Created, nil
}

// Updated returns the Updated from User resolver
func (r *UserResolver) Updated(ctx context.Context) (graphql.Time, error) {
	return r.U.Updated, nil
}

// Name returns the Name from User resolver
func (r *UserResolver) Name(ctx context.Context) (string, error) {
	return r.U.Name, nil
}

// Email returns the Email from User resolver
func (r *UserResolver) Email(ctx context.Context) (string, error) {
	return r.U.Email, nil
}

// // TrendingConnection field represents a campaign connection
// func (r *UserResolver) TrendingConnection(ctx context.Context, args connectionArgs) (*campaignConnectionResolver, error) {

// }

// // TrendingConnection field represents a campaign connection
// func (r *UserResolver) FriendsConnection(ctx context.Context, args connectionArgs) (*UserConnectionResolver, error) {

// }

// // TrendingConnection field represents a campaign connection
// func (r *UserResolver) SearchConnection(ctx context.Context, args searchConnectionArgs) (*itemConnectionResolver, error) {

// }

// // TrendingConnection field represents a campaign connection
// func (r *UserResolver) ScanConnection(ctx context.Context, args scanConnectionArgs) (*itemConnectionResolver, error) {

// }

// type searchConnectionArgs struct {
// 	connectionArgs
// 	input string
// }
