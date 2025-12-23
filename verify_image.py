import requests
import os

BASE_URL = "http://localhost:3000/api"
EMAIL = "test_image_user@example.com"
PASSWORD = "password123"

def get_token():
    # Try login first
    try:
        resp = requests.post(f"{BASE_URL}/users/login", json={"email": EMAIL, "password": PASSWORD})
        if resp.status_code == 200:
            print("Logged in successfully.")
            return resp.json()["data"]["token"]
    except Exception as e:
        print(f"Login failed: {e}")

    # Try signup
    print("User not found or login failed, trying signup...")
    resp = requests.post(f"{BASE_URL}/users", json={"username": "test_image_user", "email": EMAIL, "password": PASSWORD})
    if resp.status_code == 201:
        print("Signed up successfully.")
        return resp.json()["data"]["token"]
    elif resp.status_code == 200:
         # In case create returns 200 for existing user (based on my reading of handler)
         print("User existed (signup returned 200), logging in...")
         # Although the handler logic seemed to return user data without token for existing user?
         # "res := dto.UserResponse{...}" -> No token in that specific existing-check branch?
         # Wait, looking at CreateUser code:
         # if existingUser -> returns StatusOK, NO token in response.
         # So we MUST call login if signup returns 200 or fails.
         pass
    
    # Retry login
    resp = requests.post(f"{BASE_URL}/users/login", json={"email": EMAIL, "password": PASSWORD})
    if resp.status_code == 200:
        return resp.json()["data"]["token"]
    
    raise Exception(f"Could not get token. Status: {resp.status_code}, Body: {resp.text}")

def create_dummy_image():
    with open("test_image.jpg", "wb") as f:
        f.write(b"fake image data")

def cleanup_dummy_image():
    if os.path.exists("test_image.jpg"):
        os.remove("test_image.jpg")

def test_upload_delete(token):
    headers = {"Authorization": f"Bearer {token}"}
    
    # Upload
    print("Uploading image...")
    files = {"image": ("test_image.jpg", open("test_image.jpg", "rb"), "image/jpeg")}
    data = {"game_name": "Test Game 123"}
    resp = requests.post(f"{BASE_URL}/upload", headers=headers, files=files, data=data) 
    
    if resp.status_code != 201:
        print(f"Upload failed: {resp.status_code} {resp.text}")
        return False
    
    url = resp.json()["url"]
    print(f"Upload success! URL: {url}")
    
    # Extract filename from URL
    # URL format: /images/[filename]
    filename = url.split("/")[-1]
    
    # Verify file exists on server (we are on the "server" machine)
    if os.path.exists(f"./images/{filename}"):
        print(f"File verification passed: ./images/{filename} exists.")
    else:
        print(f"File verification FAILED: ./images/{filename} does not exist.")
        return False

    # Delete
    print(f"Deleting image {filename}...")
    resp = requests.delete(f"{BASE_URL}/images/{filename}", headers=headers)
    
    if resp.status_code != 200:
        print(f"Delete failed: {resp.status_code} {resp.text}")
        return False
        
    print("Delete success!")
    
    # Verify file is gone
    if not os.path.exists(f"./images/{filename}"):
        print("File deletion verification passed.")
    else:
        print("File deletion verification FAILED: File still exists.")
        return False
        
    return True

if __name__ == "__main__":
    try:
        create_dummy_image()
        token = get_token()
        if test_upload_delete(token):
            print("ALL TESTS PASSED")
        else:
            print("TESTS FAILED")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        cleanup_dummy_image()
