package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func dumpObj(key string, obj interface{}, w io.Writer) error {
	switch v := obj.(type) {
	case float64:
		fmt.Fprintf(w, "%s: %v\n", key, v)
	case string:
		fmt.Fprintf(w, "%s: %v\n", key, v)
	case []interface{}:
		for i := range v {
			subkey := fmt.Sprintf("%s[%d]", key, i)
			dumpObj(subkey, v[i], w)
		}
	case map[string]interface{}:
		if key != "" {
			key = key + "."
		}
		for k := range v {
			subkey := fmt.Sprintf("%s%s", key, k)
			dumpObj(subkey, v[k], w)
		}
	default:
		return fmt.Errorf("Unable to decode %v", v)
	}
	return nil
}

func main() {
	var fp io.Reader
	var err error
	if len(os.Args) > 2 {
		log.Fatalf("Usage: %s [PATH]", os.Args[0])
	}

	fp = os.Stdin
	if len(os.Args) == 2 {
		path := os.Args[1]
		fp, err = os.Open(path)
		if err != nil {
			log.Fatalf("Error while reading file %s: %v", path, err)
		}
	}

	var jData interface{}
	for {
		err = json.NewDecoder(fp).Decode(&jData)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unable to decode stream: %v", err)
		}
	}
	err = dumpObj("", jData, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
