# Usando uma imagem base do Go
FROM golang:1.22.4

# Diretório de trabalho
WORKDIR /app

# Copiar os arquivos para o contêiner
COPY . .

# Instalar as dependências e construir o binário
RUN go mod download
RUN go build -o main .

# Expor a porta usada pela API
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
