from slackclient import SlackClient


BOT_NAME = 'sworker'
SLACK_BOT_TOKEN = '<secret>'


if __name__ == '__main__':
    api_call = SlackClient(SLACK_BOT_TOKEN).api_call('users.list')
    if api_call.get('ok'):
        # retrieve all users so we can find our bot
        users = api_call.get('members')
        for user in users:
            if 'name' in user and user.get('name') == BOT_NAME:
                print('Bot ID for %s is %s' % (user['name'], user['id']))
    else:
        print('Failed')
