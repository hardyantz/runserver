package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runserver/shared"
	"strings"
	"time"

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

	if !s.ValidateToken("4e8sEXQX1F9xLTuA9uxouXLo") {
		return c.String(http.StatusUnauthorized, "invalid token")
	}

	params := &slack.Msg{Text: s.Text}

	err = CallTravis(c.Request().Context(), params.Text)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("start deploy branch %v", params.Text))
}

func CallTravis(c context.Context, branch string) error {
	req := shared.NewRequest(2, 30*time.Second)

	travisUrl := "https://api.travis-ci.com/repo/hardyantz%2Frunserver/requests"

	params := map[string]interface{}{
		"request": map[string]string{"branch": branch},
	}
	str, _ := json.Marshal(params)
	byteReq := strings.NewReader(string(str))

	headers := map[string]string{
		echo.HeaderContentType:   echo.MIMEApplicationJSON,
		echo.HeaderAccept:        echo.MIMEApplicationJSON,
		"Travis-API-Version":     "3",
		echo.HeaderAuthorization: "token SQ5Z466aaRH67-ZkLa_Clw",
	}

	var target interface{}

	_, err := req.Do(c, http.MethodPost, travisUrl, byteReq, &target, headers)
	if err != nil {
		return err
	}

	return nil
}
