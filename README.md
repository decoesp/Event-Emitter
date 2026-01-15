# Coding Challenge: EventEmitter

## O problema

Crie um package `eventemitter` usando go.

Com ele é possível escutar e emitir eventos.

## API esperada:
- `Listen`: Deve adicionar um listner para o eventName;
- `Emit`: Deve emitir eventos, trigando seus listeners;
- `RemoveEvent`: Deve remover os listeners desse eventName;
- Eventos podem ser emitidos com dados, que serão repassados para os listeners;
- `ListenOnce`: Semelhante ao `listen`, entretanto o callback será trigado apenas uma vez;
- `Reset`: Deve remover todos listeners de todos os eventos;

## Exemplo:

```go
doSomething := func(data interface{}) {
    //received!
    log.Println(data)
}

doSomethingElse := func(data interface{}) {
    //received!
    //storage.Save(...)
}

eventemitter.Listen("my-awesome-event", doSomething)
eventemitter.Listen("my-awesome-event", doSomethingElse)

eventemitter.Emit("my-awesome-event", nil)
eventemitter.Emit("my-awesome-event", "foo")

eventemitter.RemoveEvent("my-awesome-event")
```

## Expectativa de solução

- Construir a lib eventemitter, para passar nos testes automatizados.
- Os primeiros testes são os mais importantes.
- Existem bibliotecas e artigos na internet que solucionam o problema, mas elas não devem ser usadas.
- Conte com a ajuda dos entrevistadores para qualquer dúvida.
- Tudo bem nunca ter visto esse problema antes, pergunte a vontade.
- O live coding é mais sobre como resolver o problema, e menos sobre conhecimento prévio.
- Não há resposta correta.
- Problemas podem ser resolvidos de diferentes formas.
- Fique a vontade para trazer sua solução. =)

## Arquivo de tests
Junto ao desafio, há um arquivo de testes para verificar o funcionamento do código.

## Rodando os testes

```bash
go test
```
# Event-Emitter
