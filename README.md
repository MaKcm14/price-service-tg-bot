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
  This mode enables time-based product searches, working as a notification system that provides search results.  

## How to Install and Start  

Since this service depends on **Price Service**, some additional steps are required to use the bot.  

For a simpler installation and setup process, refer to the **best-price-project-deployment** guide.  

## Note  

To start the bot, create a **.env** file in the root of the source directory with the following structure:  

```
BOT_TOKEN="your_bot_token_here"

DSN="postgresql://user:pwd@ip:5432/dbName"

PRICE_SERVICE_SOCKET="socket"
```

## Examples  

(Provide some usage examples here)

### P.S.  

You need to obtain the **bot token** from [**@BotFather**](https://t.me/BotFather).  

## Technology Stack  

- PostgreSQL  
- Redis  
- Docker  
- Kafka *(in progress)*  
- Telegram API  
