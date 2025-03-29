# Price Service Telegram Bot

Prices are so unstable… You can't track them in every market.  
Trust this bot to do the job, and you'll forget about daily analysis.  

It integrates with the **Price Service Project**, making searches easier and more precise.  

---

**In brief:** *This service is a UI for the Price Service Project.*  

---

## Main Description  

The Telegram bot has several key functions:  

- **Best Price Search Mode**:  
  This mode helps you find products with the best prices across different markets.  

- **Favorites Mode**:  
  This mode allows you to manage your favorite products and interact with them (add and remove items).  

- **Tracking Mode**:  
  This mode enables time-based product searches, working as a notification system that provides search results every 24 hours.  

## How to Install and Start  

Since this service depends on **Price Service**, some additional steps are required to use the bot.  

For a simpler installation and setup process, refer to the [best-price-project-deployment](https://github.com/MaKcm14/best-price-project-deployment) guide.  

## Configuring the .env

To start the bot, create a **.env** file in the root of the source directory with the structure as in the **.env_example** 

or **copy the .env_default into the .env with your_bot_token**.

The **.env_default** already configured for using.

You can customize your .env files according to your infrastructure.

### Note
Some changes in the .env file can lead to the necessity of changes in the main docker-compose files of deployment.

## Examples  

![Creating the query in the best price mode](https://github.com/MaKcm14/price-service-tg-bot/blob/develop/docs/best-price-1.png) ![Getting the products](https://github.com/MaKcm14/price-service-tg-bot/blob/develop/docs/best-price-2.png) ![Getting the favorite products](https://github.com/MaKcm14/price-service-tg-bot/blob/develop/docs/favrorites-1.png)

### P.S.  

You need to obtain the **bot token** from [**@BotFather**](https://t.me/BotFather).  

## Technology Stack  

- PostgreSQL  
- Redis  
- Docker  
- Kafka
- Telegram API  
