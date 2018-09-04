package integration_test

import (
	. "github.com/zenaton/examples-go/integration"

	"testing"

	. "github.com/onsi/gomega"
	"os"
	"os/exec"
	"time"
	"fmt"
	"io/ioutil"
)

type entry struct {
	directory string
	context   string
	it        string
	output    string
	err       string
	sleep     int //time to wait for the file to be writen to
}

var table = []entry{
	//{
	//	directory: "event",
	//	context:   "with an event based workflow",
	//	it:        "should handle events",
	//	output:    "Task C\nTask A\nTask B\n",
	//	err:       "",
	//	sleep:     11,
	//},
	//{
	//	directory: "asynchronous",
	//	context:   "with an asychronous workflow",
	//	it:        "should perform tasks asynchronously",
	//	output:    "Task B\nTask A\n",
	//	err:       "",
	//	sleep:     9,
	//},
	{
		directory: "parallel",
		context:   "with a parallel workflow",
		it:        "should handle tasks in parallel",
		output:    "Task B\nTask A\nTask C\n",
		err:       "",
		sleep:     11,
	},
	{
		directory: "recursive",
		context:   "with a recursive workflow",
		it:        "should handle tasks recursively",
		output:    "0\nIteration: 1\n0\nIteration: 2\n0",
		err:       "",
		sleep:     12,
	},
	//{
	//	directory: "sequential",
	//	context:   "with a sequential workflow",
	//	it:        "should handle tasks sequentially",
	//	output:    "Task A\nTask B\n",
	//	err:       "",
	//	sleep:     12,
	//},
	//{
	//	directory: "version",
	//	context:   "with a version workflow",
	//	it:        "should handle versioned workflows",
	//	output:    "Task D\nTask C\nTask B\nTask A\n",
	//	err:       "",
	//	sleep:     10,
	//},
	//{
	//	directory: "wait",
	//	context:   "with a wait workflow",
	//	it:        "should wait before a task",
	//	output:    "Task A\nTask D\n",
	//	err:       "",
	//	sleep:     14,
	//},
	//{
	//	directory: "waitevent",
	//	context:   "with a waitevent workflow",
	//	it:        "should wait for either an event or a given time",
	//	output:    "Task B\n",
	//	err:       "",
	//	sleep:     6,
	//},
}

//var _ = BeforeSuite(func() {
//	for _, entry := range table {
//		entry := entry //gotcha!
//		err := SetUpTestDirectories(entry.directory)
//		g.Expect(err).NotTo(HaveOccurred())
//
//		err = Copy("../"+entry.directory+"/main.go", entry.directory+"/main.go")
//		g.Expect(err).NotTo(HaveOccurred())
//
//		err = CDIntoDir(entry.directory)
//		g.Expect(err).NotTo(HaveOccurred())
//
//		err = os.Setenv("ZENATON_APP_ENV", "dev-"+entry.directory)
//		g.Expect(err).NotTo(HaveOccurred())
//		client.SetEnv()
//
//		err = Listen()
//		g.Expect(err).NotTo(HaveOccurred())
//
//		err = CDIntoDir("..")
//		g.Expect(err).NotTo(HaveOccurred())
//	}
//})

func TestSetup(t *testing.T) {
	g := NewGomegaWithT(t)
	for _, entry := range table {
		entry := entry //gotcha!
		err := SetUpTestDirectories(entry.directory)
		g.Expect(err).NotTo(HaveOccurred())

		err = Copy("../"+entry.directory+"/main.go", entry.directory+"/main.go")
		g.Expect(err).NotTo(HaveOccurred())

		envFile := entry.directory + ".env"
		err = ChangeClient(entry.directory+"/main.go", envFile)
		g.Expect(err).NotTo(HaveOccurred())

		created, err := AddEnv(envFile)
		g.Expect(err).NotTo(HaveOccurred())

		if created {
			err = WriteAppEnv(envFile, entry.directory)
			g.Expect(err).NotTo(HaveOccurred())
		}

		bootFileName := entry.directory + "boot.go"

		err = AddBoot(bootFileName, envFile)
		g.Expect(err).NotTo(HaveOccurred())

		err = Listen(envFile, bootFileName, entry.directory)
		g.Expect(err).NotTo(HaveOccurred())
	}
}

func TestExamples(t *testing.T) {
	for _, entry := range table {
		entry := entry //gotcha!

		t.Run("", func(st *testing.T) {
			g := NewGomegaWithT(st)
			st.Parallel()
			entry := entry
			st.Log(entry.context)
			{
				st.Log(entry.it)
				{
					errFile, err := os.OpenFile(entry.directory+"/zenaton.err", os.O_RDWR, 0660)
					switch err.(type) {
					case *os.PathError:
						//this is ok
						errFile, err = os.Create(entry.directory + "/zenaton.err")
						g.Expect(err).ToNot(HaveOccurred())

					default:
						g.Expect(err).ToNot(HaveOccurred())

						//clear the files
						err = errFile.Truncate(0)
						g.Expect(err).ToNot(HaveOccurred())
						_, err = errFile.Seek(0, 0)
						g.Expect(err).ToNot(HaveOccurred())
					}
					defer errFile.Close()

					outFile, err := os.OpenFile(entry.directory+"/zenaton.out", os.O_RDWR, 0660)
					switch err.(type) {
					case *os.PathError:
						//this is ok
						outFile, err = os.Create(entry.directory + "/zenaton.out")
						g.Expect(err).ToNot(HaveOccurred())
					default:
						g.Expect(err).ToNot(HaveOccurred())

						err = outFile.Truncate(0)
						g.Expect(err).ToNot(HaveOccurred())
						_, err = outFile.Seek(0, 0)
						g.Expect(err).ToNot(HaveOccurred())
					}
					defer outFile.Close()

					cmd := exec.Command("go", "run", "main.go")
					cmd.Dir = entry.directory
					out, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Println("out3: ", string(out))
					}
					g.Expect(err).ToNot(HaveOccurred())

					time.Sleep(time.Duration(entry.sleep) * time.Second)

					errLog, err := ioutil.ReadAll(errFile)
					g.Expect(err).ToNot(HaveOccurred())
					outLog, err := ioutil.ReadAll(outFile)
					g.Expect(err).ToNot(HaveOccurred())

					g.Expect(string(errLog)).To(Equal(entry.err))
					g.Expect(string(outLog)).To(Equal(entry.output))
				}
			}
		})
	}
}
