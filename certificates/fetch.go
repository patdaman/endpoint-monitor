package fetchCertificate

import "crypto/tls"
import "crypto/sha1"
import "crypto/x509"
import "fmt"
import "encoding/pem"
import "os"
import "time"
import "bufio"
import "strings"
import "net"

func storeCertificate(cert *x509.Certificate) {
	hash := sha1.New()
	hash.Write(cert.Raw)
	filename := fmt.Sprintf("%X", hash.Sum(nil))
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, _ := os.Create(filename)
		defer f.Close()
		block := &pem.Block{"CERTIFICATE", nil, cert.Raw}
		fmt.Println("\x1b[32;1madded", filename, "\x1b[0m")
		pem.Encode(f, block)
	}
}

func analyzeDomain(domain string) {
	fmt.Println("analyzing", domain)
	dialer := net.Dialer{}
	dialer.Timeout = 10 * time.Second
	conn, err := tls.DialWithDialer(&dialer, "tcp", domain, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[31;1mfailed to connect to %s", domain), err, "\x1b[0m")
		return
	}
	defer conn.Close()
	for _, cert := range conn.ConnectionState().PeerCertificates {
		storeCertificate(cert)
	}
}

func analyzeDomains(queue chan string) {
	for {
		domain := <-queue
		analyzeDomain(domain)
	}
}

// func main() {
func fetchCertificate() {
	cs := make(chan string)
	for i := 0; i < 80; i++ {
		go analyzeDomains(cs)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			line = line + ":443"
		}
		cs <- line
	}
	time.Sleep(2 * time.Second)
}
