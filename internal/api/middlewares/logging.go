package middlewares

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        logrus.WithFields(logrus.Fields{
            "method": c.Request.Method,
            "path":   c.Request.URL.Path,
        }).Info("Request received")

        c.Next()
    }
}
