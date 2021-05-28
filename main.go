package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Livro struct {
	Id     int    `json:"id"`
	Title  string `json:"titulo"`
	Author string `json:"autor"`
}

var livros []Livro = []Livro{
	Livro{
		Id:     1,
		Title:  "Codigo Limpo",
		Author: "Robert C Martin",
	},
	Livro{
		Id:     2,
		Title:  "Desenvolvedor Limpo",
		Author: "Robert C Martin",
	},
	Livro{
		Id:     3,
		Title:  "Dom Casmurro",
		Author: "Machado de Assis",
	},
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bem vindo")
}

func listarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(livros)
}

func cadastrarLivro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}

	var novoLivro Livro
	json.Unmarshal(body, &novoLivro)
	novoLivro.Id = len(livros) + 1

	livros = append(livros, novoLivro)

	encoder := json.NewEncoder(w)
	encoder.Encode(novoLivro)
}

func routeLivros(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		listarLivros(w, r)
	} else if r.Method == "POST" {
		cadastrarLivro(w, r)
	}
}

func configHandle() {
	http.HandleFunc("/", getHello)
	http.HandleFunc("/livros", routeLivros)
}

func configServer() {
	port := "1337"
	configHandle()

	fmt.Println("Rodando na porta " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	configServer()
}
