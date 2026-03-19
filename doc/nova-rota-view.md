# Como Criar uma Nova View/Rota de Prompts (Super Simplificado!)

Agora o sistema de Prompts é totalmente genérico usando abstrações no Backend e Frontend.
Para criar uma nova tela (ex: `IzaHomol`, `BiaProd`, etc), você só precisa seguir **3 Passos Simples**:

## Passo 1: Adicione o Proxy no Backend
Abra o arquivo `backend/cmd/main.go` e procure pelo comentário `// Register specific subsystem proxies`. Você verá as rotas já existentes lá.
Adicione **uma linhazinha** com a sua nova rota usando a função `registerSystemProxy`.

**Exemplo para Iza Homol:**
```go
registerSystemProxy(
  mux, 
  authMW, 
  middleware.RateLimitMiddleware(apiLimiter), 
  "/api/proxy/iza/homol/",      // Prefixo da URL que o frontend vai bater (Sempre inicie com /api/proxy/)
  "iza_homol",                  // Nome do sistema (Permissão no Banco/Presence)
  "IZA_HOMOL_API_BASE",         // Chave de ENV contendo a URL base upstream
  "IZA_HOMOL_API_USER",         // Chave de ENV contendo o Usuário (se tiver auth)
  "IZA_HOMOL_API_PASS",         // Chave de ENV contendo a Senha (se tiver auth)
)
```

## Passo 2: Adicione a Variável de Ambiente
No seu arquivo `.env` do backend, adicione as variáveis correspondentes aos nomes que você declarou no passo anterior:

```env
IZA_HOMOL_API_BASE=http://ip-do-servidor/api/iza/homol
# Se não houver basic auth, pode omitir USER e PASS do .env.
```

## Passo 3: Adicione a UI no Vue Router (Frontend)
Abra o arquivo `frontend/src/router/index.js` e adicione a nova rota apontando para o componente genérico `PromptManager.vue`, apenas passando as **props** de configuração!

```javascript
{
  path: '/iza-homol',       // Rota que vai aparecer no navegador
  name: 'IzaHomol',         // Nome identificador pro Vue
  component: () => import('@/views/PromptManager.vue'),
  props: { 
    title: 'IZA HOMOL',           // Título gigante na Nav bar!
    systemId: 'iza_homol',        // Nome usado pro Websocket do Presence (Quem tá online)
    apiBase: '/api/proxy/iza/homol' // O exato prefixo que você registrou no Backend (passo 1)
  },
  meta: { 
    requiresAuth: true, 
    requiresSystem: 'iza_homol'   // Permissão do usuário 
  }
}
```

**Pronto! A tela inteira de gerenciamento de prompts já vai estar funcional e com WebSocket configurado.**
Nenhum arquivo `.vue` extra precisa ser criado.
Nenhum arquivo `.js` na pasta `api` precisa ser criado.
Nenhum arquivo handler novo no Backend precisa ser criado.
