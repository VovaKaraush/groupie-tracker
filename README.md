# Groupie-tracker

> Petite application en Go pour suivre les artistes (projet pédagogique).

## Description

Groupie-tracker est une application backend minimale écrite en Go qui sert des pages web et expose une API interne pour récupérer des informations d'artistes. Ce README présente l'installation, l'exécution et les bonnes pratiques de développement.

## Technologies

- Go (module Go, layout simple)
- HTML/CSS/JS pour l'interface web (dossier `web`)

## Prérequis

- Go 1.18+ installé

## Installation

1. Cloner le dépôt :

   git clone <https://github.com/VovaKaraush/groupie-tracker.git>
   cd groupie-tracker

2. Récupérer les dépendances (si nécessaire) :

   go mod tidy

## Exécution

Pour lancer l'application en mode développement :

   go run main.go
   ou
   go run .
   ou
   lancer l'éxécutable

Pour compiler :

   go build -o groupie-tracker
   ./groupie-tracker

## Structure du projet

- `main.go` : point d'entrée de l'application
- `handlers/handlers.go` : gestionnaires HTTP
- `modules/` : logique applicative et points d'accès API
- `filters/` : fonctions utilitaires et filtres
- `web/` : fichiers statiques (HTML, CSS, JS)
- `models.go` : définitions des structures métier

Conserver une séparation claire entre la logique HTTP (handlers), la logique métier (modules) et les utilitaires (filters) facilite les tests et la maintenance.

## Bonnes pratiques recommandées

- Formater le code : `gofmt` / `go fmt`
- Vérifier les problèmes : `go vet`
- Utiliser un linter (ex. `golangci-lint`) pour des règles cohérentes
- Gérer les erreurs explicitement et les propager quand nécessaire
- Utiliser `context.Context` pour les requêtes et les opérations longues
- Tests unitaires : écrire des tests pour la logique métier (dossier `modules`)
- Ne pas exposer directement les dépendances globales (préférer l'injection)

## Développement

- Ajouter des routes ou handlers : modifier `handlers/handlers.go`
- Ajouter la logique métier : modifier `modules/api.go`
- Rafraîchir les assets web dans `web/`

Exemple d'exécution locale rapide :

   go run main.go

## Tests

Si des tests existent, exécuter :

   go test ./...

## Contribuer

- Ouvrir une issue pour discuter d'une fonctionnalité
- Soumettre des pull requests petites et ciblées
- Respecter les conventions de code et les tests

## Licence

Ce projet n'inclut pas encore de licence explicite.

## Contact

Voir les membres du projet Groupie-Tracker sur GitHub.
