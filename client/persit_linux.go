package client

// Persist makes sure that the executable is run after a reboot
func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		persistAdmin(path)
	} else {
		persistUser(path)
	}
}

// persistAdmin persistence using admin priviliges
func persistAdmin(path string) {

}

// persistUser persistence using user priviliges
func persistUser(path string) {

}
