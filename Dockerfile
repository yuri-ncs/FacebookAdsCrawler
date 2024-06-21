# Use a imagem oficial do Golang como imagem base
FROM golang:1.21-alpine AS builder

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o arquivo go.mod e go.sum para baixar as dependências
COPY go.mod go.sum ./

# Baixe as dependências
RUN go mod download

# Copie o código fonte para o diretório de trabalho no contêiner
COPY . .

# Compile o aplicativo
RUN go build -o main .

# Inicie um novo estágio de construção
FROM alpine:latest

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o executável compilado do estágio anterior
COPY --from=builder /app/main .

# Exponha a porta 8080 para acesso externo
EXPOSE 8080

# Execute o aplicativo quando o contêiner for iniciado
CMD ["./main"]