package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Structura pentru informatii despre sistemul de operare
type OSInfo struct {
	Nume           string json:"nume"
	Versiune       string json:"versiune"
	Arhitectura    string json:"arhitectura"
	DataInstalarii string json:"data_instalarii"
	Licenta        string json:"licenta"
}

// Functie pentru a obtine informatii despre sistemul de operare
func getOSInfo() (*OSInfo, error) {
	osInfo := &OSInfo{}

	// Obtine numele si versiunea sistemului de operare
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "ver")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		osInfo.Nume = strings.TrimSpace(string(out))

		// Obtine arhitectura sistemului de operare
		osInfo.Arhitectura = runtime.GOARCH

		// Obtine data instalarii sistemului de operare (pentru Windows)
		cmd = exec.Command("cmd", "/c", "wmic os get InstallDate /VALUE")
		out, err = cmd.Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "=")
		if len(lines) > 1 {
			osInfo.DataInstalarii = strings.TrimSpace(lines[1])
		}

		// Obtine licenta sistemului de operare (daca este cazul)
		// Implementati aici logica specifica pentru licenta Windows
		osInfo.Licenta = "N/A"

	case "linux":
		// Implementati logica pentru Linux
		osInfo.Nume = "Linux" // Exemplu simplificat pentru numele sistemului de operare
		osInfo.Versiune = "N/A"
		osInfo.Arhitectura = runtime.GOARCH
		osInfo.DataInstalarii = "N/A"
		osInfo.Licenta = "N/A"

	case "darwin":
		// Implementati logica pentru macOS
		cmd := exec.Command("sw_vers", "-productName")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		osInfo.Nume = strings.TrimSpace(string(out))

		cmd = exec.Command("sw_vers", "-productVersion")
		out, err = cmd.Output()
		if err != nil {
			return nil, err
		}
		osInfo.Versiune = strings.TrimSpace(string(out))

		osInfo.Arhitectura = runtime.GOARCH
		osInfo.DataInstalarii = "N/A" // Exemplu simplificat pentru macOS
		osInfo.Licenta = "N/A"

	default:
		return nil, fmt.Errorf("Sistem de operare neacceptat: %s", runtime.GOOS)
	}

	return osInfo, nil
}

// Structura pentru informatii despre hardware
type HardwareInfo struct {
	Procesor      string json:"procesor"
	Nuclee        int    json:"nuclee"
	FireExecutie  int    json:"fire_executie"
	Frecventa     string json:"frecventa"
	MemorieRAM    string json:"memorie_ram"
	TipStocare    string json:"tip_stocare"
	CapacitateHDD string json:"capacitate_hdd"
	PlacaDeBaza   string json:"placa_de_baza"
	PlacaVideo    string json:"placa_video"
}

