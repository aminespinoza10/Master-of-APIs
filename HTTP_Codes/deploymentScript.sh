#!/bin/bash

#Go creation
mkdir GoCodes
cat <<EOT >> GoCodes/main.go
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
EOT

#Python creation
mkdir PythonCodes
cat <<EOT >> PythonCodes/app.py
print("Hello, World!")
EOT

#Dotnet creation
dotnet new webapi -o DotnetCodes

# How to run the servers
# dotnet (in the NetCodes directory)
# dotnet run

# Go (in the GoCodes directory)
# go run main.go

# Python (in the PythonCodes directory)
# python3 -m venv venv && source venv/bin/activate
# /Master-of-APIs/.venv/bin/python -m uvicorn app:app --reload
# python -m uvicorn app:app --reload
