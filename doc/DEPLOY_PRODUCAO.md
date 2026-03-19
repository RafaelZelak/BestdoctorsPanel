# Guia de Deploy em Produção - BestDoctors Panel

## 📋 Pré-requisitos

### No Servidor AWS VPS

- Ubuntu 20.04 LTS ou superior
- Mínimo: 2 CPU cores, 4GB RAM, 20GB SSD
- Docker e Docker Compose instalados
- Acesso SSH configurado (apenas chave, sem senha)

### Registro e Configuração

- **Domínio registrado** (ex: `bestdoctors.com.br`)
- **DNS configurado** - Apontando para o IP do servidor AWS
- **Portas abertas** no AWS Security Group: 22 (SSH), 80 (HTTP), 443 (HTTPS)

---

## 🚀 Passo a Passo de Deploy

### 1. Preparar o Servidor AWS

```bash
# Conectar ao servidor via SSH
ssh -i sua-chave.pem ubuntu@SEU_IP_AWS

# Atualizar sistema
sudo apt update && sudo apt upgrade -y

# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ubuntu

# Instalar Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Relogar para aplicar permissões do Docker
exit
ssh -i sua-chave.pem ubuntu@SEU_IP_AWS
```

### 2. Configurar Firewall (UFW)

```bash
# Habilitar firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

### 3. Clonar o Projeto no Servidor

```bash
# Criar diretório para aplicação
mkdir -p ~/apps
cd ~/apps

# Clonar repositório (ou fazer upload via SCP)
git clone https://github.com/seu-usuario/bestdoctors_panel.git
cd bestdoctors_panel
```

### 4. Gerar Senhas de Produção

```bash
# Tornar script executável
chmod +x scripts/generate-secrets.sh

# Gerar senhas aleatórias e criar .env.production
./scripts/generate-secrets.sh
```

**⚠️ IMPORTANTE:** Anote as senhas geradas! Guarde em um local seguro (gerenciador de senhas).

### 5. Configurar .env.production

```bash
# Editar arquivo de produção
nano .env.production
```

**Configurações OBRIGATÓRIAS que você DEVE alterar:**

```bash
# Seu domínio (SEM https://)
TRAEFIK_DOMAIN=bestdoctors.com.br

# Banco de dados (se externo)
PG_HOST=seu-db.rds.amazonaws.com  # ou IP do servidor de banco
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=SENHA_GERADA_PELO_SCRIPT  # Já configurado
PG_DATABASE=bestdoctors_prod
PG_SSLMODE=require  # IMPORTANTE: 'require' em produção!

# Twilio (para WhatsApp)
TWILIO_URL=https://api.twilio.com/2010-04-01/Accounts/SEU_ACCOUNT_SID/Messages.json
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=seu_token_aqui
```

**Configurações que JÁ ESTÃO PREENCHIDAS (geradas automaticamente):**

- `SUPERADMIN_PASSWORD` - Já gerado
- `REDIS_PASSWORD` - Já gerado

### 6. Configurar Email para SSL (Let's Encrypt)

```bash
# Editar configuração do Traefik
nano traefik/traefik.yml
```

Altere a linha:

```yaml
email: admin@example.com # TODO: Replace with your email
```

Para:

```yaml
email: seu-email@gmail.com # Email REAL para notificações SSL
```

### 7. Verificar Configuração do DNS

**ANTES de fazer deploy, confirme que o DNS está configurado:**

```bash
# Testar se domínio aponta para o servidor
ping bestdoctors.com.br

# Deve retornar o IP do seu servidor AWS
# Se não retornar, aguarde propagação do DNS (até 48h, geralmente 1-2h)
```

### 8. Fazer Deploy em Produção

```bash
# Tornar script executável
chmod +x start.prod.sh

# Executar deploy
./start.prod.sh
```

**O script vai:**

1. ✅ Validar todas as variáveis de ambiente
2. ✅ Verificar se portas 80/443 estão livres
3. ✅ Criar diretório para certificados SSL
4. ✅ Construir imagens Docker otimizadas
5. ✅ Iniciar serviços
6. ✅ Aguardar certificado SSL (pode levar 2-5 minutos)
7. ✅ Verificar saúde dos serviços

### 9. Acompanhar os Logs

Em outro terminal, acompanhe o processo:

```bash
# Ver logs de todos os serviços
docker-compose -f docker-compose.prod.yml logs -f

# Ver apenas Traefik (para verificar SSL)
docker-compose -f docker-compose.prod.yml logs -f traefik

# Ver apenas backend
docker-compose -f docker-compose.prod.yml logs -f backend
```

**O que observar:**

- `"Certificate obtained"` - Certificado SSL obtido ✅
- `"Server listening on :9002"` - Backend iniciado ✅
- Sem erros de conexão com banco de dados ✅

---

## 🔒 Verificações Pós-Deploy

### 1. Testar HTTPS

Acesse no navegador:

- `https://bestdoctors.com.br` - Deve carregar com HTTPS (cadeado verde)
- `https://bestdoctors.com.br/admin` - Painel administrativo

### 2. Verificar Certificado SSL

Teste em: https://www.ssllabs.com/ssltest/

- Digite: `bestdoctors.com.br`
- Aguarde análise (2-3 minutos)
- **Objetivo:** Nota A ou A+

### 3. Testar Login Administrativo

1. Acesse: `https://bestdoctors.com.br/admin`
2. Usuário: `admin`
3. Senha: A gerada pelo script `generate-secrets.sh`

### 4. Verificar Headers de Segurança

Teste em: https://securityheaders.com

- Digite: `bestdoctors.com.br`
- **Objetivo:** Nota A

---

## 🔧 Configuração de Backups Automáticos

