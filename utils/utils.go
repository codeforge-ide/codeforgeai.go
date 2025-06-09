package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrintJSON(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
}

// Add other utility functions as needed.
