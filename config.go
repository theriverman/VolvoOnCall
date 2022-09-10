package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Configuration struct {
	Username string
	Password string
	Region   string
	URL      string
}

func (c *Configuration) LoadFromFile(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		// logger here
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	for _, line := range text {
		if strings.HasPrefix(line, "#") {
			continue
		}
		tuple := strings.Split(line, ": ")
		switch {
		case tuple[0] == "username":
			c.Username = tuple[1]
		case tuple[0] == "password":
			c.Password = tuple[1]
		case tuple[0] == "region":
			c.Region = tuple[1]
		case tuple[0] == "url":
			c.URL = tuple[1]
		default:
			fmt.Println("invalid case:", tuple[1])
		}
	}
	return nil
}

func (c *Configuration) WriteToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := fmt.Sprintf("username: %s\npassword: %s\n", c.Username, c.Password)
	if c.Region != "" {
		s = s + fmt.Sprintf("region: %s\n", c.Region)
	}
	if c.URL != "" {
		s = s + fmt.Sprintf("url: %s\n", c.URL)
	}
	if _, err = f.WriteString(s); err != nil {
		return err
	}
	return nil
}
