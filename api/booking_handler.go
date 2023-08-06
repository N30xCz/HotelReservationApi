package api

import (
	"context"
	"net/http"

	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(booking)

}
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "Not Authorized",
		})
	}
	return c.JSON(booking)
}
func (h *BookingHandler) HandleCancleBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(ctx)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return ctx.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "Not Authorized",
		})
	}

	if err = h.store.Booking.UpdateBooking(context.Background(), id, bson.M{"canceled": true}); err != nil {
		return err
	}

	return ctx.JSON(genericResp{
		Type: "Msg",
		Msg:  "Booking was updated",
	})

}
