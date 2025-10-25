# PUZ File Format

Original Source: [https://code.google.com/archive/p/puz/wikis/FileFormat.wiki](https://code.google.com/archive/p/puz/wikis/FileFormat.wiki)

**Shorts are little-endian two byte integers**

**Strings are null terminated, Even if a string is empty, its trailing null still appears in the file.**

## Header

| Component            | Length (In Bytes) | Type       | Description                                                                               |
| :------------------- | :---------------- | :--------- | :---------------------------------------------------------------------------------------- |
| Checksum             | 2                 | Short      | Overall file checksum                                                                     |
| File Magic           | 12                | Strings    | Constant string: "ACROSS&DOWN"                                                            |
| CIB Checksum         | 2                 | Short      | Refer to checksums.go for implementation                                                  |
| Masked Low Checksum  | 4                 | Byte Array | Refer to checksums.go for implementation                                                  |
| Masked High Checksum | 4                 | Byte Array | Refer to checksums.go for implementation                                                  |
| Version              | 4                 | String     | e.g. "X.X"                                                                                |
| Reserved             | 2                 | N.A.       | Unused in most cases                                                                      |
| Scrambled Checksum   | 2                 | Short      | In scrambled puzzles, a checksum of the real solution (Defined Later). Otherwise, 0x0000. |
| Reserved             | 12                | N.A        | Unused in most cases                                                                      |
| Width                | 1                 | Byte       | The width of the board                                                                    |
| Height               | 1                 | Byte       | The height of the board                                                                   |
| Number of Clues      | 2                 | Short      | The number of clues for the board                                                         |
| Bitmask              | 2                 | Short      | A bitmask. _Operations Unknown_                                                           |
| Scrambled Tag        | 2                 | Short      | 0 for unscrambled puzzles. Nonzero (often 4) for scrambled puzzles.                       |

_The header will always be 52 bytes long_

## Solution and State Boards

**A Board is as a single string of ASCII, with one character per square of the board beginning at the top-left and scanning in reading order, left to right then top to bottom.**

_If a player works on a puzzle and then saves their game, the squares they've filled are stored in the state. Otherwise the state is all blank squares and contains a subset of the information in the solution._

| Component | Length (In Bytes) | Type  | Description                                                                          |
| :-------- | :---------------- | :---- | :----------------------------------------------------------------------------------- |
| Solution  | Width \* Height   | Board | The answer to the crossword, non-playable (ie: black) squares are denoted by '.'     |
| State     | Width \* Height   | Board | Empty squares are stored as '-', non-playable (ie: black) squares are denoted by '.' |

## Strings Section

**Clues are arranged numerically. When two clues have the same number, the Across clue comes before the Down clue.**

| Component | Length (In Bytes) | Type         | Description                |
| :-------- | :---------------- | :----------- | :------------------------- |
| Title     | Variable          | String       | Title of the crossword     |
| Author    | Variable          | String       | Author of the crossword    |
| Copyright | Variable          | String       | Copyright of the crossword |
| Clues     | Variable          | String Array | List of clues              |
| Notes     | Variable          | String       | Notes for the crossword    |

## Extra Sections

| Section Name | Description                                    |
| :----------- | :--------------------------------------------- |
| GRBS         | where rebuses are located in the solution      |
| RTBL         | contents of rebus squares, referred to by GRBS |
| LTIM         | timer data                                     |
| GEXT         | circled squares, incorrect and given flags     |
| RUSR         | user-entered rebus squares                     |

In official puzzles, the sections always seem to come in this order, when they appear. It is not known if the ordering is guaranteed. The GRBS and RTBL sections appear together in puzzles with rebuses. However, sometimes a GRBS section with no rebus squares appears without an RTBL, especially in puzzles that have additional extra sections.

The extra sections all follow the same general format, with variation in the data they contain.

| Component | Length (In Bytes) | Type                        | Description                                                                                    |
| :-------- | :---------------- | :-------------------------- | :--------------------------------------------------------------------------------------------- |
| Title     | 4                 | String (No null terminator) | The name of the section, these are given in the previous table                                 |
| Length    | 2                 | Short                       | The length of the data section, in bytes, not counting the null terminator                     |
| Checksum  | 2                 | Short                       | A checksum of the data section                                                                 |
| Data      | Variable          | Variable                    | The data, which varies in format but is always terminated by null and has the specified length |

### Rebus Data (GRBS)

**A Rebus Board is as a list of bytes with one byte per square of the board beginning at the top-left and scanning in reading order, left to right then top to bottom.**
The Byte for each square of the board indicate whether or not the square is a rebus.

- **0** Indicates a non-rebus square
- **1 + i** Indicates a rebus square, the solution of which is given by the key **i** in the RTBL section. A number may appear multiple times if there are multiple rebus squares wit the same solution.

_If a square is a rebus, only the first letter will be given by the solution board and only the first letter of any fill will be given in the user state board._

### Rebus Table Data (RTBL)

Contains the solutions for any rebus squares

Each solution is an ASCII string made up of a number, a colon, a string, and a semicolon.

The number is the key that the GRBS section uses to refer to a solution (it is one less than the number that appears in the corresponding rebus square). The number is always two characters long, if it is only one digit, the first character is a space.

**Example:** " 0:HEART; 1:DIAMOND;17:CLUB;23:SPADE;"

_The keys do not need to be in a consecutive order, but in official puzzles they are often in ascending order._

### Timer Data (LTIM)

Contains the amount of time the solver has used, and whether the timer is running or stopped

The data is stored as an ASCII string separated by a comma. First comes the number of seconds elapsed, then "0" if the timer is running and "1" if it is stopped.

**Example:** If the timer were stopped at 42 seconds when the puzzle was saved, the LTIM data section would contain: "42,1"

### Style Atributes (GEXT)

**The GEXT data section is another "board" of one byte per square.**

Each byte is a bitmask indicating that some style attributes are set

| Byte | Meaning                                    |
| :--- | :----------------------------------------- |
| 0x10 | The square was previously marked incorrect |
| 0x20 | The square is currently marked incorrect   |
| 0x40 | The contents were given                    |
| 0x80 | The square is circled                      |

None, some, or all of these bits may be set for each square. It is possible that they have reserved other values.

### User Rebus Data (RUSR)

Users rebus entries, uses the same format as RTBL