### Configurar Backup Diário

```bash
# Tornar script executável
chmod +x scripts/backup.sh

# Testar backup manual
./scripts/backup.sh

# Configurar cron para backup diário às 2h da manhã
crontab -e
```

Adicione esta linha:

```
0 2 * * * cd /home/ubuntu/apps/bestdoctors_panel && ./scripts/backup.sh >> /var/log/bestdoctors_backup.log 2>&1
```

### Configurar Upload para S3 (Opcional)

Se quiser enviar backups para AWS S3:

```bash
# Instalar AWS CLI
sudo apt install awscli -y

# Configurar credenciais
aws configure
```

Editar `.env.production` e adicionar:

```bash
AWS_ACCESS_KEY_ID=sua_key
AWS_SECRET_ACCESS_KEY=seu_secret
AWS_S3_BUCKET=bestdoctors-backups
AWS_REGION=us-east-1
```

---

## 📊 Monitoramento

### Comandos Úteis

```bash
# Status dos containers
docker-compose -f docker-compose.prod.yml ps

# Verificar uso de recursos
docker stats

# Logs em tempo real
docker-compose -f docker-compose.prod.yml logs -f

# Reiniciar serviço específico
docker-compose -f docker-compose.prod.yml restart backend

# Parar tudo
docker-compose -f docker-compose.prod.yml down

# Iniciar novamente
docker-compose -f docker-compose.prod.yml up -d
```

### Verificar Saúde do Sistema

```bash
# Testar endpoint de saúde
curl https://bestdoctors.com.br/health

# Deve retornar: OK
```

---

## 🔄 Atualizar Aplicação (Deploy de Nova Versão)

```bash
# 1. Fazer backup antes de atualizar
./scripts/backup.sh

# 2. Baixar nova versão
git pull origin main

# 3. Rebuild e restart (zero downtime)
docker-compose -f docker-compose.prod.yml up -d --build --no-deps backend
docker-compose -f docker-compose.prod.yml up -d --build --no-deps frontend

# 4. Verificar se atualizou
docker-compose -f docker-compose.prod.yml logs -f backend
```

---

## 🛡️ Segurança Adicional Recomendada

### 1. Fail2Ban (Proteção contra Brute Force SSH)

```bash
sudo apt install fail2ban -y
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

### 2. Atualização Automática de Segurança

```bash
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure --priority=low unattended-upgrades
```

### 3. Rotação de Senhas (A cada 90 dias)

```bash
# Gerar novas senhas
./scripts/generate-secrets.sh

# Atualizar .env.production
nano .env.production

# Restart aplicação
docker-compose -f docker-compose.prod.yml restart
```

---

## 🚨 Troubleshooting

### Certificado SSL Não Foi Obtido

```bash
# Verificar logs do Traefik
docker-compose -f docker-compose.prod.yml logs traefik

# Causas comuns:
# 1. DNS não está apontando para o servidor
# 2. Portas 80/443 bloqueadas no firewall
# 3. Email inválido em traefik.yml

# Forçar renovação
rm traefik/acme.json
touch traefik/acme.json
chmod 600 traefik/acme.json
docker-compose -f docker-compose.prod.yml restart traefik
```

### Serviço Não Inicia

```bash
# Ver erros específicos
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml logs frontend

# Verificar se variáveis de ambiente estão corretas
docker-compose -f docker-compose.prod.yml config
```

### Erro de Conexão com Banco de Dados

```bash
# Testar conexão manual
docker run -it --rm postgres:alpine psql -h SEU_PG_HOST -U postgres -d bestdoctors_prod

# Verificar:
# 1. PG_HOST está correto
# 2. Firewall do banco permite conexão do servidor
# 3. Credenciais estão corretas
```

---

## ✅ Checklist Final

Antes de considerar 100% em produção:

- [ ] Domínio aponta para servidor AWS
- [ ] SSL certificado obtido (cadeado verde no navegador)
- [ ] Nota A+ em SSL Labs
- [ ] Nota A em SecurityHeaders.com
- [ ] Login administrativo funciona
- [ ] Backup automático configurado (cron)
- [ ] Fail2ban instalado
- [ ] Firewall configurado (apenas 22, 80, 443)
- [ ] Senhas de desenvolvimento alteradas
- [ ] Logs monitorados (sem erros críticos)
- [ ] Email de notificação SSL configurado
- [ ] Documentação de procedures criada
- [ ] Plano de disaster recovery testado

---

## 📞 Suporte Pós-Deploy

### Logs Importantes

```bash
# Logs do sistema
sudo journalctl -u docker -f

# Logs da aplicação
docker-compose -f docker-compose.prod.yml logs -f

# Logs do Nginx (frontend)
docker-compose -f docker-compose.prod.yml exec frontend tail -f /var/log/nginx/access.log
```

### Contatos de Emergência

- **Let's Encrypt Status**: https://letsencrypt.status.io
- **Docker Status**: https://status.docker.com
- **AWS Status**: https://status.aws.amazon.com

---

## 🎯 Resumo Executivo

**Tempo estimado total:** 2-4 horas (incluindo propagação DNS)

**Arquivos que você VAI EDITAR:**

1. `.env.production` - Domínio e credenciais
2. `traefik/traefik.yml` - Email para SSL

**Comandos principais:**

```bash
./scripts/generate-secrets.sh  # Gerar senhas
nano .env.production           # Configurar domínio
nano traefik/traefik.yml       # Configurar email
./start.prod.sh                # Deploy!
```

**Após deploy:**

- Aplicação disponível em `https://seu-dominio.com.br`
- SSL automático (Let's Encrypt)
- Backups configurados
- Monitoramento ativo
- Segurança hardened
