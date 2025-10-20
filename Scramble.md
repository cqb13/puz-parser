# Scramble

Big thanks to Brian Raiter for [figuring out](https://www.muppetlabs.com/~breadbox/txt/acre.html) this [algorithm](https://www.muppetlabs.com/~breadbox/txt/scramble-c.txt).

## Scrambling Algorithm

1. The letters are copied from the solution into a buffer. The solution is read column-wise, i.e. going from top to bottom, then from left to right. Black squares '.' are skipped over, so that the buffer contains only letters A through Z.

2. The scrambled checksum is calculated at this point

3. The letters in the buffer are replaced with numbers in the range 0 to 25, inclusive (with A becoming 0, B becoming 1, etc).

4. The buffer contents are then arranged into a (notional) table 16 columns wide. The table is filled column-wise, but starting with the rightmost column (i.e. going from top to bottom, then right to left).

5. Successive digits of the key are added to the letters (mod 26), one digit per letter, moving column-wise through the table. Rows are then shifted from the top of the table to the end, the number of rows being equal to the first digit of the key. If the number of letters in the buffer is even, then each row shifted is also individually rotated right, with the rightmost cell shifting around to the leftmost.

6. The notional table is changed from 16 columns to 8 columns (without actually altering the buffer contents), and step 5 is then repeated, with the second digit of the key being used instead of the first in the row shift.

7. Step 5 is repeated two more times, with a 4-column table and then a 2-column table, and with the remaining two digits of the key controlling the row shift.

8. The numbers in the buffer are turned back into letters.

9. The buffer is copied back into the solution, in the same order in which it was originally copied out (top to bottom, left to right, skipping over black squares).

## Unscrambling Algorithm

1. **Copy the letters into a buffer**
   The letters are copied from the puzzle’s fill (solution) into a buffer, column by column (top to bottom, left to right), skipping over black squares (`.`).
   The buffer should contain only uppercase letters A–Z.

2. **Convert letters to numbers**
   Each letter is replaced by a number from 0 to 25 (`A` = 0, `B` = 1, etc).

3. **Reverse the four rounds of scrambling**
   The unscrambling process repeats four rounds — one for each digit of the key, starting from the **last digit** and moving backward.

   For each round `k` (where `k = 3, 2, 1, 0`):
   - Compute the “column width” as `n = 1 << (4 - k)` (i.e., 16, 8, 4, and 2 for successive rounds).
   - If `n` is larger than the total number of letters, reduce it slightly so that it stays within range.

   Then, two operations are performed in each round:

   **a. Undo the row shifts**
   - The buffer is rotated in the _opposite direction_ of the scrambling phase.
   - The number of rotations equals the value of the current key digit (`key[k]`).
   - If the total number of letters in the buffer is even, each shifted block is also “rotated left” (undoing the per-row right rotation from scrambling).

   **b. Undo the per-letter shifts**
   - Each number in the buffer has the corresponding key digit _subtracted_ (mod 26).
   - The same key digits are repeated in order across the buffer (i.e., `key[i % 4]`).

4. **Undo the column rearrangement**
   - The scrambled buffer was originally rearranged as if it had been written into a 16-column table, filled from right to left.
   - This step reverses that mapping, restoring the original column order.

5. **Convert numbers back to letters**
   - Each number (0–25) is converted back to a letter (`A` + value).

6. **Verify checksum**
   - Compute the checksum of the unscrambled buffer.
   - If it doesn’t match the stored fill checksum, the provided key is incorrect.

7. **Copy the buffer back into the puzzle grid**
   - The letters are written back into the puzzle in the same order they were originally read (top to bottom, left to right, skipping over black squares).
   - Mark the puzzle as unscrambled.
