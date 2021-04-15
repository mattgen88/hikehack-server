package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mattgen88/haljson"
	"gorm.io/gorm"
)

type UserDataKey string
type Role string

// Claims holds claims for a token
type Claims struct {
	Username string
	jwt.StandardClaims
}

// AuthMiddleware wraps something requiring auth in the form of a jwt
func AuthMiddleware(handler http.Handler, jwtKey string, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		root := haljson.NewResource()
		root.Self(r.URL.Path)

		var success bool

		// Snag JWT, verify, validate or redirect to auth endpoint
		cookie, err := r.Cookie("access.jwt")
		if err != nil {
			log.Println(err)
			success = false
		} else {
			success, ctx = validateToken(ctx, cookie, jwtKey)
		}

		if !success {
			log.Println("access.jwt failed to validate")
			// Try refresh.jwt
			cookie, err := r.Cookie("refresh.jwt")
			if err != nil {
				success = false
			} else {
				success, ctx = validateToken(ctx, cookie, jwtKey)
				if success {
					log.Println("refresh.jwt validated, updating access.jwt")
					if val := ctx.Value(UserDataKey("user_data")); val != nil {
						mapClaims := val.(jwt.MapClaims)
						if username, ok := mapClaims["Username"]; ok {
							now := time.Now()

							accessExpires := now.Add(time.Minute * 5)

							// Create the Claims
							accessClaims := Claims{
								username.(string),
								jwt.StandardClaims{
									ExpiresAt: accessExpires.Unix(),
									Issuer:    "test",
								},
							}

							ctx = context.WithValue(ctx, UserDataKey("user_data"), accessClaims)

							accessCookie, accessErr := CreateJwt("access.jwt", accessExpires, &accessClaims, jwtKey)
							if accessErr != nil {
								log.Println(accessErr)
								success = false
							} else {
								http.SetCookie(w, accessCookie)
							}
						} else {
							log.Println("No username")
							success = false
						}
					} else {
						log.Println("No user_data")
						success = false
					}
				} else {
					log.Println("refresh.jwt failed to validate")
				}
			}

		}

		if !success {
			root.Data["error"] = "Access denied"
			w.WriteHeader(http.StatusForbidden)
			json, err := json.Marshal(root)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(json)
			return
		}
		handler.ServeHTTP(w, r.WithContext(ctx))

	})
}

func validateToken(ctx context.Context, cookie *http.Cookie, jwtKey string) (bool, context.Context) {
	success := true

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		log.Println(err)
		success = false
	}

	if !token.Valid {
		log.Println("not valid")
		success = false
	}

	if _, ok := err.(*jwt.ValidationError); ok {
		log.Println("validation error")
		success = false
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		log.Println("bad claims")
		success = false
	}

	ctx = context.WithValue(ctx, UserDataKey("user_data"), token.Claims.(jwt.MapClaims))

	return success, ctx
}

func CreateJwt(name string, expire time.Time, claims jwt.Claims, jwtKey string) (*http.Cookie, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		Expires:  expire,
		SameSite: http.SameSiteNoneMode,
	}

	return &cookie, nil
}
