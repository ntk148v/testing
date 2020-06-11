import asyncio
import time
import datetime as dt
from slackclient import SlackClient

from config import BOT_ID, BOT_TOKEN


AT_BOT = '<@%s>' % BOT_ID
EXAMPLE_COMMAND = 'do'
slack_client = SlackClient(BOT_TOKEN)
READ_WEBSOCKET_DELAY = 1


def handle_command(command, channel):
    """
        Receives commands directed at the bot and determines if they
        are valid commands. If so, then acts on the commands. If not,
        returns back what it needs for clarification.
    """
    response = "Not sure what you mean. Use the *" + EXAMPLE_COMMAND + \
               "* command with numbers, delimited by spaces."
    if command.startswith(EXAMPLE_COMMAND):
        response = "Sure...write some more code then I can do that!"
    slack_client.api_call("chat.postMessage", channel=channel,
                          text=response, as_user=True)


def parse_slack_output(slack_rtm_output):
    """
        The Slack Real Time Messaging API is an events firehose.
        this parsing function returns None unless a message is
        directed at the Bot, based on its ID.
    """
    output_list = slack_rtm_output
    if output_list and len(output_list) > 0:
        for output in output_list:
            if output and 'text' in output and AT_BOT in output['text']:
                # return text after the @ mention, whitespace removed
                return output['text'].split(AT_BOT)[1].strip().lower(), \
                    output['channel']
    return None, None


def listen_sync():
    if slack_client.rtm_connect():
        print("StarterBot connected and running!")
        READ_WEBSOCKET_DELAY = 1  # 1 second delay between reading from firehose
        while True:
            command, channel = parse_slack_output(slack_client.rtm_read())
            if command and channel:
                handle_command(command, channel)
            time.sleep(READ_WEBSOCKET_DELAY)
    else:
        print("Connection failed. Invalid Slack token or bot ID?")


@asyncio.coroutine
def listen_async():
    yield from asyncio.sleep(1)
    try:
        slack_client.rtm_connect()
        info = slack_client.rtm_read()
        if len(info) > 0:
            if 'text' in info[0] and 'channel' in info[0]:
                print(info)
                resp = dt.datetime.strftime(dt.datetime.now(), '%H:%M:%S')
                channel = info[0]['channel']
                slack_client.rtm_send_message(channel, resp)
    except Exception as e:
        pass
    finally:
        asyncio.async(listen_async())


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    try:
        asyncio.async(listen_async())
        loop.run_forever()
    except KeyboardInterrupt:
        pass
    finally:
        print('Loop.close()')
        loop.close()
