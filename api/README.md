# The API Directory

This directory contains files necessary to quickly and easily deploy WIP-go as a serverless function on Vercel.

For more information, check out this link: https://vercel.com/docs/functions/serverless-functions/runtimes/go

## Vercel Hook

`vercel_hook.go` is served at the following address: https://wip.cloudpups.dev/api/vercel_hook

## Necessary Environment Variables

* `app_id` = the App ID from the GitHub App Registration
* `webhook_secret_key` = the Webhook Secret from the GitHub App Registration
* `private_key` = a Private Key from the GitHub App Registration (simply copy+paste the content of the private key)