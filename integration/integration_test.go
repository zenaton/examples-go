package integration_test

import (
	. "github.com/zenaton/examples-go/integration"

	//. "github.com/onsi/ginkgo"
	"os"

	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	{
		directory: "event",
		context:   "with an event based workflow",
		it:        "should handle events",
		output:    "Task C\nTask A\nTask B\n",
		err:       "",
		sleep:     8,
	},
	{
		directory: "asynchronous",
		context:   "with an asychronous workflow",
		it:        "should perform tasks asynchronously",
		output:    "Task B\nTask A\n",
		err:       "",
		sleep:     7,
	},
	{
		directory: "parallel",
		context:   "with a parallel workflow",
		it:        "should handle tasks in parallel",
		output:    "Task B\nTask A\nTask C\n",
		err:       "",
		sleep:     9,
	},
	{
		directory: "recursive",
		context:   "with a recursive workflow",
		it:        "should handle tasks recursively",
		output:    "0\nIteration: 1\n0\nIteration: 2\n0",
		err:       "",
		sleep:     8,
	},
	{
		directory: "sequential",
		context:   "with a sequential workflow",
		it:        "should handle tasks sequentially",
		output:    "Task A\nTask B\n",
		err:       "",
		sleep:     10,
	},
	{
		directory: "version",
		context:   "with a version workflow",
		it:        "should handle versioned workflows",
		output:    "Task D\nTask C\nTask B\nTask A\n",
		err:       "",
		sleep:     8,
	},
	{
		directory: "wait",
		context:   "with a wait workflow",
		it:        "should wait before a task",
		output:    "Task A\nTask D\n",
		err:       "",
		sleep:     12,
	},
	{
		directory: "waitevent",
		context:   "with a waitevent workflow",
		it:        "should wait for either an event or a given time",
		output:    "Task B\n",
		err:       "",
		sleep:     6,
	},
}

var _ = Describe("Integration", func() {

	It("should listen", func() {
		err := Listen()
		Expect(err).NotTo(HaveOccurred())
	})

	for _, entry := range table {
		entry := entry //gotcha!

		Context(entry.context, func() {
			It(entry.it, func() {

				path, err := os.Getwd()
				fmt.Println("path, err", path, err)

				Expect(err).NotTo(HaveOccurred())

				out, err := exec.Command("go", "run", "../"+entry.directory+"/main.go").CombinedOutput()
				if err != nil {
					fmt.Println("out3: ", string(out))
				}
				Expect(err).ToNot(HaveOccurred())

				errFile, err := os.OpenFile("../zenaton.err", os.O_RDWR, 0660)
				switch err.(type) {
				case *os.PathError:
					//this is ok
					errFile, err = os.Create("../zenaton.err")
					Expect(err).ToNot(HaveOccurred())

				default:
					Expect(err).ToNot(HaveOccurred())

					//clear the files
					err = errFile.Truncate(0)
					Expect(err).ToNot(HaveOccurred())
					_, err = errFile.Seek(0, 0)
					Expect(err).ToNot(HaveOccurred())
				}
				defer errFile.Close()

				outFile, err := os.OpenFile("../zenaton.out", os.O_RDWR, 0660)
				switch err.(type) {
				case *os.PathError:
					//this is ok
					outFile, err = os.Create("../zenaton.out")
					Expect(err).ToNot(HaveOccurred())
				default:
					Expect(err).ToNot(HaveOccurred())

					err = outFile.Truncate(0)
					Expect(err).ToNot(HaveOccurred())
					_, err = outFile.Seek(0, 0)
					Expect(err).ToNot(HaveOccurred())
				}
				defer errFile.Close()

				time.Sleep(time.Duration(entry.sleep) * time.Second)
				errLog, err := ioutil.ReadAll(errFile)
				Expect(err).ToNot(HaveOccurred())
				outLog, err := ioutil.ReadAll(outFile)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(errLog)).To(Equal(entry.err))
				Expect(string(outLog)).To(Equal(entry.output))
			})
		})
	}
})
