package users

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/errors"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
)

func GetPageData(w http.ResponseWriter, r *http.Request) *types.PageData {
    data := &types.PageData{}

    getUserData(r, data)

	return data
}

func getUserData(r *http.Request, data *types.PageData) {
    gothUser, err := auth.GetSessionUser(r)
    if err != nil {
        if err == errors.ErrorNotAuthenticated {
            log.Printf("%s", err)
            return
        }
        panic(err)
    }

    dbUser, err := database.Queries.GetUserByEmail(r.Context(), gothUser.Email)
    if err != nil {
        panic(err)
    }

    group, err := database.Queries.GetGroupInfo(r.Context(), dbUser.Group)

    data.Profile.User = dbUser
    data.Profile.UserColor = group.Color
    data.Profile.UserGroup = group.Displayname.String
}

func VerifyUser(context context.Context, inUser *goth.User) {
    _, err := database.Queries.GetUserByEmail(context, inUser.Email)
    if err != nil {
        if err == pgx.ErrNoRows {
            registerUser(context, inUser)
            return
        }
        panic(err)
    }
}

func registerUser(context context.Context, inUser *goth.User) {
    params := repository.AddUserParams{}

    params.Email = inUser.Email
    params.Group = "default"
    params.Image = inUser.AvatarURL
    params.Created = time.Now()
    params.Name = inUser.Name

    switch inUser.Provider {
    case "google":
        params.Username = strings.ReplaceAll(strings.ToLower(inUser.Name), " ", "")
    case "github":
        params.Username = strings.ReplaceAll(strings.ToLower(inUser.NickName), " ", "")
    }

    err := database.Queries.AddUser(context, params)
    if err != nil {
        panic(err)
    }
}
