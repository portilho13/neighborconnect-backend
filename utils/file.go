package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func CreateUploadsFolder() error {
	_, err := os.Stat("./uploads")
	if err != nil {
		err = os.Mkdir("./uploads", 0777)
		if err != nil {
			return err
		}

		err = os.Mkdir("./uploads/events", 0777)
		if err != nil {
			return err
		}

		err = os.Mkdir("./uploads/listing", 0777)
		if err != nil {
			return err
		}

		err = os.Mkdir("./uploads/category", 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetApiUrl() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("API_IP")

}
