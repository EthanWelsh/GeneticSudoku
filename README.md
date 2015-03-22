# GeneticSudoku
Basic genetic algorithm written in GoLang to solve Sudoku puzzles!

A genetic algorithm is an algorithm which tries to mimick the process of natural selection in order to solve a problem. Like in nature, each solution is made up of several genes. These genes are mated with one another and undergo mutations in an ongoing evolution process until a solution has been reached.

This algorithm starts by populating a 'population' of random solutions to the given board. From there, each of these solutions is graded on a set of criteria such as the number of placed tiles and how many completed rows, columns, and boxes are in each solution. After the grades have been calculated, the algorithm uses these grades to 'mate' two of the solutions together. The resulting solution of the mating process will have some of the genes from each of its parents. Genes with higher grades are more likely to be selected for mating. Lastly, the new population is subjected to random mutations which will very rarely change random genes within each solution.

* Thanks to Sean (see iph) for donating his Board class for me to use!
