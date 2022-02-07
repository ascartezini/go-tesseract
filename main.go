package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
	"unsafe"
)

var httpClient *http.Client

func main() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	http.HandleFunc("/", obterImagemDaUrlHandler)
	http.HandleFunc("/erros", obterErros)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func catalogarErro(mensagem string) {
	errosAntigos := obterMensagensDosErros()
	bytes := []byte(time.Now().Format("2006-01-02 15:04:05") + " - " + mensagem + "\n" + errosAntigos)
	ioutil.WriteFile("erros.txt", bytes, 0)
}

func obterMensagensDosErros() (texto string) {
	binario, erro := ioutil.ReadFile("erros.txt")
	if erro != nil {
		texto = ""
		return
	}

	return *(*string)(unsafe.Pointer(&binario))
}

func obterErros(w http.ResponseWriter, r *http.Request) {
	erros := obterMensagensDosErros()

	var b []byte
	b = append(b, erros...)

	w.Write(b)
}

func obterImagemDaUrlHandler(w http.ResponseWriter, r *http.Request) {
	imagens, ok := r.URL.Query()["image"]
	if !ok || len(imagens[0]) < 1 {
		json.NewEncoder(w).Encode("Informe o link da imagem. Exemplo: http://localhost:3000/?image=hello.png")
		return
	}

	imagem := string(imagens[0])
	res := obterImagemDaUrl(imagem)

	json.NewEncoder(w).Encode(res)
}

func obterImagemDaUrl(imagem string) string {
	nomeDaImagem := gerarNomeDaImagem(8)
	downloadDaImagem(imagem, nomeDaImagem)
	imageText, err := executarOCR(nomeDaImagem)
	go os.Remove(nomeDaImagem)

	if err != nil {
		return "an error occurred " + err.Error() + os.Getenv("PATH")
	}

	return imageText
}

func executarOCR(nomeDaImagem string) (string, error) {
	var tesseract []byte
	tesseract = append(tesseract, "tesseract "...)
	tesseract = append(tesseract, nomeDaImagem...)
	tesseract = append(tesseract, " stdout -l por"...)

	c := *(*string)(unsafe.Pointer(&tesseract))
	out, err := exec.Command("sh", "-c", c).Output()

	if err != nil {
		catalogarErro("erro na imagem: " + nomeDaImagem)
		return "", err
	}

	return *(*string)(unsafe.Pointer(&out)), nil
}

func gerarNomeDaImagem(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var buf []byte
	var num rune = 1

	l := len(letterRunes)
	for i := 0; i < n; i++ {
		num = letterRunes[rand.Intn(l)]
		buf = append(buf, string(num)...)
	}

	buf = append(buf, ".jpg"...)

	return *(*string)(unsafe.Pointer(&buf))
}

func downloadDaImagem(url string, nomeDaImagem string) {
	output, _ := os.Create(nomeDaImagem)
	defer output.Close()

	response, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("Erro de download", url, "-", err)
		return
	}
	defer response.Body.Close()

	io.Copy(output, response.Body)
}
