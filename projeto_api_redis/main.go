package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func main() {
	// Conectando ao Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Endereço padrão do Redis
		Password: "",           // Senha (se necessário)
		DB:       0,            // Número do banco de dados
	})

	// Verificando a conexão com o Redis
	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		fmt.Println("Erro ao conectar ao Redis:", err)
		return
	}

	// Definindo o handler para a rota "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Incrementando o contador de acessos no Redis
		err := incrementAccessCount()
		if err != nil {
			http.Error(w, "Erro ao incrementar o contador de acessos", http.StatusInternalServerError)
			return
		}

		// Escrevendo a resposta
		fmt.Fprintf(w, "Bem-vindo à minha API HTTP em Go! Quantidade de acessos: %d", getAccessCount())
	})

	// Definindo o endereço e a porta onde o servidor irá escutar
	endereco := ":8080"

	// Iniciando o servidor
	fmt.Printf("Servidor escutando na porta: %s\n", endereco)
	err = http.ListenAndServe(endereco, nil)
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %s\n", err)
	}
}

// Função para incrementar o contador de acessos no Redis
func incrementAccessCount() error {
	_, err := redisClient.Incr(redisClient.Context(), "access_count").Result()
	return err
}

// Função para obter o contador de acessos do Redis
func getAccessCount() int64 {
	val, err := redisClient.Get(redisClient.Context(), "access_count").Int64()
	if err != nil {
		fmt.Println("Erro ao obter o contador de acessos:", err)
		return 0
	}
	return val
}
