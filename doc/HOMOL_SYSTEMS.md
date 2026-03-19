## Guia para criar um novo domínio / system de *Homol*

Este guia descreve passo a passo como criar um novo **domínio (system) de homologação** seguindo o padrão já usado para `digesac_homol` e `bia_homol`.

O fluxo completo envolve:

- **Backend (Go)**: novo handler de proxy e registro na `main`.
- **Frontend (Vue 3)**: nova view, rota e módulo de API.
- **Infra**: variáveis de ambiente e `docker-compose.prod.yml`.

> Exemplos aqui usam um system fictício chamado `novo_homol`. Substitua pelos nomes reais do seu system.

---

## 1. Backend (Go)

### 1.1. Criar o handler de proxy

Crie um arquivo em `backend/routes/` seguindo o padrão de `digesac_proxy.go` e `bia_homol_proxy.go`.  
Exemplo: `backend/routes/novo_homol_proxy.go`:

```go
package routes

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "strings"
)

// NovoHomolProxyHandler forwards all /bestdoctors/novo/homol/* requests to the upstream
// NOVO_HOMOL_URL server, injecting Basic Auth credentials from environment.
func NovoHomolProxyHandler(w http.ResponseWriter, r *http.Request) {
  upstreamBase := os.Getenv("NOVO_HOMOL_URL")
  if upstreamBase == "" {
    http.Error(w, "NOVO_HOMOL_URL not configured", http.StatusInternalServerError)
    return
  }

  // Strip the /bestdoctors/novo/homol prefix so only the tail appends to NOVO_HOMOL_URL.
  upstreamPath := strings.TrimPrefix(r.URL.Path, "/bestdoctors/novo/homol")
  upstreamURL := fmt.Sprintf("%s%s", strings.TrimRight(upstreamBase, "/"), upstreamPath)
  if r.URL.RawQuery != "" {
    upstreamURL += "?" + r.URL.RawQuery
  }

  proxyReq, err := http.NewRequest(r.Method, upstreamURL, r.Body)
  if err != nil {
    http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
    return
  }

  if contentType := r.Header.Get("Content-Type"); contentType != "" {
    proxyReq.Header.Set("Content-Type", contentType)
  }

  apiUser := os.Getenv("NOVO_HOMOL_API_USER")
  apiPass := os.Getenv("NOVO_HOMOL_API_PASS")
  if apiUser != "" || apiPass != "" {
    proxyReq.SetBasicAuth(apiUser, apiPass)
  }

  resp, err := http.DefaultClient.Do(proxyReq)
  if err != nil {
    http.Error(w, "upstream request failed", http.StatusBadGateway)
    return
  }
  defer resp.Body.Close()

  w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
  w.WriteHeader(resp.StatusCode)
  io.Copy(w, resp.Body)
}
```

> Importante: siga a mesma estratégia de **prefixo protegido** usada em `bia_homol`:  
> - URL pública do painel: `/seu-systemo-homol` (rota do Vue).  
> - API interna do painel: `/bestdoctors/novo/homol/...` → passa pelo backend Go.  
> - API externa upstream real: `NOVO_HOMOL_URL` (por ex. `http://189.45.140.206/api/novo/homol`).

### 1.2. Registrar o handler na `main.go`

No arquivo `backend/cmd/main.go`, crie um `ServeMux` específico e registre o prefixo no `mux` principal, seguindo o padrão atual:

```go
// Digesac
digesacMux := http.NewServeMux()
digesacMux.HandleFunc("/", routes.DigesacProxyHandler)
mux.Handle("/digesac/homol/", middleware.RateLimitMiddleware(apiLimiter)(
  authMW(middleware.SystemMiddleware("digesac_homol")(digesacMux)),
))

// Bia Homol
biaMux := http.NewServeMux()
biaMux.HandleFunc("/", routes.BiaHomolProxyHandler)
mux.Handle("/bestdoctors/bia/homol/", middleware.RateLimitMiddleware(apiLimiter)(
  authMW(middleware.SystemMiddleware("bia_homol")(biaMux)),
))

// Novo Homol (exemplo)
novoMux := http.NewServeMux()
novoMux.HandleFunc("/", routes.NovoHomolProxyHandler)
mux.Handle("/bestdoctors/novo/homol/", middleware.RateLimitMiddleware(apiLimiter)(
  authMW(middleware.SystemMiddleware("novo_homol")(novoMux)),
))
```

