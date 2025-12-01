# BestDoctors Panel - Independent Docker Setup

Sistema de painel BestDoctors funcionando de forma independente com Docker Compose.

## 🚀 Quick Start

### Pré-requisitos
- Docker
- Docker Compose
- Arquivo `.env` configurado (veja seção abaixo)

### Iniciar o Sistema (Produção)

```bash
# Build e start de todos os serviços
docker-compose up --build -d

# Ver logs
docker-compose logs -f

# Parar os serviços
docker-compose down
```

Acesse o frontend em: **http://localhost**

### Modo Desenvolvimento

Para desenvolvimento com hot reload:

```bash
# Start em modo desenvolvimento
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Frontend estará em: http://localhost:5173
# Backend em: http://localhost:9002
```

## 📦 Arquitetura

### Serviços

1. **Backend** (Go)
   - Porta: `9002`
   - API REST com 8 endpoints
   - Conexão com PostgreSQL/Supabase
   - Health check: `/health`

2. **Frontend** (Vue 3 + Vite)
   - Porta: `80` (produção) ou `5173` (dev)
   - SPA com nginx em produção
   - Proxy para backend via nginx/Vite

### Networking

Os serviços se comunicam através da rede Docker `bestdoctors-network`:
- Frontend → Backend: `http://backend:9002`
- Health checks automáticos
- Restart automático em caso de falha

## ⚙️ Configuração

### Arquivo `.env`

O arquivo `.env` na raiz do projeto já contém todas as variáveis necessárias:

**Variáveis Obrigatórias para Backend:**
```env
# Supabase/PostgreSQL (OBRIGATÓRIO)
SUPRABASE_PGSQL=postgresql://user:pass@host:5432/db
SUPABASE_URL=http://your-supabase-url
SUPABASE_KEY=your-supabase-key
SUPRABASE_DB_HOST=host
SUPRABASE_DB_PORT=5432
SUPRABASE_DB_NAME=postgres
SUPRABASE_DB_USER=postgres
SUPRABASE_DB_PASSWORD=password
```

**Variáveis Opcionais:**
```env
# Cache local (opcional)
DB_NAME=gateway_db
DB_USER=gateway
DB_PASSWORD=password
DB_HOST=192.168.15.220
DB_PORT=5432
DB_SCHEMA=gateway_schema

# Twilio (opcional)
TWILIO_ACCOUNT_SID=your-sid
TWILIO_AUTH_TOKEN=your-token
```

## 🔧 Comandos Úteis

```bash
# Rebuild apenas um serviço
docker-compose up --build backend

# Ver status dos serviços
docker-compose ps

# Ver logs de um serviço específico
docker-compose logs -f frontend

# Executar comando no container
docker-compose exec backend sh

# Remover tudo (containers, networks, volumes)
docker-compose down -v

# Rebuild completo (sem cache)
docker-compose build --no-cache
```

## 🏗️ Estrutura do Projeto

```
bestdoctors_panel/
├── backend/
│   ├── cmd/
│   │   └── main.go          # Entry point com CORS
│   ├── internal/
│   │   └── db/
│   │       └── db.go        # Database connections
│   ├── routes/              # 8 API handlers
│   ├── models/              # Data models
│   ├── Dockerfile           # Multi-stage build
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── api/            # API client
│   │   ├── components/     # Vue components
│   │   └── composables/    # Vue composables
│   ├── Dockerfile          # Production (nginx)
│   ├── Dockerfile.dev      # Development (Vite)
│   ├── nginx.conf          # Nginx config
│   ├── vite.config.js      # Vite config
│   └── package.json
├── docker-compose.yml      # Produção
├── docker-compose.dev.yml  # Desenvolvimento
├── .env                    # Variáveis de ambiente
└── README.md
```

## 🔍 Endpoints da API

Todos os endpoints estão disponíveis em `http://localhost:9002/bestdoctors/`:

- `GET /bestdoctors/sessionphone` - Lista sessões
- `PATCH /bestdoctors/sessionphone/active` - Toggle AI
- `GET /bestdoctors/chathistory` - Histórico de chat
- `GET /bestdoctors/sessiondelta` - Sessões com novas mensagens
- `GET /bestdoctors/metrics/session` - Métricas de sessão
- `GET /bestdoctors/metrics/abandonment` - Taxa de abandono
- `GET /bestdoctors/metrics/flowdepth` - Profundidade do fluxo
- `GET /bestdoctors/metrics/reengagement` - Taxa de reengajamento
- `POST /bestdoctors/sendmessage` - Enviar mensagem
- `POST /bestdoctors/report` - Gerar relatórios (JSON/CSV/PDF/XLSX)
- `GET /health` - Health check

## 🐛 Troubleshooting

### Backend não conecta ao banco

Verifique:
1. .env tem `SUPRABASE_PGSQL` configurado
2. Host do banco é acessível do container Docker
3. Logs: `docker-compose logs backend`

### Frontend não carrega

Verifique:
1. Backend está rodando: `docker-compose ps`
2. Sem erros CORS no console do browser
3. Nginx está servindo: `docker-compose logs frontend`

### Problemas de CORS

O backend já tem CORS habilitado para `*`. Se ainda houver problemas:
1. Limpe o cache do browser
2. Verifique se está acessando pela porta correta
3. Veja logs do backend para requests OPTIONS

## 📝 Notas Importantes

- ✅ **Sem JWT**: Este código não utiliza autenticação JWT
- ✅ **Independente**: Não depende de gateway externo
- ✅ **Base Path**: Removido `/template_bestdoctors/`, agora serve em `/`
- ✅ **CORS**: Habilitado para qualquer origem
- ✅ **Health Checks**: Ambos serviços têm health checks configurados

## 🎯 Próximos Passos

Para produção:
1. Configure variáveis de ambiente específicas para produção
2. Ajuste CORS para domínios específicos (em `backend/cmd/main.go`)
3. Configure SSL/TLS se necessário
4. Considere usar secrets do Docker para credenciais
5. Configure backup do banco de dados
