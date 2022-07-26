package analytics

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type CookieSession struct {
	now time.Time
	ttl int
	sessions.Session
}

func GetCookieSession(c *gin.Context, now time.Time, ttl int) *CookieSession {
	session := &CookieSession{
		now:     now,
		ttl:     ttl,
		Session: sessions.Default(c),
	}

	session.setCookieExpiry()

	return session
}

func (c *CookieSession) AppendToCookie(cookieName, item string) bool {
	cookieItems := c.getCookieItems(cookieName)

	if _, exists := cookieItems[item]; !exists {
		cookieItems[item] = &struct{}{}
		c.setCookieItems(cookieName, cookieItems)

		return true
	}

	return false
}

func (c *CookieSession) Save() error {
	return c.Session.Save()
}

func (c *CookieSession) setCookieExpiry() {
	createdAt := c.Session.Get("created_at")
	if createdAt == nil {
		c.Session.Options(sessions.Options{
			MaxAge:   c.ttl,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		c.Session.Set("created_at", c.now.Unix())
	} else {
		updatedMaxAge := -1

		secondsSinceCookieCreation := int(c.now.Sub(time.Unix(createdAt.(int64), 0).UTC()).Seconds())
		if secondsSinceCookieCreation < c.ttl {
			updatedMaxAge = c.ttl - secondsSinceCookieCreation
		}

		c.Session.Options(sessions.Options{
			MaxAge:   updatedMaxAge,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})
	}
}

func (c *CookieSession) getCookieItems(cookieName string) map[string]*struct{} {
	result := make(map[string]*struct{})

	cookieData := c.Session.Get(cookieName)
	if cookieData != nil {
		items := strings.Split(strings.TrimSpace(cookieData.(string)), ",")
		for _, item := range items {
			if len(strings.TrimSpace(item)) == 0 {
				continue
			}

			result[item] = &struct{}{}
		}
	}

	return result
}

func (c *CookieSession) setCookieItems(cookieName string, items map[string]*struct{}) {
	arr := make([]string, 0, len(items))
	for item := range items {
		arr = append(arr, item)
	}

	c.Session.Set(cookieName, strings.Join(arr, ","))
}
