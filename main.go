package main

import (
    "net/http"
    "./src/job_controller"
)

func main() {
    http.HandleFunc("/", job_controller.StartNewJobSet)
    http.ListenAndServe(":8000", nil)
}
