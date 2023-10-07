package whatever

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	heritage "github.com/mobroco/heritage/pkg"
	"gopkg.in/yaml.v3"
)

func OpenPacket() (packet heritage.Packet) {
	if raw, err := os.ReadFile("garden.yaml"); err == nil {
		if err := yaml.Unmarshal(raw, &packet); err == nil {
			fmt.Println("using garden.yaml")
			return packet
		}
	}
	panic("?")
}

func OpenLocalPacket(locals ...string) (packet heritage.Packet) {
	packet = OpenPacket()
	for _, l := range locals {
		switch l {
		case "fqdn":
			for i, thing := range packet.Seeds {
				fqdnParts := strings.Split(thing.FQDN, ".")
				fqdnParts[len(fqdnParts)-1] = "x:3000"
				thing.FQDN = strings.Join(fqdnParts, ".")
				packet.Seeds[i] = thing
			}
		case "to":
			for i := range packet.Seeds {
				if len(packet.Seeds[i].Email.To) > 0 {
					packet.Seeds[i].Email.To = []string{"jmm@hey.com"}
				}
			}
		case "smtp":
			for i := range packet.Seeds {
				packet.Seeds[i].SMTP = heritage.SMTP{}
			}
		case "bucket":
			for i := range packet.Seeds {
				packet.Seeds[i].Bucket = heritage.Bucket{}
			}
		case "csp":
			for i := range packet.Seeds {
				packet.Seeds[i].ContentSecurityPolicy = heritage.ContentSecurityPolicy{}
			}
		case "cdn":
			for i := range packet.Seeds {
				packet.Seeds[i].CDN = heritage.CDN{}
			}
		case "login":
			for i := range packet.Seeds {
				if packet.Seeds[i].Login.RedirectURL != "" {
					packet.Seeds[i].Login.RedirectURL = fmt.Sprintf("http://" + packet.Seeds[i].FQDN + "/auth/callback")
				}
			}
		case "cookie":
			for i := range packet.Seeds {
				packet.Seeds[i].Cookie.Secure = false
				packet.Seeds[i].Cookie.Domain = fmt.Sprintf("http://" + packet.Seeds[i].FQDN)
			}
		case "stripe":
			for i := range packet.Seeds {
				packet.Seeds[i].Checkout.Stripe.SuccessURL = fmt.Sprintf("http://" + packet.Seeds[i].FQDN + "/checkout-success")
				packet.Seeds[i].Checkout.Stripe.CancelURL = fmt.Sprintf("http://" + packet.Seeds[i].FQDN + "/checkout-cancel")
			}
		}
	}
	return packet
}

func PeekAt(a any) {
	raw, _ := json.MarshalIndent(a, "", " ")
	fmt.Println(string(raw))
}

func Encrypt(plaintext, key string) string {
	iv := []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	plainText := []byte(plaintext)
	encryptor := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	encryptor.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText)
}

func Decrypt(cipherStr, key string) string {
	iv := []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	cipherText, err := base64.StdEncoding.DecodeString(cipherStr)
	if err != nil {
		panic(err)
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	decrypter.XORKeyStream(plainText, cipherText)
	return string(plainText)
}
