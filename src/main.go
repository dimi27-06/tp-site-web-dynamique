package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Produit struct {
	Id          int
	Nom         string
	Description string
	Prix        int
	Reduction   int
	Image       string
	Lareduc     bool
}

func (p Produit) PrixFinal() int {
	if p.Reduction > 0 {
		return p.Prix * (100 - p.Reduction) / 100
	}
	return p.Prix
}

var (
	produits = []Produit{
		{Id: 1, Nom: " PULL A CAPUCHE  ", Description: "Pull  confortable", Prix: 129, Reduction: 20, Image: "/static/img/products/19A.webp", Lareduc: true},
		{Id: 2, Nom: " PULL BLEU MARINE", Description: "Pull ", Prix: 119, Reduction: 10, Image: "/static/img/products/21A.webp", Lareduc: true},
		{Id: 3, Nom: " PULL NOIR", Description: "Pull noir classique", Prix: 99, Reduction: 0, Image: "/static/img/products/22A.webp", Lareduc: false},
		{Id: 4, Nom: " PULL JAUNE VERT", Description: "Hoodie vert jaune", Prix: 139, Reduction: 15, Image: "/static/img/products/16A.webp", Lareduc: true},
		{Id: 5, Nom: " PANTALON  JEAN S", Description: "Jean", Prix: 149, Reduction: 5, Image: "/static/img/products/34B.webp", Lareduc: true},
		{Id: 6, Nom: " PANTALON CARGO ", Description: "Cargo", Prix: 199, Reduction: 25, Image: "/static/img/products/33B.webp", Lareduc: true},
	}
	nextID = 7
)

func main() {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "home", produits); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/produit", func(w http.ResponseWriter, r *http.Request) {
		idProduit := r.FormValue("id")
		produitId, err := strconv.Atoi(idProduit)
		if err != nil {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		for _, product := range produits {
			if product.Id == produitId {
				if err := temp.ExecuteTemplate(w, "produit", product); err != nil {
					http.Error(w, "Erreur Template produit", http.StatusInternalServerError)
				}
				return
			}
		}

	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			nom := r.FormValue("nom")
			description := r.FormValue("description")
			prixStr := r.FormValue("prix")

			if nom == "" || description == "" || prixStr == "" {
				http.Error(w, "Champs manquants", http.StatusBadRequest)
				return
			}

			prix, err := strconv.ParseInt(prixStr, 10, 64)
			if err != nil {
				http.Error(w, "Prix invalide", http.StatusBadRequest)
				return
			}

			produit := Produit{
				Id:          nextID,
				Nom:         nom,
				Description: description,
				Prix:        int(prix),
			}

			produits = append(produits, produit)
			nextID++

			http.Redirect(w, r, fmt.Sprintf("/produit?id=%d", produit.Id), http.StatusSeeOther)
			return
		}

		if err := temp.ExecuteTemplate(w, "add", nil); err != nil {
			http.Error(w, "Erreur Template add", http.StatusInternalServerError)
		}
	})

	fileServer := http.FileServer(http.Dir("./../assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}
