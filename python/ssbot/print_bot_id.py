from slackclient import SlackClient

from config import BOT_NAME, BOT_TOKEN


if __name__ == '__main__':
    api_call = SlackClient(BOT_TOKEN).api_call('users.list')
    if api_call.get('ok'):
        # retrieve all users so we can find our bot
        users = api_call.get('members')
        for user in users:
            if 'name' in user and user.get('name') == BOT_NAME:
                print('Bot ID for %s is %s' % (user['name'], user['id']))
    else:
        print('Failed')
