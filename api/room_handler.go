package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}
type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate" bson:"fromDate"`
	TillDate   time.Time `json:"tillDate" bson:"tillDate"`
	NumPersons int       `json:"numPersons" bson:"numPersons"`
}

func (p *BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}
func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var bookParams BookRoomParams
	if err := c.BodyParser(&bookParams); err != nil {
		return err
	}
	if err := bookParams.Validate(); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal seerver error",
		})
	}
	booking := types.Booking{
		RoomID:     roomID,
		UserID:     user.ID,
		FromDate:   bookParams.FromDate,
		TillDate:   bookParams.TillDate,
		NumPersons: bookParams.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", booking)
	return c.JSON(inserted)
}
