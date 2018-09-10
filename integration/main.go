package integration

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func SetUpTestDirectories(dir string) error {
	err := os.Mkdir(dir, 0777)
	if err != nil && strings.Contains(err.Error(), "file exists") {
		return nil
	}
	return err
}

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

func ChangeClient(fileName string, envFile string) error {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	contents = bytes.Replace(contents, []byte("_ "), []byte(""), -1)
	contents = bytes.Replace(contents, []byte("examples-go/client"), []byte("examples-go/integration/client"), -1)
	contents = bytes.Replace(contents, []byte("\nfunc main() {"), []byte("func init(){client.SetEnv(\""+envFile+"\")}\nfunc main() {"), -1)
	return ioutil.WriteFile(fileName, contents, 0660)
}

func AddEnv(specificEnv string) (bool, error) {
	_, err := os.OpenFile("./client/"+specificEnv, os.O_RDWR, 0660)
	switch err.(type) {
	case *os.PathError:
		//this is ok
		err = Copy("./client/.env", "./client/"+specificEnv)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, nil
	}
}

func AddBoot(specificBoot string, envFileName string) error {
	err := os.Mkdir("./boot/"+specificBoot, 0777)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}

	_, err = os.OpenFile("./boot/"+specificBoot+"/"+specificBoot, os.O_RDWR, 0660)
	switch err.(type) {
	case *os.PathError:
		//this is ok
		err = Copy("./boot/boot.go", "./boot/"+specificBoot+"/"+specificBoot)
		if err != nil {
			return err
		}
	default:
	}

	contents, err := ioutil.ReadFile("./boot/" + specificBoot + "/" + specificBoot)
	contents = bytes.Replace(contents, []byte(`client.SetEnv("")`), []byte(`client.SetEnv("`+envFileName+`")`), -1)
	return ioutil.WriteFile("./boot/"+specificBoot+"/"+specificBoot, contents, 0660)
}

func WriteAppEnv(file, env string) error {

	envFile, err := os.OpenFile("./client/"+file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = envFile.WriteString("ZENATON_APP_ENV=dev-" + env)
	return err
}

func Listen(specificEnv string, bootFile string, executionDir string) error {

	args := []string{"listen", "--env=../client/" + specificEnv, "--boot=../boot/" + bootFile + "/" + bootFile}

	cmd := exec.Command("zenaton", args...)
	cmd.Dir = executionDir

	out, err := cmd.CombinedOutput()
	fmt.Println("out1: ", string(out))

	//try again
	if err != nil {
		out, err = cmd.CombinedOutput()
		fmt.Println("out2: ", string(out))
	}

	return err
}
