package domain

import (
	"fmt"
	"strings"
	"time"
)

// Airplane represents an aircraft that can be assigned to flights
type Airplane struct {
	ID       string `json:"id"`
	Model    string `json:"model"`
	Capacity int    `json:"capacity"`
}

// NewAirplane creates a new Airplane instance
func NewAirplane(id string, model string, capacity int) *Airplane {
	return &Airplane{
		ID:       id,
		Model:    model,
		Capacity: capacity,
	}
}

// Crew represents a crew member for a flight
type Crew struct {
	Name     string `json:"name"`
	Position string `json:"position"`
}

// Flight represents an airplane flight
type Flight struct {
	FlightNumber    string            `json:"flight_number"`
	DepartureCity   string            `json:"departure_city"`
	DestinationCity string            `json:"destination_city"`
	DepartureTime   time.Time         `json:"departure_time"`
	ArrivalTime     time.Time         `json:"arrival_time"`
	FlightCapacity  int               `json:"flight_capacity"` // Total capacity of the flight
	AvailableSeat   int               `json:"available_seat"`  // Available seats
	CrewMembers     []Crew            `json:"crew_members"`
	SeatList        map[string]bool   `json:"seat_list"` // key=seat number, value=available(true)/occupied(false)
}

// NewFlight creates a new Flight instance
func NewFlight(flightNumber, departureCity, destinationCity string, departureTime, arrivalTime time.Time, availableSeat int) *Flight {
	flight := &Flight{
		FlightNumber:    flightNumber,
		DepartureCity:   departureCity,
		DestinationCity: destinationCity,
		DepartureTime:   departureTime,
		ArrivalTime:     arrivalTime,
		FlightCapacity:  availableSeat,
		AvailableSeat:   availableSeat,
		CrewMembers:     []Crew{},
		SeatList:        make(map[string]bool),
	}
	flight.generateSeatList()
	return flight
}

// GenerateSeatList creates the initial seat map for the flight
func (f *Flight) generateSeatList() {
	row := 1
	seatLetters := []rune{'A', 'B', 'C', 'D'} 
	seatIndex := 0

	for i := 0; i < f.FlightCapacity; i++ {
		seatNumber := fmt.Sprintf("%d%c", row, seatLetters[seatIndex])
		f.SeatList[seatNumber] = true // true means the seat is available

		seatIndex++
		if seatIndex >= len(seatLetters) {
			seatIndex = 0
			row++
		}
	}
}

// AddCrewMember adds a new crew member to the flight
func (f *Flight) AddCrewMember(crew Crew) {
	f.CrewMembers = append(f.CrewMembers, crew)
}

// AssignCrew assigns multiple crew members to the flight
func (f *Flight) AssignCrew(crewMembers []Crew) {
	f.CrewMembers = crewMembers
}

