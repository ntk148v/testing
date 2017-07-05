import asyncio

import aiohttp

from config import DEBUG, TOKEN


async def call(method, data=None, token=TOKEN):
    """Slack API call"""
    with aiohttp.ClientSession() as session:
        form = aiohttp.FormData(data or {})
        form.add_field('token', token)
        async with session.post('https://slack.com/api/{}' . format(method),
                                data=form) as response:
            assert 200 == response.status, ('{} with {} failed.' \
                                            . format(method, data))
            return await response.json()


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.set_debug(DEBUG)
    response = loop.run_until_complete(call('auth.test'))
    loop.close()

    assert response['ok']
    print(response)
