package middleware

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"

	sr "github.com/BlockTeam4Boys/digitaldocs/internal/session/repository"
	"github.com/BlockTeam4Boys/digitaldocs/internal/token"
)

type AuthMiddleware struct {
	sessionRepo           *sr.RedisTokenRepository
	accessTokenCookieName string
}

func NewAuthMiddleware(sessionRepo *sr.RedisTokenRepository, accessTokenCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo:           sessionRepo,
		accessTokenCookieName: accessTokenCookieName,
	}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessCookie, err := r.Cookie(a.accessTokenCookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		accessToken := accessCookie.Value
		// check access token
		accessClaims, err := token.AccessTokenInfo(accessToken, a.sessionRepo.GetUserKey)
		sessionID := accessClaims.SessionID
		// refresh token already confirmed
		if accessClaims.ParentRefreshTokenHash == "" {
			ctx := createContext(accessClaims.UserID, sessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		parentHash := accessClaims.ParentRefreshTokenHash
		currentHash := accessClaims.CurrentRefreshTokenHash
		// try to confirm new refresh token
		for i := 0; i < 5; i++ {
			err = a.sessionRepo.CheckAndSetRefreshToken(sessionID, parentHash, currentHash)
			switch err {
			case nil:
				// the removal of parent refresh token means that the new refresh token is confirmed
				accessClaims.ParentRefreshTokenHash = ""
				key, err := a.sessionRepo.GetUserKey(accessClaims.UserID)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				accessToken, err := token.New(accessClaims, []byte(key))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:     a.accessTokenCookieName,
					Value:    accessToken,
					Expires:  time.Unix(accessClaims.ExpiresAt, 0),
					HttpOnly: true,
				})
				ctx := createContext(accessClaims.UserID, sessionID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			case sr.TokensNotEqualErr:
				// token is stolen
				a.sessionRepo.DeleteRefreshToken(sessionID)
				w.WriteHeader(http.StatusUnauthorized)
				return
			case redis.TxFailedErr:
				continue
			default:
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	})
}
