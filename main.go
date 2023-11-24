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

//	func NewGame(word string) Hangman {
//		return Hangman{
//			Message:     "Welcome back to Hangman!",
//			Display:     updateDisplay("", word),
//			LettersUsed: "",
//			TriesLeft:   6,
//			GameOver:    false,
//			Word:        word,
//		}
//	}
func NewGame(wordList []string) Hangman {
	rand.Seed(time.Now().UnixNano())
	word := wordList[rand.Intn(len(wordList))]

	return Hangman{
		Message:     "Welcome back to Hangman!",
		Display:     strings.Repeat("_ ", len(word)),
		LettersUsed: "",
		TriesLeft:   10,
		GameOver:    false,
		Word:        word,
		Gameswon:    0,
	}
}
func (game *Hangman) verif() string {
	if game.TriesLeft < 1 {
		return "lose"
	} else if !strings.Contains(game.Display, "_") {
		return "win"
	} else {
		return "guess again"
	}
}
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

	game.LettersUsed += guess
	if !game.updateDisplay([]rune(guess)[0]) {
		game.TriesLeft--
	}
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

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

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
	// Serve static files from the "static" directory.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("difficulty").ParseFiles("templates/difficulty.html"))
		tmpl.Execute(w, nil)
	})
	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		liste := r.FormValue("difficulty")
		// Load word list from file
		wordList, err := ReadWords("wordlists/" + liste + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		hangman = NewGame(wordList)
		isPlaying = true
		http.Redirect(w, r, "/hangman", 301)
	})
	http.HandleFunc("/lose", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("lose").ParseFiles("templates/lose.html"))
		tmpl.ExecuteTemplate(w, "lose", hangman)
	})
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("win").ParseFiles("templates/win.html"))
		tmpl.ExecuteTemplate(w, "win", hangman)
	})
	// Start web server
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
		//test
	})
	//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		tmpl := template.Must(template.ParseFiles("index.html"))
	//		hangman := NewGame("hangman")
	//		tmpl.Execute(w, hangman)
	//	})

	http.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		guess := r.FormValue("letter")
		if len(guess) != 1 {
			http.Error(w, "Please provide a letter.", http.StatusBadRequest)
			return
		}

		hangman.Guess(guess)

		http.Redirect(w, r, "/hangman", http.StatusSeeOther)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
