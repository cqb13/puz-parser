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
