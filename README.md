# go-ail-examples

<img align="right" width="159px" src="https://digital.aadhyarupam.com/images/AadhyarupamRound_Logo_150x150.gif">

This Go based project of Aadhyarupam Innovators demonstrate the code examples for building microservices, integration with cloud services (Google Cloud Firestore), application configuration management (Viper) etc.

You can refer our videos on YouTube channel here: https://www.youtube.com/channel/UC0uB6NjgFG3OvRNkXA24NnA

# Quick Start
go run main.go

# Pre-requisite
To execute this project, you require "Go" installed on your system.
Refer this link to download and install go language(golang): https://go.dev/doc/install

# Setup Firestore
To setup firestore, follow below steps:
1. Create Firestore project in Google Cloud Platform
2. Create service account key and download key as JSON file on your system
3. Set environment variable “GOOGLE_APPLICATION_CREDENTIALS”
    Example: Set  GOOGLE_APPLICATION_CREDENTIALS = “/filepath/Service-Account-Key.json”​
4. Create/Update "appconfig.json" file and place it along with main.go or under "resources" directory. Ensure to update the "projectid" property in the file.
    {
        "application": "go-ail-examples",
        "version": "1.0",
        "description": "Configuration Properties for AIL Examples application",
        "projectid": "GCP_FIRESTORE_PROJECT_ID",
        "log": {
            "debug": true,
            "trace": true
        }
    }
