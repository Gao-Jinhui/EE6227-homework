package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var RandSeed []int = []int{181, 637, 962, 726, 680, 199, 478, 541, 540, 266, 121, 963, 294, 798, 673, 185, 440, 91, 497, 892, 97, 200, 978, 894, 373, 74, 119, 419, 939, 214, 252, 190, 333, 27, 932, 200, 661, 4, 7, 406, 586, 314, 26, 882, 405, 426, 673, 583, 697, 116, 716, 284, 89, 817, 504, 779, 512, 606, 939, 36, 498, 247, 532, 731, 537, 503, 859, 331, 967, 891, 154, 335, 725, 834, 639, 657, 880, 922, 82, 516, 430, 954, 463, 721, 488, 533, 135, 143, 173, 964, 569, 10, 576, 800, 46, 211, 871, 619, 187, 570, 171, 947, 852, 523, 808, 12, 789, 875, 852, 845, 674, 874, 237, 935, 728, 899, 649, 791, 136, 516, 151, 243, 884, 276, 11, 594, 360, 998, 240, 565, 699, 471, 17, 78, 493, 268, 602, 476, 261, 997, 613, 197, 377, 654, 135, 759, 284, 110, 488, 696, 653, 480, 671, 285, 402, 743, 973, 132, 529, 172, 940, 159, 749, 135, 359, 407, 961, 210, 864, 684, 67, 778, 911, 943, 709, 793, 874, 957, 1, 881, 432, 796, 311, 127, 726, 849, 108, 938, 345, 329, 333, 678, 608, 834, 154, 532, 921, 42, 159, 169}
var Generation = 1

var CurPos = 58 //student number
var MutationThreshold = 0.2
var EvolutionNum = 3
var MaxMating = 5

type Person struct {
	Chromosome       []int
	BinChromosome    []string
	NormalChromosome []float64
	Fitness          float64
	RankFit          float64
	Generation       int
}

func CreatePerson(chromosome []int) Person {
	p := Person{Chromosome: chromosome, NormalChromosome: make([]float64, 3)}
	for pos := range p.Chromosome {
		p.NormalChromosome[pos] = float64(p.Chromosome[pos])/1023*40 - 20
		p.BinChromosome = append(p.BinChromosome, string(Dec2Bin(p.Chromosome[pos])))
	}
	p.Fitness = GetFitness(p.NormalChromosome)
	p.Generation = Generation
	return p
}

func main() {
	population := make([]Person, 10)
	rand.Seed(time.Now().Unix())
	fmt.Println("matric number: ", CurPos)
	fmt.Println("MutationThreshold: ", MutationThreshold)
	fmt.Println("EvolutionNum: ", EvolutionNum)
	fmt.Println("MaxMating: ", MaxMating)
	for pos := range population {
		gene := make([]int, 3)
		for p := range gene {
			gene[p] = rand.Intn(1024)
		}
		cur := CreatePerson(gene)
		population[pos] = cur
	}
	for i := 0; i < EvolutionNum; i++ {
		Generation++
		newPopulation := make([]Person, 10)
		sum := GetFitnessSum(population)
		preSum := []float64{}
		temp := 0.0
		for pos := range population {
			population[pos].RankFit = 1 / (population[pos].Fitness * sum)
			temp += population[pos].RankFit
			preSum = append(preSum, temp)
		}
		if i == EvolutionNum-1 {
			PrintInfo(population)
			break
		}
		PrintInfo(population)
		for mating := 0; mating < MaxMating; mating++ {
			parent1 := population[GetParent(preSum)]
			parent2 := population[GetParent(preSum)]
			child1, child2 := GetChildren(parent1, parent2)
			population = append(population, child1)
			population = append(population, child2)
			fmt.Println("mating ", mating+1)
			fmt.Println("parent 1 chromosome: ", parent1.BinChromosome)
			fmt.Println("parent 2 chromosome: ", parent2.BinChromosome)
			fmt.Println("child 1 chromosome: ", child1.BinChromosome)
			fmt.Println("child 2 chromosome: ", child2.BinChromosome)
		}
		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness < population[j].Fitness
		})
		copy(newPopulation, population[:10])
		population = newPopulation
	}
}

func GetFitness(normal []float64) float64 {
	res := -20 * math.Exp(-0.2*math.Sqrt(0.333333)*
		(math.Pow(normal[0], 2)+math.Pow(normal[1], 2)+math.Pow(normal[2], 2)))
	res -= math.Exp((0.3333333) * (math.Cos(2*math.Pi*normal[0]) + math.Cos(2*math.Pi*normal[1]) + math.Cos(2*math.Pi*normal[2])))
	res = res + 20 + math.E
	return res
}

func GetRandNum() float64 {
	l := len(RandSeed)
	p := float64(RandSeed[CurPos]) / 1000
	CurPos = (CurPos + 1) % l
	return p
}

func GetParent(preSum []float64) int {
	cur := GetRandNum()
	for pos := range preSum {
		if cur <= preSum[pos] {
			return pos
		}
	}
	return -1
}

func GetChildren(parent1, parent2 Person) (Person, Person) {
	c1, c2 := GetChromoByCrossover(parent1, parent2)
	c1 = GetChromoByMutation(c1)
	c2 = GetChromoByMutation(c2)
	return CreatePerson(c1), CreatePerson(c2)
}

func GetChromoByCrossover(p1, p2 Person) ([]int, []int) {
	crossIndex := 30 * GetRandNum()
	return CrossByIndex(p1.Chromosome, p2.Chromosome, int(crossIndex)/10)
}

func CrossByIndex(s1, s2 []int, index int) ([]int, []int) {
	res1, res2 := make([]int, 3), make([]int, 3)
	copy(res1[:index+1], s1[:index+1])
	copy(res1[index+1:], s2[index+1:])
	copy(res2[:index+1], s2[:index+1])
	copy(res2[index+1:], s2[index+1:])
	return res1, res2
}

func GetChromoByMutation(source []int) []int {
	binaryChromo := make([][]byte, 3)
	intChromo := make([]int, 3)
	for pos := range binaryChromo {
		binaryChromo[pos] = Dec2Bin(source[pos])
	}
	for pos := range binaryChromo {
		for i := range binaryChromo[pos] {
			p := GetRandNum()
			if p <= MutationThreshold {
				if binaryChromo[pos][i] == '0' {
					binaryChromo[pos][i] = '1'
				} else {
					binaryChromo[pos][i] = '0'
				}
			}
		}
	}
	for pos := range binaryChromo {
		intChromo[pos] = Str2DEC(string(binaryChromo[pos]))
	}
	return intChromo
}

func Dec2Bin(source int) []byte {
	b := strconv.FormatInt(int64(source), 2)
	return []byte(fmt.Sprintf("%010s", b))
}

func GetFitnessSum(p []Person) float64 {
	sum := 0.0
	for pos := range p {
		sum += 1 / p[pos].Fitness
	}
	return sum
}

func PrintInfo(p []Person) {
	fmt.Println(">>>>>  generation ", Generation-1, " list  <<<<<")
	for pos := range p {
		fmt.Printf("%+v\n", p[pos])
	}
	fmt.Println(">>>>>  generation ", Generation-1, " list  <<<<<")
}

func Str2DEC(s string) (num int) {
	l := len(s)
	for i := l - 1; i >= 0; i-- {
		num += (int(s[l-i-1]) & 0xf) << uint8(i)
	}
	return
}