// Functie pentru a obtine informatii despre hardware
func getHardwareInfo() (*HardwareInfo, error) {
	hardwareInfo := &HardwareInfo{}

	// Obtine informatii despre procesor
	cmd := exec.Command("cmd", "/c", "wmic cpu get Name,NumberOfCores,NumberOfLogicalProcessors,MaxClockSpeed /FORMAT:LIST")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			switch strings.TrimSpace(fields[0]) {
			case "Name":
				hardwareInfo.Procesor = strings.TrimSpace(fields[1])
			case "NumberOfCores":
				hardwareInfo.Nuclee, _ = strconv.Atoi(strings.TrimSpace(fields[1]))
			case "NumberOfLogicalProcessors":
				hardwareInfo.FireExecutie, _ = strconv.Atoi(strings.TrimSpace(fields[1]))
			case "MaxClockSpeed":
				hardwareInfo.Frecventa = strings.TrimSpace(fields[1]) + " MHz"
			}
		}
	}

	// Obtine informatii despre memoria RAM
	cmd = exec.Command("cmd", "/c", "wmic memorychip get Capacity /FORMAT:LIST")
	out, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	lines = strings.Split(strings.TrimSpace(string(out)), "\n")
	var totalRAM uint64
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			if len(fields) > 1 {
				capacity, _ := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
				totalRAM += capacity
			}
		}
	}
	hardwareInfo.MemorieRAM = fmt.Sprintf("%d GB", totalRAM/(1024*1024*1024))

	// Obtine informatii despre stocare (in acest exemplu, pentru HDD)
	cmd = exec.Command("cmd", "/c", "wmic diskdrive get Model,Size /FORMAT:LIST")
	out, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	lines = strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			switch strings.TrimSpace(fields[0]) {
			case "Model":
				hardwareInfo.TipStocare = strings.TrimSpace(fields[1])
			case "Size":
				sizeBytes, _ := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
				hardwareInfo.CapacitateHDD = fmt.Sprintf("%d GB", sizeBytes/(1024*1024*1024))
			}
		}
		break // Se obtin doar informatiile despre primul dispozitiv de stocare
	}

	// Obtine informatii despre placa de baza
	cmd = exec.Command("cmd", "/c", "wmic baseboard get Manufacturer,Product /FORMAT:LIST")
	out, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	lines = strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			switch strings.TrimSpace(fields[0]) {
			case "Manufacturer":
				hardwareInfo.PlacaDeBaza = strings.TrimSpace(fields[1])
			case "Product":
				hardwareInfo.PlacaDeBaza += " " + strings.TrimSpace(fields[1])
			}
		}
		break // Se obtin doar informatiile despre prima placa de baza
	}

	// Obtine informatii despre placa video
	cmd = exec.Command("cmd", "/c", "wmic path win32_videocontroller get Name /FORMAT:LIST")
	out, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	lines = strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			if len(fields) > 1 {
				hardwareInfo.PlacaVideo = strings.TrimSpace(fields[1])
			}
		}
		break // Se obtin doar informatiile despre prima placa video
	}

	return hardwareInfo, nil
}

// Structura pentru informatii despre software (programe instalate)
type SoftwareInfo struct {
	ProgrameInstalate []ProgramInfo json:"programe_instalate"
}

// Structura pentru informatii despre un program instalat
type ProgramInfo struct {
	Nume          string json:"nume"
	Versiune      string json:"versiune"
	Producator    string json:"producator"
	DataInstalare string json:"data_instalare"
	Licenta       string json:"licenta"
	// Alte informatii despre program
}

// Functie pentru a obtine informatii despre programele instalate
func getInstalledPrograms() ([]ProgramInfo, error) {
	var programs []ProgramInfo

	// Implementati metoda pentru a obtine lista de programe instalate pe sistemul de operare
	// Iata un exemplu simplificat pentru Windows, care utilizeaza registrele de sistem
	// pentru a obtine informatiile necesare. Trebuie adaptata pentru diferite sisteme de operare.

	// Exemplu: Windows
	cmd := exec.Command("cmd", "/c", "wmic product get Name,Version,Vendor,InstallDate /FORMAT:LIST")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var program ProgramInfo
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			switch strings.TrimSpace(fields[0]) {
			case "Name":
				program.Nume = strings.TrimSpace(fields[1])
			case "Version":
				program.Versiune = strings.TrimSpace(fields[1])
			case "Vendor":
				program.Producator = strings.TrimSpace(fields[1])
			case "InstallDate":
				program.DataInstalare = strings.TrimSpace(fields[1])
			}
		} else {
			// Finalizam informatiile pentru un program si adaugam in lista
			if program.Nume != "" {
				programs = append(programs, program)
				program = ProgramInfo{}
			}
		}
	}
	// Adaugam ultimul program (daca exista)
	if program.Nume != "" {
		programs = append(programs, program)
	}

	return programs, nil
}

