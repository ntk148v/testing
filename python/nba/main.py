from nba_api.stats.endpoints import playercareerstats

# Nikola Jokić
career = playercareerstats.PlayerCareerStats(player_id='203999')

# pandas data frames (optional: pip install pandas)
print(career.get_data_frames()[0])

# json
# print(career.get_json())

# dictionary
# print(career.get_dict())

from nba_api.live.nba.endpoints import scoreboard

# today's scoreboard
games = scoreboard.ScoreBoard()

# json
print(games.get_json())
