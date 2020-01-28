package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
)

var pathIngredients string = "./datafile/ingredients.csv"
var pathRecipes string = "./datafile/recipes.csv"

// ingredients variable
type ResponseList struct {
	data [][]string
}

type ResponseIngre struct {
	Title string
	UseBy string
}

type MyArr struct {
	Items []ResponseIngre
}

func (box *MyArr) AddItem(item ResponseIngre) {
	box.Items = append(box.Items, item)
}

// recipes variable
type ResponseRecipe struct {
	title string
	ingredients []string
}

func check(e error, c *gin.Context) {
	if e != nil {
		c.JSON(400, gin.H{
			"response": 400,
			"message": "error",
		})
	}
}

func ingredients(c *gin.Context) {
	iFile, err := os.Open(pathIngredients)
	check(err, c)
	scanner := bufio.NewScanner(bufio.NewReader(iFile))
	counter := 1
	box := MyArr{}
	for scanner.Scan() {
		counter++
		line := scanner.Text()
		str := strings.Split(line, ",")
		item := ResponseIngre{
			Title: str[1],
			UseBy: str[2],
		}
		box.AddItem(item)
	}

	defer iFile.Close()

	c.JSON(200, box.Items)
}

func recipes(c *gin.Context) {
	var gNil = c.Query("ingredients")
	splits := strings.Split(gNil, ",")
	fmt.Println(gNil);
	fmt.Println(splits[0]);

	//iFileIng, errIng := os.Open(pathIngredients)
	//check(errIng, c)

	iFileRec, errRec := os.Open(pathRecipes)
	check(errRec, c)

	/*scannerIng := bufio.NewScanner(bufio.NewReader(iFileIng))
	counterIng := 0
	boxIng := [][]string
	for scannerIng.Scan() {
		lineIng := scannerIng.Text()
		fmt.Println(lineIng)
		strIng := strings.Split(lineIng, ",")
		boxIng[counterIng] = strIng
		counterIng++
	}
	fmt.Println(boxIng)*/

	scannerRec := bufio.NewScanner(bufio.NewReader(iFileRec))
	counterRec := 1
	for scannerRec.Scan() {
		counterRec++
		lineRec := scannerRec.Text()
		strRec := strings.Split(lineRec, ",")
		box := []string{}

		/*for scannerIng.Scan() {
			counterIng++
			lineIng := scannerIng.Text()
			strIng := strings.Split(lineIng, ",")

			if(strings.TrimSpace(strIng[3]) == strRec[0]) {
				box = append(box, strIng[1])
			}
		}*/

		item := ResponseRecipe{
			title: strRec[1],
			ingredients: box,
		}
		fmt.Println(item)
	}
}

func create(c *gin.Context) {
	title := c.PostForm("title")
	usedBy := c.PostForm("usedBy")
	types := c.PostForm("types")
	recId := c.PostForm("ingreId")
	var path string;

	if(types == "ingredients") {
		path = pathIngredients
	} else if(types == "recipes") {
		path = pathRecipes
	} else {
		path = ""
	}

	iFile, err := os.Open(path)
	check(err, c)
	scanner := bufio.NewScanner(bufio.NewReader(iFile))
	counter := 1
	for scanner.Scan() {
		counter++
	}
	counters := strconv.Itoa(counter)

	defer iFile.Close()

	fmt.Println(recId)

	vile := counters+", "+title+", "+usedBy+", "+recId+"\n"

	if(types == "ingredients") {
		path = pathIngredients
	} else if(types == "recipes") {
		path = pathRecipes
	}

	fileR, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	check(err, c)

	n, err := fileR.WriteString(vile)
	check(err, c)

	defer fileR.Close()

	if(n == 0) {
		c.JSON(400, gin.H{
			"response": 400,
			"message": "error",
		})
	} else {
		c.JSON(200, gin.H{
			"response": 200,
			"message": "success",
		})
	}
}

func deleted(c *gin.Context) {
	types := c.PostForm("types")
	var path string;
	if(types == "ingredients") {
		path = pathIngredients
	} else if(types == "recipes") {
		path = pathRecipes
	} else {
		path = ""
	}

	iFile, err := os.Open(path)
	check(err, c)

	reader := bufio.NewReader(iFile)

	for {
		pos,_ := iFile.Seek(0, 1)
		fmt.Println("Position in file is: %d", pos)
		bytes, _, _ := reader.ReadLine()

		if (len(bytes) == 0) {
		break
		}

		lineString := string(bytes)
		if(lineString == "two") {
             iFile.Seek(int64(-(len(lineString))), 1)
             //file.WriteString("This is a test.")
         }

         fmt.Printf(lineString + "\n")
	}

	defer iFile.Close()

	c.JSON(200, gin.H{
		"response": 200,
		"message": "success",
	})
}

func main() {
	router := gin.Default()
	router.GET("/ingredients", ingredients)
	router.GET("/recipes", recipes)
	router.POST("/create", create)
	router.POST("/delete", deleted)
	router.Run(":8181")
}
