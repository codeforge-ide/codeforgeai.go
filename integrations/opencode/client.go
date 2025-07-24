package opencode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://api.opencode.ai" // Replace with actual base URL if different

// Shared error types
var (
	ErrMessageAborted = fmt.Errorf("message aborted")
	ErrProviderAuth   = fmt.Errorf("provider auth error")
	ErrUnknown        = fmt.Errorf("unknown error")
)

// --- Event ---
func ListEvents() (string, error) {
	resp, err := http.Get(baseURL + "/event")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// --- App ---
func GetApp() (string, error) {
	resp, err := http.Get(baseURL + "/app")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func InitApp() (string, error) {
	resp, err := http.Post(baseURL+"/app/init", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func LogApp(params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/log", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func GetModes() (string, error) {
	resp, err := http.Get(baseURL + "/mode")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func GetProviders() (string, error) {
	resp, err := http.Get(baseURL + "/config/providers")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// --- Find ---
func FindFiles(params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/find/file", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func FindSymbols(params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/find/symbol", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func FindText(params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/find", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// --- File ---
func ReadFile(params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/file", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func FileStatus() (string, error) {
	resp, err := http.Get(baseURL + "/file/status")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// --- Config ---
func GetConfig() (string, error) {
	resp, err := http.Get(baseURL + "/config")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// --- Session ---
func CreateSession() (string, error) {
	resp, err := http.Post(baseURL+"/session", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func ListSessions() (string, error) {
	resp, err := http.Get(baseURL + "/session")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func DeleteSession(id string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", baseURL+"/session/"+id, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func AbortSession(id string) (string, error) {
	resp, err := http.Post(baseURL+"/session/"+id+"/abort", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func ChatSession(id string, params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/session/"+id+"/message", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func InitSession(id string, params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/session/"+id+"/init", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func GetSessionMessages(id string) (string, error) {
	resp, err := http.Get(baseURL + "/session/" + id + "/message")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func RevertSession(id string, params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/session/"+id+"/revert", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func ShareSession(id string) (string, error) {
	resp, err := http.Post(baseURL+"/session/"+id+"/share", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SummarizeSession(id string, params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/session/"+id+"/summarize", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func UnrevertSession(id string) (string, error) {
	resp, err := http.Post(baseURL+"/session/"+id+"/unrevert", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func UnshareSession(id string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", baseURL+"/session/"+id+"/share", nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// --- Tui ---
func AppendPromptTui(params map[string]interface{}) (string, error) {
	b, _ := json.Marshal(params)
	resp, err := http.Post(baseURL+"/tui/append-prompt", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func OpenHelpTui() (string, error) {
	resp, err := http.Post(baseURL+"/tui/open-help", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
