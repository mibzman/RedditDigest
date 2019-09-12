# Reddit Digest

## A simple tool to not browse reddit all the time

This tool sends some number of reddit posts from given subreddits at given intervals.

*   Hottest posts every day
*   Top weekly every week
*   Top monthly every month

### Getting started:

1.  Download & Build
2.  Fill out `sample.config.example`
3.  Rename the config file to something else.
4.  `./RedditDigest <your.config>`

Running the program will send an email to your email address with the configured subs. If the system has an acurate clock, weekly and monthly content will be included on the days specefied.

The intent is for the program to be run once a day, it will not self-schedule at this time

### Config

*   UserEmail: The email you want to recieve digests
*   RedditData: Info needed for `github.com/turnage/graw`, see [this tutorial](https://turnage.gitbooks.io/graw/content/chapter1.html)
*   EmailerConfig: Something to send emails, probably gmail
*   WeeklyWeekday: The day of the week on which to send the weekly portion of the digest
*   MonthlyDay: The day of the month on which to send the monthly portion of the digest

### ToDo:

-   Convert custom json file configs to .env
