# Ageni: No-Code AI Agent Platform for DEFAI

Ageni is an open-source platform that empowers anyone to build and deploy AI agents without writing a single line of code. Designed for both simplicity and power, AgenOS combines intuitive no-code tools with advanced blockchain capabilities, enabling users to create agents that go beyond text responses to execute direct on-chain crypto transactions.

Ageni is proudly developed and maintained by REVOX, as part of our commitment to fostering innovation in AI and blockchain technologies.

![create_agent6_test](https://github.com/user-attachments/assets/e4a5ad19-eeab-4760-a24e-07590c69e0f4)

#### Key Features:
1. **No-Code Development**: Build, configure, and deploy AI agents through an intuitive, user-friendly interface.
2. **DEFAI Integration**: Leverage DeFi plugins to create agents capable of performing secure and automated on-chain transactions for trading, staking, lending, and more.
3. **Customizable Workflows**: Define detailed agent behaviors and logic to suit your unique needs.
4. **Open and Collaborative**: As an open-source project, Agenix invites contributions from developers worldwide, fostering innovation and a vibrant community.

## Local Backend Deployment
### Backend service prerequisites
```
Mysql8.0 or above
Go 1.22.0 or above
```
### Enter the backend service folder
```
cd ageni-backend
```
### Create and import database data
```
# Log in to your own MySQL database and execute the create database command
create database studio;
# exist mysql shell and initilize database from command line (replace the mysql credential to your own credential)
mysql -u root -p -h 127.0.0.1 studio < init_db.sql
```
### Modify configuration
```
# Copy a prod.toml file in the conf directory, named local.toml.
cp -rf conf/prod.toml conf/local.toml
```
### Replace your own openAI key, database account password and JWT token key in local.toml 
```
env="dev"
# Fill in your own OpenAI key
openai_key=""
bot_tabs=["Most used", "News", "Web3", "Search", "Meme", "Entertainment"]
plugin_tabs=["Most used", "Efficiency Tools", "Web3", "Entertainment", "Game", "Lifestyle", "Education"]
[log]
appName ="studio-service"
[http_server]
addr="0.0.0.0:8080"
readTimeout="1m"
writeTimeout="1m"
[studio_db]
debug_mode=true
driver_name="mysql"
# Fill in your own MySQL address and credentials
dsn="{user}:{password}@({ip}:{port})/{database}?charset=utf8mb4&parseTime=True&loc=Local"
max_idle_conns=100
max_open_conns=100
max_lifetime="60s"
max_idle_time="60s"
[jwt_token]
#Fill in your own jwt_token
secretKey=""
expire=2592000
```
### Start project
```
go run main.go
```
## Frontend Deployment
## Front end service prerequisites
```
node 18.17.0
```
### Switch to the front-end project folder
```
cd ../ageni-frontend
```
### Modify .env in ageni-frontend
```
# change VITE_PROXY to match the backend service, if you deploy it locally, it should be http://127.0.0.1:8080/
VITE_PROXY=http://127.0.0.1:8080/
# Add your WalletConnect Application ID here
VITE_WALLET_CONNECT_KEY=YOUR_WALLET_CONNECT_PRODUCT_ID
```
### (optional) If you want to customize the startup port, modify "port" in the vite.config.ts file

### Install dependencies
```
npm install
```
### Start project
```
npm run dev
```
