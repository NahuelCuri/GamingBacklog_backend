import requests
import json

try:
    response = requests.get('http://localhost:3000/api/games')
    games = response.json()['data']
    
    target_game = None
    for game in games:
        if game['title'] == 'Game with image':
            target_game = game
            # Don't break, we want the latest if duplicates exist (though ID implies order usually)
            # Actually, let's just find the one with the specific title.
    
    if target_game:
        print(f"Found game: {target_game['title']}")
        print(f"Cover URL: {target_game['cover_url']}")
        print(f"Full Data: {json.dumps(target_game, indent=2)}")
    else:
        print("Game 'Game with image' not found.")
        print("Available titles:")
        for game in games:
            print(f"- {game['title']}")

except Exception as e:
    print(f"Error: {e}")
