package main

import (
	"bufio"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Hangman représente l'état du jeu.
type Hangman struct {
	Message     string
	Display     string
	LettersUsed string
	TriesLeft   int
	GameOver    bool
	Word        string
	Gameswon    int
}

var hangman Hangman
var isPlaying bool
var win bool

// NewGame initialise une nouvelle partie avec un mot aléatoire.
func NewGame(wordList []string) Hangman {
	rand.Seed(time.Now().UnixNano())
	word := wordList[rand.Intn(len(wordList))]

	return Hangman{
		Message:     "Bienvenue dans le Pendu!",
		Display:     strings.Repeat("_ ", len(word)),
		LettersUsed: "",
		TriesLeft:   10,
		GameOver:    false,
		Word:        word,
		Gameswon:    0,
	}
}

// verif vérifie l'état du jeu et renvoie "win", "lose" ou "guess again".
func (game *Hangman) verif() string {
	if game.TriesLeft < 1 {
		return "lose"
	} else if !strings.Contains(game.Display, "_") {
		return "win"
	} else {
		return "guess again"
	}
}

// Guess traite la supposition du joueur et met à jour l'état du jeu en conséquence.
func (game *Hangman) Guess(guess string) {
	// Vérifier si la lettre a déjà été utilisée
	if strings.Contains(game.LettersUsed, guess) {
		game.Message = "Erreur : Lettre déjà utilisée !"
		return
	}

	// Vérifier si la saisie est une lettre
	if len(guess) != 1 || !isLetter(guess[0]) {
		game.Message = "Erreur : Veuillez entrer une seule lettre !"
		return
	}

	// Mettre à jour les lettres utilisées
	game.LettersUsed += guess

	// Mettre à jour l'affichage et les essais restants
	if !game.updateDisplay([]rune(guess)[0]) {
		game.TriesLeft--
	}

	// Vérifier l'état du jeu
	game.GameOver = game.TriesLeft < 1

	switch game.verif() {
	case "lose":
		game.Message = "Dommage, vous avez Perdu!"
		isPlaying = false
		win = false
	case "win":
		game.Message = "Bravo, vous avez Gagné!"
		isPlaying = false
		win = true
		game.Gameswon++
	case "guess again":
		game.Message = "Essayez encore!"
	}
}

// isLetter vérifie si le caractère est une lettre.
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// updateDisplay met à jour l'affichage du mot masqué.
func (game *Hangman) updateDisplay(char rune) bool {
	var correct bool
	arrayDisplay := []rune(game.Display)
	for i, r := range game.Word {
		if r == char {
			arrayDisplay[i*2] = char
			correct = true
		}
	}
	game.Display = string(arrayDisplay)
	return correct
}

// ReadWords lit les mots à partir d'un fichier.
func ReadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, nil
}

// main est le point d'entrée du programme.
func main() {
	// Serve static files from the "static" directory.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("difficulty").ParseFiles("templates/difficulty.html"))
		tmpl.Execute(w, nil)
	})

	// Initialisation du jeu
	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		liste := r.FormValue("difficulty")
		// Charger la liste de mots depuis le fichier
		wordList, err := ReadWords("wordlists/" + liste + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		hangman = NewGame(wordList)
		isPlaying = true
		http.Redirect(w, r, "/hangman", 301)
	})

	// Page de défaite
	http.HandleFunc("/lose", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("lose").ParseFiles("templates/lose.html"))
		tmpl.ExecuteTemplate(w, "lose", hangman)
	})

	// Page de victoire
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("win").ParseFiles("templates/win.html"))
		tmpl.ExecuteTemplate(w, "win", hangman)
	})

	// Page principale du jeu
	http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").ParseFiles("templates/index.html"))
		if !isPlaying {
			if win {
				http.Redirect(w, r, "/win", 301)
				return
			} else {
				http.Redirect(w, r, "/lose", 301)
				return
			}
		}
		tmpl.Execute(w, hangman)
	})

	// Gestion de la devinette
	http.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		guess := r.FormValue("letter")
		if len(guess) != 1 {
			http.Error(w, "Veuillez fournir une lettre.", http.StatusBadRequest)
			return
		}

		hangman.Guess(guess)

		http.Redirect(w, r, "/hangman", http.StatusSeeOther)
	})

	// Démarrer le serveur web sur le port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

