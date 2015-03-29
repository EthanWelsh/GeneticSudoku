# GeneticSudoku

A genetic algorithm is an algorithm which tries to mimick the process of natural selection in order to solve a problem. Like in nature, each solution is made up of several genes. These genes are mated with one another and undergo mutations in an ongoing evolution process until a solution has been reached.

This algorithm starts by populating a 'population' of random solutions to the given board. From there, each of these solutions is graded on a set of criteria such as the number of placed tiles and how many completed rows, columns, and boxes are in each solution. After the grades have been calculated, the algorithm uses these grades to 'mate' two of the solutions together. The resulting solution of the mating process will have some of the genes from each of its parents. Genes with higher grades are more likely to be selected for mating. Lastly, the new population is subjected to random mutations which will very rarely change random genes within each solution. 

In order to solve the sudoku puzzle, I chose to represent each row in the board as a permutation of {1, 2, 3, 4, 5, 6, 7, 8, 9}, taking care to preserve the given values in each row. Crossover will cause two boards to exchange a row with one another. Mutations will cause a random row within the board to be set to another permutation.

## Results

Sudoku is a considerably hard task to solve using a genetic algorithm. This is mostly due to the presence of countless numbers of false local maxima which cause the algorithm to arrive upon a good, albeit incorrect solution. In order to solve this issue, I tried experimenting with the crossover, mutation rates in order to maintain population diversity in the name of converging more slowly and by that means avoiding local maxima. The result can solve even the hardest 4x4 boards in very few generations without trouble, but struggles with all but the easiest 9x9 boards. In the future, I hope to implement tournament selection in order to improve the algorithm's performance on larger boards.

*** Thanks to Sean Myers (see [iph](https://github.com/iph)) for donating his Board class for me to use! ***
