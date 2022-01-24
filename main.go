package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type userData struct {
	Email            string
	FirstName        string
	LastName         string
	Phone            string
	IdentifyDocument string
	Latitude         float32
	Longitude        float32
	Address1         string
	Address2         string
	City             string
	ZipCode          string
	StateName        string
	Country          string
}

func main() {
	var csvErrors2 map[string][][]string
	csvErrors2 = make(map[string][][]string)

	csvFile, err := os.Open("users.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {

		//TODO implement the validate(email, etc) found in our services

		if len(line[0]) <= 0 || len(line[0]) > 255 {
			csvErrors2["error in email: "] = append(csvErrors2["error in email: "], line)
			continue
		}

		latitude, err := strconv.ParseFloat(line[5], 32)
		if err != nil {
			csvErrors2["error in latitude: "] = append(csvErrors2["error in latitude: "], line)
			continue
		}

		longitude, err := strconv.ParseFloat(line[6], 32)
		if err != nil {
			csvErrors2["error in latitude: "] = append(csvErrors2["error in latitude: "], line)
			continue
		}

		user := userData{
			Email:            line[0],
			FirstName:        line[1],
			LastName:         line[2],
			Phone:            line[3],
			IdentifyDocument: line[4],
			Latitude:         float32(latitude),
			Longitude:        float32(longitude),
			Address1:         line[7],
			Address2:         line[8],
			City:             line[9],
			ZipCode:          line[10],
			StateName:        line[11],
			Country:          line[12],
		}
		fmt.Println(user.Email + " " + user.FirstName + " " + user.LastName)
	}

	defer createCSVError(csvErrors2)
}

func createCSVError(csvErrors map[string][][]string) {

	csvFile, err := os.Create("users_errors.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for in, empRow := range csvErrors {
		for _, strings := range empRow {
			strings[0] = in + strings[0]
			_ = csvwriter.Write(strings)
		}
	}

	csvwriter.Flush()
	csvFile.Close()
}
