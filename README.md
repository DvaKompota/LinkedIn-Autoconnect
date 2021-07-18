# LinkedIn Autoconnect application
Automatically expand your LinkedIn network by connecting with people that work in the companies you are interested in
<br/>

## Installation and running instructions:
#### 1. To install the application on your local machine:
- clone the repo from GitHub to your local machine
- run setup.sh file in your local repo directory:<br/>
`./setup.sh`<br/>
- create a `credentials.py` file in the `data` directory a put your LinkedIn credentials there as follows:<br/>
`email = "your.email@gmail.com"`<br/>
`password = "y0uR_pA$$worD"`
#### 2. Running scripts:
- to revoke previously sent invites, older than X (default: 1 month) run:<br/>
`./modules/withdraw_old_invites.py`
- to send invites using LinkedIn search by companies from `config.py` run:<br/>
`./modules/subscribe_from_search.py`
- if you ran out of searches, you still can send invites from profile pages of your 1st circle contacts, who work in companies from `config.py`:<br/>
`./modules/subscribe_from_profiles.py`
<br/>

---
## Things to be added:
#### 1. Single application runner file with simple text interface
#### 2. Proper docstrings, in-line comments, and configuration documentation to provide easier transition for further development
