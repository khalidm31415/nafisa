import json
import requests

with open('generate_sample/signup_data_male.json') as f:
    signup_data_males = json.load(f)

with open('generate_sample/signup_data_female.json') as f:
    signup_data_females = json.load(f)

for s in signup_data_males + signup_data_females:
    selfie_with_id_card_url = f'https://image.com/selfie-idcard-{s["username"]}'
    photo_urls = [f'https://image.com/foto-{s["username"]}-{i}' for i in range(1,4)]
    s['selfieWithIDCardURL'] = selfie_with_id_card_url
    s['photoUrls'] = photo_urls
    s['password'] = 'kopinikmatnyamandilambung'
    payload = json.dumps(s)
    r = requests.post('http://localhost:3141/auth/signup', data=payload)

new_admin_data = {
    "username": "pak_admin",
    "password": "kopinikmatnyamandilambung",
    "isVerificationAdmin": True,
    "isDiscussionAdmin": True
}
payload = json.dumps(new_admin_data)
headers = {'Authorization': 'Bearer kopinikmatnyamandilambung'}
r = requests.post('http://localhost:3141/auth/new-admin', data=payload, headers=headers)


admin_login = {
    "username": "pak_admin",
    "password": "kopinikmatnyamandilambung",
}
headers = {'Content-Type': 'application/json'}
payload = json.dumps(admin_login)
r = requests.post('http://localhost:3141/auth/login', data=payload, headers=headers)

cookies = {'token': r.json()['token']}
r = requests.get('http://localhost:3141/verification/unverified-users', cookies=cookies)
unverified_users = r.json()

for u in unverified_users:
    body = {'userId': u['userId']}
    payload = json.dumps(body)
    r = requests.post('http://localhost:3141/verification/verify', data=payload, cookies=cookies)
