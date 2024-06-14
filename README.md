# Snake Go

## Description du projet

Le projet consiste en un Snake en Golang comportant un menu disposant de plusieurs options pour changer les paramètres de la partie. Il y a également une musique de fond, des bruitages et l'affichage des meilleurs scores

## Equipe

- Florent Weltmann
- Dantin Durand

## Commande pour lancer le projet

- Rendez-vous dans le répertoire du projet.
- Lancer la commande `go run .`. Le programme va se lancer !

## Les menus

Au lancement du jeu, différentes options vous seront proposés :

- Commencer le jeu, accéder aux crédits ou quitter le jeu.
- Ensuite en commençant le jeu, vous devrez entrer votre nom pour garder une trace des meilleurs scores.
- Vous pourrez ensuite choisir le mode de jeu : Le mode classique (1 vie, pas d'obstacle) ou bien le mode challenge (plusieurs vies, des obstacles)
- Vous pourrez ensuite choisir la difficulté : Qui change la vitesse du snake selon la difficulté (plus le niveau de difficulté est facile, plus le snake sera lent au début), et si vous êtes en mode challenge changera également le nombre de vies et d'obstacles.

## Le jeu

Vous devrez manger de la nourriture pour que le serpent grossisse, au fur et à mesure, le serpent ira de plus en plus vite et il ne devra pas se rentrer dedans ni toucher un obstacle ni toucher les murs sinon il perdra une vie ou si il n'en a plus Game Over.

## Stack

- Golang
- ebiten

## Musiques utilisées

Music: "8 Bit Adventure" By HeatleyBros
https://www.youtube.com/watch?v=Wsw-86zjb8I
