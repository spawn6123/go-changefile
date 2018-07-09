package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	dstpath  string
	dstfiles []string
	dstnames []string
	srcpath  string
	srcfiles []string
	srcnames []string
	fail     error
	tt       = time.Now().Format("20060102")
	tt1      = time.Now().Format("20060102150405")
)

func init() {
	if err := os.Mkdir(`./log`, os.ModePerm); err != nil {
		fmt.Println("Log folder already exists!")
	} else {
		fmt.Println("Log folder has been created!")
	}
}

func main() {
	srcpath = ckinputempty("Please enter the source file location:")
	dstpath = ckinputempty("Please enter the destination file location:")

	t := time.Now()                            //Program time start
	duration := time.Duration(5) * time.Second //5sec

	srcfiles, srcnames, fail = getfilelist(srcpath)
	if fail != nil {
		stoppro(fail, duration, 911)
	}

	dstfiles, dstnames, fail = getfilelist(dstpath)
	if fail != nil {
		stoppro(fail, duration, 911)
	}

	dochange(srcnames, srcfiles, dstnames, dstfiles)

	elapsed := time.Since(t).String() //Program execution time
	write2log("copy execution time: "+elapsed, tt, "std")
	fmt.Println("copy execution time: " + elapsed + ",換檔完成，5秒後程式自動關閉......")
	time.Sleep(duration) //waiting 5sec close
}

func readinput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return scanner.Text()
}

func write2log(a, t, fn string) {
	filename, err := os.OpenFile("./log/"+t+"_"+fn+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file open fail: %v", err)
	}
	defer filename.Close()
	if fn == "bk" {
		filename.WriteString(a + "\n")
	} else {
		log.SetOutput(filename)
		log.Println(a)
	}

}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func getfilelist(a string) ([]string, []string, error) {
	var ff, nn []string
	err := filepath.Walk(a,
		func(a string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				ff = append(ff, a)
				nn = append(nn, info.Name())
			}
			return nil
		})
	if err != nil {
		return ff, nn, err
	} else if len(nn) <= 0 {
		return ff, nn, errors.New(a + "(*There are no files in the specified folder*)")
	}
	return ff, nn, err
}

func ckinputempty(p string) string {
	valid := true
	var a string
	for valid {
		fmt.Print(p)
		a = readinput()
		if a == "" {
			fmt.Println("*Path cannot be blank*")
		} else if _, err := os.Stat(a); os.IsNotExist(err) {
			fmt.Println("*Unable to find the specified path*")
		} else {
			valid = false
		}
	}
	return a
}

func dochange(srcname, srcfile, dstname, dstfile []string) {
	for index := 0; index < len(srcname); index++ {
		for i, v := range dstname {
			if srcnames[index] == v {
				err := os.Rename(dstfile[i], dstfile[i]+"."+tt1)
				if err != nil {
					write2log(err.Error(), tt, "error")
				} else {
					write2log(dstfile[i]+"."+tt1, tt1, "bk")
					err = Copy(srcfile[index], dstfile[i])
					if err != nil {
						write2log(err.Error(), tt, "error")
					} else {
						write2log(dstfile[i], tt1, "cp")
					}
				}
			}
		}
	}
}

func stoppro(e error, t time.Duration, n int) {
	write2log(e.Error(), tt, "error")
	fmt.Println("An exception occurs and the program closes automatically after 5 seconds:", e)
	time.Sleep(t)
	os.Exit(n)
}
