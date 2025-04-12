package main

import (
	"fmt"
	"path/filepath"
	"time"

	// "golang-airplane/internal/components/airplane"
	"golang-airplane/internal/components/flight"
	"golang-airplane/internal/core/domain"
	"golang-airplane/internal/storage/json"
	"golang-airplane/internal/utils"
)

// App represents the main application
type App struct {
	flightService      *flight.Service
	reservationService *flight.ReservationService
	validation         *utils.ValidationService
	dataManager        *utils.DataManager
}

func main() {
	// Initialize storage, repositories, and services
	dataDir := filepath.Join(".", "data") // Store data in ./data directory
	
	// Setup storage
	storage := json.NewStorage(dataDir)
	flightRepo := json.NewFlightRepository(storage)
	reservationRepo := json.NewReservationRepository(storage)
	
	// Setup services
	flightService := flight.NewService(flightRepo, reservationRepo)
	reservationService := flight.NewReservationService(flightRepo, reservationRepo)
	validation := utils.NewValidationService()
	dataManager := utils.NewDataManager(dataDir)
	
	// Create app
	app := &App{
		flightService:      flightService,
		reservationService: reservationService,
		validation:         validation,
		dataManager:        dataManager,
	}
	
	// Run the app
	app.run()
}

// run starts the application main loop
func (app *App) run() {
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("|             AIRLINE MANAGEMENT SYSTEM                     |")
	fmt.Println("+-----------------------------------------------------------+")
	
	menu := []string{
		"Add a Flight",
		"Book a Flight",
		"Check-in",
		"Assign Crew to Flight",
		"Display All Flights",
		"Display Reservations of a Flight",
		"Exit",
	}
	
	for {
		// Display menu
		fmt.Println("\nMenu Options:")
		for i, option := range menu {
			fmt.Printf("%d. %s\n", i+1, option)
		}
		
		// Get user choice
		choice := app.validation.GetInteger("\nPlease select an option: ", "Invalid input. Please enter a number.", 1, len(menu))
		
		// Handle the choice
		switch choice {
		case 1:
			app.addFlightMenu()
		case 2:
			app.bookFlightMenu()
		case 3:
			app.checkInMenu()
		case 4:
			app.assignCrewMenu()
		case 5:
			app.displayAllFlightsMenu()
		case 6:
			app.displayFlightReservationsMenu()
		case 7:
			fmt.Println("Exiting program. Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// addFlightMenu handles adding a new flight
func (app *App) addFlightMenu() {
	fmt.Println("\n--- Add Flight ---")
	
	for {
		flightNumber := app.validation.GetString("Enter flight number (Must be Fxxxx and no space): ", 
			"Flight number should match the format Fxxxx", false)
		
		if !app.validation.ValidateFlightNumber(flightNumber) {
			fmt.Println("Flight number must be in the format Fxxxx (e.g., F1234)")
			continue
		}
		
		// Check if flight already exists
		flight, err := app.flightService.GetFlight(flightNumber)
		if err == nil && flight != nil {
			fmt.Println("Flight number already exists. Please enter a unique flight number.")
			continue
		}
		
		departureCity := app.validation.GetString("Enter departure city: ", "Departure city cannot be empty", false)
		destinationCity := app.validation.GetString("Enter destination city: ", "Destination city cannot be empty", false)
		
		fmt.Println("+------------------------------------------------------------------------------------------+")
		fmt.Println("|<> Please enter departure and arrival times based on our instruction:                     |")
		fmt.Println("|1. Time cannot be empty.                                                                  |")
		fmt.Println("|2. Time format dd/mm/yyyy-HH:mm (day/month/year-Hour:minutes)                             |")
		fmt.Println("|3. Times must be later than current time (initial flight time) 24 hours.                  |")
		fmt.Println("|4. Duration of commercial flight must be between 30 minutes and 24 hours.                 |")
		fmt.Println("+------------------------------------------------------------------------------------------+")
		
		var departureTime, arrivalTime time.Time
		
		for {
			departureTime = app.validation.GetDate("Enter departure time (format dd/MM/yyyy-HH:mm): ",
				"Please follow our format and input realistic times, try again", "02/01/2006-15:04", false)
			
			arrivalTime = app.validation.GetDate("Enter arrival time (format dd/MM/yyyy-HH:mm): ",
				"Please follow our format and input realistic times, try again", "02/01/2006-15:04", false)
			
			if app.validation.ValidateDates(departureTime, arrivalTime) {
				break
			}
		}
		
		availableSeat := app.validation.GetInteger("Enter available seats for flight: ",
			"Available seats of commercial flight must be between 36 and 853 and cannot be empty, try again!", 36, 853)
		
		flight, err = app.flightService.AddFlight(flightNumber, departureCity, destinationCity, departureTime, arrivalTime, availableSeat)
		if err != nil {
			fmt.Printf("Error adding flight: %v\n", err)
			continue
		}
		
		fmt.Println("Flight information:")
		fmt.Println(flight)
		
		if !app.validation.CheckYesOrNo("Do you want to add another flight? \nChoose 'Y' for YES || Choose 'N' for NO : ") {
			break
		}
	}
}

// bookFlightMenu handles booking a flight
func (app *App) bookFlightMenu() {
	fmt.Println("\n--- Book Flight ---")
	
	for {
		location := app.validation.GetString("Please input a location (departure or destination): ", 
			"Location cannot be empty", false)
		
		date := app.validation.GetDate("Please enter date (departure or arrival) to find flight (dd/mm/yyyy): ", 
			"Please follow our format and input realistic times, try again", "02/01/2006", false)
		
		// Search for flights
		flights, err := app.flightService.SearchFlights(location, date)
		if err != nil {
			fmt.Printf("Error searching flights: %v\n", err)
			continue
		}
		
		if len(flights) == 0 {
			fmt.Println("No flights found for the given location and date.")
			return
		}
		
		fmt.Println("Found Flight(s):")
		
		// Display flights
		fmt.Println("+-----+--------------+--------------------+--------------------+--------------------+--------------------+--------------------+--------------------+")
		fmt.Println("|Index|Flight number |   Departure City   | Destination City   |   Departure time   |    Arrival time    |    Available Seat  |   Flight Duration  |")
		fmt.Println("+-----+--------------+--------------------+--------------------+--------------------+--------------------+--------------------+--------------------+")
		
		for i, flight := range flights {
			fmt.Printf("|%5d| %-12s | %-18s | %-18s | %-18s | %-18s | %-18d | %-18s |\n", i+1,
				flight.FlightNumber, flight.DepartureCity, flight.DestinationCity,
				flight.DepartureTime.Format("02/01/2006"), flight.ArrivalTime.Format("02/01/2006"),
				flight.AvailableSeat, flight.GetDuration())
			fmt.Print("+-----+--------------+--------------------+--------------------+--------------------+--------------------+--------------------+--------------------+\n")
		}
		
		// Select a flight
		selectedIndex := app.validation.GetInteger("Select a flight by entering the corresponding number at the 'index' column of each flight: ",
			"Invalid selection, please try again", 1, len(flights))
		
		selectedFlight := flights[selectedIndex-1]
		
		// Check available seats
		if selectedFlight.AvailableSeat <= 0 {
			fmt.Println("Available slots for the flight are running out. Cannot add a reservation.")
			return
		}
		
		// Enter customer information
		name := app.validation.GetString("Enter name: ", "Name cannot be empty", false)
		address := app.validation.GetString("Enter address: ", "Address cannot be empty", false)
		phoneNumber := app.validation.GetLong("Enter phone number: ", "Phone number must be a valid number", false)
		idCardNumber := app.validation.GetLong("Enter identity card number: ", "ID card number must be a valid number", false)
		
		// Create reservation
		reservation, err := app.reservationService.BookFlight(name, address, phoneNumber, idCardNumber, selectedFlight.FlightNumber)
		if err != nil {
			fmt.Printf("Error booking flight: %v\n", err)
			continue
		}
		
		fmt.Printf("Reservation ID: %s added successfully.\nReservation ID is required for check-in progress, selecting a seat, and receiving a boarding pass\n\n",
			reservation.ReservationID)
		fmt.Println(reservation)
		fmt.Println("When you go to the airport, please select the 'Flight check-in' option to choose your seat and receive your boarding pass.\n")
		
		if !app.validation.CheckYesOrNo("Do you want to create another reservation? \nChoose 'Y' for YES || Choose 'N' for NO : ") {
			break
		}
	}
}

// checkInMenu handles the check-in process
func (app *App) checkInMenu() {
	fmt.Println("\n--- Check-In ---")
	
	for {
		reservationID := app.validation.GetString("Please input reservation ID: ", "Reservation ID cannot be empty", false)
		
		// Find the reservation
		reservation, err := app.reservationService.GetReservation(reservationID)
		if err != nil {
			fmt.Printf("No such reservation ID found: %v\n", err)
			return
		}
		
		// Find the flight
		flight, err := app.flightService.GetFlight(reservation.ReservationFlightNumber)
		if err != nil {
			fmt.Printf("No such flight found for this reservation: %v\n", err)
			return
		}
		
		if reservation.CheckedIn {
			fmt.Println("This reservation has already been checked in. Please try with another reservation.")
			return
		}
		
		// Display available seats and let the user select one
		fmt.Println("Please choose your seat on this journey:")
		app.displaySeatsMap(flight)
		
		seatNumber := app.validation.GetString("Enter the seat number you want to choose: ", 
			"Seat number cannot be empty", false)
		
		// Perform check-in
		err = app.reservationService.CheckIn(reservationID, seatNumber)
		if err != nil {
			fmt.Printf("Error checking in: %v\n", err)
			continue
		}
		
		// Get updated reservation
		reservation, _ = app.reservationService.GetReservation(reservationID)
		
		// Display boarding pass
		fmt.Println(reservation.BoardingPassToString(flight))
		
		if !app.validation.CheckYesOrNo("Do you want to get another boarding pass? \nChoose 'Y' for YES || Choose 'N' for NO : ") {
			break
		}
	}
}

// assignCrewMenu handles assigning crew to a flight
func (app *App) assignCrewMenu() {
	fmt.Println("\n--- Assign Crew to Flight ---")
	
	for {
		flightNumber := app.validation.GetString("Enter flight number (Fxxxx and no space): ", 
			"Flight number should match the format Fxxxx", false)
		
		if !app.validation.ValidateFlightNumber(flightNumber) {
			fmt.Println("Flight number must be in the format Fxxxx (e.g., F1234)")
			continue
		}
		
		// Check if flight exists
		flight, err := app.flightService.GetFlight(flightNumber)
		if err != nil {
			fmt.Printf("Flight number does not exist: %v\n", err)
			return
		}
		
		// Check if flight already has crew assigned
		if len(flight.CrewMembers) > 0 {
			fmt.Println("This flight already has a crew assigned.")
			return
		}
		
		// Input crew members
		crewMembers := app.inputCrew()
		
		// Assign crew to flight
		err = app.flightService.AssignCrew(flightNumber, crewMembers)
		if err != nil {
			fmt.Printf("Error assigning crew: %v\n", err)
			continue
		}
		
		// Display confirmation
		flight, _ = app.flightService.GetFlight(flightNumber) // Get updated flight info
		fmt.Printf("Crew of flight %s added successfully\n\n", flightNumber)
		
		// Display crew list
		fmt.Println("+-------------------+-------------------+")
		fmt.Println("|        Name       |     Position      |")
		fmt.Println("+-------------------+-------------------+")
		for _, crew := range flight.CrewMembers {
			fmt.Printf("|%-18s |%-18s |\n", crew.Name, crew.Position)
			fmt.Println("+-------------------+-------------------+")
		}
		
		if !app.validation.CheckYesOrNo("Do you want to add crew for another flight? \nChoose 'Y' for YES || Choose 'N' for NO : ") {
			break
		}
	}
}

// inputCrew handles the input of crew members
func (app *App) inputCrew() []domain.Crew {
	crewList := []domain.Crew{}
	pilotCount := 0
	attendantCount := 0
	groundStaffCount := 0
	
	fmt.Println("+-------------------------------------Assign-Crew------------------------------------------+")
	fmt.Println("|<> Please assign crew base on our instruction:                                            |")
	fmt.Println("|1. Must have at least one crew member for each position (Pilot, Attendant, Ground Staff). |")
	fmt.Println("|2. Maximum 2 pilots allowed                                                               |")
	fmt.Println("|3. Crew members must be less or equal to the initial quantity of crew members             |")
	fmt.Println("+------------------------------------------------------------------------------------------+")
	fmt.Println()
	
	maxCrewMember := app.validation.GetInteger("Please input quantity of crew members: ", 
		"Quantity must be a positive integer", 3, 100)
	
	for len(crewList) < maxCrewMember {
		name := app.validation.GetString("Please input name of crew member (Enter 'Q' if you want to stop input new crew members): ",
			"Name should not be blank", false)
		
		if name == "Q" || name == "q" {
			break
		}
		
		choice := app.validation.GetInteger("Input their position (1. Pilot - 2. Attendant - 3. Ground Staff): ",
			"Must be an integer between 1 and 3", 1, 3)
		
		var position string
		switch choice {
		case 1:
			if pilotCount >= 2 {
				fmt.Println("Maximum 2 pilots allowed.")
				continue
			}
			position = "Pilot"
			pilotCount++
		case 2:
			position = "Attendant"
			attendantCount++
		case 3:
			position = "Ground Staff"
			groundStaffCount++
		default:
			fmt.Println("Invalid position")
			continue
		}
		
		crewList = append(crewList, domain.Crew{Name: name, Position: position})
	}
	
	// Check if we have at least one of each position
	if pilotCount < 1 || attendantCount < 1 || groundStaffCount < 1 {
		fmt.Println("You must have at least one crew member for each position (Pilot, Attendant, Ground Staff).")
		return []domain.Crew{}
	}
	
	return crewList
}

// displayAllFlightsMenu displays all flights sorted by departure time
func (app *App) displayAllFlightsMenu() {
	fmt.Println("\n--- All Flights ---")
	
	flights, err := app.flightService.ListAllFlights()
	if err != nil {
		fmt.Printf("Error retrieving flights: %v\n", err)
		return
	}
	
	if len(flights) == 0 {
		fmt.Println("No flights found.")
		return
	}
	
	// Display flights in a table format
	fmt.Println("+--------------+--------------------+--------------------+--------------------+--------------------+--------------------+--------------------+")
	fmt.Println("|Flight number |   Departure City   | Destination City   |   Departure time   |    Arrival time    |    Available Seat  |   Flight Duration  |")
	fmt.Println("+--------------+--------------------+--------------------+--------------------+--------------------+--------------------+--------------------+")
	
	for _, flight := range flights {
		fmt.Printf("| %-12s | %-18s | %-18s | %-18s | %-18s | %-18d | %-18s |\n",
			flight.FlightNumber, flight.DepartureCity, flight.DestinationCity,
			flight.DepartureTime.Format("02/01/2006-15:04"), flight.ArrivalTime.Format("02/01/2006-15:04"),
			flight.AvailableSeat, flight.GetDuration())
		fmt.Print("+--------------+--------------------+--------------------+--------------------+--------------------+--------------------+--------------------+\n")
	}
}

// displayFlightReservationsMenu shows all reservations for a selected flight
func (app *App) displayFlightReservationsMenu() {
	fmt.Println("\n--- Flight Reservations ---")
	
	flightNumber := app.validation.GetString("Enter flight number (Must be Fxxxx and no space): ", 
		"Flight number should match the format Fxxxx", false)
	
	if !app.validation.ValidateFlightNumber(flightNumber) {
		fmt.Println("Flight number must be in the format Fxxxx (e.g., F1234)")
		return
	}
	
	// Check if flight exists
	_, err := app.flightService.GetFlight(flightNumber)
	if err != nil {
		fmt.Println("Flight is not found")
		return
	}
	
	// Get reservations for the flight
	reservations, err := app.reservationService.GetReservationsForFlight(flightNumber)
	if err != nil {
		fmt.Printf("Error retrieving reservations: %v\n", err)
		return
	}
	
	fmt.Printf("Reservations for flight %s:\n", flightNumber)
	
	if len(reservations) == 0 {
		fmt.Println("No reservation in flight yet")
		return
	}
	
	// Display reservations in a table format
	fmt.Println("+--------------+--------------------+--------------------+--------------------+--------------------+------------------------------------+")
	fmt.Println("|ReservationID |        Name        |     Phone Number   |   ID Card Number   |    Seat Number     |                  Address           |")
	fmt.Println("+--------------+--------------------+--------------------+--------------------+--------------------+------------------------------------+")
	
	for _, reservation := range reservations {
		seatLocation := reservation.SeatLocation
		if seatLocation == "" {
			seatLocation = "   X"
		}
		
		fmt.Printf("| %-12s | %-18s | %-18d | %-18d | %-18s |%-36s|\n",
			reservation.ReservationID, reservation.Name, reservation.PhoneNumber,
			reservation.IdentityCardNumber, seatLocation, reservation.Address)
		fmt.Print("+--------------+--------------------+--------------------+--------------------+--------------------+------------------------------------+\n")
	}
}

// displaySeatsMap displays the seating layout for a flight
func (app *App) displaySeatsMap(flight *domain.Flight) {
	maxSeatsPerRow := 4 // Maximum seats per row
	seatIndex := 0
	row := 1
	
	fmt.Println("=====================================================================")
	
	for seat, available := range flight.SeatList {
		if seatIndex == 0 {
			fmt.Printf("Row %-4d|", row)
		}
		
		if available {
			fmt.Printf(" %-6s |", seat)
		} else {
			fmt.Printf(" %s(x) |", seat)
		}
		
		seatIndex++
		
		if seatIndex == 2 {
			fmt.Print("                       |")
		}
		
		if seatIndex >= maxSeatsPerRow {
			fmt.Println()
			seatIndex = 0
			row++
		}
	}
	
	// Ensure we end with a newline
	if seatIndex != 0 {
		fmt.Println()
	}
	
	fmt.Println("=====================================================================")
}