import json


def get_colors(data):
    colors = {}
    for color in data['colors'].values():
        if color.startswith("#"):
            colors[color] = color
    for item in data['tokenColors']:
        for color in item['settings'].values():
            if color.startswith("#"):
                colors[color] = color
    return colors


with open('vscolors.json') as json_file:
    data = json.load(json_file)
    colors = get_colors(data)
    with open('output.json', 'w') as outfile:
        json.dump(colors, outfile)
