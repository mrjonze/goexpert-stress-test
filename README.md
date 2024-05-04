# goexpert-stress-test
Resposta do desafio técnico Stress Test da pós Go Expert.

Para buildar a imagem docker do projeto, execute o comando abaixo:
```shell
docker build -t <nome_da_imagem> -f Dockerfile .
```

Para rodar o container da imagem docker do projeto, execute o comando abaixo:
```shell
docker run <nome_da_imagem> -url=<url_de_preferencia> -concurrency=<quantidade_de_requests_simultaneos> -requests=<quantidade_total_de_requests>
``` 