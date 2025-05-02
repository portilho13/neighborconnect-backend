package utils

import "os"

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
