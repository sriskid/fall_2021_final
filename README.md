# Visualizing A Cell and Cellular Processes

This is a project designed to show a cell, the organelles with a small description, as well as a game to show the cellular processes.

## Installation

Make sure you have SDL2 binding for Go installed as well as its supplemental packages SDL_ttf and SDL_mix. 
Visit: https://github.com/veandco/go-sdl2 and view installation.

## Usage
By running the main.go of the final_project folder, you should be able to use the program. It is entirely UI based and so there is no need for data input. Click around and have fun.

## Purpose and Changes
The purpose of this project was to use graphics to show the concept of the cell for educational purposes. The original purpose was to show the different parts of the cell and then show an animation for different cellular processes. However, I realized that animation required using different sprites and I was unsure how to get the sprites and whether the amount of time for the project was going to be devoted to drawing rather than coding. Therefore, I changed the project from animation to an RPG-like game.

## Game
You have the option to play a game depicting protein synthesis (transcription to translation). However, creating a game is a difficult task and therefore I was only able to show one process. However, the structure of the code allows for the possibility of including more cellular processes. In addition, the game has many bugs in it, but none that make it unplayable. The game is playable as the player can move around, pick up items (proteins, transcription factors, enzymes) that have descriptions to them. The player will have directions as to what to do and how what they are doing relates to the cellular process. Although there are a lot of changes that can be made, it is still a good task for the semester. 

## Playing the Game
You can move your player around using the arrow keys. You will see various items laid around representing either a description or a protein or something. There is a screen in the bottom left that will tell you what the item is when you pick it up. To pick up an item stand on top of it, and either press T or on the bottom right corner it will show you the item and then you can click on it. This will put the item into the inventory, which you can access using I. Furthermore, I have implemented a codon table that is useful for the translation step, which can be accessed by pressing C. Picking up a book will give you directions as to what to do or about the cellular process and will be displayed at the top. If you step on a staircase you will move to a different level/part of the cell. To go through the door, face them and press the arrow key into them to open the door. To drop an item, open your inventory and drag-click and item out of the inventory box and it will appear at the position of the player. To exit the game, you have to close the window. 

## Further additions

Many further additions can be made. Animations could be added to have the player interact using the genome sequences, proteins, and other molecules in their "inventory". Furthermore, the cellular processes (particularly translation) can be better explained and animated. Additional cellular processes could be applied with a different map and different items/directions. Some more user interactions rather than picking up items and reading through item descriptions would be nice, however I realized it was beyond my capabilities. With regards to the cell model, the pictures could have better resolution and better quality though it is still readable. 

## Concluding thoughts
The task I proposed, and the project submitted are different in that I proposed to use animation to show cellular processes, but instead I ended up making a game to show a cellular process that could be expanded to fit more. The game is fully playable without issues and offers some amount of fun and learning experience. While it is not what I originally proposed, and while there is a lot I can go better with the program, I still think it is a good project. Thank you very much and enjoy!

