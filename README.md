# pg-shazam

Limits:

-   Max 5 ms on proxy

Todo:

-   Adjust parser
-   Ensure table state on shazam db in all pgs
-   Distribued kv state
-   Separar error e fatal e direcionar para o front o erro correto ou abortar a conexão com um fatal e msg de erro do proto3
-   Tratar erro de timeout com o pg

https://github.com/mackerelio/go-osstat

## How it works

1. Inicia raft, eleição e sync das informações
2. Inicia conexão com os postgres
3. Synca as informações com os postgres
4. Inicial pool de conexões com os postgres
5. Inicia handler de conexões
6. Inicia o servidor tcp
7. Espera conexões
8. Nova conexão
    1. Comunicação inicial
    2. Validação de autenticação
    3. Pronto para query















https://github.com/alecthomas/participle
https://github.com/tshprecher/antlr_psql/blob/master/antlr4/PostgreSQLLexer.g4
https://github.com/pgcodekeeper/pgcodekeeper/blob/master/apgdiff/antlr-src/SQLParser.g4
https://github.com/postgres/postgres/blob/master/src/backend/parser/gram.y
https://github.com/antlr/grammars-v4
