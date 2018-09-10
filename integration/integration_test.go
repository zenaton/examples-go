package integration_test

import (
	. "github.com/zenaton/examples-go/integration"

	"testing"

	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/gomega"
)

type entry struct {
	directory     string
	context       string
	it            string
	fileRunOutput string
	fileRunErr    bool
	outFile       string
	errFile       string
	sleep         int //time to wait for the file to be written to
}

var table = []entry{
	{
		directory: "testFindWorkflow",
		context:   "with a workflow that has an ID method",
		it:        "should be able to find a workflow and dispatch it again",
		outFile:   "out:  testTaskReturn\nerr:  <nil>\nout:  testTaskReturn\nerr:  <nil>\n",
		errFile:   "",
		sleep:     11,
	},
	{
		directory: "testReturnFromTaskInsideTask",
		context:   "with a task that lanches another task",
		it:        "the outer task should be able to get the return value of the inner task",
		outFile:   "err:  <nil>\nout:  test return\nerr:  <nil>\n",
		errFile:   "",
		sleep:     10,
	},
	{
		directory: "testReturnError",
		context:   "with a task that returns an error",
		it:        "the workflow should exit and the error should be logged",
		outFile:   "",
		errFile:   "testTaskError\n",
		sleep:     5,
	},
	{
		directory:  "testUnserializableHandlerType",
		context:    "with workflow handler type that is unserializable",
		it:         "the application should panic",
		fileRunErr: true,
		fileRunOutput: "panic: handler type 'WithTask' must be able to be unmarshaled from json. " +
			"json: cannot unmarshal object into Go struct field Task.Handler of type interfaces.Handler",
	},
	{
		directory: "event",
		context:   "with an event based workflow",
		it:        "should handle events",
		outFile:   "Task C\nTask A\nTask B\n",
		errFile:   "",
		sleep:     13,
	},
	{
		directory: "wait",
		context:   "with a wait workflow",
		it:        "should wait before a task",
		outFile:   "Task A\nTask D\n",
		errFile:   "",
		sleep:     16,
	},
	{
		directory: "waitevent",
		context:   "with a waitevent workflow",
		it:        "should wait for either an event or a given time",
		outFile:   "Task A\n",
		errFile:   "",
		sleep:     10,
	},
	{
		directory: "asynchronous",
		context:   "with an asychronous workflow",
		it:        "should perform tasks asynchronously",
		outFile:   "Task B\nTask A\n",
		errFile:   "",
		sleep:     11,
	},
	{
		directory: "parallel",
		context:   "with a parallel workflow",
		it:        "should handle tasks in parallel",
		outFile:   "Task B\nTask A\nTask C\n",
		errFile:   "",
		sleep:     15,
	},
	//{
	//	directory: "recursive",
	//	context:   "with a recursive workflow",
	//	it:        "should handle tasks recursively",
	//	outFile:   "0\nIteration: 1\n0\nIteration: 2\n0",
	//	errFile:   "",
	//	sleep:     16,
	//},
	{
		directory: "sequential",
		context:   "with a sequential workflow",
		it:        "should handle tasks sequentially",
		outFile:   "Task A\nTask B\n",
		errFile:   "",
		sleep:     15,
	},
	{
		directory: "version",
		context:   "with a version workflow",
		it:        "should handle versioned workflows",
		outFile:   "Task D\nTask C\nTask B\nTask A\n",
		errFile:   "",
		sleep:     10,
	},
}

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

					logFile, err := os.OpenFile(entry.directory+"/zenaton.log", os.O_RDWR, 0660)
					switch err.(type) {
					case *os.PathError:
						//this is ok
					default:
						g.Expect(err).ToNot(HaveOccurred())

						err = logFile.Truncate(0)
						g.Expect(err).ToNot(HaveOccurred())
						_, err = logFile.Seek(0, 0)
						g.Expect(err).ToNot(HaveOccurred())
					}
					defer logFile.Close()

					cmd := exec.Command("go", "run", "main.go")
					cmd.Dir = entry.directory

					out, err := cmd.CombinedOutput()

					g.Expect(string(out)).To(ContainSubstring(entry.fileRunOutput))

					if entry.fileRunErr {
						g.Expect(err).To(HaveOccurred())
					} else {
						g.Expect(err).ToNot(HaveOccurred())
					}

					time.Sleep(time.Duration(entry.sleep) * time.Second)

					errLog, err := ioutil.ReadAll(errFile)
					g.Expect(err).ToNot(HaveOccurred())
					outLog, err := ioutil.ReadAll(outFile)
					g.Expect(err).ToNot(HaveOccurred())

					g.Expect(string(errLog)).To(Equal(entry.errFile))
					g.Expect(string(outLog)).To(Equal(entry.outFile))
				}
			}
		})
	}
}
