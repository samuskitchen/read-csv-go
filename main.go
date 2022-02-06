package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strconv"
)

const (
	fileName      = "users.csv"
	errorFilename = "users_errors.csv"
)

type userData struct {
	Email            string
	FirstName        string
	LastName         string
	Phone            string
	IsPhoneVerified  bool
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
	var csvErrors map[string][][]string = make(map[string][][]string)

	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			log.Fatalf("error close file %s:", err)
		}
	}(csvFile)

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		var isPhoneVerified bool = false
		var latitude32 float32
		var longitude32 float32

		isPhoneVerified, latitude32, longitude32 = lineValidation(line, csvErrors, isPhoneVerified, latitude32, longitude32)

		if len(csvErrors) > 0 {
			continue
		}

		user := userData{
			Email:            line[0],
			FirstName:        line[1],
			LastName:         line[2],
			Phone:            line[3],
			IsPhoneVerified:  isPhoneVerified,
			IdentifyDocument: line[4],
			Latitude:         latitude32,
			Longitude:        longitude32,
			Address1:         line[7],
			Address2:         line[8],
			City:             line[9],
			ZipCode:          line[10],
			StateName:        line[11],
			Country:          line[12],
		}

		fmt.Println(user.Email + " " + user.FirstName + " " + user.LastName)
	}

	if len(csvErrors) > 0 {
		err = createCSVError(csvErrors)
		if err != nil {
			log.Fatal(err)
		}

		//We delete the file when we finish processing it
		err = os.Remove(errorFilename)
		if err != nil {
			log.Fatalf("error writing record to file: %s", err)
		}
	}
}

func lineValidation(line []string, csvErrors map[string][][]string, isPhoneVerified bool, latitude32 float32, longitude32 float32) (bool, float32, float32) {
	// validate email required
	if len(line[0]) <= 0 {
		csvErrors["email is required: "] = append(csvErrors["email is required: "], line)
	}

	// validate length email
	if len(line[0]) <= 0 || len(line[0]) > 255 {
		csvErrors["length email error: "] = append(csvErrors["length email error: "], line)
	}

	// validate format email
	if !validFormatEmail(line[0]) {
		csvErrors["incorrect format email: "] = append(csvErrors["incorrect format email: "], line)
	}

	// validate first name required
	if len(line[1]) <= 0 {
		csvErrors["first name is required: "] = append(csvErrors["first name is required: "], line)
	}

	// validate length first name
	if len(line[1]) <= 0 || len(line[1]) > 255 {
		csvErrors["first name length error: "] = append(csvErrors["first name length error: "], line)
	}

	// validate last name required
	if len(line[2]) <= 0 {
		csvErrors["last name is required: "] = append(csvErrors["last name is required: "], line)
	}

	// validate length last name
	if len(line[2]) <= 0 || len(line[2]) > 255 {
		csvErrors["last name length error: "] = append(csvErrors["last name length error: "], line)
	}

	// validate if contain a phone number
	if len(line[3]) > 0 {
		isPhoneVerified = true
	}

	// validate identity document required
	if len(line[4]) <= 0 {
		csvErrors["identify document is required: "] = append(csvErrors["identify document is required: "], line)
	}

	// validate length identity document
	if len(line[4]) <= 0 || len(line[4]) > 40 {
		csvErrors["error identity document length: "] = append(csvErrors["error identity document length: "], line)
	}

	// validate latitude required
	if len(line[5]) <= 0 {
		csvErrors["latitude is required: "] = append(csvErrors["latitude is required: "], line)
	}

	// validate format latitude
	latitude, err := strconv.ParseFloat(line[5], 32)
	if err != nil {
		csvErrors["error in latitude is not a float number: "] = append(csvErrors["error in latitude is not a float number: "], line)
	}
	latitude32 = float32(latitude)

	// validate longitude required
	if len(line[6]) <= 0 {
		csvErrors["longitude is required: "] = append(csvErrors["longitude is required: "], line)
	}

	// validate format longitude
	longitude, err := strconv.ParseFloat(line[6], 32)
	if err != nil {
		csvErrors["error in longitude is not a float number: "] = append(csvErrors["error in longitude is not a float number: "], line)
	}
	longitude32 = float32(longitude)

	// validate city required
	if len(line[9]) <= 0 {
		csvErrors["city is required: "] = append(csvErrors["city is required: "], line)
	}

	// validate length city
	if len(line[9]) <= 0 || len(line[9]) > 255 {
		csvErrors["error city length: "] = append(csvErrors["error city length: "], line)
	}

	// validate country required
	if len(line[12]) <= 0 {
		csvErrors["country is required: "] = append(csvErrors["country is required: "], line)
	}

	// validate length country
	if len(line[12]) <= 0 || len(line[12]) > 2 {
		csvErrors["country length error : "] = append(csvErrors["country length error : "], line)
	}

	return isPhoneVerified, latitude32, longitude32
}

func validFormatEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func createCSVError(csvErrors map[string][][]string) error {
	csvFile, err := os.Create(errorFilename)
	if err != nil {
		return fmt.Errorf("failed creating file: %s", err)
	}

	csvWriter := csv.NewWriter(csvFile)
	for errorLine, lines := range csvErrors {
		for _, line := range lines {
			line[0] = errorLine + line[0]

			err = csvWriter.Write(line)
			if err != nil {
				return fmt.Errorf("error writing record to file: %s", err)
			}
		}
	}

	csvWriter.Flush()
	err = csvFile.Close()
	if err != nil {
		return fmt.Errorf("error close csv file: %s", err)
	}

	//Add logic to upload the file to S3 or send an email

	return nil
}
