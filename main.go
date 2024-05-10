package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Rooms struct {
	ID     string
	XCoord int
	YCoord int
}
type Links struct {
	Room1ID string
	Room2ID string
}
type Ants struct {
	ID     string
	XCoord int
	YCoord int
}

func main() {
	if len(os.Args) != 2 {
		println("Wrong input")
		return
	}
	input := os.Args[1]
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ilkSatir string
	var ilkSatirAlindi bool
	var startpoint, endpoint Rooms

	var digerNoktalar []Rooms
	var linksSlice []Links

	for scanner.Scan() {
		satir := scanner.Text()
		if !ilkSatirAlindi {
			ilkSatir = scanner.Text()
			ilkSatirAlindi = true
		}
		parcalar := strings.Fields(satir)

		if len(parcalar) > 0 {

			if satir == "##start" {
				scanner.Scan()
				fmt.Sscanf(scanner.Text(), "%s %d %d", &startpoint.ID, &startpoint.XCoord, &startpoint.YCoord)
			} else if satir == "##end" {
				scanner.Scan()
				fmt.Sscanf(scanner.Text(), "%s %d %d", &endpoint.ID, &endpoint.XCoord, &endpoint.YCoord)
			}
			if len(parcalar) == 3 {
				// Diğer noktaların bilgilerini al ve slice'e ekle
				var room Rooms
				fmt.Sscanf(scanner.Text(), "%s %d %d", &room.ID, &room.XCoord, &room.YCoord)
				digerNoktalar = append(digerNoktalar, room)
			} else if len(parcalar) == 1 && len(parcalar[0]) == 3 {
				// Bağlantıları al ve slice'e ekle
				var link Links
				fmt.Sscanf(scanner.Text(), "%s-%s", &link.Room1ID, &link.Room2ID)
				linksSlice = append(linksSlice, link)
			}
		}
	}
	numbersofants, err := strconv.Atoi(ilkSatir)
	if err != nil {
		return
	}
	for i := 0; i < numbersofants; i++ {
		ants := "L" + strconv.Itoa(i+1)
		fmt.Println(ants)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okuma hatası", err)
	}
	fmt.Println("Start Point:", startpoint)
	fmt.Println("End Point:", endpoint)
	fmt.Println("Other Points:", digerNoktalar)
	fmt.Println("Links:", linksSlice)
}
