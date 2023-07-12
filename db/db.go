package db

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "Hotel-Reservation"
	testDBNAME = "hotel-reservation-test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
