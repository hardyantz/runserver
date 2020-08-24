package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/slack-go/slack"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8009"
	}

	slackAuth := os.Getenv("SLACK_VERIFICATION_TOKEN")
	fmt.Println(slackAuth)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/slack/build", PushDeploy)

	e.Logger.Fatal(e.Start(":" + port))
}

func PushDeploy(c echo.Context) error {
	s, err := slack.SlashCommandParse(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		return c.String(http.StatusUnauthorized, "invalid token")
	}

	params := &slack.Msg{Text: s.Text}

	return c.String(http.StatusOK, fmt.Sprintf("deployed branch %v", params.Text))
}
