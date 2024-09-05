package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go main()
	os.Exit(m.Run())
}

func TestServerHappyPath(t *testing.T) {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "Hello, Server\n")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("could not read from connection: %v", err)
	}

	if message != "Hello, Server\n" {
		t.Errorf("unexpected response: got %v want %v", message, "Hello, Server\n")
	}
}

func TestServerNoPortProvided(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	out, _ := bufio.NewReader(r).ReadString('\n')
	os.Stdout = rescueStdout

	if out != "Please provide a port number!\n" {
		t.Errorf("unexpected output: got %v want %v", out, "Please provide a port number!\n")
	}
}

func TestServerInvalidPort(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "invalidPort"}

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	out, _ := bufio.NewReader(r).ReadString('\n')
	os.Stdout = rescueStdout

	if out == "Listening on :invalidPort\n" {
		t.Errorf("server started with invalid port: got %v", out)
	}
}

func TestServerConnectionRefused(t *testing.T) {
	conn, err := net.Dial("tcp", ":8082")
	if err == nil {
		conn.Close()
		t.Fatal("connection should have been refused")
	}
}

func TestServerHandleConnectionError(t *testing.T) {
	ln, err := net.Listen("tcp", ":8083")
	if err != nil {
		t.Fatalf("could not start server: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatalf("could not accept connection: %v", err)
		}
		conn.Close()
	}()

	conn, err := net.Dial("tcp", ":8083")
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	conn.Close()

	// Wait for the server to handle the connection and close it
	time.Sleep(100 * time.Millisecond)
}
