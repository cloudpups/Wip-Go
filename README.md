<div align="center">

# WIP-go

A GitHub App for Do-Not-Merge and Work-In-Progress functionality! (...which is different than Draft PRs)

🚧 This is a work in progress!!

[![MediumLink](https://img.shields.io/badge/Read%20about%20me%20on%20-Medium-lightgrey?style=flat-square)][medium] [![CodeQL](https://github.com/cloudpups/Wip-Go/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/cloudpups/Wip-Go/actions/workflows/github-code-scanning/codeql)

</div>

> **❗ For those following the [Medium Article][medium]**
> 
> If you are coming here from the [Medium article][medium], note that you should use the source from the [Source for Medium Article](https://github.com/cloudpups/Wip-Go/releases/tag/medium_article) Release to minimize any issues. The app, and the permissions it requires, has changed since it was first written.

## What does this do?

Do you have a pull request that is not a draft, has passed review, and yet for some reason it **still** should not be merged? Look no further, the WIP-go bot is here for you! Use the default labels supported by this bot (`"dnm", "do not merge", "do-not-merge", "wip", "work in progress", "work-in-progress"`), or supply a list of up to 15 custom labels for this bot to recognize (**not** case-sensitive)!

Whenever such a label is added to a pull request, this bot will add a failing status check. If this status check is `required` via branch protections, then that pull request will be blocked from being merged (baring no admin overrides, of course)!

## Why does this exist?

Mainly as an experiment, but one that is production worthy! There is another, more mature, WIP bot out there. I advise checking it out: https://github.com/marketplace/wip.

To learn more about this bot, check out the corresponding Medium article: [LINK](https://medium.com/@JoshuaTheMiller/go-go-github-apps-d1b4bb26812b).

## Running the App

This app requires you to configure a GitHub Application Registration *beforehand* as you will need to use the Private Key and App ID from your registration for this GitHub App to work! The necessary permissions and events for the GitHub App are documented in the [GitHub Application Manifest (app.yaml)](./app.yaml).

### Configuration

To run, the application must be configured. Such configuration can be set either via a `config.yaml` file present in the same directory as the executable (this can actually be configured), or by adding values to the processes's environment variables. 

To keep things simple, let's stick with using the `config.yaml` file:

```yaml
# This is all the configuration needed to run this bot!
github:
  # While this should NOT be required, it makes for easier following of the Medium article!
  v3_api_url: "https://api.github.com/"    
  app:
    # Note: integration_id is referring to the App ID of your GitHub App Registration
    integration_id: 000000
    # You can simply paste the contents of your private key as a multiline string!
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      .....content
      .....content
      .....content
      -----END RSA PRIVATE KEY-----
```

🔍 Note: all configurable options and their corresponding documentation can be found in the [app config file](./app.go).

### Starting the app

Assuming you are running this from the source code, the following steps are all that is needed:

1. Ensure that you have Go installed on your machine: https://go.dev/learn/
2. Configure your GitHub App Registration as per the #Configuration section.
3. Create a `config.yaml` file somewhere on your machine (as you will see, it *technically* does not need to be named `config.yaml`).     
    * ❗ Make sure to populate the file with the values from your GitHub App registration!
4. Run the following command: `go run . run --config-file ./path-to-your-config.yaml
    * By passing the `-config-file` argument, you can place the config file at an arbitrary location. You can even give it an arbitrary name! Just ensure that it is actually valid YAML.

You now have your own instance of the WIP-go bot!

If you are running your bot somewhere where GitHub cannot reach, consider using Smee.io to forward webhook events. Do note that you SHOULD NOT use Smee.io for production OR in the case where your repositories should not have any data exposed to the public (i.e., repository names). To learn more about how to use Smee.io with this app, refer to the #Forwarding-Webhook-Events section of this README.

### Forwarding Webhook Events

> It should be stressed that Smee.io **must not** be used for production use cases, nor for any situations where the repositories in question are sensitive in nature. Be responsible in your decision making process. ~ @JoshuaTheMiller

With that out of the way, getting started with Smee.io is fairly easy. Just follow these instructions: https://github.com/probot/smee-client

After it is installed, you can start forwarding events by simply running the following command:

```sh
smee --url https://smee.io/SomeRandomValues --target http://127.0.0.1:8080/api/github/hook
```

* `url` is the URL of the Smee endpoint that was generated by Smee for you.
    * To generate a new URL (or *Channel*, as Smee.io calls them), go here and click *Start a new channel*: https://smee.io/
    * ❗ **IMPORTANT** you MUST update your GitHub App registration so that the `Webhook URL` property is pointing towards this URL. Otherwise the events will never make it to your app! Again- do not use Smee.io for production use cases.
* `target` is the URL (including protocol, port, and path) for your event handling endpoint.

[medium]: https://medium.com/@JoshuaTheMiller/go-go-github-apps-d1b4bb26812b
