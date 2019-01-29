# txtcrusher
### Installation
1. Create Pastebin account and obtain unique developer API key:
https://pastebin.com/api#1
2. Paste your developer API key to `$HOME/.config/config.json`
```json
{
    "pastebin": {
        "api_dev_key": "your_developer_key",
        "api_user_key": "your_user_key"
    }
}
```
3. (Optional) For some operation you need obtain user API key. You can generate it, following a link provided by pastebin.com:
https://pastebin.com/api/api_user_key.html
or using txtcrusher:
```bash
txtcrusher -k username password
```
And add it to the configuration file mentioned above.

### Usage
`txtcrusher` supports all provided Pastebin API.
##### Getting raw paste output of any "public" and "unlisted" pastes
Pastebin provides option to getting any "public" or "unlisted" paste even without `config.json` file.
```bash
txtcrusher -p PASTE_ID
```
##### Creating a new paste
You can upload paste with following:
```bash
txtcrusher [-guest] [-title TITLE] [-format FORMAT] [-expire EXPIREDATE] [-mod MODIFICATOR] -c TEXT
```
Flags with `[]` are optional. Info about them you can get there:
[about format](https://pastebin.com/api#5)
[about expire date](https://pastebin.com/api#6)
[about modificators](https://pastebin.com/api#7)
##### Deleting a paste created by a user
```bash
txtcrusher -d PASTE_ID
```
##### Listing pastes created by a user
```bash
txtcrusher -l RESULT_LIMIT
RESULT_LIMIT must be in range, where min value is 1, max value is 1000
```
##### Getting a user information and settings
```bash
txtcrusher -i
```
##### Getting raw paste output of users pastes including "private" pastes
```bash
txtcrusher -g PASTE_ID
```