// Functie pentru a obtine informatii despre securitate
func getSecurityInfo() (string, error) {
	// Implementati metoda pentru a obtine informatii despre statusul de securitate al sistemului
	// Informatiile pot include statusul antivirusului, firewall-ului, etc.

	// Exemplu simplificat pentru Windows
	cmd := exec.Command("cmd", "/c", "wmic /namespace:\\\\root\\SecurityCenter2 path AntiVirusProduct get displayName /FORMAT:LIST")
	out, err := cmd.Output()
	if err != nil {
		return "N/A", err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			fields := strings.Split(line, "=")
			if len(fields) > 1 {
				return strings.TrimSpace(fields[1]), nil
			}
		}
	}
	return "N/A", nil
}

// Structura pentru informatii despre utilizator
type UserInfo struct {
	NumeUtilizator string json:"nume_utilizator"
	GrupUtilizator string json:"grup_utilizator"
	// Alte informatii despre utilizator
}

// Functie pentru a obtine informatii despre utilizatorul curent
func getCurrentUserInfo() (*UserInfo, error) {
	userInfo := &UserInfo{}

	// Implementati metoda pentru a obtine informatii despre utilizatorul curent
	// Aceasta poate include numele utilizatorului, grupul utilizatorului, etc.

	// Exemplu simplificat pentru Windows
	cmd := exec.Command("cmd", "/c", "whoami")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	userInfo.NumeUtilizator = strings.TrimSpace(string(out))

	// Pentru grupul utilizatorului, se poate face o interogare suplimentara sau folosind pachete/librarii aditionale

	return userInfo, nil
}

func main() {
	// Creare structura pentru informatiile complete despre sistem
	systemInfo := make(map[string]interface{})

	// Obtinere informatii despre sistemul de operare
	osInfo, err := getOSInfo()
	if err != nil {
		fmt.Printf("Eroare la obtinerea informatiilor despre sistemul de operare: %v\n", err)
		osInfo = &OSInfo{} // Initializare structura goala
	}
	systemInfo["sistem_de_operare"] = osInfo

	// Obtinere informatii despre hardware
	hardwareInfo, err := getHardwareInfo()
	if err != nil {
		fmt.Printf("Eroare la obtinerea informatiilor despre hardware: %v\n", err)
		hardwareInfo = &HardwareInfo{} // Initializare structura goala
	}
	systemInfo["hardware"] = hardwareInfo

	// Obtinere informatii despre programele instalate
	installedPrograms, err := getInstalledPrograms()
	if err != nil {
		fmt.Printf("Eroare la obtinerea informatiilor despre programele instalate: %v\n", err)
		installedPrograms = []ProgramInfo{} // Initializare slice gol
	}
	softwareInfo := &SoftwareInfo{
		ProgrameInstalate: installedPrograms,
	}
	systemInfo["software"] = softwareInfo

	// Obtinere informatii despre securitate
	securityInfo, err := getSecurityInfo()
	if err != nil {
		fmt.Printf("Eroare la obtinerea informatiilor despre securitate: %v\n", err)
		securityInfo = "N/A" // Initializare cu valoare default
	}
	systemInfo["securitate"] = securityInfo

	// Obtinere informatii despre utilizatorul curent
	userInfo, err := getCurrentUserInfo()
	if err != nil {
		fmt.Printf("Eroare la obtinerea informatiilor despre utilizator: %v\n", err)
		userInfo = &UserInfo{} // Initializare structura goala
	}
	systemInfo["utilizator"] = userInfo

	// Generare fisier JSON cu informatiile complete
	jsonData, err := json.MarshalIndent(systemInfo, "", "    ")
	if err != nil {
		fmt.Printf("Eroare la generarea JSON: %v\n", err)
		return
	}

	// Scriere fisier JSON
	fileName := "system_info.json"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Eroare la crearea fisierului JSON: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Printf("Eroare la scrierea fisierului JSON: %v\n", err)
		return
	}

	fmt.Printf("Informatii sistem salvate cu succes in fisierul %s\n", fileName)
}
