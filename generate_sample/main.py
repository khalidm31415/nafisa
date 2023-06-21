import json
import requests
import http

with open('generate_sample/signup_data_male.json') as f:
    signup_data_males = json.load(f)

with open('generate_sample/signup_data_female.json') as f:
    signup_data_females = json.load(f)

for s in signup_data_males + signup_data_females:
    credential = {
        "username": s["username"],
        "password": "kopinikmatnyamandilambung",
    }
    headers = {'Content-Type': 'application/json'}
    payload = json.dumps(credential)
    res = requests.post('http://localhost:3141/auth/signup', data=payload, headers=headers)

    payload = json.dumps(credential)
    sess = requests.Session()
    res = sess.post('http://localhost:3141/auth/login', data=payload, headers=headers)

    selfie_with_id_card_url = f'https://image.com/selfie-idcard-{s["username"]}'
    photo_urls = [f'https://image.com/foto-{s["username"]}-{i}' for i in range(1,4)]
    s['selfieWithIDCardURL'] = selfie_with_id_card_url
    s['photoUrls'] = photo_urls
    payload = json.dumps(s)
    jar = http.cookies.SimpleCookie(res.headers["Set-Cookie"])
    cookies = {'token':jar['token'].value}
    res = sess.post('http://localhost:3141/profile', data=payload, headers=headers, cookies=cookies)

new_admin_data = {
    "oauthGmail": "khalidmuhammad1100@gmail.com",
    "username": "admin",
    "password": "kopinikmatnyamandilambung",
    "isVerificationAdmin": True,
    "isDiscussionAdmin": True
}
payload = json.dumps(new_admin_data)
headers = {'Authorization': 'Bearer kopinikmatnyamandilambung'}
r = requests.post('http://localhost:3141/auth/new-admin', data=payload, headers=headers)


admin_login = {
    "username": "admin",
    "password": "kopinikmatnyamandilambung",
}
headers = {'Content-Type': 'application/json'}
payload = json.dumps(admin_login)
res = requests.post('http://localhost:3141/auth/login', data=payload, headers=headers)

jar = http.cookies.SimpleCookie(res.headers["Set-Cookie"])
cookies = {'token':jar['token'].value}
res = requests.get('http://localhost:3141/verification/unverified-users', cookies=cookies)
unverified_users = res.json()

for u in unverified_users:
    body = {'userId': u['userId']}
    payload = json.dumps(body)
    requests.post('http://localhost:3141/verification/verify', data=payload, cookies=cookies)
