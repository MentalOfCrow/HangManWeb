La Structure du Projet Sur VsCode Sera la suivant : 

/HANGMANWEBCOMPLET
│
├── /static
│   ├── style.css
│   └── images/
│   ├── 10.png
│   ├── 9.png
│   ├── 8.png
│   ├── 7.png
│   ├── 6.png
│   ├── 5.png
│   ├── 4.png
│   ├── 3.png
│   ├── 2.png
│   ├── 1.png
│   ├── 0.png
│   ├── HangMan.jpg
│   └── Win.jpg
│
├── /templates
│   ├── difficulty.html
│   ├── index.html
│   ├── lose.html
│   └── win.html
│ 
├── /wordlists
│   ├── OptionADictionnaire.txt
│   ├── OptionBFacile.txt
│   ├── OptionCMoyen.txt
│   ├── OptionDDifficile.txt
│   └── OptionEHalloween.txt
│
├── main.go
└── README.md



Projet hangman : Vous devez mettre en place une IHM orienter web dans l’objectif de permettre à la majorité des
utilisateurs de pouvoir utiliser votre jeu « Hangman CLI ». 

Attention il faudra apporter à votre
projet précèdent quelques modifications puisqu’il faudra que l’utilisateur puisse choisir une
difficulté (avec des malus ? , des bonus ? , des mots plus complexe ?, une langue étrangère ? à
vous de voir ça !), mais également lui afficher des messages de victoire, ou bien de défaite selon
le nombre de points ou bien le niveau de difficulté avec la possibilité de directement pouvoir
rejouer. 

Il faudra penser à sécuriser les différentes routes mais également les entrées utilisateur
(uniquement vérifier la validité des données entrées) pour éviter les hacker. Vous êtes fan de
cybersécurité ? Alors commençons ! Vous pouvez une fois ceci fait améliorer l’IHM en ajoutant
des fonctionnalités mais attentions à l'hors sujet !


LES CONTRAINTES : Utiliser le projet HANGMAN CLI comme module
Mettre en place un serveur HTTP et des routes écrites avec le langage Golang
Utiliser des modèles de vue
Utiliser uniquement des librairies natives à Golang
Interdiction d’utilisation de Framework HTML/JS/CSS / 


LE DÉROULÉ  LE TRAVAIL
Effectuer le projet en binôme avec le ou la partenaire de votre choix, penser à le renseigner dans le
fichier Excel présent dans le canal Teams « Classe INFO B1 ».
LES LIVRABLES
Le code source du projet, héberger sur un repository Github avec un readme rédiger avec les sections
suivantes : une présentation rapide du projet (thème, objectif…), lancement du projet (liste des
commandes à exécuter, lien à mettre dans le navigateur…), explication d’utilisation du projet.
Un fichier Execl avec la gestion de projet contenant la liste des tâches, la répartition, l’ordre dans lequel
elles ont été exécuter puis enfin un petit retour de chaque membre de l’équipe (sur les difficultés
techniques, de l’équipe…) juste quelques lignes. Le fichier est à mettre sur le Git.LES CRITÈRES
D'ÉVALUATION : Les points indiqués ci-dessous sont les éléments qui vont permettre d'évaluer votre projet
hangman-web.


Implémentation d’un serveur HTTP et
utilisation de multiples routes (minimum 4)
Utilisation de Templates avec des
conditions et des boucles
Implémentation de formulaires et gestion
des données de celui-ci sur des routes
séparées
Utilisation de fichiers statiques (CSS, image)
Utilisation de module
Organisation des fichiers et dossiers du projet
Un code source clair, optimisé, commenté et
versionné (Grâce à Git Hub)
Respect des contraintes et des fonctionnalités
demandées
Design de l’interface web (respect des
fondamentaux UX/UI)
Les rendues (répartitions des tâches et readme) / QUELQUES
PISTES ? : Mettre en place
l’environnement de
développement
Implémenter le serveur HTTP,
suivit de ses routes pour traiter
les données et distribuer les vues
Préparer les Templates
GOHTML


Tester le fonctionnement du serveur
(vérifier que les données utilisateurs
sont bien valides, les routes soient
bien distribuées et sécurisées) Intégrer le module
Hangman CLI
Tester et valider
(PS : Pour la disposition des fichiers : sous VsCode je peux pas avoir de fichiers exemple.gohtml par contre on peux exemple.html) 

Important ! 

Probleme de commit pour abdallah ( il avait casser son pc ) au niveau des commits on a faits 50% chaqu'un en terme de travaille . 

