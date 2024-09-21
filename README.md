# Padel Game Organizer Bot

[![Go](https://img.shields.io/badge/Go-1.19-blue.svg)](https://golang.org)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot-API-blue)](https://core.telegram.org/bots/api)

This Telegram bot is written in Go and helps padel enthusiasts organize games effortlessly. With this bot, users can create new games, view active games, and join games with other players. It supports multiple commands for smooth game management.

## Features

- **/help**: Displays a help message with available commands and their usage.
- **/new**: Allows users to create a new game by specifying the date, time, place, and players' skill level.
- **/games**: Shows all active games in the chat and lets users join existing games.

## Table of Contents

- [Getting Started](#getting-started)
- [Commands](#commands)
- [Installation](#installation)

## Getting Started

### Prerequisites

- [Go 1.19+](https://golang.org/doc/install)
- [Telegram Bot API Token](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
- Access to a Telegram chat where the bot will be used.

### Telegram Bot Token

To get started, you'll need a Telegram bot token from [BotFather](https://t.me/BotFather). Create a new bot, and BotFather will give you the token that is required for the bot to interact with Telegram's API.

## Commands

The following commands are available for users in the chat:

### /help
Displays the help message with details on how to use the bot.

### /new
Creates a new padel game by specifying the game details. You will need to provide the following parameters:
- **Game date**: The date of the game (format: MM-DD).
- **Game time**: The time of the game (format: HH:MM).
- **Game location**: Location of the game.
- **Players level**: Skill level of the players (e.g., Beginner, Intermediate, Advanced).

### /games
Lists all active padel games in the chat. Users can see game details and join a game.

## Installation

To install and run the bot locally, follow these steps:

1. Clone the repository:
```bash
git clone https://github.com/yourusername/padel-bot.git
cd padel-bot
```
2. Create a .env file in the root of the project and add your Telegram bot token:
```
PADEL_BOT_TOKEN=your-telegram-bot-token
```
3. Install Go dependencies:
```
go mod tidy
```
5. Run the bot:
```
go run main.go
```
