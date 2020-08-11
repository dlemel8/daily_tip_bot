# Daily Tip Bot
This bot is designed to run on Heruku free plan account.

![Design](design.jpg?raw=true "Design")

Since free apps sleep automatically after 30 mins of inactivity, this bot contains 2 services:
1. Interactive service wakes up on Slack event or first command, serve user commands and schedule tips for them on a given hour of day.
2. Scheduled service wakes up some external scheduler event, and send tips to all relevant users that schedule tip for that hour.  

Separation between services is done using sub command. 

### Slack Setup
* Create new app
* Add Bot Token Scopes:
  - commands
  - chat:write
  - users:read
* Define Slash Commands:
  - Command: /get-tip
    - Request URL: will be filled later
    - Short Description: get an immediate tip from selected tips source
    - Usage Hint: [topic]
  - Command: /list-topics
    - Request URL: will be filled later
    - Short Description: list all tip topics
    - Usage Hint: 
  - Command: /schedule-tip
    - Request URL: will be filled later
    - Short Description: schedule an hour in day (0-23) to send a tip 
    - Usage Hint: [hour] [topic]
* Install the app in your workspace
 
### Environment Variables
* PORT - Local listening port (interactive service only). Commands are served under /slack endpoint.
* SLACK_SIGNING_SECRET - App Signing Secret (interactive service only)
* SLACK_BOT_TOKEN - App Bot User OAuth Access Token (starts with xoxb)
* DATABASE_URL - PostgreSql connection string

### Local Development & Testing
* Run testing DB and ngrok using `docker-compose up -d`
* Open ngrok management website (local port 4040) and copy tunnel http url to all Slash commands Request URL field.
* Run interactive service using `daily_tip_bot` or `daily_tip_bot web_server`, where:
  - PORT is set to 8080
  - DATABASE_URL is set to postgres://test:test123@localhost:5432/scheduled_tips?sslmode\=disable
* Run scheduled service using `daily_tip_bot scheduled_tips_sender` (use same DATABASE_URL)

### Heroku Deployment
* Create new app
* Configure resources:
  - Heroku Postgres
    - Plan : hobby-dev
  - Heroku Scheduler
    - Schedule: Every 10 minutes
    -  Run Command: `bin/daily_tip_bot scheduled_tips_sender`
* Connect to Github and set Automatic deploys
* Set all env vars under Config Vars (PORT is automatically set by Heroku)
* Copy app url to all Slash commands Request URL field.