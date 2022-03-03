from datetime import datetime, timedelta

import requests

base_url = "https://data.nba.net/prod/v1"
yesterday = (datetime.today() - timedelta(days=1)).strftime("%Y%m%d")
resp = requests.get(f'{base_url}/{yesterday}/scoreboard.json')
games = resp.json()['games']
gameIds = [g['gameId'] for g in games]

resp = requests.get(f'{base_url}/{yesterday}/{gameIds[0]}_boxscore.json')
first_game = resp.json()['basicGameData']
hTeam = first_game['hTeam']
vTeam = first_game['vTeam']
print(f"{hTeam['triCode']} {hTeam['score']}-{vTeam['score']} {vTeam['triCode']}")