Pontos importantes:

- O `SystemMiddleware("novo_homol")` **define o nome do system** usado para permissão de acesso (ver seção do frontend sobre `requiresSystem`).
- Use sempre um prefixo que **passe pelo backend** e não conflite com outro app (seguindo o padrão `/bestdoctors/<system>/homol/`).

---

## 2. Frontend (Vue 3)

### 2.1. Criar o módulo de API (`frontend/src/api`)

Use como base `frontend/src/api/digesac.js` e `frontend/src/api/bia.js`.  
Exemplo: `frontend/src/api/novo.js`:

```js
// src/api/novo.js
// Use backend-prefixed path to avoid clashing with other apps
const API_BASE = '/bestdoctors/novo/homol';

export async function fetchPromptsList() {
  const res = await fetch(`${API_BASE}/prompts?content=false`, {
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch prompts list: ${res.statusText}`);
  }
  return res.json();
}

export async function fetchPromptDetails(name) {
  const res = await fetch(`${API_BASE}/prompts/${name}`, {
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch prompt ${name}: ${res.statusText}`);
  }
  return res.json();
}

export async function updatePrompt(name, promptContent) {
  const res = await fetch(`${API_BASE}/prompts/${name}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ prompt: promptContent }),
  });
  if (!res.ok) {
    throw new Error(`Failed to update prompt ${name}: ${res.statusText}`);
  }
  return res.json();
}

