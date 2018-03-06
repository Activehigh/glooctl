package file_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"

	"path/filepath"

	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
	. "github.com/solo-io/gloo/pkg/secretwatcher/file"
)

var _ = Describe("FileSecretWatcher", func() {
	var (
		dir   string
		file  string
		err   error
		watch secretwatcher.Interface
	)
	BeforeEach(func() {
		dir, err = ioutil.TempDir("", "filesecrettest")
		Must(err)
		file = filepath.Join(dir, "secrets.yml")
		watch, err = NewSecretWatcher(dir, time.Millisecond)
		Must(err)
	})
	AfterEach(func() {
		log.Debugf("removing " + dir)
		os.RemoveAll(dir)
	})
	Describe("watching file", func() {
		Context("an invalid structure is written to a file", func() {
			It("sends an error on the Error() channel", func() {
				invalidData := []byte("]]foo: bar")
				err = ioutil.WriteFile(file, invalidData, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-watch.Secrets():
					Fail("config was received, expected error")
				case err := <-watch.Error():
					Expect(err).To(HaveOccurred())
				case <-time.After(time.Second * 1):
					Fail("expected err to have occurred before 1s")
				}
			})
		})
		Context("no secrets wanted", func() {
			It("doesnt send anything on any channel", func() {
				missingSecrets := map[string]map[string][]byte{"another-key": {"foo": []byte("bar"), "baz": []byte("qux")}}
				data, err := json.Marshal(missingSecrets)
				Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(file, data, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-watch.Secrets():
					Fail("secretmap was received, expected timeout")
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 1):
					// passed
				}
			})
		})
		Context("want secrets that the file doesn't contain", func() {
			It("sends an error on the Error() channel", func() {
				missingSecrets := map[string]map[string][]byte{"another-key": {"foo": []byte("bar"), "baz": []byte("qux")}}
				data, err := json.Marshal(missingSecrets)
				Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(file, data, 0644)
				Expect(err).NotTo(HaveOccurred())
				go watch.TrackSecrets([]string{"this key really should not be in the secretmap"})
				select {
				case <-watch.Secrets():
					Fail("secretmap was received, expected error")
				case err := <-watch.Error():
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("secretmap not found"))
				case <-time.After(time.Second * 1):
					Fail("expected err to have occurred before 1s")
				}
			})
		})
		Context("a valid config is written to a file", func() {
			It("sends a corresponding secretmap on Secrets()", func() {
				secretMap := NewTestSecrets()
				yml, err := yaml.Marshal(secretMap)
				Must(err)
				err = ioutil.WriteFile(file, yml, 0644)
				Must(err)
				var key string
				for k := range secretMap {
					key = k
					break
				}
				go watch.TrackSecrets([]string{key})
				select {
				case parsedSecrets := <-watch.Secrets():
					Expect(parsedSecrets).To(Equal(secretMap))
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 5):
					Fail("expected new secrets to be read in before 1s")
				}
			})
		})
	})
})
