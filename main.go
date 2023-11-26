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

// Hangman représente l'état du jeu Hangman
type Hangman struct {
	Message     string // Message affiché à l'utilisateur
	Display     string // Affichage du mot avec les lettres devinées et les espaces
	LettersUsed string // Lettres déjà utilisées par le joueur
	TriesLeft   int    // Nombre de tentatives restantes
	GameOver    bool   // Indique si le jeu est terminé
	Word        string // Mot à deviner
	Gameswon    int    // Nombre de parties gagnées
}

var hangman Hangman // Variable pour suivre l'état du jeu actuel
var isPlaying bool  // Indique si une partie est en cours
var win bool        // Indique si le joueur a gagné la partie actuelle

// NewGame initialise une nouvelle partie Hangman
func NewGame(wordList []string) Hangman {
	rand.Seed(time.Now().UnixNano())
	word := wordList[rand.Intn(len(wordList))]

	return Hangman{
		Message:     "Bienvenue dans Hangman !",
		Display:     strings.Repeat("_ ", len(word)),
		LettersUsed: "",
		TriesLeft:   10,
		GameOver:    false,
		Word:        word,
		Gameswon:    0,
	}
}

// verif vérifie l'état actuel du jeu (gagné, perdu ou deviner à nouveau)
func (game *Hangman) verif() string {
	if game.TriesLeft < 1 {
		return "lose"
	} else if !strings.Contains(game.Display, "_") {
		return "win"
	} else {
		return "guess again"
	}
}

// Guess prend en charge la tentative du joueur pour deviner une lettre
func (game *Hangman) Guess(guess string) {
	// Convertie les lettres essayer en majuscules
	guess = strings.ToUpper(guess)
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

	game.LettersUsed += guess
	if !game.updateDisplay([]rune(guess)[0]) {
		game.TriesLeft--
	}
	game.GameOver = game.TriesLeft < 1

	switch game.verif() {
	case "lose":
		game.Message = "Dommage, vous avez perdu !"
		isPlaying = false
		win = false
	case "win":
		game.Message = "Bravo, vous avez gagné !"
		isPlaying = false
		win = true
		game.Gameswon++
	case "guess again":
		game.Message = "Essayez encore !"
	}
}

// isLetter vérifie si le caractère est une lettre
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// updateDisplay met à jour l'affichage du mot en fonction de la lettre devinée
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

// ReadWords lit les mots depuis un fichier spécifié
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

func main() {
	// Serveur de fichiers statiques depuis le répertoire "static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("difficulty").ParseFiles("templates/difficulty.html"))
		tmpl.Execute(w, nil)
	})

	// Initialiser une nouvelle partie
	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		liste := r.FormValue("difficulty")
		// Charger la liste de mots depuis un fichier
		wordList, err := ReadWords("wordlists/" + liste + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		// Initialiser une nouvelle partie
		hangman = NewGame(wordList)
		isPlaying = true
		// Redirection vers la page du jeu
		http.Redirect(w, r, "/hangman", http.StatusMovedPermanently)
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
			// Redirection vers la page de victoire ou de défaite
			if win {
				http.Redirect(w, r, "/win", http.StatusMovedPermanently)
				return
			} else {
				http.Redirect(w, r, "/lose", http.StatusMovedPermanently)
				return
			}
		}
		tmpl.Execute(w, hangman)
	})

	// Gérer les tentatives du joueur pour deviner une lettre
	http.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		guess := r.FormValue("letter")
		if len(guess) != 1 {
			http.Error(w, "Veuillez fournir une lettre.", http.StatusBadRequest)
			return
		}

		// Gérer la tentative du joueur
		hangman.Guess(guess)

		// Redirection vers la page principale du jeu
		http.Redirect(w, r, "/hangman", http.StatusSeeOther)
	})

	// Démarrer le serveur web
	log.Fatal(http.ListenAndServe(":8080", nil))
}
