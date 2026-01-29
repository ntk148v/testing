package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	uuid, err := GetMachineUUID()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(uuid)
}

// GetMachineUUID returns a canonical RFC-4122 UUID
// equivalent to what Veeam uses internally.
func GetMachineUUID() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxUUID()
	case "windows":
		return getWindowsUUID()
	case "darwin":
		return getMacUUID()
	default:
		return getFallbackUUID()
	}
}

// ---------- Linux ----------

func getLinuxUUID() (string, error) {
	data, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return "", err
	}

	uuid := strings.TrimSpace(string(data))
	return transformSMBIOSUUID(uuid)
}

// ---------- Windows ----------

func getWindowsUUID() (string, error) {
	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-Command",
		"(Get-WmiObject Win32_ComputerSystemProduct).UUID",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// ---------- macOS ----------

func getMacUUID() (string, error) {
	cmd := exec.Command(
		"ioreg",
		"-rd1",
		"-c",
		"IOPlatformExpertDevice",
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				return parts[3], nil
			}
		}
	}

	return "", errors.New("UUID not found")
}

// ---------- Fallback ----------

func getFallbackUUID() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}

	h := sha1.Sum([]byte(host))
	return formatUUID(h[:16]), nil
}

// ---------- SMBIOS → RFC-4122 ----------

func transformSMBIOSUUID(dmiUUID string) (string, error) {
	clean := strings.ReplaceAll(dmiUUID, "-", "")
	if len(clean) != 32 {
		return "", errors.New("invalid UUID")
	}

	b, err := hex.DecodeString(clean)
	if err != nil {
		return "", err
	}

	// Reverse first 3 fields (SMBIOS endianness)
	reverse(b[0:4]) // time_low
	reverse(b[4:6]) // time_mid
	reverse(b[6:8]) // time_hi_and_version

	return formatUUID(b), nil
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func formatUUID(b []byte) string {
	return strings.ToLower(
		hex.EncodeToString(b[0:4]) + "-" +
			hex.EncodeToString(b[4:6]) + "-" +
			hex.EncodeToString(b[6:8]) + "-" +
			hex.EncodeToString(b[8:10]) + "-" +
			hex.EncodeToString(b[10:16]),
	)
}

