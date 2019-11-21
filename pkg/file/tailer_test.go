package file_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sjeandeaux/access-log-monitor/pkg/file"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Suite")
}

var _ = Describe("Tailer", func() {

	Describe("Tailering the file", func() {
		Context("With a correct file", func() {
			It("should receive lines", func() {
				name := emulateWriteInFileForTail(context.TODO())
				tailer, err := file.NewDefaultTailer(name)
				Ω(err).Should(BeNil())

				lines, errs := tailer.Tail(context.TODO())

				var data string
				Eventually(lines).Should(Receive(&data))
				Ω(data).Should(Equal("tail"))
				Eventually(errs).ShouldNot(Receive())
			})
		})

		Context("With a correct file and closed context", func() {
			It("should close the channels", func() {
				name := emulateWriteInFileForTail(context.TODO())
				tailer, err := file.NewDefaultTailer(name)
				Ω(err).Should(BeNil())

				ctx := context.Background()
				ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
				defer cancel()
				lines, errs := tailer.Tail(ctx)
				Eventually(lines).Should(BeClosed())
				Eventually(errs).Should(BeClosed())
			})
		})
	})
})

// emulateWriteInFileForTail it creates a temporary file
// a goroutine will write inside each 500 * time.Millisecond
func emulateWriteInFileForTail(ctx context.Context) string {
	tmpfile, err := ioutil.TempFile("", "tailer")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer os.Remove(tmpfile.Name()) // at the end, it can remove the fail
		ticker := time.NewTicker(500 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Fprintln(tmpfile, "tail")
			}
		}

	}()
	return tmpfile.Name()
}
