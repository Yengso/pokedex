# Pokedex CLI

A command-line Pokedex built in Go as part of the [boot.dev](https://boot.dev) curriculum. Explore the Pokemon world, search for Pokemon, catch them, and inspect your collection — all from your terminal.

## Motivation

This was my first Go project, built to learn the fundamentals of the language through something fun and practical. It covers structs, JSON parsing, functions, HTTP requests, and working with external APIs. The project also serves as a foundation for continued development as I grow as an aspiring DevOps engineer.

## Quick Start

**Prerequisites:** Go 1.22 or higher

```bash
git clone https://github.com/Yengso/pokedex.git
cd pokedex
go build -o pokedexcli
./pokedexcli
```

## Usage

Once running, you'll see the `Pokedex >` prompt. Available commands:

| Command | Description |
|---|---|
| `help` | Show all available commands |
| `map` | Display the next list of location areas |
| `mapb` | Display the previous list of location areas |
| `explore <area>` | List Pokemon found in a location area |
| `catch <pokemon>` | Attempt to catch a Pokemon |
| `inspect <pokemon>` | View details of a caught Pokemon |
| `pokedex` | List all Pokemon you have caught |
| `exit` | Exit the program |

**Example:**
```
Pokedex > explore pallet-town
Pokedex > catch pikachu
Pokedex > inspect pikachu
```

## Contributing

This is a personal learning project, but suggestions and feedback are welcome. Feel free to open an issue or submit a pull request.
