# Use uma imagem base do Go
FROM golang:1.21

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Copie os arquivos do projeto para o diretório de trabalho
COPY . .

# Navegue para o diretório src
WORKDIR /app/src

# Baixe as dependências e compile o binário
RUN go mod tidy
RUN go build -o /app/app

# Defina a porta em que a aplicação irá rodar
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["/app/app"]
