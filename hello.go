package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5

func main() {

	exibeIntroducao()

	for {
		exibeMenu()

		comando := leComando() //comando vai receber o retorna dessa funçao

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando.")
			os.Exit(-1) //-1 = ocorreu um erro no código
		}
	}

}

func exibeIntroducao() {
	nome := "Juliana" //declarando e atribuindo valor a uma var
	versao := 1.1
	fmt.Println("Olá, sra.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	//menu
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair")
}

// essa funçao vai retornar um int
func leComando() int {
	//capturando oq o usuário vai escrever
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento() { //função c/ múltiplo retorno
	fmt.Println("Monitorando...")

	sites := leSitesDoArquivo()

	//executando o loop X vezes
	for i := 0; i < monitoramentos; i++ {
		//range -> retorna o indice(i) , e o elemento que ele ta passando(sites)
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site)
			testaSite(site)
		}
		time.Sleep(delay * time.Second) //monitorando a cada 5s
		fmt.Println("")
	}

	fmt.Println("")

}

func testaSite(site string) {
	// ignorando a variavel de erro por enquanto
	resp, err := http.Get(site) //fzd uma requisição GET p/ site informado

	if err != nil { //se der erro -> mostra a msg de erro
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 { //detectando se o status-code é 200OK
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registraLog(site, true) //true -> site online
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status Code:", resp.StatusCode)
		registraLog(site, false)
	}

}

// func -> ler cada linha do arquivo.txt e extrai cada linha
func leSitesDoArquivo() []string { //retorna um Slice de string
	var sites []string

	//abrindo o arquivo p/ dps ser lido
	arquivo, err := os.Open("sites.txt")

	if err != nil { //se der erro -> mostra a msg de erro
		fmt.Println("Ocorreu um erro:", err)
	}

	//lendo linha a linha
	leitor := bufio.NewReader(arquivo)

	// for p/ ler todas as linhas
	for {
		linha, err := leitor.ReadString('\n') //ler ate '\n' e ja retorna uma string
		linha = strings.TrimSpace(linha)      //removendo os espaços desnecessarios
		sites = append(sites, linha)          //colocando a linha dentro do Slice de sites

		if err == io.EOF { //encontra o final do arquivo -> sai do loop
			break
		}
	}

	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	//quando abrir o arquivo -> posso escrever ou ler -> e caso o arquivo nao exista, crie ele ; 0666(permissao)
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	// formataçao da data e hora(go formata por numeros especificos) - nome do site - se ta online/offline
	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n") //convertendo o status p/ string
	arquivo.Close()
}

func imprimeLogs() {
	//abrindo o arquivo
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo)) //convertendo os bytes p/ string
}
