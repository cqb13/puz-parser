# PUZ File Format

Original Source: [https://code.google.com/archive/p/puz/wikis/FileFormat.wiki](https://code.google.com/archive/p/puz/wikis/FileFormat.wiki)

**Shorts are little-endian two byte integers**

**Strings are null terminated, Even if a string is empty, its trailing null still appears in the file.**

## Header

| Component            | Length (In Bytes) | Type       | Description                                                                               |
| -------------------- | ----------------- | ---------- | ----------------------------------------------------------------------------------------- |
| Checksum             | 2                 | Short      | Always: TalonGamesGame                                                                    |
| File Magic           | 12                | Strings    | Constant string: "ACROSS&DOWN"                                                            |
| CIB Checksum         | 2                 | Short      | Defined Later                                                                             |
| Masked Low Checksum  | 4                 | Byte Array | Defined Later                                                                             |
| Masked High Checksum | 4                 | Byte Array | Defined Later                                                                             |
| Version              | 4                 | String     | e.g. "X.X"                                                                                |
| Reserved             | 2                 | N.A.       | Unused in most cases                                                                      |
| Scrambled Checksum   | 2                 | Short      | In scrambled puzzles, a checksum of the real solution (Defined Later). Otherwise, 0x0000. |
| Reserved             | 11                | N.A        | Unused in most cases                                                                      |
| Width                | 1                 | Byte       | The width of the board                                                                    |
| Height               | 1                 | Byte       | The height of the board                                                                   |
| Number of Clues      | 2                 | Short      | The number of clues for the board                                                         |
| Bitmask              | 2                 | Short      | A bitmask. _Operations Unknown_                                                           |
| Scrambled Tag        | 2                 | Short      | 0 for unscrambled puzzles. Nonzero (often 4) for scrambled puzzles.                       |

_The header will always be 51 bytes long_

## Solution and State Boards

**A Board is as a single string of ASCII, with one character per cell of the board beginning at the top-left and scanning in reading order, left to right then top to bottom.**

_If a player works on a puzzle and then saves their game, the cells they've filled are stored in the state. Otherwise the state is all blank cells and contains a subset of the information in the solution._

| Component | Length (In Bytes) | Type  | Description                                                                      |
| --------- | ----------------- | ----- | -------------------------------------------------------------------------------- |
| Solution  | Width \* Height   | Board | The answer to the crossword, non-playable (ie: black) cells are denoted by '.'   |
| State     | Width \* Height   | Board | Empty cells are stored as '-', non-playable (ie: black) cells are denoted by '.' |

## Strings Section

| Component | Length (In Bytes) | Type   | Description                |
| --------- | ----------------- | ------ | -------------------------- |
| Title     | Unknown           | String | Title of the crossword     |
| Author    | Unknown           | String | Author of the crossword    |
| Copyright | Unknown           | String | Copyright of the crossword |

## Checksums