export async function syncAllPrompts() {
  const res = await fetch(`${API_BASE}/prompts/sync`, {
    method: 'PUT',
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(`Failed to sync all prompts: ${res.statusText}`);
  }
  return res.json();
}

export async function syncSinglePrompt(name) {
  const res = await fetch(`${API_BASE}/prompts/sync/${name}`, {
    method: 'PUT',
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(`Failed to sync prompt ${name}: ${res.statusText}`);
  }
  return res.json();
}
```

### 2.2. Criar a view (tela de gerenciamento)

Use como base `frontend/src/views/DigesacHomol.vue` e `frontend/src/views/BiaHomol.vue`.  
Na prática, você pode:

- Duplicar uma das views existentes (`DigesacHomol.vue` ou `BiaHomol.vue`).
- Ajustar:
  - Título e textos (ex.: de `DIGESAC HOMOL` para `NOVO HOMOL`).
  - Import do módulo de API (`@/api/novo`).

Exemplo de trecho chave do script:

```js
// Dentro de NovoHomol.vue
import { fetchPromptsList, fetchPromptDetails, updatePrompt, syncAllPrompts, syncSinglePrompt } from '@/api/novo';
```

O restante da lógica (lista de arquivos `.md`, edição, sync, etc.) pode ser reaproveitada integralmente.

### 2.3. Registrar a rota no router

No arquivo `frontend/src/router/index.js`, adicione uma nova rota apontando para a view criada e configurando o `requiresSystem` igual ao nome de system que você usou no backend:

```js
const routes = [
  {
    path: '/digesac-homol',
    name: 'DigesacHomol',
    component: () => import('@/views/DigesacHomol.vue'),
    meta: { requiresAuth: true, requiresSystem: 'digesac_homol' },
  },
  {
    path: '/bia-homol',
    name: 'BiaHomol',
    component: () => import('@/views/BiaHomol.vue'),
    meta: { requiresAuth: true, requiresSystem: 'bia_homol' },
  },
  {
    path: '/novo-homol',
    name: 'NovoHomol',
    component: () => import('@/views/NovoHomol.vue'),
    meta: { requiresAuth: true, requiresSystem: 'novo_homol' },
  },
];
```

> O valor de `requiresSystem` **precisa bater** com o nome passado ao `SystemMiddleware("novo_homol")` no backend, e com o campo `system` do usuário retornado por `getCurrentUser()`.

### 2.4. Ajustar redirecionamentos de login (opcional)

Em `frontend/src/views/Login.vue` e `frontend/src/router/index.js`, existe lógica que redireciona o usuário para um system que ele tem acesso (por exemplo, `bia_homol` ou `digesac_homol`).  
Para incluir o novo system nesse fallback, basta seguir o mesmo padrão:

- No mapa de nomes exibidos.
- Nos `if (systems.includes('...'))` que escolhem a rota padrão.

---

## 3. Variáveis de ambiente

### 3.1. `.env` (usado pelo `docker-compose.prod.yml`)

Adicione as variáveis específicas do novo system, seguindo o padrão do Digesac e da Bia:

```env
# DIGESAC
DIGESAC_API_BASE=http://189.45.140.206/api/digesac/homol
DIGESAC_API_USER=homol
DIGESAC_API_PASS=********

# BIA HOMOL
BIA_HOMOL_URL=http://189.45.140.206/api/bia/homol
BIA_HOMOL_API_USER=homol
BIA_HOMOL_API_PASS=********

# NOVO HOMOL (exemplo)
NOVO_HOMOL_URL=http://189.45.140.206/api/novo/homol
NOVO_HOMOL_API_USER=homol
NOVO_HOMOL_API_PASS=********
```

> Sempre valide a URL e credenciais do upstream diretamente com `curl` antes de subir o painel.

Exemplo de teste:

```bash
curl -X GET "http://189.45.140.206/api/novo/homol/prompts/router.md" \
  -u homol:********
```

---

## 4. Docker / Traefik (produção)

### 4.1. `docker-compose.prod.yml` – backend

O serviço `backend` precisa expor as envs para o container, usando os mesmos nomes que o código Go lê:

```yaml
backend:
  environment:
    - PORT=9002
    - ENVIRONMENT=production
    # ...
    # DIGESAC
    - DIGESAC_API_BASE=${DIGESAC_API_BASE}
    - DIGESAC_API_USER=${DIGESAC_API_USER}
    - DIGESAC_API_PASS=${DIGESAC_API_PASS}
    # BIA HOMOL
    - BIA_HOMOL_URL=${BIA_HOMOL_URL}
    - BIA_HOMOL_API_USER=${BIA_HOMOL_API_USER}
    - BIA_HOMOL_API_PASS=${BIA_HOMOL_API_PASS}
    # NOVO HOMOL (exemplo)
    - NOVO_HOMOL_URL=${NOVO_HOMOL_URL}
    - NOVO_HOMOL_API_USER=${NOVO_HOMOL_API_USER}
    - NOVO_HOMOL_API_PASS=${NOVO_HOMOL_API_PASS}
```

> **Atenção:** o nome à esquerda (por ex. `NOVO_HOMOL_URL`) deve ser exatamente o que o `os.Getenv()` usa no código Go.

### 4.2. Roteamento Traefik

O Traefik já está configurado para:

- Enviar **tudo que começa com `/bestdoctors/`, `/auth/` e `/admin/`** para o backend Go.
- Enviar o host `chat.setuptecnologia.com.br` (sem prefixo especial) para o frontend Vue.

Por isso, use sempre prefixos de API do tipo:

- `/bestdoctors/<system>/homol/...`

Em vez de reutilizar diretamente o mesmo prefixo da SPA de outro app (por exemplo, não usar `/bia/homol/...` se já existe um app WebChat Bia nesse path).

---

## 5. Checklist rápido para um novo system

1. **Backend**
   - [ ] Criar `routes/<system>_proxy.go` com `*_URL`, `*_API_USER`, `*_API_PASS`.
   - [ ] Registrar `ServeMux` e `mux.Handle("/bestdoctors/<system>/homol/", ...)` em `cmd/main.go`.
2. **Frontend**
   - [ ] Criar `src/api/<system>.js` com `API_BASE = '/bestdoctors/<system>/homol'`.
   - [ ] Criar view `src/views/<SystemHomol>.vue` reutilizando a tela de prompts.
   - [ ] Adicionar rota em `src/router/index.js` com `requiresSystem: '<system>_homol'`.
   - [ ] (Opcional) Ajustar redirecionamentos de login para incluir o novo system.
3. **Infra**
   - [ ] Adicionar variáveis no `.env`: `<SYSTEM>_HOMOL_URL`, `<SYSTEM>_HOMOL_API_USER`, `<SYSTEM>_HOMOL_API_PASS`.
   - [ ] Referenciar essas variáveis no `docker-compose.prod.yml` do `backend`.
   - [ ] Testar o upstream real com `curl`.
   - [ ] Buildar e subir `backend` e `frontend` em produção.
4. **Validação final**
   - [ ] Acessar `https://chat.setuptecnologia.com.br/<system>-homol`.
   - [ ] Conferir no DevTools → Network que as chamadas vão para `/bestdoctors/<system>/homol/...` e retornam JSON, não HTML.

