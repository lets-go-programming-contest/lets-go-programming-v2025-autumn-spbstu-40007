package wifi

type Repository interface {
	GetWifiStatus(id int) (string, error)
}

func WiFiHandle(repo Repository, id int) string {
	status, err := repo.GetWifiStatus(id)
	if err != nil {
		return "Error: Database Connection"
	}
	if status == "active" {
		return "WiFi is ON"
	}
	return "WiFi is OFF"
}

func SimpleCheck(ssid string) bool {
	return ssid != ""
}
