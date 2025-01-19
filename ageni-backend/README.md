## project introduction
### REVOX Studio is a cutting-edge low-code development platform at the core of the REVOX ecosystem, designed to bridge the gap between artificial intelligence (AI) and Web3. It provides a powerful and user-friendly environment for developers to build, deploy, and manage decentralized AI applications. With a rich library of plugins and an integrated agent marketplace, REVOX Studio empowers developers to create sophisticated AI agents and applications that seamlessly interact with Web3 protocols, enabling groundbreaking use cases.
#### Key Features:
1. **Low-Code Development**: Simplifies the creation of AI-Web3 applications, making it accessible to developers of all levels.
2. **Agent Marketplace**: Facilitates the sharing and monetization of AI agents, fostering collaboration and innovation.
3. **Integration with Web3**: Supports AI agents capable of executing on-chain transactions and interacting with smart contracts.
4. **Rich Plugin Ecosystem**: Offers pre-built modules for enhanced functionality and faster development cycles.
####  Core Benefits:
- **AI and Web3 Connectivity**: Enables AI to participate in decentralized transactions and utilize blockchain-based resources.
- **Developer Empowerment**: Provides tools to experiment and scale applications with ease, unlocking the potential of decentralized AI.
- **Ecosystem Growth**: Contributes to the broader REVOX mission of creating a unified AI-Web3 ecosystem.
## Now we are focusing on open-source studio products!
## Local deployment
### Preconditions
```
Mysql8.0 or above
Go 18.0 or above
```
### Create and switch a directory to store studit backend
```
mkdir studit-backend & cd studit-backend
```
### Download this front-end project
```
git clone  https://github.com/readonme/open-studio-backend.git
```
### Create and import database data
```
##Execute the command to create a database
create database studio;
#Import data
mysql -u root -p -h 127.0.0.1 studio < init_db.sql
```
### Modify configuration
```
#Copy a prod.toml file in the conf directory, named local. toml.
cp -rf conf/prod.toml conf/local.toml
```
### Replace your own openAI key, database account password in local.toml JWT token key，
```
vim conf/local.toml
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
# Fill in your own MySQL address 和 account password
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
