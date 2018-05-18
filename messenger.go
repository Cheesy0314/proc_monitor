package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"net/smtp"
	"log"
	"strings"
	"os"
	"time"
)

func sendAlert(message string) {
	json := "{\"error\":\"processor down\", \"list\":\"" + message + " \"}"

	var jsonStr = []byte(json)
	url,success := os.LookupEnv("ALERT_SYSTEM")
	shouldLog,_ := os.LookupEnv("SHOULD_LOG")
	if !success {
		log.Fatal("Can't load URL to send alert to.")
	}
	req,err := http.NewRequest("POST",url,bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept","application/json")
	req.Header.Set("Content-Type","application/json")

	if err != nil {
		log.Printf("An error occurred when trying to create request: %s\n",err)
		return
	}

	client := &http.Client{}
	client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Fatal request error occurred: %s\n",err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if shouldLog == "true" {
		//Log resp from server
		log.Printf("Response Status: %s\n", resp.Status)
		log.Printf("Response Body: %s\n\n", string(body))
	}

	shouldEmail,_ := os.LookupEnv("SHOULD_EMAIL")
	if shouldEmail == "true" {
		sendEmail(message)
	}
}

func sendEmail(message string) {
	addrs, _ := os.LookupEnv("RECIPIENTS")
	emailAddr,found := os.LookupEnv("SENDER")

	if !found {
		log.Println("Could not load email, error occurred")
		return
	}

	pass,_:= os.LookupEnv("EMAIL_PASS")
	to := strings.Split(addrs,",")
	msg := "From: " + emailAddr + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: Critical Java Server Failure\n\n" +
		"The following java processes are down: \n" + strings.Join(strings.Split(message, ", "), "\n") + "\n\nThis is an automated message\r\n"
	a := smtp.PlainAuth("", emailAddr,pass,"smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:465", a,  emailAddr, to, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s\n", err)
	}
}
