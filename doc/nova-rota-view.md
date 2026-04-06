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

## Passo 4: Registrar o Nome Amigável no Login
Para que o novo sistema apareça com um nome bonito no "Painel de Sistemas" (ao invés do ID técnico), você precisa atualizar o arquivo `frontend/src/views/Login.vue`.

Procure a função `formatSystemName` e adicione o mapeamento:
```javascript
function formatSystemName(sys) {
  const names = {
    // ... outros sistemas
    'iza_homol': 'IZA (Homologação)'   // O ID que você criou e o Nome que o usuário vê
  }
}
```

E na função `handleSystemSelection`, adicione o redirecionamento:
```javascript
function handleSystemSelection(sys) {
  // ... outros ifs
  } else if (sys === 'iza_homol') {
    router.push('/iza-homol')         // A rota que você criou no Passo 3
  }
}
```

## Passo 5: Mapear Variáveis no Docker Compose
Mesmo que as variáveis estejam no `.env`, o Docker não as passa automaticamente para dentro do container do backend. Você precisa mapeá-las nos arquivos `docker-compose.yml` (local) e `docker-compose.prod.yml` (produção).

Abra os arquivos e procure pela seção `backend` -> `environment`. Adicione as novas variáveis:

```yaml
      # ... outras variáveis
      - IZA_HOMOL_API_BASE=${IZA_HOMOL_API_BASE}
      - IZA_HOMOL_API_USER=${IZA_HOMOL_API_USER}
      - IZA_HOMOL_API_PASS=${IZA_HOMOL_API_PASS}
```

**Pronto! Agora o fluxo completo do Backend até a UI está configurado.**
Infelizmente esse passo ainda é manual para garantir que você tenha controle total sobre os nomes exibidos na tela de Login.
