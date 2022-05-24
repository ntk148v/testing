import asyncio
import httpx
import time


# Synchronous requests
start_time = time.time()
def run_sync():
    client = httpx.Client(verify=False)

    for number in range(1, 151):
        url = f'https://pokeapi.co/api/v2/pokemon/{number}'
        resp = client.get(url)
        pokemon = resp.json()
        print(pokemon['name'])

run_sync()
print("--- Synchronous requests: %s seconds ---" % (time.time() - start_time))


# Asynchronous requests
start_time = time.time()


async def get_pokemon(client, url):
        resp = await client.get(url)
        pokemon = resp.json()

        return pokemon['name']


async def run_async():
    async with httpx.AsyncClient(verify=False) as client:
        tasks = []
        for number in range(1, 151):
            url = f'https://pokeapi.co/api/v2/pokemon/{number}'
            tasks.append(asyncio.ensure_future(get_pokemon(client, url)))

        original_pokemon = await asyncio.gather(*tasks)
        for pokemon in original_pokemon:
            print(pokemon)

asyncio.run(run_async())
print("--- Asynchronous requests: %s seconds ---" % (time.time() - start_time))
