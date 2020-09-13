import os

from pyrogram import Client, filters

api_id = os.getenv("API_ID")
api_hash = os.getenv("API_HASH")
bot_token = os.getenv("BOT_TOKEN")

app = Client(
    "my_bot",
    api_id=api_id,
    api_hash=api_hash,
    bot_token=bot_token
)

@app.on_message(filters.command(["start"]))
def start(client, message):
    message.reply_text("Hi!")

@app.on_message(filters.text & filters.private)
def echo(client, message):
    message.reply_text(message.text)

app.run()