// GetDuration returns the flight duration as a string in format "hh:mm"
func (f *Flight) GetDuration() string {
	duration := f.ArrivalTime.Sub(f.DepartureTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// String returns a string representation of the flight
func (f *Flight) String() string {
	var sb strings.Builder
	sb.WriteString("+---------------------------------------------------------------+\n")
	sb.WriteString("|                      FLIGHT INFORMATION                       |\n")
	sb.WriteString("+---------------------------+-----------------------------------+\n")
	sb.WriteString(fmt.Sprintf("| Flight Number             | %-27s       |\n", f.FlightNumber))
	sb.WriteString(fmt.Sprintf("| Departure City            | %-27s       |\n", f.DepartureCity))
	sb.WriteString(fmt.Sprintf("| Destination City          | %-27s       |\n", f.DestinationCity))
	sb.WriteString(fmt.Sprintf("| Departure Time            | %-27s       |\n", f.DepartureTime.Format("02/01/2006-15:04")))
	sb.WriteString(fmt.Sprintf("| Arrival Time              | %-27s       |\n", f.ArrivalTime.Format("02/01/2006-15:04")))
	sb.WriteString(fmt.Sprintf("| Available Seat            | %-27d       |\n", f.AvailableSeat))
	sb.WriteString(fmt.Sprintf("| Flight duration           | %-27s       |\n", f.GetDuration()))
	sb.WriteString("+---------------------------+-----------------------------------+\n")

	// Crew Information
	if len(f.CrewMembers) > 0 {
		sb.WriteString("| Crew Members:                                                 |\n")
		sb.WriteString("+-------------------+-------------------+\n")
		sb.WriteString("|        Name       |     Position      |\n")
		sb.WriteString("+-------------------+-------------------+\n")
		
		for _, crew := range f.CrewMembers {
			sb.WriteString(fmt.Sprintf("|%-18s |%-18s |\n", crew.Name, crew.Position))
			sb.WriteString("+-------------------+-------------------+\n")
		}
	} else {
		sb.WriteString("| No crew assigned yet                                          |\n")
		sb.WriteString("+---------------------------------------------------------------+\n")
	}

	return sb.String()
}

// Reservation represents a flight booking
type Reservation struct {
	ReservationID           string    `json:"reservation_id"`
	Name                    string    `json:"name"`
	Address                 string    `json:"address"`
	PhoneNumber             int64     `json:"phone_number"`
	IdentityCardNumber      int64     `json:"identity_card_number"`
	ReservationFlightNumber string    `json:"reservation_flight_number"`
	SeatLocation            string    `json:"seat_location"`
	CheckedIn               bool      `json:"checked_in"`
	ReservationTime         time.Time `json:"reservation_time"`
}

// NewReservation creates a new Reservation
func NewReservation(name, address string, phoneNumber, identityCardNumber int64, flightNumber string) *Reservation {
	return &Reservation{
		ReservationID:           generateReservationID(),
		Name:                    name,
		Address:                 address,
		PhoneNumber:             phoneNumber,
		IdentityCardNumber:      identityCardNumber,
		ReservationFlightNumber: flightNumber,
		CheckedIn:               false,
		ReservationTime:         time.Now(),
	}
}

// Global counter for reservation IDs
var reservationIDCounter = 0

// generateReservationID generates a unique reservation ID
func generateReservationID() string {
	reservationIDCounter++
	return fmt.Sprintf("R%04d", reservationIDCounter)
}

// CheckIn marks the reservation as checked in
func (r *Reservation) CheckIn() {
	r.CheckedIn = true
}

// String returns a string representation of the reservation
func (r *Reservation) String() string {
	var sb strings.Builder
	sb.WriteString("+------------------------------------------------------------+\n")
	sb.WriteString("|                    RESERVATION DETAILS                     |\n")
	sb.WriteString("+-------------------------+----------------------------------+\n")
	sb.WriteString(fmt.Sprintf("| Reservation ID          | %-30s |\n", r.ReservationID))
	sb.WriteString(fmt.Sprintf("| Name                    | %-30s |\n", r.Name))
	sb.WriteString(fmt.Sprintf("| Address                 | %-30s |\n", r.Address))
	sb.WriteString(fmt.Sprintf("| Phone Number            | %-30d |\n", r.PhoneNumber))
	sb.WriteString(fmt.Sprintf("| ID Card Number          | %-30d |\n", r.IdentityCardNumber))
	sb.WriteString(fmt.Sprintf("| Flight Number           | %-30s |\n", r.ReservationFlightNumber))
	if r.SeatLocation != "" {
		sb.WriteString(fmt.Sprintf("| Seat Location           | %-30s |\n", r.SeatLocation))
	} else {
		sb.WriteString("| Seat Location           | Not Assigned                   |\n")
	}
	checkInStatus := "No"
	if r.CheckedIn {
		checkInStatus = "Yes"
	}
	sb.WriteString(fmt.Sprintf("| Checked In              | %-30s |\n", checkInStatus))
	sb.WriteString("+-------------------------+----------------------------------+\n")
	return sb.String()
}

// BoardingPassToString generates a boarding pass for the reservation
func (r *Reservation) BoardingPassToString(flight *Flight) string {
	var sb strings.Builder
	sb.WriteString("+---------------------------------------------------------------+\n")
	sb.WriteString("|                        BOARDING PASS                          |\n")
	sb.WriteString("+---------------------------+----------------------------------+\n")
	sb.WriteString(fmt.Sprintf("| Passenger Name           | %-30s |\n", r.Name))
	sb.WriteString(fmt.Sprintf("| Flight                   | %-30s |\n", r.ReservationFlightNumber))
	sb.WriteString(fmt.Sprintf("| From                     | %-30s |\n", flight.DepartureCity))
	sb.WriteString(fmt.Sprintf("| To                       | %-30s |\n", flight.DestinationCity))
	sb.WriteString(fmt.Sprintf("| Date                     | %-30s |\n", flight.DepartureTime.Format("02/01/2006")))
	sb.WriteString(fmt.Sprintf("| Time                     | %-30s |\n", flight.DepartureTime.Format("15:04")))
	sb.WriteString(fmt.Sprintf("| Seat                     | %-30s |\n", r.SeatLocation))
	sb.WriteString("+---------------------------+----------------------------------+\n")
	sb.WriteString("|                 THANK YOU FOR FLYING WITH US                 |\n")
	sb.WriteString("+---------------------------------------------------------------+\n")
	return sb.String()
}