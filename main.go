package main

//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "bufio"
    "time"
    "sync"
    "os"
    "mvdan.cc/xurls"
)


var(
    wg          sync.WaitGroup
    FromFile    bool
    input       string
    ManualURL   string
    FileName    string
    LINKS       []string
)

func ExtractURLSFromFile(filename string) ([]string){
    file, err := os.Open(filename)
    defer file.Close()
    if err != nil{
        fmt.Println(err)
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)

    var URLS []string


    for scanner.Scan(){
        URLS = append(URLS,scanner.Text())
    }

    return URLS
}


func CheckIfcontainsString(ARRAY *[]string, str string) bool{
    for _, i := range *ARRAY{
        if i == str { return true }
    }
    return false
}


func FindURLS(URL string, LIST *[]string ){
    if FromFile == true{
        defer wg.Done()
    }

    resp,err := http.Get(URL)
    defer resp.Body.Close()

    if err != nil {
        fmt.Println(err)
        wg.Done()
    }

    body,_ := ioutil.ReadAll(resp.Body)

    rxRelaxed := xurls.Strict()
    links := rxRelaxed.FindAllString(string(body),-1)

    var counter uint64 = 0
    for _,link := range links{
        if CheckIfcontainsString(LIST,link) != true{
                *LIST = append(*LIST,link)
                counter++
            }
        }
    if counter != 0{
        fmt.Printf("Found %v unique link(s)\n",counter)
    }

}



func SaveStringArrayToFile(filename string, ARRAY *[]string){
    f, err := os.Create(filename)
    defer f.Close()
    if err != nil {
        fmt.Println(err)
        f.Close()
    }

    for _,data := range *ARRAY{
        _,err := f.WriteString(data + "\n")
        if err != nil{
            fmt.Println(err)
        }
    }
}



func main(){

    fmt.Println(
        "  _        _         _                      \n",
        "| |      (_)       | |                     \n",
        "| |       _  ____  | |  _   ____   ____    \n",
        "| |      | ||  _ \\ | | / ) / _  ) / ___)  \n",
        "| |_____ | || | | || |< ( ( (/ / | |       \n",
        "|_______)|_||_| |_||_| \\_) \\____)|_|    \n\n")

    LinksPointer := &LINKS

    // input
    fmt.Printf("Do you want to read URLs from file ? (y/n) : ")
    fmt.Scanln(&input)
    if input == "y" || input == "Y"{
        FromFile = true
    }else if input == "n" || input == "N"{
        FromFile = false
    }else{
        FromFile = false
        fmt.Println("Invalid input. Not reading from file")
    }

    if FromFile == false{ //Serve without file
        fmt.Println("Enter URL : ")
        fmt.Scanln(&ManualURL)

        t0 := time.Now()

        FindURLS(ManualURL,LinksPointer)
        fmt.Printf("\n\n%v links in total \n\n",len(LINKS))

        if len(LINKS) != 0{
            SaveStringArrayToFile("Output.txt",LinksPointer)
            fmt.Println("Saved as \"Output.txt\"")
        }

        t1 := time.Now()
        fmt.Printf("\nTook %v ",t1.Sub(t0))

    }else if FromFile == true{ //Serve with a file
        fmt.Println("Filename : ")
        fmt.Scanln(&FileName)
        URLS := ExtractURLSFromFile(FileName)
        t0 := time.Now()


        for _,link := range URLS{
            wg.Add(1)
            go FindURLS(link,LinksPointer)
            time.Sleep(time.Millisecond * 250) // A reasonable time for server resting
        }
        wg.Wait()

        fmt.Printf("\n\n%v links in total \n\n",len(LINKS))
        if len(LINKS) != 0{
            SaveStringArrayToFile("Output.txt",LinksPointer)
            fmt.Println("Saved as \"Output.txt\"")
        }



        t1 := time.Now()
        fmt.Printf("\nTook %v ",t1.Sub(t0))


    }


    fmt.Scanln()
}
