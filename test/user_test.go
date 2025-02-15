package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	helper "github.com/yuugure-aikouka/kyoto-common/test/helper"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

var (
	PARTNERS = 1
	NONE     = 2
	LIKE     = 3
	LIKED    = 4
)

func TestGetPartners(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setupTest()

		// create 10 users
		users := []db.User{}
		for i := 0; i < 10; i++ {
			users = append(users, createRandomUser(t))
		}

		// assign randomly partners to the first user
		partnerCount := 0
		partnershipStatus := make([]int, 10)
		for i := 1; i < 10; i++ {
			partnershipStatus[i] = int(utils.RandomInt(1, 2))
			if partnershipStatus[i] == PARTNERS {
				establishPartnership(t, users[0], users[i])
				partnerCount += 1
			}
		}

		// make the request with the first user's id
		url := fmt.Sprintf("/v1/users/%d/partners", users[0].ID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		responseBody, err := helper.UnmarshalResponseBody[[]db.ListPartnersRow](res.Body)
		require.Nil(t, err)

		partnerIDs := []int32{}
		for _, v := range *responseBody.Data {
			partnerIDs = append(partnerIDs, v.ID)
		}

		// validate
		require.Equal(t, res.Code, http.StatusOK)
		require.Equal(t, len(partnerIDs), partnerCount)
		for i := 1; i < 10; i++ {
			if partnershipStatus[i] == PARTNERS {
				require.Contains(t, partnerIDs, users[i].ID)
			} else {
				require.NotContains(t, partnerIDs, users[i].ID)
			}
		}
	})

	t.Run("User not found", func(t *testing.T) {
		setupTest()

		// make the request with a non existing user
		url := fmt.Sprintf("/v1/users/%d/partners", 1)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("Invalid ID Param", func(t *testing.T) {
		setupTest()

		url := "/v1/users/hello/partners"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, res.Code, http.StatusBadRequest)
	})
}

func TestGetPotentialPartners(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setupTest()

		// create 10 users
		users := []db.User{}
		for i := 0; i < 10; i++ {
			users = append(users, createRandomUser(t))
		}

		// for users 2-10, set its status randomly to
		// like/liked/partners/none with user 1
		potentialsCount := 0
		partnershipStatus := make([]int, 10)
		for i := 1; i < 10; i++ {
			partnershipStatus[i] = int(utils.RandomInt(1, 4))
			switch partnershipStatus[i] {
			case LIKE:
				initiatePartnership(t, users[0], users[i])
				break
			case LIKED:
				initiatePartnership(t, users[i], users[0])
				potentialsCount += 1
				break
			case PARTNERS:
				establishPartnership(t, users[0], users[i])
				break
			case NONE:
				potentialsCount += 1
			}
		}

		// make the request with the first user's id
		url := fmt.Sprintf("/v1/users/%d/potential-partners", users[0].ID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		responseBody, err := helper.UnmarshalResponseBody[[]db.ListPotentialPartnersRow](res.Body)
		require.Nil(t, err)

		potentialIDs := []int32{}
		for _, v := range *responseBody.Data {
			potentialIDs = append(potentialIDs, v.ID)
		}

		// validate
		require.Equal(t, res.Code, http.StatusOK)
		require.Equal(t, len(potentialIDs), potentialsCount)
		for i := 1; i < 10; i++ {
			if partnershipStatus[i] == LIKED || partnershipStatus[i] == NONE {
				require.Contains(t, potentialIDs, users[i].ID)
			} else {
				require.NotContains(t, potentialIDs, users[i].ID)
			}
		}
	})

	t.Run("User not found", func(t *testing.T) {
		setupTest()

		// make the request with a non existing user
		url := fmt.Sprintf("/v1/users/%d/potential-partners", 1)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("Invalid ID Param", func(t *testing.T) {
		setupTest()

		url := "/v1/users/i-am-string/potential-partners"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, res.Code, http.StatusBadRequest)
	})
}

func createRandomUser(t *testing.T) db.User {
	user, err := store.CreateUser(context.Background(), db.CreateUserParams{
		Username:    utils.RandomString(8),
		DisplayName: utils.RandomString(8),
		AvatarUrl:   "",
	})

	require.Nil(t, err)

	return user
}

func initiatePartnership(t *testing.T, user1, user2 db.User) {
	_, err := store.CreatePartnership(context.Background(), db.CreatePartnershipParams{
		UserID1: user1.ID,
		UserID2: user2.ID,
	})
	require.Nil(t, err)
}

func establishPartnership(t *testing.T, user1, user2 db.User) {
	initiatePartnership(t, user1, user2)

	err := store.UpdatePartnershipStatus(context.Background(), db.UpdatePartnershipStatusParams{
		UserID1: user1.ID,
		UserID2: user2.ID,
		Status:  "accepted",
	})
	require.Nil(t, err)
}
