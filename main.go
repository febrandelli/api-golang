package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func pegarId(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)

	id, err := vars["livroId"]

	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return 999, errors.New("Id formato incorreto")
	}

	return strconv.Atoi(id)
}

func listarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(livros)
}

func cadastrarLivro(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var novoLivro Livro
	json.Unmarshal(body, &novoLivro)
	novoLivro.Id = len(livros) + 1

	livros = append(livros, novoLivro)

	encoder := json.NewEncoder(w)
	encoder.Encode(novoLivro)
}

func apagarLivro(w http.ResponseWriter, r *http.Request) {
	id, err := pegarId(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indiceLivro := -1
	for indice, livro := range livros {
		if livro.Id == id {
			indiceLivro = indice
		}
	}

	livros = append(livros[0:indiceLivro], livros[indiceLivro+1:]...)

	w.WriteHeader(http.StatusNoContent)
}

func modificarLivro(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id, errId := pegarId(w, r)

	if errId != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indiceLivro := -1
	for indice, livro := range livros {
		if livro.Id == id {
			indiceLivro = indice
		}
	}

	if indiceLivro < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, errBody := ioutil.ReadAll(r.Body)
	if errBody != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var livroModificado Livro
	erroJson := json.Unmarshal(body, &livroModificado)
	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	livroModificado.Id = id

	livros[indiceLivro] = livroModificado

	json.NewEncoder(w).Encode(livros[indiceLivro])
}

func buscarLivro(w http.ResponseWriter, r *http.Request) {

	id, err := pegarId(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, livro := range livros {
		if livro.Id == id {
			json.NewEncoder(w).Encode(livro)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func configHandle(roteador *mux.Router) {
	roteador.HandleFunc("/livros", listarLivros).Methods("GET")
	roteador.HandleFunc("/livros/{livroId}", buscarLivro).Methods("GET")
	roteador.HandleFunc("/livros", cadastrarLivro).Methods("POST")
	roteador.HandleFunc("/livros/{livroId}", apagarLivro).Methods("DELETE")
	roteador.HandleFunc("/livros/{livroId}", modificarLivro).Methods("PUT")
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func configServer() {
	port := "1337"
	roteador := mux.NewRouter().StrictSlash(true)
	roteador.Use(jsonMiddleware)
	configHandle(roteador)

	fmt.Println("Rodando na porta " + port)
	log.Fatal(http.ListenAndServe(":"+port, roteador))
}

func main() {
	configServer()
}
